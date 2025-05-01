#!/bin/bash
swag init \
    --output docs/swagger \
    --generalInfo cmd/server/main.go \
    --v3.1 \
    --parseDependency \
    --parseInternal
