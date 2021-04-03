@ECHO OFF

protoc --go_out=. --go-grpc_out=. proto/session_service.proto
protoc --go_out=. --go-grpc_out=. proto/user.proto
