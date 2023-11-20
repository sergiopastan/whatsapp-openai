run:
	@export $(shell cat .env | xargs) && go run -race .

.PHONY: run
