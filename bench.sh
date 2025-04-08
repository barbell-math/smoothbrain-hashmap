#!/bin/bash

echo "Default Implementation"
go test -bench=. -benchmem ./...

echo "simd128 Implementation"
go test -tags=simd128 -bench=. -benchmem ./...

echo "simd256 Implementation"
go test -tags=simd256 -bench=. -benchmem ./...
