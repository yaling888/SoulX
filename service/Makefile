NAME=SoulService
BINDIR=../bin
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s -buildid='

WINDOWS_ARCH_LIST = \
	windows-386 \
	windows-amd64 \
	windows-arm64 \
	windows-arm32v7

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/386/$(NAME)-$@.exe

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/amd64/$(NAME)-$@.exe

windows-arm64:
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/arm64/$(NAME)-$@.exe

windows-arm32v7:
	GOARCH=arm GOOS=windows GOARM=7 $(GOBUILD) -o $(BINDIR)/arm32v7/$(NAME)-$@.exe

releases: $(WINDOWS_ARCH_LIST)
