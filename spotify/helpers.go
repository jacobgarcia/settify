package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Pallinder/go-randomdata"
	"github.com/jacobgarcia/settify/transport"
)

func operation(token string, firstPlaylist string, secondPlaylist string, c Client, fn method) (*NewPlaylistResponse, error) {
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

	// Get playlist with applied operation
	op, err := fn(first, second)
	if err != nil {
		return nil, err
	}

	if len(op) == 0 {
		nestedError := transport.NestedError{
			Status:  204,
			Message: "Playlists doesn't have anything in common",
		}
		errResponse := transport.IntersectError{
			Error: nestedError,
		}

		resp, err := json.Marshal(errResponse)

		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", resp)
	}

	results := Playlists{
		Items: op,
		Total: len(op),
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

	fmt.Printf("%+v\n", newPlaylist)
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