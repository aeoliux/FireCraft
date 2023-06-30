GO ?= go

all: mcAuth firecraft

%: bin/%/main.go
	${GO} build -o $@ $^

clean:
	${GO} clean
	rm -f mcAuth firecraft