# Gin and Tonic
A go api server for fun


### Compile the programs

    make all

### Create sqlite database file

    touch data/sqlite_database.db

### Run the migrations

    ./bin/migration_client up-all

###  Run the web server

     ./bin/api_server

### Create a new user

    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com", "birthday": "2025-10-12"}' \
    localhost:8080/v1.0/user

### Get a user

    curl -X GET -H "Content-Type: application/json" -d '{"id": 1}' localhost:8080/v1.0/user

### Update a user

    curl -X PUT \
    -H "Content-Type: application/json" \
    -d '{"id": 1, "first_name": "Sam", "last_name": "Rachal", "email": "sam.rachal@gmail.com", "birthday": "1990-06-15"}' \
    localhost:8080/v1.0/user

### Delete a user

    curl -X DELETE -H "Content-Type: application/json" -d '{"id": 1}' localhost:8080/v1.0/user

### Get all users

    curl -X GET -H "Content-Type: application/json" localhost:8080/v1.0/users

### Get all users with age

    curl -X GET -H "Content-Type: application/json" localhost:8080/v1.0/users_with_age

### Get age stats

    curl -X GET -H "Content-Type: application/json" localhost:8080/v1.0/age_stats