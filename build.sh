#!/bin/bash
if [[ "$OSTYPE" == "linux-gnu" ]]; then
	sed -i 's/\"syscall\"/\/\/ \"syscall\"/g' main.go
	sed -i 's/cmd.SysProcAttr/\/\/ cmd.SysProcAttr/g' main.go
	astilectron-bundler
	sed -i 's/\/\/ \"syscall\"/\"syscall\"/g' main.go
	sed -i 's/\/\/ cmd.SysProcAttr/cmd.SysProcAttr/g' main.go
elif [[ "$OSTYPE" == "darwin19" || "$OSTYPE" == "darwin17" ]]; then
	sed -i 's/\"syscall\"/\/\/ \"syscall\"/g' main.go
	sed -i 's/cmd.SysProcAttr/\/\/ cmdSysProcAttr/g' main.go
    	astilectron-bundler
	sed -i 's/\/\/ \"syscall\"/\"syscall\"/g' main.go
else
    echo "$OSTYPE is not a supported platform"
fi
