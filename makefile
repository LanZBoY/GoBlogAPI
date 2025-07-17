
swag:
	@echo "Generating Swagger File..."
	cd ./app && swag init

test:
	@echo "Start to Run Test"
	cd ./app && go test ./...