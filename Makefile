test:
	go test ./...

test-cover:
	go tool cover -html cover.out

test-cover-cli:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out
