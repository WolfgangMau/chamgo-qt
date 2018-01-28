#!/bin/bash
# crossbuild-information taken from https://github.com/therecipe/qt/wiki


echo "Linux"
docker pull therecipe/qt:linux
qtdeploy -docker build linux

echo "windows64"
docker pull therecipe/qt:windows_64_static
qtdeploy -docker build windows_64_static
mv deploy/windows deploy/windows64

echo "windows32"
docker pull therecipe/qt:windows_32_static
qtdeploy -docker build windows_32_static
mv deploy/windows deploy/windows32

echo "darwin"
qtdeploy