@echo off
setlocal EnableDelayedExpansion
set "CMD=%~dp0"
for %%I in ("%CMD%\..") do set "SRC=%%~fI"

:: Permisos de admin
net session >nul 2>&1
if %errorlevel% neq 0 (
    powershell -Command "Start-Process -FilePath '%~f0' -Verb RunAs"
    exit /b
)

:: Binarios
set "INSTALL_DIR=%ProgramFiles%\TTetris"
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

copy /Y "%SRC%\bin\windows.exe" "%INSTALL_DIR%\tetris.exe"
copy /Y "%SRC%\shell\uninstall.cmd" "%INSTALL_DIR%\uninstall.cmd"

:: PATH del SISTEMA
for /f "tokens=2*" %%A in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH 2^>nul') do set "SYSPATH=%%B"
echo !SYSPATH! | findstr /I /C:"%INSTALL_DIR%" >nul
if errorlevel 1 (
    reg add "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH /t REG_EXPAND_SZ /d "!SYSPATH!;%INSTALL_DIR%" /f >nul
)

:: Resources
set "DATA_DIR=%INSTALL_DIR%\data"
if not exist "%DATA_DIR%" mkdir "%DATA_DIR%"

copy /Y "%SRC%\resources\icon.ico" "%DATA_DIR%\"
xcopy /Y /E /I "%SRC%\resources\audio\" "%DATA_DIR%\audio\"

:: Config
set "CONFIG_DIR=%AppData%\TTetris"
if not exist "%CONFIG_DIR%" mkdir "%CONFIG_DIR%"

if not exist "%CONFIG_DIR%\config.json" copy /Y "%SRC%\resources\config.json" "%CONFIG_DIR%\config.json"

:: Acceso directo
if not exist "%ProgramData%\Microsoft\Windows\Start Menu\Programs\TTetris" mkdir "%ProgramData%\Microsoft\Windows\Start Menu\Programs\TTetris"
powershell -NoProfile -Command ^
  "$s=(New-Object -COM WScript.Shell).CreateShortcut('%ProgramData%\Microsoft\Windows\Start Menu\Programs\TTetris\Tetris.lnk');" ^
  "$s.TargetPath   = '%INSTALL_DIR%\tetris.exe';" ^
  "$s.IconLocation = '%INSTALL_DIR%\data\icon.ico';" ^
  "$s.WindowStyle  =1;" ^
  "$s.Save()"