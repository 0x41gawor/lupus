#!/bin/bash

# Set OPA server URL
OPA_URL="http://192.168.56.111:8282"

# Step 1: Create the policy
echo "Creating policy for mc_server..."
curl --location --request PUT "$OPA_URL/v1/policies/mc_server" \
--header 'Content-Type: text/plain' \
--data 'package mc_server

# Rule to calculate the new license value for CPU only
cpu = cpu_new_license {
    cpu_new_license := ceil(input.in_use * 1.2)
}

# Rule to calculate the new license value for RAM only
ram = ram_new_license {
    ram_new_license := ceil(input.in_use * 1.4)
}
'
echo -e "\nPolicy created."

# Step 2: Query the RAM rule
echo "Querying RAM license calculation..."
curl --location "$OPA_URL/v1/data/mc_server/ram" \
--header 'Content-Type: application/json' \
--data '{
    "input": {
            "in_use": 8,
            "license": 8
    }
}'
echo -e "\nRAM query complete."

# Step 3: Query the CPU rule
echo "Querying CPU license calculation..."
curl --location "$OPA_URL/v1/data/mc_server/cpu" \
--header 'Content-Type: application/json' \
--data '{
    "input": {
            "in_use": 8,
            "license": 8
    }
}'
echo -e "\nCPU query complete."
