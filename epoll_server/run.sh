#!/usr/bin/env bash

unzip -j /tmp/data/data.zip '*.json' -d data

cp /tmp/data/options.txt data

./go/bin/highloadcup2017 80 data