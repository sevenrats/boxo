# Changelog

All notable changes to this project will be documented in this file.

Note:
* The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
* This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).
( More notes about versioning and our release policies are [here](./RELEASE.md).

## Legend
The following emojis are used to highlight certain changes:
* 🛠 - BREAKING CHANGE.  Action is required if you use this functionality.
* ✨ - Noteworthy change to be aware of.

## [Unreleased]

- ✨ The gateway templates were updated to provide better features for users and gateway implementers:
  - New human-friendly error messages.
  - Updated, higher-definition icons in directory listings.
  - Customizable menu items next to "About IPFS" and "Install IPFS".
  - Valid DAG-CBOR and DAG-JSON blocks now provide a preview, where links can be followed.

## [0.8.0] - 2023-04-05
### Added

- ✨ Migrated repositories into Boxo (https://github.com/sevenrats/boxo/pull/220)
  - github.com/ipfs/interface-go-ipfs-core => ./coreiface
  - github.com/ipfs/go-pinning-service-http-client => ./pinning/remote/client
  - github.com/ipfs/go-path => ./path
  - github.com/ipfs/go-namesys => ./namesys
  - github.com/ipfs/go-mfs => ./mfs
  - github.com/ipfs/go-ipfs-provider => ./provider
  - github.com/ipfs/go-ipfs-pinner => ./pinning/pinner
  - github.com/ipfs/go-ipfs-keystore => ./keystore
  - github.com/ipfs/go-filestore => ./filestore
  - github.com/ipfs/go-ipns => ./ipns
  - github.com/ipfs/go-blockservice => ./blockservice
  - github.com/ipfs/go-ipfs-chunker => ./chunker
  - github.com/ipfs/go-fetcher => ./fetcher
  - github.com/ipfs/go-ipfs-blockstore => ./blockstore
  - github.com/ipfs/go-ipfs-posinfo => ./filestore/posinfo
  - github.com/ipfs/go-ipfs-util => ./util
  - github.com/ipfs/go-ipfs-ds-help => ./datastore/dshelp
  - github.com/ipfs/go-verifcid => ./verifcid
  - github.com/ipfs/go-ipfs-exchange-offline => ./exchange/offline
  - github.com/ipfs/go-ipfs-routing => ./routing
  - github.com/ipfs/go-ipfs-exchange-interface => ./exchange
  - github.com/ipfs/go-unixfs => ./ipld/unixfs
  - github.com/ipfs/go-merkledag => ./ipld/merkledag
  - github.com/ipld/go-car => ./ipld/car
- ✨ Added a migration tool to aid in migrating from the migrated repositories to Boxo, see the documentation here: https://github.com/sevenrats/boxo/blob/main/README.md#migrating-to-boxo (https://github.com/sevenrats/boxo/pull/226)
  - Added a check to ensure the migration tool is only run in a Git repository (with an optional override flag)
- ✨ Added tracing and metrics to the refactored gateway for its IPFS backend


### Changed

- Removed a mention of "bitswap" in blockservice debug logs
- Changed the Bitswap message package from "bitswap.message.pb" to "bitswap.message.v1.pb" to avoid protobuf panics due to duplicate registration with [go-bitswap](https://github.com/ipfs/go-bitswap) (https://github.com/sevenrats/boxo/pull/212)
- ✨ Remove a busyloop in blockservice getBlocks by removing batching when caching (https://github.com/sevenrats/boxo/pull/232)

### Deprecated

None

### Removed

None

### Fixed

- Ensure dag-cbor/json codecs are registered in the gateway handler (https://github.com/sevenrats/boxo/pull/223)
- ✨ Refactor the Gateway API to operate on higher level semantics (https://github.com/sevenrats/boxo/pull/176)
- Fixed a panic in the gateway handler when returning errors (https://github.com/sevenrats/boxo/pull/255)

### Security

None
