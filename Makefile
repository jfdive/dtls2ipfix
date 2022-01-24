DIST_DIR= ./
OUTPUT_DTLS2IPFIX   := $(DIST_DIR)dtls2ipfix

.PHONY: all
all: build

.PHONY: test
test:
	@echo testing code
	go test ./...

.PHONY: clean
clean:
	rm  -f $(OUTPUT_DTLS2IPFIX)

.PHONY: build
build:
	go build -o $(OUTPUT_DTLS2IPFIX) -gcflags="all=-N -l" ./dtls2ipfix.go

