GOPATH="$(shell echo $$GOPATH)"
BINDIR = bin

BINFILE=go

APPS = go-auth
all: $(APPS)

$(APPS): $(BINDIR)
	@echo "build $@"
	@export
	@echo "gopath $(GOPATH)"
	@GOPATH=$(GOPATH) $(BINFILE) build -o $(BINDIR)/$@ cmd/*.go

test:
	@GOPATH=$(GOPATH) go test ./...

$(BINDIR):
	mkdir -p $(BINDIR)
