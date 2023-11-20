run:
	@export $(shell cat .env | xargs) && go run main.go -race
