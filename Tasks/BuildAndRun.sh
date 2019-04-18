#!/bin/bash
if go build; then
    echo "Build Successful"
    ./ParallelLife
else
    echo "Build Error"
fi