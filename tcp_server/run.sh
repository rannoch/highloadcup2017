#!/usr/bin/env bash

cp /tmp/data data
unzip -j data/data.zip '*.json' -d data

./go/bin/highloadcup2017 80 data