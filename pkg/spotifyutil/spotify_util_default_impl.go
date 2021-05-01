package spotifyutil

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type spotifyUtilDefaultImpl struct {
	userRepo      *user.UserRepo
	authenticator *spotify.Authenticator
	state         string
}

func NewSpotifyUtilDefaultImpl(userRepo *user.UserRepo,
	redirectURI string, state string) SpotifyUtil {

	scopes := []string{
		spotify.ScopeUserReadPrivate,
		spotify.ScopeUserReadRecentlyPlayed,
	}
	authenticator := spotify.NewAuthenticator(redirectURI, scopes...)

	return &spotifyUtilDefaultImpl{
		userRepo:      userRepo,
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

	sUser, err := client.CurrentUser()
	if err != nil {
		return err
	}

	dUser, err := (*s.userRepo).GetByID(sUser.ID)
	if err != nil {
		dUser = user.User{
			CreatedAt:       time.Now().UTC(),
			LastPlayedItems: make([]spotify.RecentlyPlayedItem, 0),
		}
	}

	dUser.UserInfo = *sUser
	dUser.Token = *token
	dUser.UpdatedAt = time.Now().UTC()
	dUser.Active = true
	(*s.userRepo).InsertOrUpdate(dUser)

	return nil
}
