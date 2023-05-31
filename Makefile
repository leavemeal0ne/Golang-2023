.PHONY: run
run :
	go run .\cmd\web\

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: test_cover
test_cover:
			go test -cover ./...
.DEFAULT.GOAL := run