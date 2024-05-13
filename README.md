# GoBank

A simple Bank microservice application built on top of Golang used as a simple demo on how to utilize new feature in Mux standard package Go 1.22

## Running the App/Service

1. `docker-compose up`

## Testing the Endpoints

1. Create / Register a new Account
```
curl -X POST -H 'Content-Type: application/json' -d '{"first_name": "M Yauri M", "last_name": "Attamimi", "account_number": "7830371235", "balance": 500000}' localhost:8080/api/accounts
```
2. Fetch all Accounts
```
curl -X GET -H 'Accept: application/json' localhost:8080/api/accounts
```