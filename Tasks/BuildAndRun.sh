#!/bin/bash
if go build; then
    echo "Build Successful"
    ./ParallelLife -ll=clean
else
    echo "Build Error"
fi