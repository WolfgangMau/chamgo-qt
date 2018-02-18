#!/bin/bash
 crossbuild-information taken from https://github.com/therecipe/qt/wiki
rm -rf deploy/*
#
echo "Linux"
docker pull therecipe/qt:linux
qtdeploy -docker build linux
mkdir deploy/linux/config
cp config/config.yaml deploy/linux/config/
#zip -r deploy/chamgo-qt-linux.zip deploy/linux
#
#echo "windows64"
##docker pull therecipe/qt:windows_64_static
#qtdeploy -docker build windows_64_static
#mv deploy/windows deploy/windows64
#zip -r deploy/chamgo-qt-win64.zip deploy/windows64
#
echo "windows32"
docker pull therecipe/qt:windows_32_static
qtdeploy -docker build windows_32_static
mv deploy/windows deploy/windows32
mkdir deploy/windows/config
cp config/config.yaml deploy/windows/config/
#zip -r deploy/chamgo-qt-win32.zip deploy/windows32

echo "darwin"
qtdeploy
mkdir deploy/darwin/chamgo-qt.app/Contents/MacOS/config
cp config/config.yaml deploy/darwin/chamgo-qt.app/Contents/MacOS/config/
cp -a bin macros  /Users/wolfgang/Development/workspace-go/src/github.com/WolfgangMau/chamgo-qt/deploy/darwin/chamgo-qt.app/Contents/MacOS/

#zip -r deploy/chamgo-qt-osx.zip deploy/darwin
