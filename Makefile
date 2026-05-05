.PHONY: build serve clean

build:
	@go run ./cmd/build

serve: build
	@cd dist && python3 -m http.server 8000

clean:
	@rm -rf dist
