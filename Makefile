run:
	go run main.go

build:
	go build -o APIProbe.exe main.go

zip:
	rm -f APIProbe.zip && go build && "/c/Program Files/7-Zip/7z.exe" a APIProbe.zip APIProbe.exe && rm APIProbe.exe

lint:
	golangci-lint run ./...

test:
	go test -v ./...

cover:
	go test -v ./... -coverprofile=./coverage/coverage.out ./...
	go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
