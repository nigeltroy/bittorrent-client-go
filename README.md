# bittorrent-client-go
A simple BitTorrent client in Go.

## Project progress

### Components of program

- [ ] CLI
- [ ] GTK UI
- [ ] API
    - [ ] Torrent
        - [ ] URL/file
        - [ ] Info hash/Magnet link
    - [ ] Tracker
    - [ ] Peers
    - [ ] Client
- [ ] HTTP API

### BitTorrent protocol spec fields
Information on the following fields can be found in the [BitTorrent protocol spec (or a spec change request)](https://www.bittorrent.org/beps/bep_0003.html) or in the spec in the [BitTorrent Wiki](https://wiki.theory.org/BitTorrentSpecification). A simplified version of this program only needs to support some non-optional fields.

#### Metainfo
- [ ] info
    - [ ] piece length
    - [ ] pieces
    - [ ] private
    - Single file mode
        - [ ] name
        - [ ] length
        - [ ] md5sum
    - Multiple file mode
        - [ ] name
        - [ ] files
            - [ ] length
            - [ ] path
            - [ ] md5sum
- [ ] announce
- [ ] announce-list
- [ ] creation date
- [ ] comment
- [ ] created by
- [ ] encoding

#### Trackers

##### Request
- [ ] info_hash
- [ ] peer_id
- [ ] port
- [ ] uploaded
- [ ] downloaded
- [ ] left
- [ ] compact
- [ ] no_peer_id
- [ ] event
- [ ] ip
- [ ] numwant
- [ ] key
- [ ] trackerid

##### Response
- [ ] failure reason
- [ ] warning message
- [ ] interval
- [ ] min interval
- [ ] tracker id
- [ ] complete
- [ ] incomplete
- [ ] peers
    - [ ] Dictionary model
        - [ ] peer id
        - [ ] ip
        - [ ] port
    - [ ] Binary model

#### Peers
