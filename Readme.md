https://grpc.io/docs/languages/go/quickstart/

brew install protobuf

mkdir -p gen/go

cd protos
protoc -I proto proto/sso/sso.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

go get github.com/pavelfire/protostu

go get github.com/ilyakaznacheev/cleanenv

====Run======
cd ssotu
go run ./cmd/ssotu/main.go --config=./config/local.yaml

go get google.golang.org/grpc

shift + alt + F --- formating
cmd + . --- import func


cd /Users/pvdo/VScodeProjects/authtuz/ssotu
# подтянуть зависимости и обновить go.sum
go mod tidy
# (необязательно, но можно явно)
go get github.com/pavelfire/protostu@v0.1.2
# проверить сборку
go build ./cmd/ssotu


go get golang.org/x/crypto
go get github.com/golang-jwt/jwt/v5

go test ./internal/lib/jwt/... -v из каталога ssotu.

go get github.com/golang-migrate/migrate/v4