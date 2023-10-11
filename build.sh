#!/bin/bash
cd frontend
npm install -g pnpm
pnpm install
pnpm build
cp -r dist ../server/static/files

cd ../cmd/authlibskin
go build -trimpath -ldflags "-w -s" -tags="redis,sqlite" -o out/authlibskin
