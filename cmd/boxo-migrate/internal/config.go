package migrate

import (
	"encoding/json"
	"fmt"
	"io"
)

type Config struct {
	ImportPaths map[string]string
	Modules     []string
}

var DefaultConfig = Config{
	ImportPaths: map[string]string{
		"github.com/ipfs/go-bitswap":                     "github.com/sevenrats/boxo/bitswap",
		"github.com/ipfs/go-ipfs-files":                  "github.com/sevenrats/boxo/files",
		"github.com/ipfs/tar-utils":                      "github.com/sevenrats/boxo/tar",
		"github.com/ipfs/interface-go-ipfs-core":         "github.com/sevenrats/boxo/coreiface",
		"github.com/ipfs/go-unixfs":                      "github.com/sevenrats/boxo/ipld/unixfs",
		"github.com/ipfs/go-pinning-service-http-client": "github.com/sevenrats/boxo/pinning/remote/client",
		"github.com/ipfs/go-path":                        "github.com/sevenrats/boxo/path",
		"github.com/ipfs/go-namesys":                     "github.com/sevenrats/boxo/namesys",
		"github.com/ipfs/go-mfs":                         "github.com/sevenrats/boxo/mfs",
		"github.com/ipfs/go-ipfs-provider":               "github.com/sevenrats/boxo/provider",
		"github.com/ipfs/go-ipfs-pinner":                 "github.com/sevenrats/boxo/pinning/pinner",
		"github.com/ipfs/go-ipfs-keystore":               "github.com/sevenrats/boxo/keystore",
		"github.com/ipfs/go-filestore":                   "github.com/sevenrats/boxo/filestore",
		"github.com/ipfs/go-ipns":                        "github.com/sevenrats/boxo/ipns",
		"github.com/ipfs/go-blockservice":                "github.com/sevenrats/boxo/blockservice",
		"github.com/ipfs/go-ipfs-chunker":                "github.com/sevenrats/boxo/chunker",
		"github.com/ipfs/go-fetcher":                     "github.com/sevenrats/boxo/fetcher",
		"github.com/ipfs/go-ipfs-blockstore":             "github.com/sevenrats/boxo/blockstore",
		"github.com/ipfs/go-ipfs-posinfo":                "github.com/sevenrats/boxo/filestore/posinfo",
		"github.com/ipfs/go-ipfs-util":                   "github.com/sevenrats/boxo/util",
		"github.com/ipfs/go-ipfs-ds-help":                "github.com/sevenrats/boxo/datastore/dshelp",
		"github.com/ipfs/go-verifcid":                    "github.com/sevenrats/boxo/verifcid",
		"github.com/ipfs/go-ipfs-exchange-offline":       "github.com/sevenrats/boxo/exchange/offline",
		"github.com/ipfs/go-ipfs-routing":                "github.com/sevenrats/boxo/routing",
		"github.com/ipfs/go-ipfs-exchange-interface":     "github.com/sevenrats/boxo/exchange",
		"github.com/ipfs/go-merkledag":                   "github.com/sevenrats/boxo/ipld/merkledag",
		"github.com/ipld/go-car":                         "github.com/sevenrats/boxo/ipld/car",

		// Pre Boxo rename
		"github.com/ipfs/go-libipfs/gateway":               "github.com/sevenrats/boxo/gateway",
		"github.com/ipfs/go-libipfs/bitswap":               "github.com/sevenrats/boxo/bitswap",
		"github.com/ipfs/go-libipfs/files":                 "github.com/sevenrats/boxo/files",
		"github.com/ipfs/go-libipfs/tar":                   "github.com/sevenrats/boxo/tar",
		"github.com/ipfs/go-libipfs/coreiface":             "github.com/sevenrats/boxo/coreiface",
		"github.com/ipfs/go-libipfs/unixfs":                "github.com/sevenrats/boxo/ipld/unixfs",
		"github.com/ipfs/go-libipfs/pinning/remote/client": "github.com/sevenrats/boxo/pinning/remote/client",
		"github.com/ipfs/go-libipfs/path":                  "github.com/sevenrats/boxo/path",
		"github.com/ipfs/go-libipfs/namesys":               "github.com/sevenrats/boxo/namesys",
		"github.com/ipfs/go-libipfs/mfs":                   "github.com/sevenrats/boxo/mfs",
		"github.com/ipfs/go-libipfs/provider":              "github.com/sevenrats/boxo/provider",
		"github.com/ipfs/go-libipfs/pinning/pinner":        "github.com/sevenrats/boxo/pinning/pinner",
		"github.com/ipfs/go-libipfs/keystore":              "github.com/sevenrats/boxo/keystore",
		"github.com/ipfs/go-libipfs/filestore":             "github.com/sevenrats/boxo/filestore",
		"github.com/ipfs/go-libipfs/ipns":                  "github.com/sevenrats/boxo/ipns",
		"github.com/ipfs/go-libipfs/blockservice":          "github.com/sevenrats/boxo/blockservice",
		"github.com/ipfs/go-libipfs/chunker":               "github.com/sevenrats/boxo/chunker",
		"github.com/ipfs/go-libipfs/fetcher":               "github.com/sevenrats/boxo/fetcher",
		"github.com/ipfs/go-libipfs/blockstore":            "github.com/sevenrats/boxo/blockstore",
		"github.com/ipfs/go-libipfs/filestore/posinfo":     "github.com/sevenrats/boxo/filestore/posinfo",
		"github.com/ipfs/go-libipfs/util":                  "github.com/sevenrats/boxo/util",
		"github.com/ipfs/go-libipfs/datastore/dshelp":      "github.com/sevenrats/boxo/datastore/dshelp",
		"github.com/ipfs/go-libipfs/verifcid":              "github.com/sevenrats/boxo/verifcid",
		"github.com/ipfs/go-libipfs/exchange/offline":      "github.com/sevenrats/boxo/exchange/offline",
		"github.com/ipfs/go-libipfs/routing":               "github.com/sevenrats/boxo/routing",
		"github.com/ipfs/go-libipfs/exchange":              "github.com/sevenrats/boxo/exchange",

		// Unmigrated things
		"github.com/ipfs/go-libipfs/blocks": "github.com/ipfs/go-block-format",
		"github.com/sevenrats/boxo/blocks":       "github.com/ipfs/go-block-format",
	},
	Modules: []string{
		"github.com/ipfs/go-bitswap",
		"github.com/ipfs/go-ipfs-files",
		"github.com/ipfs/tar-utils",
		"gihtub.com/ipfs/go-block-format",
		"github.com/ipfs/interface-go-ipfs-core",
		"github.com/ipfs/go-unixfs",
		"github.com/ipfs/go-pinning-service-http-client",
		"github.com/ipfs/go-path",
		"github.com/ipfs/go-namesys",
		"github.com/ipfs/go-mfs",
		"github.com/ipfs/go-ipfs-provider",
		"github.com/ipfs/go-ipfs-pinner",
		"github.com/ipfs/go-ipfs-keystore",
		"github.com/ipfs/go-filestore",
		"github.com/ipfs/go-ipns",
		"github.com/ipfs/go-blockservice",
		"github.com/ipfs/go-ipfs-chunker",
		"github.com/ipfs/go-fetcher",
		"github.com/ipfs/go-ipfs-blockstore",
		"github.com/ipfs/go-ipfs-posinfo",
		"github.com/ipfs/go-ipfs-util",
		"github.com/ipfs/go-ipfs-ds-help",
		"github.com/ipfs/go-verifcid",
		"github.com/ipfs/go-ipfs-exchange-offline",
		"github.com/ipfs/go-ipfs-routing",
		"github.com/ipfs/go-ipfs-exchange-interface",
		"github.com/ipfs/go-libipfs",
	},
}

func ReadConfig(r io.Reader) (Config, error) {
	var config Config
	err := json.NewDecoder(r).Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("reading and decoding config: %w", err)
	}
	return config, nil
}
