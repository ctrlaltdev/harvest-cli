build:
	cd src; go build -o ../bin/thi .

compile:
	cd src; GOOS=darwin GOARCH=amd64 go build -o ../bin/thi-darwin-amd64 .
	cd src; GOOS=darwin GOARCH=arm64 go build -o ../bin/thi-darwin-arm64 .
	cd src; GOOS=linux GOARCH=amd64 go build -o ../bin/thi-linux-amd64 .
	cd src; GOOS=windows GOARCH=amd64 go build -o ../bin/thi-windows-amd64.exe .

	tar czf bin/thi-darwin-amd64.tar.gz bin/thi-darwin-amd64
	tar czf bin/thi-darwin-arm64.tar.gz bin/thi-darwin-arm64
	tar czf bin/thi-linux-amd64.tar.gz bin/thi-linux-amd64
	tar czf bin/thi-windows-amd64.exe.tar.gz bin/thi-windows-amd64.exe

	sha256sum bin/thi-darwin-amd64.tar.gz > bin/thi-darwin-amd64.tar.gz.sha256
	sha256sum bin/thi-darwin-arm64.tar.gz > bin/thi-darwin-arm64.tar.gz.sha256
	sha256sum bin/thi-linux-amd64.tar.gz > bin/thi-linux-amd64.tar.gz.sha256
	sha256sum bin/thi-windows-amd64.exe.tar.gz > bin/thi-windows-amd64.exe.tar.gz.sha256
