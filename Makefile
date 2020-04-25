build-client:
	mkdir -p dist/
	go build -o dist/videoparty cmd/client/*.go
