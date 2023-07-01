GO ?= go
QTDEPLOY ?= $(shell go env GOPATH)/bin/qtdeploy

all: go.sum deploy

go.sum:
	go get -v .

deploy: clean go.sum $(QTDEPLOY)
	mkdir -p ~/go/src/github.com/zapomnij
	cp -pr $(shell pwd) $(HOME)/go/src/github.com/zapomnij/firecraft
	GO111MODULE=off $(QTDEPLOY) build desktop .
	rm -rf linux windows macos

$(QTDEPLOY):
	GO111MODULE=off go get -v github.com/therecipe/qt/cmd/...

firecraft: go.sum
	go build -v .

%: go.sum bin/%/main.go
	${GO} build -o $@ $^

clean:
	${GO} clean
	rm -rf mcAuth firecraft deploy

linux-install: deploy
	mkdir -p ~/.minecraft
	cp -pr ./deploy/linux ~/.minecraft/launcher
	desktop-file-install --dir=$(HOME)/.local/share/applications ./share/applications/firecraft.desktop