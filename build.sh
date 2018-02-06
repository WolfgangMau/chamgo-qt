#!/bin/bash
# crossbuild-information taken from https://github.com/therecipe/qt/wiki
rm -rf deploy/*
#
echo "Linux"
#docker pull therecipe/qt:linux
qtdeploy -docker build linux
#zip -r deploy/chamgo-qt-linux.zip deploy/linux
#
#echo "windows64"
##docker pull therecipe/qt:windows_64_static
#qtdeploy -docker build windows_64_static
#mv deploy/windows deploy/windows64
#zip -r deploy/chamgo-qt-win64.zip deploy/windows64
#
#echo "windows32"
##docker pull therecipe/qt:windows_32_static
#qtdeploy -docker build windows_32_static
#mv deploy/windows deploy/windows32
#zip -r deploy/chamgo-qt-win32.zip deploy/windows32
#
#echo "darwin"
#qtdeploy
#zip -r deploy/chamgo-qt-osx.zip deploy/darwin
