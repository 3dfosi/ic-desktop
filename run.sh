#!/bin/bash
if [[ "$OSTYPE" == "linux-gnu" ]]; then
	./output/linux-amd64/InstaCrypt
elif [[ "$OSTYPE" == "darwin19" || "$OSTYPE" == "darwin17" ]]; then
    	open ./output/darwin-amd64/InstaCrypt.app
else
    echo "$OSTYPE is not a supported platform"
fi
