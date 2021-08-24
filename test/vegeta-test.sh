#!/bin/sh
echo "POST http://localhost:8081/hello/10" | \
  vegeta attack -connections=100000 -duration=60s -body=test_body.txt | \
  tee results.bin | \
  vegeta report -type="hist[0ms,0.2ms,0.3ms,0.4ms,0.5ms,10ms,100ms]"