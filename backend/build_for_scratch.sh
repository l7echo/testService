#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend_wsl_linux .
