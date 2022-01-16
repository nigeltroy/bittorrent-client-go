package client

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

func (t *torrent) createRequest(announce string, uploaded int64, downloaded int64) (*request, error) {
	hasMultipleFiles := t.metainfo.info.hasMultipleFiles
	var length int64

	if hasMultipleFiles {
		length = 0
		for _, file := range t.metainfo.info.files {
			length += file.length
		}
	} else {
		length = t.metainfo.info.length
	}

	if length == 0 {
		return nil, errors.New("length has not been set or is equal to zero")
	}

	peerId := "-NG0001-"
	for i := 0; i < 12; i++ {
		peerId += strconv.Itoa(rand.Intn(10))
	}

	return &request{
		announce:   announce,
		infoHash:   t.metainfo.infoHash.urlEncodedString,
		peerId:     peerId,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       length - downloaded,
	}, nil
}

func (r *request) getTrackerResponse() (*response, error) {
	requestUrl := fmt.Sprintf(
		"%s?info_hash=%s&peer_id=%s&uploaded=%d&downloaded=%d&left=%d",
		r.announce,
		r.infoHash,
		url.QueryEscape(r.peerId),
		r.uploaded,
		r.downloaded,
		r.left,
	)

	if strings.HasPrefix(r.announce, "udp://") {
		return nil, errors.New("cannot perform UDP connection requests")
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(requestUrl)
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

func (t *torrent) tryToAnnounce() (*tracker, error) {
	// Need to clean this function up
	var couldAnnounce bool

	// Try announcing to the main announce URL first
	request, err := t.createRequest(t.metainfo.announce, 0, 0)
	if err != nil {
		return nil, err
	}

	response, err := request.getTrackerResponse()
	if err != nil {
		// Main announce URL failed, try announcing to other URLs
		for _, url := range t.metainfo.announceList {
			request, err = t.createRequest(url, 0, 0)
			if err != nil {
				return nil, err
			}

			response, _ = request.getTrackerResponse()
			if response != nil {
				couldAnnounce = true
				break
			}
		}
	} else {
		couldAnnounce = true
	}

	if couldAnnounce {
		return &tracker{
			request:  *request,
			response: *response,
		}, nil
	}
	return nil, errors.New("could not announce to any announce URL")
}
