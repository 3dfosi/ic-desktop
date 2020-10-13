#!/bin/bash

POSITIONAL=()

if [ "$#" != 2 ] && [ "$1" != "-h" ]; then
    echo -e '\nSomething is missing... Type "./setup -h" without the quotes to find out more...\n'
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
    echo -e "\ndist.sh $VERSION\n\nOPTIONS:\n\n-v: version\n-r\n"
    exit 0
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

DIR="InstaCrypt-v$VERSION-x64-Linux"
mkdir $DIR
cp ../output/linux-amd64/InstaCrypt $DIR/ 
cp ../InstaCrypt-icon.png $DIR
cp dist/install.sh $DIR/
tar -cjvf InstaCrypt-v$VERSION-x64-Linux.tar.bz2 $DIR 
