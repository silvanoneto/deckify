package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/silvanoneto/deckify/pkg/collector"
	"github.com/silvanoneto/deckify/pkg/spotifyutil"
	"github.com/silvanoneto/deckify/pkg/user"
)

func main() {
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

	userRepo := user.NewUserRepoInMemoryImpl()

	var sUtil spotifyutil.SpotifyUtil = spotifyutil.NewSpotifyUtilDefaultImpl(
		&userRepo, "http://localhost:8080/callback", "default")

	var cInstance collector.Collector = collector.NewLazyCollector(
		&userRepo, &sUtil, uint(collectorPageSize),
		uint(collectorCallIntervalTime))
	go cInstance.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Got request for: %s", r.URL.String())
		http.Redirect(w, r, sUtil.GetAuthURL(),
			http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/callback", sUtil.AuthCallback)

	addr := ":8080"
	log.Println("Deckify started on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
