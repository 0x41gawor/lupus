#!/bin/bash

# Set OPA server URL
OPA_URL="http://127.0.0.1:8383"

# Step 1: Create the policy
echo "Creating policy for mc_server..."
curl --location --request PUT "$OPA_URL/v1/policies/sayer" \
--header 'Content-Type: text/plain' \
--data 'package sayer

# Rule to explicitly check if ram is even
auth = result {
    result := input.ram % 2 == 0
}
'
echo -e "\nPolicy created."

# Step 2: Query the policy
echo "Querying policy..."
curl --location "$OPA_URL/v1/data/sayer/auth" \
--header 'Content-Type: application/json' \
--data '{
    "input": {
        "ram": 1
    }
}'
echo -e "\nQuery complete."
