.PHONY: run
run :
	go run .\cmd\web\

.PHONY: test
test:
	go test -v ./..

.DEFAULT.GOAL := run