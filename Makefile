dc-build:
	docker-compose build
dc-start:
	docker-compose up -d
dc-stop:
	docker-compose down
dc-logs:
	docker-compose logs -f
test:
	docker run --name chat-redis-test -d --rm -p 60001:6379 redis:5
	go test -v ./...
	docker stop chat-redis-test