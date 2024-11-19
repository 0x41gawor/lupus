#!/bin/bash

# Set OPA server URL
OPA_URL="http://192.168.56.111:8383"


# Step 1: Create the policy
echo "Creating policy for mc_server..."
curl --location --request PUT 'http://192.168.56.111:8383/v1/policies/sayer' \
--header 'Content-Type: text/plain' \
--data 'package sayer

# Rule to explicitly check if ram is even
auth = result {
    result := input.ram % 2 == 0
}
'
echo -e "\nPolicy created."

# Step 2: Query the policy
curl --location 'http://192.168.56.111:8383/v1/data/sayer/auth' \
--header 'Content-Type: application/json' \
--data '{
    "input":{
        "ram": 1
    }
}'
echo -e "\ Query complete."