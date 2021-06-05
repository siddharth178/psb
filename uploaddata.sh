#!/bin/bash

curl -v -H "Content-Type: application/x-protobuf" -H "Content-Encoding: snappy" -H "X-Prometheus-Remote-Write-Version: 0.1.0" --data-binary "@real-dataset.sz" "http://localhost:9201/write"
