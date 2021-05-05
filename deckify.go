package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/silvanoneto/deckify/pkg/collector"
	"github.com/silvanoneto/deckify/pkg/group"
	"github.com/silvanoneto/deckify/pkg/spotifyutil"
	"github.com/silvanoneto/deckify/pkg/stacker"
	"github.com/silvanoneto/deckify/pkg/user"
)

func main() {
	godotenv.Load()

	value := os.Getenv("DECKIFY_COLLECTOR_PAGESIZE")
	if value == "" {
		log.Fatalln("No collector page size")
	}
	collectorPageSize, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Collector page size value is invalid")
	}

	value = os.Getenv("DECKIFY_COLLECTOR_CALL_INTERVAL_TIME_IN_SECONDS")
	if value == "" {
		log.Fatalln("No collector call interval time")
	}
	collectorCallIntervalTime, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Collector call interval time value is invalid")
	}

	value = os.Getenv("DECKIFY_STACKER_PAGESIZE")
	if value == "" {
		log.Fatalln("No stacker page size")
	}
	stackerPageSize, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Stacker page size value is invalid")
	}

	value = os.Getenv("DECKIFY_STACKER_CALL_INTERVAL_TIME_IN_SECONDS")
	if value == "" {
		log.Fatalln("No stacker call interval time")
	}
	stackerCallIntervalTime, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Stacker call interval time value is invalid")
	}

	value = os.Getenv("DECKIFY_STACKER_TRACK_WINDOW_IN_DAYS")
	if value == "" {
		log.Fatalln("No stacker track window")
	}
	stackerTrackWindow, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalln("Stacker track window value is invalid")
	}

	userRepo := user.NewUserRepoInMemoryImpl()
	groupRepo := group.NewGroupRepoInMemoryImpl()

	var spotifyUtil spotifyutil.SpotifyUtil = spotifyutil.NewSpotifyUtilDefaultImpl(&userRepo, &groupRepo,
		"http://localhost:8080/callback", "default")

	var collectorInstance collector.Collector = collector.NewLazyCollector(&userRepo, &spotifyUtil,
		uint(collectorPageSize), uint(collectorCallIntervalTime))
	go collectorInstance.Start()

	var stackerInstance stacker.Stacker = stacker.NewLazyStacker(&userRepo, &groupRepo, &spotifyUtil,
		uint(stackerPageSize), uint(stackerCallIntervalTime), uint(stackerTrackWindow))
	go stackerInstance.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request for: %s", r.URL.String())
		http.Redirect(w, r, spotifyUtil.GetAuthURL(),
			http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/callback", spotifyUtil.AuthCallback)

	addr := ":8080"
	log.Println("Deckify started on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
