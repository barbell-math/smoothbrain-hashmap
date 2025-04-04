#!/bin/bash

echo "Default Implementation"
go test -bench=AgainstMapGetOnly -benchmem ./...

echo "simd256 Implementation"
go test -tags=simd256 -bench=AgainstMapGetOnly -benchmem ./...
