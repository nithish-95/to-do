all: clean build run

build:
	go build -o bin/todo .
run: 
	bin/todo
clean:
	go mod tidy
	rm bin/* || true