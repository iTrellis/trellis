#!/bin/sh
rm *.go
protoc --go_out=plugins=grpc:. *.proto
easyjson -all *.pb.go
go install