package collector

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/silvanoneto/deckify/pkg/spotifyutil"
	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type userRepoMockImpl struct{}

func (r *userRepoMockImpl) InsertOrUpdate(user.User)            {}
func (r *userRepoMockImpl) Remove(string) error                 { return nil }
func (r *userRepoMockImpl) GetByID(string) (user.User, error)   { return user.User{}, nil }
func (r *userRepoMockImpl) GetAllActive(uint, uint) []user.User { return []user.User{} }
func NewUserRepoMockImpl() user.UserRepo                        { return &userRepoMockImpl{} }

type spotifyUtilMockImpl struct{}

func (s *spotifyUtilMockImpl) AuthCallback(http.ResponseWriter, *http.Request) {}
func (s *spotifyUtilMockImpl) GetAuthURL() string                              { return "" }
func (s *spotifyUtilMockImpl) GetClient(*oauth2.Token) spotify.Client          { return spotify.Client{} }
func NewSpotifyUtilMockImpl() spotifyutil.SpotifyUtil                          { return &spotifyUtilMockImpl{} }

func TestNewLazyCollector(t *testing.T) {
	type args struct {
		userRepo              *user.UserRepo
		spotifyUtil           *spotifyutil.SpotifyUtil
		pageSize              uint
		callIntervalInSeconds uint
	}
	var (
		defaultUserRepo              user.UserRepo           = NewUserRepoMockImpl()
		defaultSpotifyUtil           spotifyutil.SpotifyUtil = NewSpotifyUtilMockImpl()
		defaultPageSize              uint                    = 10
		defaultCallIntervalInSeconds uint                    = 10
	)
	tests := []struct {
		name string
		args args
		want *lazyCollector
	}{
		{
			name: "createAnInstance",
			args: args{
				userRepo:              &defaultUserRepo,
				spotifyUtil:           &defaultSpotifyUtil,
				pageSize:              defaultPageSize,
				callIntervalInSeconds: defaultCallIntervalInSeconds,
			},
			want: &lazyCollector{
				userRepo:              &defaultUserRepo,
				spotifyUtil:           &defaultSpotifyUtil,
				pageSize:              defaultPageSize,
				callIntervalInSeconds: defaultCallIntervalInSeconds,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLazyCollector(tt.args.userRepo, tt.args.spotifyUtil, tt.args.pageSize,
				tt.args.callIntervalInSeconds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLazyCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lazyCollector_Start(t *testing.T) {
	type fields struct {
		userRepo              *user.UserRepo
		spotifyUtil           *spotifyutil.SpotifyUtil
		pageSize              uint
		callIntervalInSeconds uint
	}
	var (
		defaultUserRepo              user.UserRepo           = NewUserRepoMockImpl()
		defaultSpotifyUtil           spotifyutil.SpotifyUtil = NewSpotifyUtilMockImpl()
		defaultPageSize              uint                    = 10
		defaultCallIntervalInSeconds uint                    = 10
	)
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "runCollector",
			fields: fields{
				userRepo:              &defaultUserRepo,
				spotifyUtil:           &defaultSpotifyUtil,
				pageSize:              defaultPageSize,
				callIntervalInSeconds: defaultCallIntervalInSeconds,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lazyCollector{
				userRepo:              tt.fields.userRepo,
				spotifyUtil:           tt.fields.spotifyUtil,
				pageSize:              tt.fields.pageSize,
				callIntervalInSeconds: tt.fields.callIntervalInSeconds,
			}
			go c.Start()
		})
	}
}

func Test_lazyCollector_printNewPlayedItems(t *testing.T) {
	type fields struct {
		userRepo              *user.UserRepo
		spotifyUtil           *spotifyutil.SpotifyUtil
		pageSize              uint
		callIntervalInSeconds uint
	}
	type args struct {
		user  *user.User
		items []spotify.RecentlyPlayedItem
	}
	var (
		defaultUserRepo              user.UserRepo           = NewUserRepoMockImpl()
		defaultSpotifyUtil           spotifyutil.SpotifyUtil = NewSpotifyUtilMockImpl()
		defaultPageSize              uint                    = 10
		defaultCallIntervalInSeconds uint                    = 10
	)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "printWhenItemsListIsEmpty",
			fields: fields{
				userRepo:              &defaultUserRepo,
				spotifyUtil:           &defaultSpotifyUtil,
				pageSize:              defaultPageSize,
				callIntervalInSeconds: defaultCallIntervalInSeconds,
			},
			args: args{user: &user.User{}, items: []spotify.RecentlyPlayedItem{}},
		},
		{
			name: "printWhenItemsListIsNotEmpty",
			fields: fields{
				userRepo:              &defaultUserRepo,
				spotifyUtil:           &defaultSpotifyUtil,
				pageSize:              defaultPageSize,
				callIntervalInSeconds: defaultCallIntervalInSeconds,
			},
			args: args{
				user: &user.User{},
				items: []spotify.RecentlyPlayedItem{
					{},
					{Track: spotify.SimpleTrack{Artists: []spotify.SimpleArtist{{Name: "ABC"}}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &lazyCollector{
				userRepo:              tt.fields.userRepo,
				spotifyUtil:           tt.fields.spotifyUtil,
				pageSize:              tt.fields.pageSize,
				callIntervalInSeconds: tt.fields.callIntervalInSeconds,
			}
			c.printNewPlayedItems(tt.args.user, tt.args.items)
		})
	}
}
