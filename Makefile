GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test

all: build

test:
	$(GOTEST) -v ./...

build:
	$(GOBUILD) -o build/detector ./cmd/github-detector

linux32:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) -o build/detector_linux_i386 ./cmd/github-detector

linux64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o build/detector_linux_amd64 ./cmd/github-detector

darwin32:
	CGO_ENABLED=0 GOOS=darwin GOARCH=386 $(GOBUILD) -o build/detector_darwin_i386 ./cmd/github-detector

darwin64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o build/detector_darwin_amd64 ./cmd/github-detector

windows32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 $(GOBUILD) -o build/detector_windows_i386 ./cmd/github-detector
	
windows64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o build/detector_windows_amd64 ./cmd/github-detector

clean:
	$(GOCLEAN)
	rm -rf build