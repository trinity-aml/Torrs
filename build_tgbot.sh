#!/bin/bash
##rm -fr ./tgbotapi
##git clone --recursive https://github.com/tdlib/telegram-bot-api.git tgbotapi
cd ./tgbotapi
rm -fr ./build
mkdir ./build
cd ./build
cmake -DCMAKE_BUILD_TYPE=Release ..
cmake --build .
cp ./telegram-bot-api ../../dist