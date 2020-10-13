@echo off
title Windows Installation Script

rem echo Setting policy to allow localhost web server...
rem powershell -command "start-process -verb runAs 'CheckNetIsolation.exe' -argumentlist 'LoopbackExempt -a -n=\"Microsoft.Win32WebViewHost_cw5n1h2txyewy\"'

rem InstaCrypt
xcopy /y /s /i "InstaCrypt" "C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt"

set SCRIPT="%TEMP%\%RANDOM%-%RANDOM%-%RANDOM%-%RANDOM%.vbs"

echo Set oWS = WScript.CreateObject("WScript.Shell") >> %SCRIPT%
echo sLinkFile = "C:\Users\%USERNAME%\Desktop\InstaCrypt.lnk" >> %SCRIPT%
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> %SCRIPT%
echo oLink.TargetPath = "C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt\InstaCrypt.exe" >> %SCRIPT%
echo oLink.WorkingDirectory ="C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt" >> %SCRIPT%
echo oLink.Save >> %SCRIPT%

cscript /nologo %SCRIPT%

del %SCRIPT%
