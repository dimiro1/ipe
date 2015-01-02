default: install-debug

install-debug:
	go install -ldflags "-w" github.com/dimiro1/ipe

run-debug: install-debug
	${GOPATH}/bin/ipe --config ${GOPATH}/src/github.com/dimiro1/ipe/config-example.json

test:
	go test github.com/dimiro1/ipe
