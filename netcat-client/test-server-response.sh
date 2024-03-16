#!/bin/sh

message="Test message to server"
response=$(echo $message | nc $SERVER_IP $SERVER_PORT)

if [ "$response" == "$message" ]; then
  echo "OK: Server response is correct: message: $message, response: $response"
else
  echo "ERROR: Server response is incorrect: message: $message, response: $response"
fi