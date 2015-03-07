default: debug

debug:
	go install -ldflags "-w" github.com/keithmattix/ipe

run-debug: debug
	${GOPATH}/bin/ipe --config ${GOPATH}/src/github.com/keithmattix/ipe/config.json -logtostderr=true -v=2

test:
	go test github.com/keithmattix/ipe...

macos:
	GOOS=darwin GOARCH=amd64 go install github.com/keithmattix/ipe

linux:
	GOOS=linux GOARCH=amd64 go install github.com/keithmattix/ipe

windows:
	GOOS=windows GOARCH=amd64 go install github.com/keithmattix/ipe

raspberry:
	GOOS=linux GOARCH=arm go install github.com/keithmattix/ipe
