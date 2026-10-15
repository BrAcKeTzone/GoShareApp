@echo off
SETLOCAL

SET BINARY_NAME=goshare.exe
SET INSTALL_DIR=C:\Program Files\GoShare

echo 🔍 Checking for binary...
IF NOT EXIST "%BINARY_NAME%" (
    echo [ERROR] %BINARY_NAME% not found in the current directory. Please build it first.
    pause
    exit /b 1
)

echo 📦 Creating installation directory at %INSTALL_DIR%...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

echo 📋 Copying binary...
copy /Y "%BINARY_NAME%" "%INSTALL_DIR%\%BINARY_NAME%" >nul

echo ⚙️ Adding to system PATH...
setx /M PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1
if errorlevel 1 (
    echo [WARNING] Could not write to System PATH. Adding to User PATH instead...
    powershell -Command "[Environment]::SetEnvironmentVariable('PATH', [Environment]::GetEnvironmentVariable('PATH', 'User') + ';%INSTALL_DIR%', 'User')"
)

echo 🛠️ Generating uninstaller...
(
    echo @echo off
    echo echo [INFO] Uninstalling GoShare...
    echo rmdir /S /Q "%INSTALL_DIR%"
    echo echo [SUCCESS] GoShare has been successfully uninstalled.
    echo pause
    echo (goto) 2^>nul ^& del "%%~f0"
) > "%INSTALL_DIR%\uninstall.bat"

echo.
echo 🎉 Installation complete!
echo 👉 Open a brand-new terminal window in any folder and type: goshare
echo 👉 To uninstall, navigate to "%INSTALL_DIR%" and run "uninstall.bat" as Administrator.
pause