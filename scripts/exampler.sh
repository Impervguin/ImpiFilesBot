#!/bin/bash

python3 ./scripts/exampler.py --dirs ./config \
    --extensions .yaml .env \
    --suffix .example \
    --exclude docker-compose \
    --override
