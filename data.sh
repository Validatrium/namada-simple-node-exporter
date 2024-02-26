#!/bin/bash

> status.txt

watch -n 10 'wget -q -O - http://you_rpc:you_port/status > status.txt' &

