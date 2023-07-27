# builds and tests project via go tools
# installsolver plugins in ricrob (server) directory
all:
	@echo "update dependencies"
	go get -u ./...
	go mod tidy
	@echo "build and test"
	go build -v ./...
	go vet ./...
	golint -set_exit_status=true ./...
	staticcheck -checks all -fail none ./...
	go test ./...
	@echo "install wasm"
#	GOOS=js GOARCH=wasm go build -o ./ricrob/assets/ricrob.wasm ./ricrob/wasm
#see fsfe reuse tool (https://git.fsfe.org/reuse/tool)
	@echo "reuse (license) check"
	pipx run reuse lint

#install additional tools
tools:
#install linter
	@echo "install latest go linter version"
	go install golang.org/x/lint/golint@latest
#install staticcheck
	@echo "install latest staticcheck version"
	go install honnef.co/go/tools/cmd/staticcheck@latest
