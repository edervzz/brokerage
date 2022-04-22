# Brokerage by Oscar Eder Velázquez Pineda
Rest service to process buy & sell orders. This project requires Docker Compose to host API and Database (MySQL).

## Post methods
- `/migration` : Create database and tables.
- `/accounts` : Create new account.
- `/accounts/{id}/orders` : Create order from account.


## Running locally
1. Save and uncompress project in a locar folder

2. Build API
```
docker compose build
```

3. Execte
```
docker compose up -d
```

## API Examples
- __Create database__
    ```
    curl --location --request POST 'http://localhost:8000/migration'
    ```
- Response
    ```
    200 OK
    ```
- __Create account__
    ```
    curl --location --request POST 'http://localhost:8000/accounts' \
    --header 'Content-Type: application/json' \
    --data-raw '{
        "cash": 5000
    }'
    ```
- Response = 200 OK
    ```
    {
        "id":1,
        "cash":5000,
        "issuers":[]
    }
    ```
- __Create an Order__
    ```
    curl --location --request POST 'http://localhost:8000/accounts/1/orders' \
    --header 'Content-Type: application/json' \
    --data-raw '[
        {
            "timestamp": "1650525504",
            "operation": "BUY",
            "issuer_name": "APPL",
            "total_shares": 10,
            "total_price": 10
        }
    ]'
    ```
- Response = 200 OK
    ```
    {
        "cash":4900,
        "issuers": [
            { 
                "issuer_name":"APPL",
                "total_shares":10,
                "share_price":10,
                "business_erros":[]
            }
        ]
    }
    ```
- __Create two Orders__
    ```
    curl --location --request POST 'http://localhost:8000/accounts/1/orders' \
    --header 'Content-Type: application/json' \
    --data-raw '[
        {
            "timestamp": "1650525904",
            "operation": "BUY",
            "issuer_name": "APPL",
            "total_shares": 10,
            "total_price": 10
        },
        {
            "timestamp": "1650525904",
            "operation": "BUY",
            "issuer_name": "SBUX",
            "total_shares": 11,
            "total_price": 20
        }
    ]'
    ```
- Response = 200 OK
    ```
    {
        "cash":4360,
        "issuers": [ 
            { 
                "issuer_name":"APPL",
                "total_shares":10,
                "share_price":10,
                "business_erros":[]
            },
            {
                "issuer_name":"SBUX",
                "total_shares":11,
                "share_price":20,
                "business_erros":[]
            }
        ]
    }
    ```
- __Check coverage__
    ```
    go test ./... -cover -coverprofile cover.out
    go tool cover -html cover.out
    ```