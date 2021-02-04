package storage

type Friends struct{
Friends map[Username]map[Username]struct{}
}

func (f * Friends)  becomeFriends(username1, username2 Username) error {
	if _, ok := f.Friends[username1]; !ok { // if 1 doesnt exist in map
		f.Friends[username1] = map[Username]struct{}{}
	}

	if _, ok := f.Friends[username1][username2]; !ok { // if 2 is not a friend of 1
		f.Friends[username1][username2] = struct{}{}
		return nil
	} else {
		return errors.New("you are already friends")
	}
}

func (f * Friends) AddFriend(username1, username2 Username) error {
	err := f.becomeFriends(username1, username2)
	if err != nil {
		return err
	}
	_ = f.becomeFriends(username2, username1)
	return nil

	//if _, ok := f.Friends[username1]; ok{ 	// if 1 exists in map
	//	if _, ok := f.Friends[username1][username2]; !ok { // and 1 is not a friend of 2
	//		f.Friends[username1][username2] = struct{}{} // 1 and 2 are friends
	//		if _, ok := f.Friends[username2]; ok { // if 2 exists in map
	//			f.Friends[username2][username1] = struct{}{} // Make 2 a friend of 1
	//		} else { // if 2 doesn't exist in map
	//			f.Friends[username2] = map[Username]struct{}{}
	//			f.Friends[username2][username1] = struct{}{}
	//		}
	//	} else{
	//		return errors.New("you are already friends")
	//	}
	//} else { // if 1 doesn't exist in map
	//	f.Friends[username1] = map[Username]struct{}{}
	//	f.Friends[username1][username2] = struct{}{}
	//	//Make 2 a friend of 1
	//	// add 2 to 1
	//}
}
