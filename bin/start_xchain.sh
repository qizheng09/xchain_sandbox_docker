#!/bin/bash

cd `dirname $0`/..

if ! [ -e data/blockchain/root ]; then
    bin/xchain createChain
fi

bin/xchain
