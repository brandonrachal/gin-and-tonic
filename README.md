# Gin and Tonic
A go api server for fun

### Compile the programs

    make all

###  Run the web server

     ./bin/api_server

### Create a new user

    curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"first_name": "Brandon", "last_name": "Rachal", "email": "brandon.rachal@gmail.com"}' \
    localhost:8080/user

### Get a user

    curl -X GET -H "Content-Type: application/json" -d '{"id": 1}' localhost:8080/user

### Update a user

    curl -X PUT \
    -H "Content-Type: application/json" \
    -d '{"id": 1, "first_name": "Phillip", "last_name": "Rachal", "email": "brandon.rachal@gmail.com"}' \
    localhost:8080/user

### Delete a user

    curl -X DELETE -H "Content-Type: application/json" -d '{"id": 1}' localhost:8080/user

### Get all users

    curl -X GET -H "Content-Type: application/json" localhost:8080/users

