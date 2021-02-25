#!/bin/sh
SRCDIR=$GOPATH/src
DIR=`pwd`

echo "DIR $DIR"
echo "SRCDIR $SRCDIR"

find $DIR -path ${DIR}/vendor -prune -o -name '*.pb.go' -exec rm {} \;
find $DIR -path ${DIR}/vendor -prune -o -name '*.proto' -exec echo {} \;
find $DIR/proto -path ${DIR}/vendor -prune -o -name '*.proto' -exec protoc --go_out=plugins=grpc:${SRCDIR} -I=${DIR}/proto -I=. {} \;

