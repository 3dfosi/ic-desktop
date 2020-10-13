@echo off
title Windows Distribution Script

if not [%3]==[] (
echo Too many arguments...
goto:eof
) 

if not "%1"=="-v" (
echo Must use -v flag...
goto:eof
)

set dir=InstaCrypt-v%2-x64-Win10
echo Preparing %dir%...
mkdir %dir%
mkdir %dir%\InstaCrypt
wget -O ssh-vault_0.12.6_windows_amd64.zip https://bintray.com/nbari/ssh-vault/download_file?file_path=ssh-vault_0.12.6_windows_amd64.zip
unzip -n ssh-vault_0.12.6_windows_amd64.zip
move ssh-vault.exe %dir%\InstaCrypt
del LICENSE
del ssh-vault_0.12.6_windows_amd64.zip
copy ..\output\windows-amd64\InstaCrypt.exe %dir%\InstaCrypt 
rem copy webview.dll %dir%\InstaCrypt
rem copy WebView2Loader.dll %dir%\InstaCrypt
copy install.bat %dir%
zip -r %dir%.zip %dir%
echo All Done!



