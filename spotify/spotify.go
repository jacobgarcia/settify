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

// method is a custom type abstrction in order to pass functions as parameters
type method func(first PlaylistResponse, second PlaylistResponse) ([]Playlist, error)

// Client contains the required params to connect succesfully to Spotify API
type Client struct {
	authURL string
	URL     string
	id      string
	secret  string
}

// Service expose all endpoints as services
// This is a microservices architecture pattern
type Service interface {
	UserPlaylists(token, offset, username string) (*Playlists, error)
	Profile(token string) (*User, error)
	Playlists(token string, offset string) (*Playlists, error)
	Intersect(token, firstPlaylist, secondPlaylist, name string) (*NewPlaylistResponse, error)
	Union(token, firstPlaylist, secondPlaylist, name string) (*NewPlaylistResponse, error)
	Complement(token, firstPlaylist, secondPlaylist, name string) (*NewPlaylistResponse, error)
}

// Image specifies image urls of an object
type Image struct {
	URL string `json:"url"`
}

// PlaylistsDecoder decodes the response object from Spotify for the playlists endpoint
type PlaylistsDecoder struct {
	Items []PlaylistDecoder `json:"items"`
	Total int               `json:"total"`
}

// PlaylistDecoder contains all the objects we want to decode from the items response from Spotify
type PlaylistDecoder struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Owner  Owner   `json:"owner"`
	Public bool    `json:"public"`
	Tracks Tracks  `json:"tracks"`
	Images []Image `json:"images"`
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
	Image  string `json:"image,omitempty"`
}

// User encodes/decodes the user id for Spotify
type User struct {
	ID     string  `json:"id"`
	Name   string  `json:"display_name,omitempty"`
	Email  string  `json:"email,omitempty"`
	Images []Image `json:"images,omitempty"`
}

// PlaylistResponse contains the format for the response object
type PlaylistResponse struct {
	Reference string  `json:"href"`
	Items     []Track `json:"items"`
}

// Track is abstraction for playlist contained as track in the json from the Spotify API
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

// Playlists retrieves the playlists from the user
func (c Client) Playlists(token string, offset string) (*Playlists, error) {
	return getPlaylists(token, offset, "me", c)
}

// UserPlaylists retrieves the playlists from the user
func (c Client) UserPlaylists(token, offset, username string) (*Playlists, error) {
	username = fmt.Sprintf("users/%s", username)
	return getPlaylists(token, offset, username, c)
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
		ID:     user.ID,
		Email:  user.Email,
		Images: user.Images,
		Name:   user.Name,
	}

	return &userResponse, nil
}
