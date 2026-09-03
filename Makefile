.PHONY: build test lint fmt precommit update-charter

build:
	go build ./...

test:
	go test ./...

lint:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to be run on:"; \
		echo "$$out"; \
		exit 1; \
	fi
	go vet ./...

fmt:
	gofmt -w .

precommit:
	pre-commit run --all-files

update-charter:
	CHARTER_UPDATE_ONLY=1 bash <(curl -fsSL https://raw.githubusercontent.com/y-marui/dev-charter/main/scripts/install.sh)
