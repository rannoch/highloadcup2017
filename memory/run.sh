#!/usr/bin/env bash

unzip -j /tmp/data/data.zip '*.json' -d data

./go/bin/memory 80 data