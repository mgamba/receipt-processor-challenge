run:
	go run main.go

dev:
	go run main.go -mode development

b:
	go build

tidy:
	go mod tidy

shell:
	docker compose run --rm -it -v `pwd`\:/usr/src/app -p 3333\:3333 api /bin/sh

t:
	go test ./... -v
