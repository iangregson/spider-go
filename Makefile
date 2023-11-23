all : clean build

clean: 
	rm bin/*

build:
	go build -o bin/spider-go

test:
	go test ./crawler/... ./parser/... ./urls/... 
