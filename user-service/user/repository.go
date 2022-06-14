package user

import (
	"fmt"
)

type Repository interface {
	FindAll() (*[]User, error)
	FindById(Id int) (*User, error)
	Create(user User) (*User, error)
	Update(user User) (*User, error)
	Delete(Id int) (*User, error)
}

// var monkUsers = []User{
// 	{1, "abc123456", "Mike", "0918000111"},
// 	{2, "zxc963", "Kyle", "0927789456"},
// 	{3, "ewer456", "June", "0939456123"},
// }

type repository struct {
	monkUsers []User
}

func NewRepository() *repository {
	monkUsers := []User{
		{1, "abc123456", "Mike", "0918000111"},
		{2, "zxc963", "Kyle", "0927789456"},
		{3, "ewer456", "June", "0939456123"},
	}
	repository := repository{
		monkUsers,
	}
	return &repository
}

func (r *repository) FindAll() (*[]User, error) {
	return &r.monkUsers, nil
}

func (r *repository) FindById(Id int) (*User, error) {
	var user User
	err := r.findUserFromMonk(&user, Id)

	return &user, err
}

func (r *repository) Create(user User) (*User, error) {
	user.Id = r.monkUsers[len(r.monkUsers)-1].Id + 1
	r.monkUsers = append(r.monkUsers, user)
	return &user, nil
}

func (r *repository) Update(updateUser User) (*User, error) {
	for i, u := range r.monkUsers {
		if u.Id == updateUser.Id {
			r.monkUsers[i] = updateUser
			return &updateUser, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *repository) Delete(Id int) (*User, error) {
	for i, u := range r.monkUsers {
		if u.Id == Id {
			r.monkUsers = append(r.monkUsers[:i], r.monkUsers[i+1:]...)
			return &u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *repository) findUserFromMonk(user *User, Id int) error {
	for _, u := range r.monkUsers {
		if u.Id == Id {
			*user = u
			return nil
		}
	}

	return fmt.Errorf("user id:%d was not found", Id)
}
