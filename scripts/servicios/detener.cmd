@echo off
rem Doble clic para ejecutar detener.ps1 sin cambiar la politica de ejecucion de PowerShell.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0detener.ps1" %*
pause
