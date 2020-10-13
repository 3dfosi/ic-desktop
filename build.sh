#!/bin/bash
if [[ "$OSTYPE" == "linux-gnu" ]]; then
	astilectron-bundler
elif [[ "$OSTYPE" == "darwin19" || "$OSTYPE" == "darwin17" ]]; then
    	astilectron-bundler
else
    echo "$OSTYPE is not a supported platform"
fi
