build-sandbox-proto:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	proto/sandbox.proto

docker-build-worker:
	docker build -t worker -f dockerfiles/Dockerfile.worker .

docker-run-worker:
	docker run --privileged --name worker -p 50051:50051 -v /var/run/docker.sock:/var/run/docker.sock -v sandbox:/app/source -it worker

debug-worker:
	docker-build-worker
	docker-run-worker

docker-build-python-worker:
	docker build -t worker-python:v0.0.1 -f dockerfiles/Dockerfile.sandbox.python .
