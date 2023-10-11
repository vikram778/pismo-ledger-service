# pismo-ledger-service

Transaction service is a service which allows to create account and once the account is created can perform
transactions (purchase , withdrwal , credit voucher)

#Trnasaction service flow

1. Account needs to be created via createAccount api.
2. Once the account is created transactions can be performed.
3. Transactions can be of only certain types given in the operationTypes table.
4. Upon receiving the transaction request service checks if the account exists and then also checks for the operationType.
5. If account exists and operationType in the request is valid service proceeds with creating a transaction record in db.

#Database Migrations

 To execute the Database Migration run the following command 

    - go run cmd/migrate/main.go up


#Application startup and set up

1. To run and build auth service :
   go run cmd/app/main.go or execute run.sh script (sh run.sh)
   go build ./cmd/app

   
#Prerequisites

1. PostgresSQl


Above dependencies are added in docker-compose file .

#Starting dependencies and service in local

docker-compose -f docker-compose.yml up --build

#Configs

config .yml files are in /config directory