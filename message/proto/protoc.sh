#!/bin/sh
rm message.pb.go
protoc --go_out=plugins=grpc:. *.proto
go install