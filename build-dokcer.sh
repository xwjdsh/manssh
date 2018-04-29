#! /bin/bash
set -e 

docker build --build-arg VERSION=`git describe --tags` -t wendellsun/manssh .
