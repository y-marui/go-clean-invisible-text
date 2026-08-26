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
	git remote | grep -q '^dev-charter$$' || \
	  git remote add dev-charter https://github.com/y-marui/dev-charter
	git fetch dev-charter
	@STASHED=0; \
	if ! git diff --quiet || ! git diff --cached --quiet || [ -n "$$(git ls-files --others --exclude-standard)" ]; then \
		git stash push -u -m "update-charter"; \
		STASHED=1; \
	fi; \
	git subtree pull --prefix=docs/dev-charter dev-charter main --squash; \
	if [ "$$STASHED" = "1" ]; then git stash pop; fi
