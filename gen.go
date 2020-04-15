package grafton

// The following generate directives allow generating the legacy swagger models
//go:generate rm -rf ./generated/*
//go:generate sh -c "go run github.com/go-swagger/go-swagger/cmd/swagger generate client -f specs/connector.yaml -t generated/connector"
//go:generate sh -c "go run github.com/go-swagger/go-swagger/cmd/swagger generate client -f specs/provider.yaml -t generated/provider"
