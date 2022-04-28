protoc --gofast_out=protos/pb protos/pbdef/*.proto --proto_path=protos/pbdef
go run ../../cmd/metactl/main.go --module=github.com/jageros/metactl/example/game