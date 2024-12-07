#!/bin/bash

# Define OPA server URL
OPA_URL="http://localhost:8181"

# Step 1: Load the OPA policy
echo "Uploading policy..."
curl --location --request PUT "$OPA_URL/v1/policies/mmet" \
--header 'Content-Type: text/plain' \
--data 'package mmet

# Rule to calculate the total usage
total_usage := sum([v | city = input[_]; v = city])

# Rule to calculate the target usage for each city based on the desired percentages
target_usage[city] := target {
    percentage := data.percentages[city]
    target := total_usage * percentage / 100
}

# Rule to identify cities that are above their target usage
overloaded[city] := {
    "city": city,
    "excess": usage - target
} {
    usage := input[city]
    target := target_usage[city]
    usage > target
}

# Rule to identify cities that are below their target usage
underloaded[city] := {
    "city": city,
    "deficit": target - usage
} {
    usage := input[city]
    target := target_usage[city]
    usage < target
}

# Rule to generate move commands by pairing overloaded cities with underloaded cities
move_commands[command] {
    from := overloaded[_]
    to := underloaded[_]
    from_excess := from.excess
    to_deficit := to.deficit

    # Determine the move count based on the minimum of excess and deficit
    count := min([from_excess, to_deficit])

    # Round the count to an integer if necessary
    count_int := round(count)

    # Create the move command
    command := {
        "From": from.city,
        "To": to.city,
        "Count": count_int
    }
}'
echo -e "\nPolicy uploaded."

# Step 2: Check if the policy exists
echo "Checking if policy is loaded..."
curl --location "$OPA_URL/v1/policies/mmet"
echo -e "\n"

# Step 3: Upload data for percentages
echo "Uploading data for percentages..."
curl --location --request PUT "$OPA_URL/v1/data/percentages" \
--header 'Content-Type: application/json' \
--data '{
    "Gdansk": 25,
    "Krakow": 25,
    "Poznan": 25,
    "Warsaw": 25
}'
echo -e "\nPercentages data uploaded."

# Step 4: Check if the data exists
echo "Checking if data is loaded..."
curl --location "$OPA_URL/v1/data/percentages"
echo -e "\n"

# Step 5: Perform a test query
echo "Performing test query..."
curl --location "$OPA_URL/v1/data/mmet/move_commands" \
--header 'Content-Type: application/json' \
--data '{
    "input": {
        "Gdansk": 12,
        "Krakow": 8,
        "Poznan": 10,
        "Warsaw": 13
    }
}'
echo -e "\nTest query completed."

echo "OPA setup and test completed."
