package spotify

// Profile gets the user information
func (c Client) Profile(token string) (*User, error) {
	// Next, we need the user.id of the current session.
	// This is a requirement to create the new playlist.
	user, err := request(c.URL, "v1/me", token, nil)
	if err != nil {
		return nil, err
	}

	return user, nil
}
