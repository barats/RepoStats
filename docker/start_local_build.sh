#!/usr/bin/env bash
docker-compose -p repostats -f local_build.yml --env-file vars.env up -d --build --force-recreate