package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jacobgarcia/settify/transport"
)

// New instatiates a new API client for Spotify
func New(a, u, i, s string) *Client {
	return &Client{
		authURL: a,
		URL:     u,
		id:      i,
		secret:  s,
	}
}

// Client contains the required params to connect succesfully to Spotify API
type Client struct {
	authURL string
	URL     string
	id      string
	secret  string
}

// Service expose all endpoints as a services.
// This is a microservices architecture
type Service interface {
	Playlists(token string, offset string) (*Playlists, error)
	Intersect(token string, firstPlaylist string, secondPlaylist string) (*Playlists, error)
}

// PlaylistsDecoder decodes the response object from Spotify for the playlists endpoint
type PlaylistsDecoder struct {
	Items []PlaylistDecoder `json:"items"`
	Total int               `json:"total"`
}

// PlaylistDecoder contains all the objects we want to decode from the items response from Spotify
type PlaylistDecoder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Owner  Owner  `json:"owner"`
	Public bool   `json:"public"`
	Tracks Tracks `json:"tracks"`
}

// Owner refers to the author of a playlist
type Owner struct {
	ID string `json:"id"`
}

// Tracks refers to tracks of a playlist
type Tracks struct {
	Total int `json:"total"`
}

// Playlists encodes the response object for the playlists endpoint
type Playlists struct {
	Items []Playlist `json:"items"`
	Total int        `json:"total"`
}

// Playlist contains the response object for the playlists endpoint
type Playlist struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Scope  string `json:"scope"`
	Tracks int    `json:"tracks"`
}

type PlaylistResponse struct {
	Reference string  `json:"href"`
	Items     []Track `json:"items"`
}

type Track struct {
	Track Playlist `json:"track"`
}

var httpClient *http.Client = &http.Client{}
var data url.Values = url.Values{}

// Playlists is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Playlists(token string, offset string) (*Playlists, error) {
	uri := fmt.Sprintf("%s/v1/me/playlists?offset=%s", c.URL, offset)
	req, err := http.NewRequest("GET", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errResponse transport.IntersectError
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.IntersectError{
			Error: errResponse.Error,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var playslistsDecoder PlaylistsDecoder
	err = json.Unmarshal(response, &playslistsDecoder)

	var playlists []Playlist
	for _, playlist := range playslistsDecoder.Items {
		scope := "public"
		if !playlist.Public {
			scope = "private"
		}
		newPlaylist := Playlist{
			ID:     playlist.ID,
			Name:   playlist.Name,
			Owner:  playlist.Owner.ID,
			Tracks: playlist.Tracks.Total,
			Scope:  scope,
		}
		playlists = append(playlists, newPlaylist)
	}

	playlistsResponse := Playlists{
		Items: playlists,
		Total: playslistsDecoder.Total,
	}

	return &playlistsResponse, err
}

// Intersect is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Intersect(token string, firstPlaylist string, secondPlaylist string) (*Playlists, error) {
	uri := fmt.Sprintf("%s/v1/playlists/%s/tracks", c.URL, firstPlaylist)
	req, err := http.NewRequest("GET", uri, bytes.NewBufferString(data.Encode()))

	fmt.Println(uri)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errResponse transport.IntersectError
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.IntersectError{
			Error: errResponse.Error,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var trackResponse PlaylistResponse
	err = json.Unmarshal(response, &trackResponse)

	first := PlaylistResponse{
		Reference: trackResponse.Reference,
		Items:     trackResponse.Items,
	}

	if err != nil {
		return nil, err
	}

	// Second Playlist Hit
	uri = fmt.Sprintf("%s/v1/playlists/%s/tracks", c.URL, secondPlaylist)
	req, err = http.NewRequest("GET", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", token)

	res, err = httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	response, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		var errResponse transport.IntersectError
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.IntersectError{
			Error: errResponse.Error,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var secondPlaylistResponse PlaylistResponse
	err = json.Unmarshal(response, &secondPlaylistResponse)

	second := PlaylistResponse{
		Reference: secondPlaylistResponse.Reference,
		Items:     secondPlaylistResponse.Items,
	}

	intersection := []Playlist{}
	for _, firstItem := range first.Items {
		for _, secondItem := range second.Items {
			if firstItem.Track.ID == secondItem.Track.ID {
				playlist := Playlist{
					ID:   firstItem.Track.ID,
					Name: firstItem.Track.Name,
				}
				intersection = append(intersection, playlist)
			}
		}
	}

	results := Playlists{
		Items: intersection,
	}
	return &results, err
}
