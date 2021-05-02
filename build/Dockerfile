FROM golang:alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v -coverpkg=./pkg/... -coverprofile=coverage.out ./... \
    && go tool cover -func=coverage.out

RUN go build -o ./deckify

FROM alpine
RUN apk add ca-certificates

COPY --from=build_base /tmp/app/deckify /app/deckify

WORKDIR /app

EXPOSE 8080

CMD ["./deckify"]