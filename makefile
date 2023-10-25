db:
	docker build -t db . && docker run -p 5432:5432 --name db-container -d db

add:
	curl -i -X POST http://localhost:8080/ \
	-H 'Content-Type: application/json' \
	-d '{"name": "Dmitriy", "surname": "Ushakov", "patronymic": "Vasilevich"}'

delete:
	curl -i -X DELETE http://localhost:8080/1

