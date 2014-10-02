package models

import "testing"

func TestUser(t *testing.T) {
	var (
		user  User
		user2 User
		users []User
		err   error
	)

	user = User{Name: "user1", Email: "user1@test.com"}
	err = user.Save()
	if err != nil {
		t.Fatal(err.Error())
	}

	user3, err := GetUser(user.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if user3.Name != "user1" {
		t.Fatal("error")
	}

	user3.Name = "user3"
	user3.Save()
	user, err = GetUser(user.Id)
	if user.Name != "user3" {
		t.Fatal("error")
	}

	users, err = GetAllUsers()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(users) != 1 {
		t.Fatal("there should be 1 users, not ", len(users))
	}

	user2 = User{Name: "user2", Email: "user2@test.com"}
	err = user2.Save()
	if err != nil {
		t.Fatal(err.Error())
	}

	users, err = GetAllUsers()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(users) != 2 {
		t.Fatal("there should be 2 users, not ", len(users))
	}

	user.Delete()
	users, err = GetAllUsers()
	if len(users) != 1 {
		t.Fatal("there should be 1 users, not ", len(users))
	}

}
