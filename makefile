
swag:
	@echo "Generating Swagger File..."
	cd ./app && swag init

test:
	@echo "Start to Run Test"
	go test ./... -coverprofile coverReport.out

generate_test_report:
	@echo "Generating Report..."
	go tool cover -html=coverReport.out -o coverReport.html

test_and_genreate_report: test generate_test_report