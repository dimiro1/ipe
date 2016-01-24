default: debug

debug:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-w" github.com/dimiro1/ipe

run-debug: debug
	${GOPATH}/bin/ipe --config ${GOPATH}/src/github.com/dimiro1/ipe/config.json -logtostderr=true -v=2

test:
	GO15VENDOREXPERIMENT=1 go test -bench . `go list ./... | grep -v vendor | grep -v functional | grep -v github.com/dimiro1/ipe$$`

dev-deps:
	go get github.com/pusher/pusher-http-go
