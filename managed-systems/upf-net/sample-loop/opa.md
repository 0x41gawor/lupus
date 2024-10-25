#### 1. Prepare OPA
##### 1.1 Run OPA
On 1st terminal
```sh
docker run -p 8181:8181 openpolicyagent/opa     run --server --log-level debug
```

From now on you can use Postman or perform curl commands in terminal:
##### 1.2 Put policy
```sh
curl --location --request PUT 'http://192.168.56.111:8181/v1/policies/mmet' \
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
```

You can check if it exists by:
```sh
curl --location 'http://192.168.56.111:8181/v1/policies/mmet'
```

##### 2.3 Put data
```sh
curl --location --request PUT 'http://192.168.56.111:8181/v1/data/percentages' \
--header 'Content-Type: application/json' \
--data '{
    	"Gdansk": 25,
		"Krakow": 25,
		"Poznan": 25,
		"Warsaw": 25
}'
```

You check if it exists by:
```sh
curl --location 'http://192.168.56.111:8181/v1/data/percentages' \
--data ''
```

##### 2.4 Perform test query
```sh
curl --location 'http://192.168.56.111:8181/v1/data/mmet/move_commands' \
--header 'Content-Type: application/json' \
--data '{
    "input":{
        "Gdansk": 12,
        "Krakow": 8,
        "Poznan": 10,
        "Warsaw": 13
    }
}'
```