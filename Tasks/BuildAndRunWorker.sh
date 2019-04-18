#!/bin/bash
cd worker
if go build; then
    echo "Build Successful"
    ./worker
else
    echo "Build Error"
fi