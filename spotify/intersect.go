package spotify

func intersect(first PlaylistResponse, second PlaylistResponse) ([]Playlist, error) {
	// Then we traverse all elements and create the interesection
	// We doing it this unoptimized way so we can compare it after this O(n*m)
	intersection := []Playlist{}
	for _, firstItem := range first.Items {
		for _, secondItem := range second.Items {

			go func(firstItem Track, secondItem Track) {
				if firstItem.Track.ID == secondItem.Track.ID {
					playlist := Playlist{
						ID:   firstItem.Track.ID,
						Name: firstItem.Track.Name,
						URI:  firstItem.Track.URI,
					}
					intersection = append(intersection, playlist)
				}
			}(firstItem, secondItem)

		}
	}

	return intersection, nil
}

// Intersect is the first method will be implementing in Settify. Basically takes two playlists, and generates a new playlist containing the interesection between them.
func (c Client) Intersect(token, firstPlaylist, secondPlaylist, name string) (*NewPlaylistResponse, error) {
	return operation(token, firstPlaylist, secondPlaylist, name, c, intersect)
}
