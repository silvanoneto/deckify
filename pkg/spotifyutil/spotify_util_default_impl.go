package spotifyutil

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/silvanoneto/deckify/pkg/group"
	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type spotifyUtilDefaultImpl struct {
	userRepo      *user.UserRepo
	groupRepo     *group.GroupRepo
	authenticator *spotify.Authenticator
	state         string
}

func NewSpotifyUtilDefaultImpl(userRepo *user.UserRepo, groupRepo *group.GroupRepo, redirectURI string,
	state string) SpotifyUtil {

	scopes := []string{
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadRecentlyPlayed,
		spotify.ScopePlaylistModifyPublic,
		spotify.ScopePlaylistModifyPrivate,
	}
	authenticator := spotify.NewAuthenticator(redirectURI, scopes...)

	return &spotifyUtilDefaultImpl{
		userRepo:      userRepo,
		groupRepo:     groupRepo,
		authenticator: &authenticator,
		state:         state,
	}
}

func (s *spotifyUtilDefaultImpl) AuthCallback(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request for:", r.URL.String())

	token, err := s.authenticator.Token(s.state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Println(err)
		return
	}
	if st := r.FormValue("state"); st != s.state {
		http.NotFound(w, r)
		log.Printf("State mismatch: %s != %s\n", st, s.state)
		return
	}

	s.updateUserData(token)
	if err != nil {
		http.Error(w, "User could not be updated",
			http.StatusInternalServerError)
		log.Println(err)
		return
	}

	fmt.Fprintf(w, "Login Completed!")
}

func (s *spotifyUtilDefaultImpl) GetAuthURL() string {
	return s.authenticator.AuthURL(s.state)
}

func (s *spotifyUtilDefaultImpl) GetClient(token *oauth2.Token) spotify.Client {
	return s.authenticator.NewClient(token)
}

func (s *spotifyUtilDefaultImpl) updateUserData(token *oauth2.Token) error {
	client := s.GetClient(token)

	spotifyUser, err := client.CurrentUser()
	if err != nil {
		return err
	}

	deckifyUser, err := (*s.userRepo).GetByID(spotify.ID(spotifyUser.ID))
	if err != nil {
		deckifyUser = user.User{
			CreatedAt:    time.Now().UTC(),
			PlayedTracks: make([]spotify.RecentlyPlayedItem, 0),
		}

		newGroup := group.Group{
			Name:      "Your Deck <3",
			Users:     make(map[spotify.ID]struct{}),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Active:    true,
		}

		newGroup.Owner = spotify.ID(spotifyUser.ID)
		newGroup.Users[spotify.ID(spotifyUser.ID)] = struct{}{}

		spotifyPlaylist, err := client.CreatePlaylistForUser(spotifyUser.ID, newGroup.Name, "Created by Deckify", true)
		if err != nil {
			return err
		}

		newGroup.ID = spotifyPlaylist.ID
		(*s.groupRepo).InsertOrUpdate(newGroup)
	}

	deckifyUser.ID = spotify.ID(spotifyUser.ID)
	deckifyUser.UserInfo = *spotifyUser
	deckifyUser.Token = *token
	deckifyUser.UpdatedAt = time.Now().UTC()
	deckifyUser.Active = true
	(*s.userRepo).InsertOrUpdate(deckifyUser)

	return nil
}
