package spotifyutil

import (
	"net/http"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type SpotifyUtil interface {
	AuthCallback(http.ResponseWriter, *http.Request)
	GetAuthURL() string
	GetClient(*oauth2.Token) spotify.Client
}
