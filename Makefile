build-sandbox-proto:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	proto/sandbox.proto

docker-build-worker:
	docker build -t worker:v0.0.3 -f dockerfiles/Dockerfile.worker .

docker-build-master:
	docker build -t master:v0.0.1 -f dockerfiles/Dockerfile.master .

docker-build-python-sandbox:
	docker build -t sandbox-python:v0.0.1 -f dockerfiles/Dockerfile.sandbox.python .

docker-build-golang-sandbox:
	docker build -t sandbox-golang:v0.0.1 -f dockerfiles/Dockerfile.sandbox.golang .

docker-build-frontend:
	 docker build -t  frontend:v0.0.1 -f dockerfiles/Dockerfile.frontend .