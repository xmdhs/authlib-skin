#!/bin/bash
cd frontend
npm install -g pnpm
pnpm install
pnpm build
cp -r dist ../server/static/files