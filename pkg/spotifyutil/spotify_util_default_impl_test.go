package spotifyutil

import (
	"reflect"
	"testing"

	"github.com/silvanoneto/deckify/pkg/group"
	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
)

type userRepoMockImpl struct{}

func (r *userRepoMockImpl) InsertOrUpdate(user.User)              {}
func (r *userRepoMockImpl) Remove(spotify.ID) error               { return nil }
func (r *userRepoMockImpl) GetByID(spotify.ID) (user.User, error) { return user.User{}, nil }
func (r *userRepoMockImpl) GetAllActive(uint, uint) []user.User   { return []user.User{} }
func NewUserRepoMockImpl() user.UserRepo                          { return &userRepoMockImpl{} }

type groupRepoMockImpl struct{}

func (r *groupRepoMockImpl) InsertOrUpdate(group.Group)              {}
func (r *groupRepoMockImpl) Remove(spotify.ID) error                 { return nil }
func (r *groupRepoMockImpl) GetByID(spotify.ID) (group.Group, error) { return group.Group{}, nil }
func (r *groupRepoMockImpl) GetAllActive(uint, uint) []group.Group   { return []group.Group{} }
func (r *groupRepoMockImpl) GetAllActiveByUserID(spotify.ID, uint, uint) []group.Group {
	return []group.Group{}
}
func NewGroupRepoMockImpl() group.GroupRepo { return &groupRepoMockImpl{} }

func TestNewSpotifyUtilDefaultImpl(t *testing.T) {
	type args struct {
		userRepo    *user.UserRepo
		groupRepo   *group.GroupRepo
		redirectURI string
		state       string
	}
	var (
		defaultUserRepo      user.UserRepo   = NewUserRepoMockImpl()
		defaultGroupRepo     group.GroupRepo = NewGroupRepoMockImpl()
		defaultSpotifyScopes []string        = []string{
			spotify.ScopeUserReadPrivate,
			spotify.ScopeUserReadRecentlyPlayed,
			spotify.ScopePlaylistModifyPublic,
			spotify.ScopePlaylistModifyPrivate,
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
				userRepo:  &defaultUserRepo,
				groupRepo: &defaultGroupRepo,
			},
			want: &spotifyUtilDefaultImpl{
				userRepo:      &defaultUserRepo,
				groupRepo:     &defaultGroupRepo,
				authenticator: &defaultSpotifyAuthenticator,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpotifyUtilDefaultImpl(tt.args.userRepo, tt.args.groupRepo, tt.args.redirectURI,
				tt.args.state); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpotifyUtilDefaultImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}
