GO ?= go
QTDEPLOY ?= $(shell go env GOPATH)/bin/qtdeploy

all: deploy

deploy: $(QTDEPLOY)
	GO111MODULE=off $(QTDEPLOY) build desktop ./bin/firecraft/main.go
	rm -rf ./bin/firecraft/linux
	mv ./bin/firecraft/deploy ./deploy

$(QTDEPLOY):
	GO111MODULE=off go get -v github.com/therecipe/qt/cmd/...

%: bin/%/main.go
	${GO} build -o $@ $^

clean:
	${GO} clean
	rm -rf mcAuth firecraft deploy

linux-install: deploy
	mkdir -p ~/.minecraft
	cp -pr ./deploy/linux ~/.minecraft/launcher
	desktop-file-install --dir=$(HOME)/.local/share/applications ./share/applications/firecraft.desktop