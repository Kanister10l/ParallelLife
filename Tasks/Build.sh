#!/bin/bash
if go build; then
    echo "Build Successful"
else
    echo "Build Error"
fi