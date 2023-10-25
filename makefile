db:
	docker build -t db . && docker run -p 5432:5432 --name db-container -d db

add:
	curl -i -X POST http://localhost:8080/ \
	-H 'Content-Type: application/json' \
	-d '{"name": "Margarita", "surname": "Gamaleeva", "patronymic": "Sergeevna"}'

delete:
	curl -i -X DELETE http://localhost:8080/5

find:
	curl -X GET "http://localhost:8080?name=Dmitriy&offset=0&limit=2" | jq

findall:
	curl -X GET "http://localhost:8080?offset=0&limit=2" | jq

byid:
	curl -X GET "http://localhost:8080/6"