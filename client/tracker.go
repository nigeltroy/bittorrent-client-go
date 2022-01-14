package client

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/marksamman/bencode"
)

type peerDict struct {
	peerId string
	ip     string
	port   int64
}

type request struct {
	announce   string
	infoHash   string
	peerId     string
	uploaded   int64
	downloaded int64
	left       int64
	// compact
	// ip
	// port
	// noPeerId
	// event
	// numwant
	// key
	// trackerid
}

type response struct {
	failureReason string
	interval      int64
	peersDict     []peerDict // dictionary model
	// peersBinary   string     // binary model
	// warningMessage
	// minInterval
	// trackerId
	// complete
	// incomplete
}

type tracker struct {
	request  request
	response response
}

func (t *torrent) createRequest(
	announce string, infoHash string, uploaded int64, downloaded int64, left int64,
) (*request, error) {
	hasMultipleFiles := t.metainfo.info.hasMultipleFiles
	length := t.metainfo.info.length

	if hasMultipleFiles {
		return nil, errors.New("length requested from torrent with multiple files")
	} else if !hasMultipleFiles && length == 0 {
		return nil, errors.New("length has not been set or is equal to zero")
	}

	peerId := "-NG0001-"
	for i := 0; i < 12; i++ {
		peerId += strconv.Itoa(rand.Intn(10))
	}

	return &request{
		announce:   announce,
		infoHash:   infoHash,
		peerId:     peerId,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       length,
	}, nil
}

func (r *request) getTrackerResponse() (*response, error) {
	url := fmt.Sprintf(
		"%s?info_hash=%s&peer_id=%s&uploaded=%d&downloaded=%d&left=%d",
		r.announce,
		r.infoHash,
		url.QueryEscape(r.peerId),
		r.uploaded,
		r.downloaded,
		r.left,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	decodedResp, err := bencode.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	failureReason, requestFailed := decodedResp["failure reason"].(string)
	if requestFailed {
		return &response{
			failureReason: failureReason,
		}, nil
	}

	interval := decodedResp["interval"].(int64)

	peers := make([]peerDict, 0)
	for _, p := range decodedResp["peers"].([]interface{}) {
		peerInterface := p.(map[string]interface{})
		peerId, _ := peerInterface["peer id"].(string)
		peer := peerDict{
			peerId: peerId,
			ip:     peerInterface["ip"].(string),
			port:   peerInterface["port"].(int64),
		}
		peers = append(peers, peer)
	}

	return &response{
		interval:  interval,
		peersDict: peers,
	}, nil
}
