server:
	go run main.go
proto:
	rm -f pb/*
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=f \
    proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl
	
.PHONY: server proto evans