package user

import (
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type User struct {
	UserInfo        spotify.PrivateUser
	LastPlayedItems []spotify.RecentlyPlayedItem
	Token           oauth2.Token
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Active          bool
}
