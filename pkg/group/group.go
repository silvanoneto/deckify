package group

import (
	"time"

	"github.com/zmb3/spotify"
)

type Group struct {
	ID        spotify.ID
	Name      string
	Users     map[spotify.ID]struct{}
	Owner     spotify.ID
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
}
