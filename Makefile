run-compose:
	$(info starting services. http interface available on 8080 port)
	@docker-compose up --build -d

stop:
	@docker-compose down