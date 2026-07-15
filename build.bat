@echo off
chcp 65001 > nul

echo Сборка Base-версии...
go build -o bin/minecraft_server_analyser_base.exe ./src/base

echo Сборка Extended-версии...
go build -o bin/minecraft_server_analyser_extended.exe ./src/extended

echo Готово!
pause