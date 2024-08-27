module server

go 1.18

require (
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.34.2
	gopkg.in/yaml.v3 v3.0.1
	iproto v0.0.0-00010101000000-000000000000
)

require (
	github.com/stretchr/testify v1.9.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace iproto => ./../iproto
