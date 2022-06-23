NAME=SoulX
BINDIR=bin
VERSION=$(shell git describe --tags --always || echo "unknown version")
BUILDTIME=$(shell date -u +%F)
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-X "github.com/yaling888/soulx/constant.Version=$(VERSION)" \
                		-X "github.com/yaling888/soulx/constant.BuildTime=$(BUILDTIME)" \
                		-w -s -H=windowsgui -buildid='

WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64 \
	windows-arm64 \
	windows-arm32v7

windows-386:
	mkdir -p $(BINDIR)/386/resources
	mkdir -p $(BINDIR)/386/core
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/386/$(NAME)-$@.exe

windows-amd64:
	mkdir -p $(BINDIR)/amd64/resources
	mkdir -p $(BINDIR)/amd64/core
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/amd64/$(NAME)-$@.exe

windows-arm64:
	mkdir -p $(BINDIR)/arm64/resources
	mkdir -p $(BINDIR)/arm64/core
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/arm64/$(NAME)-$@.exe

windows-arm32v7:
	mkdir -p $(BINDIR)/arm32v7/resources
	mkdir -p $(BINDIR)/arm32v7/core
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD) -o $(BINDIR)/arm32v7/$(NAME)-$@.exe

releases: $(WINDOWS_ARCH_LIST)
