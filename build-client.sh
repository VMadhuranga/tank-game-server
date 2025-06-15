#!/bin/bash

echo "building client"
rm -rf client && cd ../tank-game-client/ && npm run build && cp -r dist ../tank-game-server/client
echo "done"
