package storage

import "errors"

type Friends struct {
	Friends map[string]map[string]struct{}
}

func (f *Friends) Add(username1, username2 string) error {
	err := f.becomeFriends(username1, username2)
	if err != nil {
		return err
	}
	_ = f.becomeFriends(username2, username1)
	return nil
}

func (f *Friends) GetFriendsOf(username string) []string {
	m := f.Friends[username]
	friends := make([]string, len(m))

	i := 0
	for k := range m {
		friends[i] = k
		i++
	}

	return friends
}

func (f *Friends) becomeFriends(username1, username2 string) error {
	if _, ok := f.Friends[username1]; !ok { // If 1 doesn't exist in the map
		f.Friends[username1] = map[string]struct{}{}
	}

	if _, ok := f.Friends[username1][username2]; !ok { // If 2 is not a friend of 1
		f.Friends[username1][username2] = struct{}{}
		return nil
	} else {
		return errors.New("you are already friends")
	}
}
