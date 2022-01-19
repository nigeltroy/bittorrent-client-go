# bittorrent-client-go
A simple BitTorrent client in Go.

## Motivation
I wanted to learn Go and async programming!
An actual fully-featured BitTorrent client is a huge undertaking (see FOSS project [qBittorrent](https://github.com/qbittorrent/qBittorrent)), so this only aims to be an extremely simplified version of one.

## Project progress

### Components of program

- [x] CLI
- [ ] GTK UI
- [ ] API
    - [x] Torrent
        - [ ] URL
        - [x] File
        - [ ] Info hash/Magnet link
    - [x] Tracker
    - [ ] Peers
    - [ ] Client

### BitTorrent protocol spec fields
Information on the following fields can be found in the [BitTorrent protocol spec (or a spec change request)](https://www.bittorrent.org/beps/bep_0003.html) or in the spec in the [BitTorrent Wiki](https://wiki.theory.org/BitTorrentSpecification). A simplified version of this program only needs to support some non-optional fields.

#### Metainfo
- [x] info
    - [x] piece length
    - [x] pieces
    - [ ] private
    - Single file mode
        - [x] name
        - [x] length
        - [ ] md5sum
    - Multiple file mode
        - [x] name
        - [x] files
            - [x] length
            - [x] path
            - [ ] md5sum
- [x] announce
- [x] announce-list
- [x] creation date
- [x] comment
- [x] created by
- [ ] encoding
- [ ] url-list ??? - I found this while trying to debug the Arch torrent, but it isn't documented in the BitTorrent spec

#### Trackers
- [x] TCP announce
- [ ] UDP announce
- [ ] DHT announce
- [ ] Async announcing
- [ ] Reannouncing every <Response.interval> seconds

##### Request
- [x] info_hash
- [x] peer_id
- [ ] port
- [x] uploaded
- [x] downloaded
- [x] left
- [ ] compact
- [ ] no_peer_id
- [ ] event
- [ ] ip
- [ ] numwant
- [ ] key
- [ ] trackerid

##### Response
- [x] failure reason
- [ ] warning message
- [x] interval
- [ ] min interval
- [ ] tracker id
- [ ] complete
- [ ] incomplete
- [x] peers
    - [x] Dictionary model
        - [x] peer id
        - [x] ip
        - [x] port
    - [ ] Binary model

#### Peers
- [x] Handshake
