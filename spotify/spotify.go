package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jacobgarcia/settify/transport"

	"github.com/Pallinder/go-randomdata"
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
	Intersect(token string, firstPlaylist string, secondPlaylist string) (*NewPlaylistResponse, error)
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
	ID     string `json:"id,omitempty"`
	Name   string `json:"name"`
	Owner  string `json:"owner,omitempty"`
	Scope  string `json:"scope,omitempty"`
	Tracks int    `json:"tracks,omitempty"`
	URI    string `json:"uri,omitempty"`
}

// User encodes/decodes the user id for Spotify
type User struct {
	ID string `json:"id"`
}

type PlaylistResponse struct {
	Reference string  `json:"href"`
	Items     []Track `json:"items"`
}

type Track struct {
	Track Playlist `json:"track"`
}

// NewPlaylistResponse is the response object when creating a new playlist
type NewPlaylistResponse struct {
	Name   string `json:"name"`
	Href   string `json:"href"`
	Tracks int    `json:"tracks"`
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

func request(url string, path string, token string, dat interface{}) (*User, error) {
	// The URL for the request
	uri := fmt.Sprintf("%s/%s", url, path)
	// Specify if the request its a GET or a POST
	method := "GET"
	var requestBody []byte
	if dat != nil {
		method = "POST"
		// Add body
		requestBody, _ = json.Marshal(dat)
	}
	// Create the request object
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	// Add an Authorization token
	req.Header.Add("Authorization", token)

	// Actually DO the request
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// Close response using defer
	defer res.Body.Close()
	// Read the body of the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Manage the response if it's not an OK Status
	if res.StatusCode > 299 || res.StatusCode < 200 {
		var errResponse transport.IntersectError
		err = json.Unmarshal(body, &errResponse)

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

	// Create the decoder object
	var user User
	// Decode the object
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	userResponse := User{
		ID: user.ID,
	}

	return &userResponse, nil
}

// Intersect is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Intersect(token string, firstPlaylist string, secondPlaylist string) (*NewPlaylistResponse, error) {
	// First we need to retrieve the first playlist tracks
	uri := fmt.Sprintf("%s/v1/playlists/%s/tracks", c.URL, firstPlaylist)
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

	var trackResponse PlaylistResponse
	err = json.Unmarshal(response, &trackResponse)

	first := PlaylistResponse{
		Reference: trackResponse.Reference,
		Items:     trackResponse.Items,
	}

	if err != nil {
		return nil, err
	}

	// Now we need to second playlist
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

	if err != nil {
		return nil, err
	}

	second := PlaylistResponse{
		Reference: secondPlaylistResponse.Reference,
		Items:     secondPlaylistResponse.Items,
	}

	// Then we traverse all elements and create the interesection
	// We doing it this unoptimized way so we can compare it after this
	intersection := []Playlist{}
	for _, firstItem := range first.Items {
		for _, secondItem := range second.Items {
			if firstItem.Track.ID == secondItem.Track.ID {
				playlist := Playlist{
					ID:   firstItem.Track.ID,
					Name: firstItem.Track.Name,
					URI:  firstItem.Track.URI,
				}
				intersection = append(intersection, playlist)
			}
		}
	}

	results := Playlists{
		Items: intersection,
		Total: len(intersection),
	}

	// Next, we need the user.id of the current session.
	// This is a requirement to create the new playlist.
	user, err := request(c.URL, "v1/me", token, nil)
	if err != nil {
		return nil, err
	}

	// Next, we create the empty playlist with a new random name
	name := randomdata.SillyName()
	newPlaylist := Playlist{
		Name: name,
	}
	uri = fmt.Sprintf("v1/users/%s/playlists", user.ID)
	playlist, err := request(c.URL, uri, token, newPlaylist)
	if err != nil {
		return nil, err
	}

	// Finally, we need to add the tracks to the playlist
	// Create an slice containing the tracks
	tracks := []string{}
	for _, item := range results.Items {
		tracks = append(tracks, item.URI)
	}
	uri = fmt.Sprintf("v1/playlists/%s/tracks", playlist.ID)
	jsonTracks := map[string][]string{
		"uris": tracks,
	}
	_, err = request(c.URL, uri, token, jsonTracks)

	if err != nil {
		return nil, err
	}

	// At the end, we just create a new response object containing the information we need
	newPlaylistResponse := NewPlaylistResponse{
		Name:   name,
		Href:   playlist.ID,
		Tracks: len(tracks),
	}

	return &newPlaylistResponse, err
}
