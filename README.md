# deckify

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/silvanoneto/go-learning.svg)](https://github.com/silvanoneto/deckify)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/silvanoneto/deckify)
[![GoReportCard example](https://goreportcard.com/badge/github.com/silvanoneto/deckify)](https://goreportcard.com/report/github.com/silvanoneto/deckify)
[![GitHub license](https://img.shields.io/github/license/silvanoneto/deckify.svg)](https://github.com/silvanoneto/deckify/main/LICENSE)

Deckify collects the recently played musics from spotify users

### How to run it locally

You'll need to install Go in your machine. Follow the instructions in [https://golang.org/doc/install](https://golang.org/doc/install) or use the package manager of your preference.

After this:
1. Download the source code;
2. Open your terminal/prompt within the project folder path;
3. Create the following environment variables:
    - SPOTIFY_ID={your spotify client id}
    - SPOTIFY_SECRET={your spotify client secret}
    - DECKIFY_COLLECTOR_PAGESIZE=100
    - DECKIFY_COLLECTOR_CALL_INTERVAL_TIME_IN_SECONDS=5
3. Run the following commands:
```sh
go mod download
go run deckify.go
```

### How to run it locally (Docker)

1. Create an .env file and set the following environment variables:
    - SPOTIFY_ID={your spotify client id}
    - SPOTIFY_SECRET={your spotify client secret}
    - DECKIFY_COLLECTOR_PAGESIZE=100
    - DECKIFY_COLLECTOR_CALL_INTERVAL_TIME_IN_SECONDS=5
2. Run the following commands:
```sh
docker pull silvanoneto/deckify
docker run --env-file .env -p 8080:8080 --rm silvanoneto/deckify
```
