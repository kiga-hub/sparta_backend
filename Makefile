PROJECT_NAME=websocket
API_ROOT=/api

.PHONY: image swagger-client swagger-server

image:
	docker build -t ${PROJECT_NAME}:dev .

run:
	docker run --rm \
	-p 8000:8000 \
	${PROJECT_NAME}:dev

