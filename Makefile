GO ?= go
QTDEPLOY ?= $(subst \,/,$(shell go env GOPATH)/bin/qtdeploy)

all: deploy

deploy: clean $(QTDEPLOY)
	mkdir -p $(shell go env GOPATH)/src/github.com/zapomnij
	rm -rf $(shell go env GOPATH)/src/github.com/zapomnij/firecraft
	cp -pr $(shell pwd) $(shell go env GOPATH)/src/github.com/zapomnij/firecraft
	GO111MODULE=off $(QTDEPLOY) build desktop .
	rm -rf linux windows darwin

$(QTDEPLOY):
	GO111MODULE=off go get -v github.com/therecipe/qt/cmd/...

firecraft:
	go build -v .

%: bin/%/main.go
	${GO} build -o $@ $^

clean:
	${GO} clean
	rm -rf mcAuth firecraft deploy

linux-install: deploy
	mkdir -p ~/.minecraft
	rm -rf ~/.minecraft/launcher
	cp -pr ./deploy/linux ~/.minecraft/launcher
	desktop-file-install --dir=$(HOME)/.local/share/applications ./share/applications/firecraft.desktop