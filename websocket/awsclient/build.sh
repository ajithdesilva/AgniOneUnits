#! /bin/bash

# This how we want version, name the binary output
VESRION=1.0.0
SOURCE=./plugin/awebsocket_client.go
BINARY=./plugin/awebsocket_client.so

# values to pass for BinVersion, GitCommitLog, GitStatus, BuildTime and BuildGoVersion"
# Version=`git describe --tags`  # git tag 1.0.1  # require tag tagged before

BuildTime=`date`
BuildGoVersion=`go version`


# Setup the -ldflags option for build 
#LDFLAGS=" -s -w -X 'awebsocket.client/src/build.Version=${VESRION}' \
LDFLAGS=" -X 'awebsocket.client/src/build.Version=${VESRION}' \
-X 'awebsocket.client/src/build.User=$(id -u -n)' \
-X 'awebsocket.client/src/build.Time=${BuildTime}' \
-X 'awebsocket.client/src/build.BuildGoVersion=${BuildGoVersion}' "


###echo ${LDFLAGS}

cd ./src
echo "clean old binaries....."
rm ${BINARY}
echo "clean old binaries ......... DONE"

echo "building plug-in....."

go build -v -buildmode=plugin -ldflags="${LDFLAGS}" -o ${BINARY} ${SOURCE}
echo "building plug-in ......... DONE"

cp ${BINARY} $1/plugins/websocket
echo "copy websocket plug-in ${BINARY} -> $1/plugins/websocket/${BINARY} DONE"
cd ..
