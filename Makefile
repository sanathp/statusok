all: vendor-build

build:
	go build

vendor-build:
	go build -mod=vendor


clean:
	rm -fr statusok
