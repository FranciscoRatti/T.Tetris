@echo off
setlocal EnableDelayedExpansion

net session >nul 2>&1
if %errorlevel% neq 0 (
    powershell -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

if exist "%ProgramData%\Microsoft\Windows\Start Menu\Programs\TTetris" rmdir /s /q "%ProgramData%\Microsoft\Windows\Start Menu\Programs\TTetris"
if exist "%AppData%\TTetris" rmdir /s /q "%AppData%\TTetris"
if exist "%ProgramFiles%\TTetris" rmdir /s /q "%ProgramFiles%\TTetris"