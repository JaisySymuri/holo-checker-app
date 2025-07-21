@echo off
echo Building holo-checker-app with icon...

:: Generate rsrc.syso inside ./cmd
rsrc -ico favicon.ico -o ./cmd/rsrc.syso
if %errorlevel% neq 0 (
    echo Failed to generate rsrc.syso
    pause
    exit /b %errorlevel%
)

:: Build the executable with windowsgui mode
go build -ldflags="-H=windowsgui" -o holo-checker-app.exe ./cmd
if %errorlevel% neq 0 (
    echo Build failed!
    pause
    exit /b %errorlevel%
)

echo Build complete: holo-checker-app.exe
pause
