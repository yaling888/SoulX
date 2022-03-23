NAME=SoulService
BINDIR=../bin
GOBUILD=go build -trimpath -ldflags '-w -s -buildid='

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe
