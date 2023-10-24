db:
	docker build -t db . && docker run -p 5432:5432 --name db-container -d db