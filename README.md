# Personality Service

Implementation of the service that receives full name via API and enriches it from open APIs
answer with the most likely age, gender and nationality and storing the data in
DB. 

## Usage
**1.** [Install golang](https://go.dev/doc/install)  
**2.** Download repository from github
```shell
go get "https://github.com/superhacker2002/personality-service"
```
**2.** Set up environment variables:
- Create .env file
- Set `PORT` variable to port on which the server will listen for incoming connections.
- Set `DATABASE_URL` variable to URL of the database where the service will store the data

**3.** Run web service from the command line:
```shell
go run cmd/maing.go
```

## API Endpoints

1. Obtaining data with name filter and pagination
```
curl -X GET "http://localhost:8080?name=Dmitriy&offset=0&limit=2"
```
2. Deleting by ID
```
   curl -i -X DELETE http://localhost:8080/5
```
3. Changing person information
```
curl -i -X PUT http://localhost:8080/10 \
    	-H 'Content-Type: application/json' \
    	-d '{"name": "Dmitriy", "surname": "NewSurname", "patronymic": "NewPatronymic"}'
```
4. Adding new person (information about them will be enriched by the server)
```
curl -i -X POST http://localhost:8080/ \
	-H 'Content-Type: application/json' \
	-d '{"name": "Margarita", "surname": "Gamaleeva", "patronymic": "Sergeevna"}'
```
