#! /bin/bash

# This how we want version, name the binary output
VESRION=1.0.0
SOURCE=./plugin/ahttpclient.go
BINARY=./plugin/ahttpclient.so

# values to pass for BinVersion, GitCommitLog, GitStatus, BuildTime and BuildGoVersion"
# Version=`git describe --tags`  # git tag 1.0.1  # require tag tagged before

BuildTime=`date`
BuildGoVersion=`go version`


# Setup the -ldflags option for build 

LDFLAGS=" -s -w -X 'ahttp.client/src/build.Version=${VESRION}' \
-X 'ahttp.client/src/build.User=$(id -u -n)' \
-X 'ahttp.client/src/build.Time=${BuildTime}' \
-X 'ahttp.client/src/build.BuildGoVersion=${BuildGoVersion}' "

echo ${LDFLAGS}

cd ./src
echo "clean old binaries....."
rm ${BINARY}
echo "clean old binaries ......... DONE"

echo "building plug-in....."

GO111MODULE=on go build -v -buildmode=plugin -ldflags="${LDFLAGS}" -o  ${BINARY} ${SOURCE}
echo "building plug-in ......... DONE"



cp ${BINARY} $1/plugins/http
echo "plugin ${BINARY} copied to $1/plugins/http "

