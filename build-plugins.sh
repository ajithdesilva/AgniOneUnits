#! /bin/bash
echo "./build_plugins.sh <BASE_PATH_TO_DEPLOY>"

echo "deployment path -- > $1"
echo ""
echo "Building plugins......"
echo ""

echo "building http plugin....."
cd http/ahttpclient
./build.sh $1
echo "building http plugin..... DONE"

cd ..
cd ..


echo "building websocket plugin....."
cd websocket/awsclient
./build.sh $1
echo "building websocket plugin.....DONE"


echo "Building plugins......COMPLETED"


