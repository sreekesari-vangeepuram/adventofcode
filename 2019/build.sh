#!/bin/bash

for i in {1..25}; do
	go build -o ./bin/Day-$i/main ./Day-$i/*/main.go	
done
