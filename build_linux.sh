#!/bin/bash

# Создаем папку bin, если её нет
mkdir -p bin

echo "Сборка Minimal-версии под Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/linux/minimal/minecraft_server_analyser_base ./src/base

echo "Сборка Base-версии под Linux..."
GOOS=linux GOARCH=amd64 go build -o bin/linux/base/minecraft_server_analyser_base ./src/base

echo "Сборка Extended-версии под Linux..."
echo "Внимание: extended версия работает нестабильно на Linux!"
GOOS=linux GOARCH=amd64 go build -o bin/linux/extended/minecraft_server_analyser_extended ./src/extended

echo "Сборка завершена!"
