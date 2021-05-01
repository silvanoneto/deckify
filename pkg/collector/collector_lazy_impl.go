package collector

import (
	"log"
	"time"

	"github.com/silvanoneto/deckify/pkg/spotifyutil"
	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
)

type lazyCollector struct {
	userRepo              *user.UserRepo
	spotifyUtil           *spotifyutil.SpotifyUtil
	pageSize              uint
	callIntervalInSeconds uint
}

func NewLazyCollector(userRepo *user.UserRepo, spotifyUtil *spotifyutil.SpotifyUtil, pageSize uint,
	callIntervalInSeconds uint) *lazyCollector {

	return &lazyCollector{
		userRepo:              userRepo,
		spotifyUtil:           spotifyUtil,
		pageSize:              pageSize,
		callIntervalInSeconds: callIntervalInSeconds,
	}
}

func (c *lazyCollector) Start() {
	for {
		c.collect()
		time.Sleep(time.Second * time.Duration(c.callIntervalInSeconds))
	}
}

func (c *lazyCollector) collect() {
	var pageNumber uint = 0
	for {
		users := (*c.userRepo).GetAllActive(pageNumber, c.pageSize)
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			if err := c.collectLastPlayedItems(&user, 50); err != nil {
				log.Println(err)
			}
			user.UpdatedAt = time.Now().UTC()
			(*c.userRepo).InsertOrUpdate(user)
		}
		pageNumber++
	}
}

func (c *lazyCollector) collectLastPlayedItems(user *user.User, limit uint) error {

	client := (*c.spotifyUtil).GetClient(&user.Token)

	options := &spotify.RecentlyPlayedOptions{
		Limit: int(limit),
	}
	if len(user.LastPlayedItems) > 0 {
		lastItem := user.LastPlayedItems[len(user.LastPlayedItems)-1]
		options.AfterEpochMs = lastItem.PlayedAt.UnixNano() /
			int64(time.Millisecond)
	}

	items, err := client.PlayerRecentlyPlayedOpt(options)
	if err != nil {
		return err
	}

	reversedItems := make([]spotify.RecentlyPlayedItem, 0, len(items))
	for i := range items {
		item := items[len(items)-1-i]
		reversedItems = append(reversedItems, item)
	}

	c.printNewPlayedItems(user, reversedItems)

	user.LastPlayedItems = append(user.LastPlayedItems, reversedItems...)
	return nil
}

func (c *lazyCollector) printNewPlayedItems(user *user.User, items []spotify.RecentlyPlayedItem) {
	for _, item := range items {
		userName := user.UserInfo.DisplayName
		trackName := item.Track.Name
		playedAt := item.PlayedAt

		var artistName string
		if len(item.Track.Artists) > 0 {
			artistName = item.Track.Artists[0].Name
		}

		log.Printf("%s heard %s, by %s, at %v", userName, trackName, artistName, playedAt)
	}
}
