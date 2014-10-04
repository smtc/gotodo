package models

import (
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	var (
		user  User
		user2 User
		users []User
		err   error
	)

	user = User{Name: "中文", Email: "user1@test.com"}
	err = user.Save()
	if err != nil {
		t.Fatal(err.Error())
	}

	user3, err := GetUser(user.Id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if user3.Name != "中文" {
		t.Fatal("error")
	}

	user3.Name = "user3"
	user3.Save()
	user4, err := GetUser(user.Id)
	if user4.Name != "user3" {
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

	ids := fmt.Sprintf("%v,%v", user.Id, user2.Id)
	names := GetMultUserName(ids)
	if names != "user3,user2" {
		t.Fatal("names should be user3,user2, not ", names)
	}

	UserDelete(user.Id)
	users, err = GetAllUsers()
	if len(users) != 1 {
		t.Fatal("there should be 1 users, not ", len(users))
	}

}
