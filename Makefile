default: debug

debug:
	GO15VENDOREXPERIMENT=1 go install -ldflags "-w" github.com/dimiro1/ipe

run-debug: debug
	${GOPATH}/bin/ipe --config ${GOPATH}/src/github.com/dimiro1/ipe/config.json -logtostderr=true -v=2

test:
	GO15VENDOREXPERIMENT=1 go test `go list ./... | grep -v vendor`

macos:
	GO15VENDOREXPERIMENT=1 GOOS=darwin GOARCH=amd64 go install github.com/dimiro1/ipe

linux:
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=amd64 go install github.com/dimiro1/ipe

windows:
	GO15VENDOREXPERIMENT=1 GOOS=windows GOARCH=amd64 go install github.com/dimiro1/ipe

raspberry:
	GO15VENDOREXPERIMENT=1 GOOS=linux GOARCH=arm go install github.com/dimiro1/ipe
