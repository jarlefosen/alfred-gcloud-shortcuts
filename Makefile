

all: dep test butler-darwin

dep:
	cd gcloud-butler && dep ensure

test:
	go test -v ./gcloud-butler/...

butler-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s -w' -o bin/butler-darwin gcloud-butler/main.go
