#!/bin/bash
sudo rm -rf data/*
sudo rm -rf provingKey.txt
sudo rm -rf verifyingKey.json
mkdir -p data
echo -n 1 > data/podNumber.txt
go run main.go
