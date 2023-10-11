#!/bin/bash
bash ./build_frontend.sh

cd cmd/authlibskin
go build -trimpath -ldflags "-w -s" -tags="redis,sqlite" -o out/authlibskin
