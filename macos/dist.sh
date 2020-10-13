#!/bin/bash

VERSION="v0.1"

POSITIONAL=()

if [ "$#" != 2 ] && [ "$1" != "-h" ]; then
    echo -e '\nSomething is missing... Type "./dist.sh -h" without the quotes to find out more...\n'
    exit 0
fi

while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -v|--version)
    VERSION="$2"
    shift # past argument
    shift # past value
    ;;
    -h|--help)
    echo -e "\n3DF install.sh $VERSION\n\nOPTIONS:\n\n-v: Version Number\n-h, --help: Help\n\n"
    exit 0
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

DIR="InstaCrypt-v$VERSION-x64-MacOS"

cp -R ../output/darwin-amd64/InstaCrypt.app .
cp Info.plist InstaCrypt.app/
cp ssh-vault InstaCrypt.app/Contents/MacOS/
./make_dmg \
    -image dmgback.png \
    -file 144,144 InstaCrypt.app \
    -symlink 416,144 /Applications \
    -convert UDBZ \
    $DIR.dmg


