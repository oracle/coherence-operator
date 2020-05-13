#!/usr/bin/env bash


kill -9 $(ps aux | grep 'operator-sdk up' | awk '{print $2}')
kill -9 $(ps aux | grep 'tee operator.out' | awk '{print $2}')
kill -9 $(ps aux | grep 'operator-local' | awk '{print $2}')
kill -9 $(ps aux | grep 'coherence-operator' | awk '{print $2}')
kill -9 $(ps aux | grep 'local-watches.yaml' | awk '{print $2}')
