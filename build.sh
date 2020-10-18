#!/bin/bash
if [[ "$OSTYPE" == "linux-gnu" ]]; then
	sed -i 's/\"syscall\"/\/\/ \"syscall\"/g' main.go
	sed -i 's/cmd.SysProcAttr/\/\/ cmd.SysProcAttr/g' main.go
	astilectron-bundler
	sed -i 's/\/\/ \"syscall\"/\"syscall\"/g' main.go
	sed -i 's/\/\/ cmd.SysProcAttr/cmd.SysProcAttr/g' main.go
elif [[ "$OSTYPE" == "darwin19" || "$OSTYPE" == "darwin17" ]]; then
	sed -i -e 's/\"syscall\"/\/\/ \"syscall\"/g' main.go
        sed -i -e 's/cmd.SysProcAttr/\/\/ cmd.SysProcAttr/g' main.go
        astilectron-bundler
        sed -i -e 's/\/\/ \"syscall\"/\"syscall\"/g' main.go
        sed -i -e 's/\/\/ cmd.SysProcAttr/cmd.SysProcAttr/g' main.go
        rm main.go-e
else
    echo "$OSTYPE is not a supported platform"
fi
