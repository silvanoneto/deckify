package spotifyutil

import (
	"reflect"
	"testing"

	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
)

type userRepoMockImpl struct{}

func (r *userRepoMockImpl) InsertOrUpdate(user.User)            {}
func (r *userRepoMockImpl) Remove(string) error                 { return nil }
func (r *userRepoMockImpl) GetByID(string) (user.User, error)   { return user.User{}, nil }
func (r *userRepoMockImpl) GetAllActive(uint, uint) []user.User { return []user.User{} }
func NewUserRepoMockImpl() user.UserRepo                        { return &userRepoMockImpl{} }

func TestNewSpotifyUtilDefaultImpl(t *testing.T) {
	type args struct {
		userRepo    *user.UserRepo
		redirectURI string
		state       string
	}
	var (
		defaultUserRepo      user.UserRepo = NewUserRepoMockImpl()
		defaultSpotifyScopes []string      = []string{
			spotify.ScopeUserReadPrivate,
			spotify.ScopeUserReadRecentlyPlayed,
		}
		defaultSpotifyAuthenticator spotify.Authenticator = spotify.NewAuthenticator("", defaultSpotifyScopes...)
	)
	tests := []struct {
		name string
		args args
		want *spotifyUtilDefaultImpl
	}{
		{
			name: "createAnInstance",
			args: args{
				userRepo: &defaultUserRepo,
			},
			want: &spotifyUtilDefaultImpl{
				userRepo:      &defaultUserRepo,
				authenticator: &defaultSpotifyAuthenticator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpotifyUtilDefaultImpl(tt.args.userRepo, tt.args.redirectURI,
				tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpotifyUtilDefaultImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}
