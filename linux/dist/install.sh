#!/bin/bash

#################################################################
### PROJECT:
### InstaCrypt
### SCRIPT:
### install.sh
### DESCRIPTION:
### Install Script for InstaCrypt
### MAINTAINED BY:
### hkdb <hkdb@3df.io>
### Disclaimer:
### This application is maintained by volunteers and in no way
### do the maintainers make any guarantees. Use at your own risk.
### ##############################################################

echo "Beginning Installation."

echo ""

if [ ! -d ~/.local/bin ]; then
    echo "~/.local/bin doesn't exist... Creating..."
    mkdir -p ~/.local/bin
fi

if [ ! -d ~/.local/share/applications ]; then
    echo "~/.local/share/applications doesn't exist... Creating..."
    mkdir -p ~/.local/share/applications   
fi

if [ ! -d ~/.local/share/icons/hicolor/256x256 ]; then
    echo "~/.local/share/icons/hicolor/256x256 doesn't exist... Creating..."
    mkdir -p ~/.local/share/icons/hicolor/256x256
fi

echo "Installing icon..."
cp InstaCrypt-icon.png ~/.local/share/icons/hicolor/256x256/

echo "Installing .desktop..."
# Generate .desktop
cat > ~/.local/share/applications/InstaCrypt.desktop <<EOF
[Desktop Entry]
Version=0.1.0
Name=InstaCrypt
Comment=Encryption Simplified
GenericName=InstaCrypt
Exec=$HOME/.local/bin/InstaCrypt
Path=$HOME/.local/bin/
Terminal=false
Type=Application
Icon=$HOME/.local/share/icons/hicolor/256x256/InstaCrypt-icon.png
StartupNotify=true
Categories=Utility;
EOF

echo "Installing binary..."
cp InstaCrypt ~/.local/bin/
chmod +x ~/.local/bin/InstaCrypt

echo ""

echo "Entering ~/local/bin..."
cd ~/.local/bin/

if [ -f "~/.local/bin/ssh-vault" ];
then
    echo "SSH-Vault doesn't exist... Installing..."
    wget https://bintray.com/nbari/ssh-vault/download_file?file_path=ssh-vault_0.12.6_linux_amd64.tar.gz
    tar -xzvf ssh-vault_0.12.6_linux_amd64.tar.gz
    mv ssh-vault_0.12.6_linux_amd64/ssh-vault .
    chmod +x ssh-vault
    rm -rf ssh-vault_0.12.6_linux_amd64
    rm ssh-vault_0.12.6_linux_amd64.tar.gz
else
    echo "SSH-Vault is already installed... Skipping..."
fi

echo ""

echo "Installation Complete. If you don't see any errors above, you are good to go! :)"
