package spotify

func complement(first PlaylistResponse, second PlaylistResponse) ([]Playlist, error) {
	// Then we traverse all elements and create the interesection
	// We doing it this unoptimized way so we can compare it after this O(n*m)
	complement := []Playlist{}
	// First we need to calculate the intersection, then remove it from B
	for _, firstItem := range first.Items {
		for index, secondItem := range second.Items {

			go func(firstItem Track, secondItem Track, index int) {
				if firstItem.Track.ID == secondItem.Track.ID {
					// Remove the element at index i from a.
					second.Items[index] = second.Items[len(second.Items)-1] // Copy last element to index i.
					second.Items[len(second.Items)-1] = Track{}             // Erase last element (write zero value).
					second.Items = second.Items[:len(second.Items)-1]       // Truncate slice.
				}
			}(firstItem, secondItem, index)
		}
	}

	for _, item := range second.Items {
		playlist := Playlist{
			ID:   item.Track.ID,
			Name: item.Track.Name,
			URI:  item.Track.URI,
		}
		complement = append(complement, playlist)
	}

	return complement, nil
}

// Complement creates a playlist containing all elements that are not in A
func (c Client) Complement(token string, firstPlaylist string, secondPlaylist string) (*NewPlaylistResponse, error) {
	return operation(token, firstPlaylist, secondPlaylist, c, complement)
}
