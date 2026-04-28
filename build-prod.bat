@echo off
echo Building holo-checker-app with icon...

:: Generate Windows resource file in the root package
rsrc -ico favicon.ico -o resource.syso
if %errorlevel% neq 0 (
    echo Failed to generate resource.syso
    pause
    exit /b %errorlevel%
)

:: Build the executable with windowsgui mode
go build -ldflags="-H=windowsgui" -o holo-checker-app.exe .
if %errorlevel% neq 0 (
    echo Build failed!
    pause
    exit /b %errorlevel%
)

echo Build complete: holo-checker-app.exe
pause
