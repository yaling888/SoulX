NAME=SoulX
BINDIR=bin
GOBUILD=CGO_ENABLED=1 go build -trimpath -ldflags '-w -s -H=windowsgui -buildid='

WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64 \
	windows-arm64 \
	windows-arm32v7

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-arm64:
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-arm32v7:
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

zip_releases=$(addsuffix .zip, $(WINDOWS_ARCH_LIST))

$(zip_releases): %.zip : %
	zip -m -j $(BINDIR)/$(NAME)-$(basename $@)-$(VERSION).zip $(BINDIR)/$(NAME)-$(basename $@).exe

all-arch: $(WINDOWS_ARCH_LIST)

releases: $(zip_releases)
