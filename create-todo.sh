#!/bin/bash
set -eu
curl -X POST localhost:8080/todos -d '{"subject": "hoge", "description": "piyo"}'