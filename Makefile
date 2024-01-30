run:
	go run .
build:
	GO111MODULE=on CGO_ENABLED=1 go build -o rmediy .
build-win:
	 GO111MODULE=on GOOS=windows CGO_ENABLED=1 go build -o rmediy.exe .
build-darwin:
	 GO111MODULE=on GOOS=darwin GOARCH=amd64 CC=clang  CXX=clang go build -o rmediy . 
install:
	 CGO_ENABLED=1 go install .