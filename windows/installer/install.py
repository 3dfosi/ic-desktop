#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import unicode_literals
from __future__ import print_function
import sys
from time import sleep
from gooey import Gooey, GooeyParser
import os
import tempfile
from pathlib import Path
import shutil

# nonbuffered_stdout = os.fdopen(sys.stdout.fileno(), 'w', 0)
# sys.stdout = nonbuffered_stdout

@Gooey(progress_regex=r"^progress: (\d+)/(\d+)$",
       progress_expr="x[0] / x[1] * 100",
       program_name='InstaCrypt Desktop Installer',
       image_dir='images',
       program_description='Password-less Encryption Simplified',
       show_restart_button=False,
       show_success_modal=False)
#       auto_start=True)

def main():
    parser = GooeyParser(prog="InstaCrypt Desktop Installer")
    args = parser.parse_args(sys.argv[1:])

    print("Beginning Installation...\n")

    # Get HomeDir
    homedir = str(Path.home())

    # Scoop and SSH
    if os.path.isfile(homedir + "\\scoop\\shims\\ssh.exe") == False:
        print("SSH is not on system. Installing scoop and then will install ssh")
        prep_env = "powershell -command 'Set-ExecutionPolicy RemoteSigned -scope CurrentUser'"
        os.system(prep_env)
        install_scoop = "powershell -command 'Invoke-Expression (New-Object System.Net.WebClient).DownloadString(\"https://get.scoop.sh\")'"
        os.system(install_scoop)
        install_ssh = "powershell -command 'scoop install openssh'"
        os.system(install_ssh)
    else:
        print("Skipping installation of scoop and ssh since it already exists...")

    print("progress: {}/{}".format(1, 4))
    print("\n")

    # InstaCrypt
    icdir = homedir + "\\AppData\\Local\\Programs\\InstaCrypt"
    if os.path.isdir(icdir) == True:
        print("Removing previous versions of InstaCrypt Desktop...")
        shutil.rmtree(icdir)
    
    shutil.copytree("InstaCrypt", homedir + "\\AppData\\Local\\Programs\\InstaCrypt")
    print("Installed InstaCrypt Desktop...\n")
    print("progress: {}/{}".format(3, 4))
    print ("\n")

    # Create InstaCrypt Desktop Shortcut
    dir = tempfile.gettempdir()
    Path(dir+"\\shortcut.vbs").touch()
    vbs = open(dir+"\\shortcut.vbs", 'a')
    vbs.write("Set oWS = WScript.CreateObject(\"WScript.Shell\")\n")
    vbs.write("sLinkFile = \"" + homedir + "\\Desktop\\InstaCrypt.lnk\"\n")
    vbs.write("Set oLink = oWS.CreateShortcut(sLinkFile)\n")
    vbs.write("oLink.TargetPath = \"" + homedir + "\\AppData\\Local\\Programs\\InstaCrypt\\InstaCrypt.exe\"\n")
    vbs.write("oLink.WorkingDirectory =\"" + homedir + "\\AppData\\Local\\Programs\\InstaCrypt\"\n")
    vbs.write("oLink.Save\n")
    vbs.close()
    os.system("cscript /nologo " + dir + "\\shortcut.vbs")
    os.remove(dir+"\\shortcut.vbs")
    print("Created desktop shortcut...")
    print("progress: {}/{}".format(4,4))
    print("\n")

    print("Installation Completed!")

if __name__ == "__main__":
    sys.exit(main())
