#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from __future__ import print_function
import sys
from time import sleep
from gooey import Gooey, GooeyParser
import os


@Gooey(progress_regex=r"^progress: (\d+)/(\d+)$",
       progress_expr="x[0] / x[1] * 100",
       program_name='InstaCrypt Desktop Installer',
       image_dir='./images',
       program_description='Password-less Encryption Simplified',
       show_restart_button=False,
       show_success_modal=False)
#       auto_start=True)

def main():
    parser = GooeyParser(prog="InstaCrypt Desktop Installer")
    args = parser.parse_args(sys.argv[1:])

    print("Beginning Installation...")

    # Scoop and SSH
    prep_env = "powershell -command 'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
    os.system(prep_env)
    install_scoop = "powershell -command 'Invoke-Expression (New-Object System.Net.WebClient).DownloadString(\"https://get.scoop.sh\")'"
    os.sytem(install_scoop)
    install_ssh = "powershell -command 'scoop install git-with-openssh'"
    os.system(install_ssh)
    print("progress: {}/{}".format(1, 4))

    # InstaCrypt
    install_ic = 'xcopy /y /s /i "InstaCrypt" "C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt"'
    os.system(install_ic)
    print("progress: {}/{}".format(3, 4))

    sc1 = 'set SCRIPT="%TEMP%\%RANDOM%-%RANDOM%-%RANDOM%-%RANDOM%.vbs"'

    sc2 = 'echo Set oWS = WScript.CreateObject("WScript.Shell") >> %SCRIPT%'
    sc3 = 'echo sLinkFile = "C:\Users\%USERNAME%\Desktop\InstaCrypt.lnk" >> %SCRIPT%'
    sc4 = 'echo Set oLink = oWS.CreateShortcut(sLinkFile) >> %SCRIPT%'
    sc5 = 'echo oLink.TargetPath = "C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt\InstaCrypt.exe" >> %SCRIPT%'
    sc6 = 'echo oLink.WorkingDirectory ="C:\Users\%USERNAME%\AppData\Local\Programs\InstaCrypt" >> %SCRIPT%'
    sc7 = 'echo oLink.Save >> %SCRIPT%'
    sc8 = 'cscript /nologo %SCRIPT%'
    sc9 = 'del %SCRIPT%'
    print("progress: {}/{}".format(4,4))

    print("Installation Completed!")

if __name__ == "__main__":
    sys.exit(main())
