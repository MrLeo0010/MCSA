@echo off
chcp 65001 > nul

echo Сборка Minimal-версии...
go build -o bin/windows/minimal/minecraft_server_analyser_minimal.exe ./src/minimal

echo Сборка Base-версии...
go build -o bin/windows/base/minecraft_server_analyser_base.exe ./src/base

echo Сборка Extended-версии...
go build -o bin/windows/extended/minecraft_server_analyser_extended.exe ./src/extended

pause
