package user

import "github.com/zmb3/spotify"

type UserRepo interface {
	InsertOrUpdate(User)
	Remove(spotify.ID) error
	GetByID(spotify.ID) (User, error)
	GetAllActive(uint, uint) []User
}
