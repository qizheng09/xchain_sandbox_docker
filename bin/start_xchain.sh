#!/bin/bash

cd `dirname $0`/..

if ! [ -e data/blockchain/xuper ]; then
    bin/xchain-cli createChain
fi

bin/xchain
