@echo off
rem Doble clic para ejecutar levantar.ps1 sin cambiar la politica de ejecucion de PowerShell.
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0levantar.ps1" %*
pause
