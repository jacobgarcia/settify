package spotify

func unify(first PlaylistResponse, second PlaylistResponse) ([]Playlist, error) {
	tracksUnion := append(first.Items, second.Items...)
	union := []Playlist{}
	for _, item := range tracksUnion {
		playlist := Playlist{
			ID:   item.Track.ID,
			Name: item.Track.Name,
			URI:  item.Track.URI,
		}
		union = append(union, playlist)
	}
	return union, nil
}

// Union merges two playlist tracks into one
func (c Client) Union(token string, firstPlaylist string, secondPlaylist string) (*NewPlaylistResponse, error) {
	return operation(token, firstPlaylist, secondPlaylist, c, unify)
}
