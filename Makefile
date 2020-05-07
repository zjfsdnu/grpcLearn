.PHONY: proto

proto:
	protoc --go_out=plugins=grpc:. proto/*.proto
	#protoc --go_out=. proto/*.proto