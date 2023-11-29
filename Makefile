
build-proto:
	protoc --go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=worker/proto \
	--go-grpc_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
	worker/proto/*.proto
load-image: