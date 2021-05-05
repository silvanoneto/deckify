package group

import "github.com/zmb3/spotify"

type GroupRepo interface {
	InsertOrUpdate(Group)
	Remove(spotify.ID) error
	GetByID(spotify.ID) (Group, error)
	GetAllActive(uint, uint) []Group
	GetAllActiveByUserID(spotify.ID, uint, uint) []Group
}
