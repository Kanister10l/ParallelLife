#!/bin/bash
if go build; then
    echo "Build Successful"
    ./ParallelLife -ll=debug
else
    echo "Build Error"
fi