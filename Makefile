# Albert Espí­n
# set variables
BINDIR = ./bin
VERSION = latest
BINARYFILE = go-ikea-server
DOCKERFILE = go-ikea-server
DOCKERREGISTRY = ea3hsp
DOCKERPUSH = $(DOCKERREGISTRY)/$(DOCKERFILE):$(VERSION)
ARCH = amd64
OS = linux

.PHONY: clean

build:
	go build -o $(BINDIR)/$(BINARYFILE) -i cmd/main.go
run:build
	$(BINDIR)/$(DOCKERFILE)
clean:
	rm -f .$(BINDIR)/$(BINARYFILE)
docker:
	env GOARCH=$(ARCH) GOOS=$(OS) go build -o $(BINDIR)/$(BINARYFILE) -i cmd/main.go
	docker build -t $(DOCKERFILE) .
docker-push:
	docker tag $(DOCKERFILE) $(DOCKERPUSH)
	docker push $(DOCKERPUSH)