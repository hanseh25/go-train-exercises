# go-train-exercises

via command line :

-- clean db 
- migrate -database ${POSTGRESQL_URL_MAIN} -path migrations down

-- initialize db
- migrate -database ${POSTGRESQL_URL_MAIN} -path migrations up

-- if you need to disable sslmode
- export POSTGRESQL_URL='postgres://testuser:testpassword@localhost:5432/passlocker?sslmode=disable'

-- full command to run 
- PSWLCKRDSN=postgres://testuser:testpassword@localhost/passlocker AUTH_USERNAME=hans AUTH_PASSWORD=password go run .

-- curl test with basic auth
- curl -k -u hans:password -d '{"url":"www.painhub.com", "username":"hans", "password":"password"}'  -X POST https://localhost:4000/save-credentials

-- list items
- curl -k -u hans:password https://localhost:4000/list-credentials

-- package
- go get github.com/hanseh25/go-password-gen   

