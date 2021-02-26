#!/bin/sh -e

docker build -t jk/coherence:1.0.0 -f operator-compatibility/target/classes/Dockerfile ./operator-compatibility