package stacker

import (
	"log"
	"sort"
	"time"

	"github.com/silvanoneto/deckify/pkg/group"
	"github.com/silvanoneto/deckify/pkg/spotifyutil"
	"github.com/silvanoneto/deckify/pkg/user"
	"github.com/zmb3/spotify"
)

type lazyStacker struct {
	userRepo              *user.UserRepo
	groupRepo             *group.GroupRepo
	spotifyUtil           *spotifyutil.SpotifyUtil
	pageSize              uint
	callIntervalInSeconds uint
	trackWindowInDays     uint
}

func NewLazyStacker(userRepo *user.UserRepo, groupRepo *group.GroupRepo, spotifyUtil *spotifyutil.SpotifyUtil,
	pageSize uint, callIntervalInSeconds uint, trackWindowInDays uint) *lazyStacker {

	return &lazyStacker{
		userRepo:              userRepo,
		groupRepo:             groupRepo,
		spotifyUtil:           spotifyUtil,
		pageSize:              pageSize,
		callIntervalInSeconds: callIntervalInSeconds,
		trackWindowInDays:     trackWindowInDays,
	}
}

func (st *lazyStacker) Start() {
	for {
		st.stack()
		time.Sleep(time.Second * time.Duration(st.callIntervalInSeconds))
	}
}

func (st *lazyStacker) stack() {
	var pageNumber uint = 0
	for {
		groups := (*st.groupRepo).GetAllActive(pageNumber, st.pageSize)
		if len(groups) == 0 {
			break
		}

		for _, group := range groups {
			if err := st.updateGroupPlaylist(&group); err != nil {
				log.Println(err)
			}
		}
		pageNumber++
	}
}

func (st *lazyStacker) updateGroupPlaylist(group *group.Group) error {
	type criteria struct {
		count      int
		lastPlayAt time.Time
	}

	tracks := make(map[spotify.ID]criteria)

	for userID := range group.Users {
		user, err := (*st.userRepo).GetByID(userID)
		if err != nil {
			log.Println(err)
			continue
		}
		if !user.Active {
			continue
		}
		for _, track := range user.PlayedTracks {
			if time.Now().UTC().AddDate(0, 0, -1*int(st.trackWindowInDays)).After(track.PlayedAt) {
				continue
			}
			trackCriteria := tracks[track.Track.ID]
			trackCriteria.count += 1
			if trackCriteria.lastPlayAt.Before(track.PlayedAt) {
				trackCriteria.lastPlayAt = track.PlayedAt
			}
			tracks[track.Track.ID] = trackCriteria
		}
	}

	tracksMapLength := len(tracks)
	if tracksMapLength == 0 {
		return nil
	}
	tracksIDs := make([]spotify.ID, 0, tracksMapLength)
	for trackID := range tracks {
		tracksIDs = append(tracksIDs, trackID)
	}
	sort.Slice(tracksIDs, func(i, j int) bool {
		if tracks[tracksIDs[i]].count == tracks[tracksIDs[j]].count {
			return tracks[tracksIDs[i]].lastPlayAt.After(tracks[tracksIDs[j]].lastPlayAt)
		}

		return tracks[tracksIDs[i]].count > tracks[tracksIDs[j]].count
	})
	if tracksMapLength > 100 {
		tracksIDs = tracksIDs[:100]
	}

	owner, err := (*st.userRepo).GetByID(group.Owner)
	if err != nil {
		return err
	}
	client := (*st.spotifyUtil).GetClient(&owner.Token)

	client.ReplacePlaylistTracks(group.ID, tracksIDs...)

	return nil
}
