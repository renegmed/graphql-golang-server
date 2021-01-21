init-project:
	go mod init lyrical-app

up:
	docker-compose up -d 

cert:
	openssl genrsa -out ./tls/server.key 2048
	openssl req -new -x509 -key ./tls/server.key -out ./tls/server.pem -days 365

test-repo:  
	go test -v ./repository/...
	
test-resolver:
	go test -v ./graph/...  

run:
	go run -race server.go 