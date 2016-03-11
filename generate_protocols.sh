#!/usr/bin/env bash
# for swift:    https://github.com/alexeyxo/protobuf-swift
# for go:       https://github.com/golang/protobuf
protoc --go_out=vmail/ vmail_proto/*.proto
protoc --swift_out=vmail_swift/ vmail_proto/*.proto