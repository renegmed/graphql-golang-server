init-project:
	go mod init lyrical-app

up:
	docker-compose up -d 

test-repo:  
	go test -v ./repository/...
	
test-resolver:
	go test -v ./graph/...  

run:
	go run -race server.go 