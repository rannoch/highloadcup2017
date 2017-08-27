#!/usr/bin/env bash

cp -R /tmp/data data/
unzip -j data/data.zip '*.json' -d data

./go/bin/highloadcup2017 80 data