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

linux-install:
	[ ! -d ./deploy/linux ] && make deploy || exit 0
	mkdir -p $(HOME)/.minecraft
	rm -rf $(HOME)/.minecraft/launcher
	cp -pr ./deploy/linux $(HOME)/.minecraft/launcher
	desktop-file-install --dir=$(HOME)/.local/share/applications ./share/applications/firecraft.desktop

linux-uninstall:
	[ -f $(HOME)/.minecraft/launcher/FireCraft -o -f $(HOME)/.minecraft/launcher/firecraft ] && \
		rm -rf $(HOME)/.minecraft/launcher && \
		rm -f $(HOME)/.local/share/applications/firecraft.desktop $(HOME)/.local/share/applications/FireCraft.desktop || exit 0