build:
	GOARCH=wasm GOOS=js go build -o test.wasm

server-main:
	go build -o server-main server/main.go

run-server: server-main
	@echo http://localhost:3000/wasm_exec.html
	./server-main
