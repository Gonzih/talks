build:
	GOARCH=wasm GOOS=js go build -o test.wasm
	@echo

server-main:
	env GO111MODULE=on go build -o server-main server/main.go

run-server: server-main
	@echo http://localhost:3000/wasm_exec.html
	./server-main

auto-build:
	reflex -r '.*\.go' make build
