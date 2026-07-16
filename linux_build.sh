#!/bin/bash

# Создаем папку bin, если её нет
mkdir -p bin

echo "Сборка Base-версии под Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/minecraft_server_analyser_base ./src/base

echo "Сборка Extended-версии под Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/minecraft_server_analyser_extended ./src/extended

echo "Сборка завершена!"
