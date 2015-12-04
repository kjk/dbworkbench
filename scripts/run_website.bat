cd website
go vet ./...

go test .

go run main.go

go build
./website
rm website

