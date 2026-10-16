include .env
export

.PHONY: run

run:
	go run -tags fts5 .