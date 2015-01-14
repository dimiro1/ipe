default: debug

debug:
	go install -ldflags "-w" github.com/dimiro1/ipe

run-debug: debug
	${GOPATH}/bin/ipe --config ${GOPATH}/src/github.com/dimiro1/ipe/config.json -logtostderr=true -v=2

test:
	go test github.com/dimiro1/ipe

macos:
	GOOS=darwin GOARCH=amd64 go install github.com/dimiro1/ipe

linux:
	GOOS=linux GOARCH=amd64 go install github.com/dimiro1/ipe

windows:
	GOOS=windows GOARCH=amd64 go install github.com/dimiro1/ipe

raspberry:
	GOOS=linux GOARCH=arm go install github.com/dimiro1/ipe
