_NAME=harvest

build:
	go build -o bin/$(_NAME) .

compile:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/$(_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -o bin/$(_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/$(_NAME)-linux-arm64 .
	# GOOS=windows GOARCH=amd64 go build -o bin/$(_NAME)-windows-amd64.exe .

	tar czf bin/$(_NAME)-darwin-amd64.tar.gz bin/$(_NAME)-darwin-amd64
	tar czf bin/$(_NAME)-darwin-arm64.tar.gz bin/$(_NAME)-darwin-arm64
	tar czf bin/$(_NAME)-linux-amd64.tar.gz bin/$(_NAME)-linux-amd64
	tar czf bin/$(_NAME)-linux-arm64.tar.gz bin/$(_NAME)-linux-arm64
	# tar czf bin/$(_NAME)-windows-amd64.exe.tar.gz bin/$(_NAME)-windows-amd64.exe

	sha256sum bin/$(_NAME)-darwin-amd64.tar.gz > bin/$(_NAME)-darwin-amd64.tar.gz.sha256
	sha256sum bin/$(_NAME)-darwin-arm64.tar.gz > bin/$(_NAME)-darwin-arm64.tar.gz.sha256
	sha256sum bin/$(_NAME)-linux-amd64.tar.gz > bin/$(_NAME)-linux-amd64.tar.gz.sha256
	sha256sum bin/$(_NAME)-linux-arm64.tar.gz > bin/$(_NAME)-linux-arm64.tar.gz.sha256
	# sha256sum bin/$(_NAME)-windows-amd64.exe.tar.gz > bin/$(_NAME)-windows-amd64.exe.tar.gz.sha256
