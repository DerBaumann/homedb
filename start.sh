#!/bin/bash

# Run this if the normal air command doesn't work!
# source ./.env.sh && air

source .env.sh

if [ $DEV_PLATFORM = "linux" ]; then
  docker start homedb-dev
fi

air
