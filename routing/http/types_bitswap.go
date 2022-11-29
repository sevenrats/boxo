package http

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multibase"
)

type BitswapReadProviderResponse struct {
	Protocol string
	ID       *peer.ID
	Addrs    []Multiaddr
}

type BitswapWriteProviderRequest struct {
	Protocol  string
	Signature string

	// this content must be untouched because it is signed and we need to verify it
	RawPayload json.RawMessage                    `json:"Payload"`
	Payload    BitswapWriteProviderRequestPayload `json:"-"`
}

type BitswapWriteProviderRequestPayload struct {
	Keys        []CID
	Timestamp   *Time
	AdvisoryTTL *Duration
	ID          *peer.ID
	Addrs       []Multiaddr
}

type tmpBWPR BitswapWriteProviderRequest

func (p *BitswapWriteProviderRequest) UnmarshalJSON(b []byte) error {
	var bwp tmpBWPR
	err := json.Unmarshal(b, &bwp)
	if err != nil {
		return err
	}

	p.Protocol = bwp.Protocol
	p.Signature = bwp.Signature
	p.RawPayload = bwp.RawPayload

	return json.Unmarshal(bwp.RawPayload, &p.Payload)
}

func (p *BitswapWriteProviderRequest) IsSigned() bool {
	return p.Signature != ""
}

func (p *BitswapWriteProviderRequest) setRawPayload() error {
	payloadBytes, err := json.Marshal(p.Payload)
	if err != nil {
		return fmt.Errorf("marshaling bitswap write provider payload: %w", err)
	}

	p.RawPayload = payloadBytes

	return nil
}

func (p *BitswapWriteProviderRequest) Sign(peerID peer.ID, key crypto.PrivKey) error {
	if p.IsSigned() {
		return errors.New("already signed")
	}

	if key == nil {
		return errors.New("no key provided")
	}

	sid, err := peer.IDFromPrivateKey(key)
	if err != nil {
		return err
	}
	if sid != peerID {
		return errors.New("not the correct signing key")
	}

	err = p.setRawPayload()
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(p.RawPayload))
	sig, err := key.Sign(hash[:])
	if err != nil {
		return err
	}

	sigStr, err := multibase.Encode(multibase.Base64, sig)
	if err != nil {
		return fmt.Errorf("multibase-encoding signature: %w", err)
	}

	p.Signature = sigStr
	return nil
}

func (p *BitswapWriteProviderRequest) Verify() error {
	if !p.IsSigned() {
		return errors.New("not signed")
	}

	if p.Payload.ID == nil {
		return errors.New("peer ID must be specified")
	}

	// note that we only generate and set the payload if it hasn't already been set
	// to allow for passing through the payload untouched if it is already provided
	if p.RawPayload == nil {
		err := p.setRawPayload()
		if err != nil {
			return err
		}
	}

	pk, err := p.Payload.ID.ExtractPublicKey()
	if err != nil {
		return fmt.Errorf("extracing public key from peer ID: %w", err)
	}

	_, sigBytes, err := multibase.Decode(p.Signature)
	if err != nil {
		return fmt.Errorf("multibase-decoding signature to verify: %w", err)
	}

	hash := sha256.Sum256([]byte(p.RawPayload))
	ok, err := pk.Verify(hash[:], sigBytes)
	if err != nil {
		return fmt.Errorf("verifying hash with signature: %w", err)
	}
	if !ok {
		return errors.New("signature failed to verify")
	}

	return nil
}

type BitswapWriteProviderResponse struct {
	AdvisoryTTL time.Duration
}