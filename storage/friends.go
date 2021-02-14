package storage

import "errors"

type Friends struct {
	Friends map[Username]map[Username]struct{}
}

func (f *Friends) becomeFriends(username1, username2 Username) error {
	if _, ok := f.Friends[username1]; !ok { // If 1 doesn't exist in the map
		f.Friends[username1] = map[Username]struct{}{}
	}

	if _, ok := f.Friends[username1][username2]; !ok { // If 2 is not a friend of 1
		f.Friends[username1][username2] = struct{}{}
		return nil
	} else {
		return errors.New("you are already friends")
	}
}

func (f *Friends) Add(username1, username2 Username) error {
	err := f.becomeFriends(username1, username2)
	if err != nil {
		return err
	}
	_ = f.becomeFriends(username2, username1)
	return nil
}

// TODO removeFriend
