# Opa for mc-server
## Run opa
```sh
docker run -p 8181:8181 openpolicyagent/opa     run --server --log-level debug
```
## Create policy
```sh
curl --location --request PUT 'http://192.168.56.111:8181/v1/policies/mc_server' \
--header 'Content-Type: text/plain' \
--data 'package mc_server

# Rule to calculate the new license value for CPU only
cpu = {"cpu": cpu_new_license} {
    cpu_new_license := ceil(input.cpu.in_use * 1.2)
}

# Rule to calculate the new license value for RAM onlylup
ram = {"ram": ram_new_license} {
    ram_new_license := ceil(input.ram.in_use * 1.2)
}
'
```

## Query
RAM:
```sh
curl --location 'http://192.168.56.111:8181/v1/data/mc_server/ram' \
--header 'Content-Type: application/json' \
--data '{
    "input": {
        "ram" : {
                "in_use": 40
        }
    }
}'
```
CPU:
```sh
curl --location 'http://192.168.56.111:8181/v1/data/mc_server/cpu' \
--header 'Content-Type: application/json' \
--data '{
    "input": {
        "cpu" : {
                "in_use": 20 
        }
    }
}'
```
