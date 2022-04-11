#!/usr/bin/env bash
docker-compose -p repostats -f dependances.yml --env-file vars.env up -d --build --force-recreate