@echo off
REM CyberLab Server — Startup Script (Windows)
REM
REM Usage: run.bat
REM Reads from .env if it exists, otherwise prompts.

setlocal enabledelayedexpansion

cd /d "%~dp0"

REM Load .env if it exists
if exist .env (
  echo ======================================== Loading configuration from .env...
  for /f "usebackq tokens=1,2 delims==" %%a in (.env) do (
    if not "%%b"=="" set "%%a=%%b"
  )
)

REM Check required environment variables
set MISSING=false
if "%JWT_SECRET%"=="" set MISSING=true
if "%DB_PASSWORD%"=="" set MISSING=true
if "%FLAG_HASH_SALT%"=="" set MISSING=true

if "%MISSING%"=="true" (
  echo ======================================== WARNING: Missing required environment variables!
  echo.
  echo Please create a .env file next to this script with:
  echo.
  echo   JWT_SECRET=your-random-secret-key
  echo   DB_PASSWORD=your-database-password
  echo   FLAG_HASH_SALT=your-random-64-char-salt
  echo.
  echo Optionally:
  echo   AI_API_KEY=sk-your-openai-api-key
  echo.

  if not exist .env (
    echo ======================================== Generating .env with random values...
    powershell -Command "$r1 = -join ((48..57)+(65..90)+(97..122) | Get-Random -Count 64 | ForEach-Object{[char]$_}); $r2 = -join ((48..57)+(65..90)+(97..122) | Get-Random -Count 64 | ForEach-Object{[char]$_}); @('JWT_SECRET='+$r1, 'DB_PASSWORD=change-me-db-password', 'FLAG_HASH_SALT='+$r2, '# AI_API_KEY=sk-your-openai-api-key') | Out-File -FilePath .env -Encoding utf8"
    echo ======================================== .env file created. Edit it to set your database password.
    echo   notepad .env
  )
  pause
  exit /b 1
)

REM Determine binary
set BINARY=cyberlab-server.exe
if not exist "%BINARY%" (
  echo ERROR: Binary not found: %BINARY%
  echo Did you extract the full release package?
  pause
  exit /b 1
)

echo.
echo ========================================
echo   CyberLab Server v1.0.0
echo ========================================
echo   API:       http://localhost:17420
echo   Swagger:   http://localhost:17420/swagger/index.html
echo   Frontend:  http://localhost:17420
echo ========================================
echo.

"%BINARY%"
