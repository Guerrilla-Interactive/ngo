main: ngo
	./ngo

ngo: main.go
	go build .

install:
	go install .

clean:
	rm ngo
