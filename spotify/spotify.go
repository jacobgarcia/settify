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
	Authenticate() (*AuthenticationResponse, error)
	Playlists(token string) (*TrackResponse, error)
	Intersect(token string, firstPlaylist string, secondPlaylist string) (*TrackResponse, error)
}

// AuthenticationResponse is the response struct from Spotify
type AuthenticationResponse struct {
	Token      string `json:"access_token"`
	Type       string `json:"token_type"`
	Expiration int    `json:"expires_in"`
	Scope      string `json:"scope"`
}

type TrackResponse struct {
	Reference string     `json:"href"`
	Items     []Playlist `json:"items"`
}

type PlaylistResponse struct {
	Reference string   `json:"href"`
	Items     []Tracks `json:"items"`
}

type Tracks struct {
	Track Playlist `json:"track"`
}

type Playlist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var httpClient *http.Client = &http.Client{}
var data url.Values = url.Values{}

// Authenticate is the method for authentication and getting a valid Spotify token
func (c Client) Authenticate() (*AuthenticationResponse, error) {
	uri := fmt.Sprintf("%s/api/token", c.authURL)

	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.id, c.secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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
		var errResponse transport.ErrorResponse
		err = json.Unmarshal(response, &errResponse)

		if err != nil {
			return nil, err
		}

		errResponse = transport.ErrorResponse{
			Message: errResponse.Message,
			Status:  res.StatusCode,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s", resp)
	}

	var authResponse AuthenticationResponse
	err = json.Unmarshal(response, &authResponse)

	auth := AuthenticationResponse{
		Token:      authResponse.Token,
		Type:       authResponse.Type,
		Expiration: authResponse.Expiration,
		Scope:      authResponse.Scope,
	}

	return &auth, err
}

// Playlists is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Playlists(token string) (*TrackResponse, error) {
	uri := fmt.Sprintf("%s/v1/me/playlists", c.URL)
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

	var trackResponse TrackResponse
	err = json.Unmarshal(response, &trackResponse)

	track := TrackResponse{
		Reference: trackResponse.Reference,
		Items:     trackResponse.Items,
	}

	return &track, err
}

// Intersect is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Intersect(token string, firstPlaylist string, secondPlaylist string) (*TrackResponse, error) {
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

	results := TrackResponse{
		Reference: secondPlaylistResponse.Reference,
		Items:     intersection,
	}
	return &results, err
}
