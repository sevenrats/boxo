package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sevenrats/boxo/blockservice"
	"github.com/sevenrats/boxo/blockstore"
	"github.com/sevenrats/boxo/examples/gateway/common"
	"github.com/sevenrats/boxo/gateway"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"
	"github.com/stretchr/testify/assert"
)

const (
	HelloWorldCID = "bafkreifzjut3te2nhyekklss27nh3k72ysco7y32koao5eei66wof36n5e"
)

func newProxyGateway(t *testing.T, rs *httptest.Server) *httptest.Server {
	blockStore := blockstore.NewBlockstore(dssync.MutexWrap(datastore.NewMapDatastore()))
	exch := newProxyExchange(rs.URL, nil)
	blockService := blockservice.New(blockStore, exch)
	routing := newProxyRouting(rs.URL, nil)

	gw, err := gateway.NewBlocksGateway(blockService, gateway.WithValueStore(routing))
	if err != nil {
		t.Error(err)
	}

	handler := common.NewHandler(gw)
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	return ts
}

func TestErrorOnInvalidContent(t *testing.T) {
	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("wrong data"))
	}))
	t.Cleanup(rs.Close)
	ts := newProxyGateway(t, rs)

	res, err := http.Get(ts.URL + "/ipfs/" + HelloWorldCID)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, res.StatusCode)
	assert.Contains(t, string(body), blocks.ErrWrongHash.Error())
}

func TestPassOnOnCorrectContent(t *testing.T) {
	rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	t.Cleanup(rs.Close)
	ts := newProxyGateway(t, rs)

	res, err := http.Get(ts.URL + "/ipfs/" + HelloWorldCID)
	assert.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.NoError(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, string(body), "hello world")
}

func TestTraceContext(t *testing.T) {
	doCheckRequest := func(t *testing.T, req *http.Request) {
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.EqualValues(t, string(body), "hello world")
	}

	const (
		traceVersion  = "00"
		traceID       = "4bf92f3577b34da6a3ce929d0e0e4736"
		traceParentID = "00f067aa0ba902b7"
		traceFlags    = "00"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := common.SetupTracing(ctx, "Proxy Test")
	assert.NoError(t, err)
	defer (func() { _ = tp.Shutdown(ctx) })()

	t.Run("Re-use Traceparent Trace ID Of Initial Request", func(t *testing.T) {
		rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The expected prefix for the traceparent header consists of the version and trace id.
			expectedPrefix := fmt.Sprintf("%s-%s-", traceVersion, traceID)
			if !strings.HasPrefix(r.Header.Get("traceparent"), expectedPrefix) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.Write([]byte("hello world"))
			}
		}))

		t.Cleanup(rs.Close)
		ts := newProxyGateway(t, rs)

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/ipfs/"+HelloWorldCID, nil)
		assert.NoError(t, err)
		req.Header.Set("Traceparent", fmt.Sprintf("%s-%s-%s-%s", traceVersion, traceID, traceParentID, traceFlags))
		doCheckRequest(t, req)
	})

	t.Run("Create New Trace ID If Not Given", func(t *testing.T) {
		rs := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In this request we are not sending a traceparent header, so a new one should be created.
			if r.Header.Get("traceparent") == "" {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.Write([]byte("hello world"))
			}
		}))

		t.Cleanup(rs.Close)
		ts := newProxyGateway(t, rs)

		req, err := http.NewRequest(http.MethodGet, ts.URL+"/ipfs/"+HelloWorldCID, nil)
		assert.NoError(t, err)
		doCheckRequest(t, req)
	})
}
