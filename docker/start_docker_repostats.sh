#!/usr/bin/env bash
docker-compose -p repostats -f pull_build.yml --env-file vars.env up -d --build --force-recreate