default: test

.PHONY: test
test:
	docker compose -f docker/docker-compose.test.yml up


# docker run -v ./:/app/src --workdir /app/src --rm -it girmes go run cmd/main.go -org girmes
.PHONY: image
image:
	docker build -f docker/Dockerfile -t girmes .
	docker run -v ./:/app/src --workdir /app/src --rm -it girmes sh
