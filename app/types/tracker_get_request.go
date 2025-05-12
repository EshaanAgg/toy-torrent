package types

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// TrackerGetRequest is a struct that holds the parameters for
// making a GET request to a HTTP tracker.
type TrackerGetRequest struct {
	TrackerURL string
	InfoHash   []byte
	PeerID     string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Compact    int
}

// Makes a GET request to a HTTP tracker to discover peers
// to download the file from.
// connectToPeers is a boolean that indicates whether to make a network connection
// to each of the peers returned by the tracker.
func (r *TrackerGetRequest) MakeRequest(connectToPeers bool) (*TrackerGetResponse, error) {
	// URL query parameters
	p := url.Values{}
	p.Add("info_hash", string(r.InfoHash))
	p.Add("peer_id", r.PeerID)
	p.Add("port", fmt.Sprintf("%d", r.Port))
	p.Add("uploaded", fmt.Sprintf("%d", r.Uploaded))
	p.Add("downloaded", fmt.Sprintf("%d", r.Downloaded))
	p.Add("left", fmt.Sprintf("%d", r.Left))
	p.Add("compact", fmt.Sprintf("%d", r.Compact))

	// Make the GET request
	url := r.TrackerURL + "?" + p.Encode()
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	response, err := NewTrackerGetResponse(data, connectToPeers)
	if err != nil {
		return nil, fmt.Errorf("error decoding tracker response: %w", err)
	}
	return response, nil
}
