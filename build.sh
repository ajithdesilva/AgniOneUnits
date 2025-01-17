#! /bin/bash

# This how we want version, name the binary output
UNIT_NAME=demohttp
VESRION=1.0.0
SOURCE="./unit/demohttp_main.go ./unit/demohttp.go"
BINARY=./unit/demohttp.so

# values to pass for BinVersion, GitCommitLog, GitStatus, BuildTime and BuildGoVersion"
# Version=`git describe --tags`  # git tag 1.0.1  # require tag tagged before

BuildTime=`date`
BuildGoVersion=`go version`


# Setup the -ldflags option for build 
LDFLAGS=" -s -w -X 'unit.demo.http/src/build.Version=${VESRION}' \
-X 'unit.demo.http/src/build.User=$(id -u -n)' \
-X 'unit.demo.http/src/build.Time=${BuildTime}' \
-X 'unit.demo.http/src/build.BuildGoVersion=${BuildGoVersion}' "

echo ${LDFLAGS}

cd ./src
echo "clean old binaries....."
rm ${BINARY}
echo "clean old binaries ......... DONE"

echo "building plug-in....."

go build -v -buildmode=plugin -ldflags="${LDFLAGS}" -o ${BINARY} ${SOURCE}

echo "building plug-in ......... DONE"


echo "Deploying ${UNIT_NAME} to $1/apps/config/${UNIT_NAME}"

cp ${BINARY} $1/apps/units

mkdir -p $1/apps/config/${UNIT_NAME}

### make the unit & config path in the app.config


sed -i  '/"path":"",/c\\t\t\t\t"path":"$1/apps/units/${BINARY}",' ./config/app.config
sed -i  '/"config":"",/c\\t\t\t\t"path":"$1/apps/config/${UNIT_NAME}/${UNIT_NAME}.config",' ./config/app.config

cp ./config/*.config $1/apps/config/${UNIT_NAME}

echo "Deploying ${UNIT_NAME} ................... DONE"
cd ..