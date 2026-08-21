module github.com/yottta/configbuddy.v2

go 1.25.7

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // TODO: replace this with https://github.com/expr-lang/expr
	github.com/jawher/mow.cli v1.2.0 // TODO: replace it with https://github.com/spf13/cobra
	github.com/stretchr/testify v1.12.1
	github.com/yottta/go-core v0.0.13
	gopkg.in/yaml.v3 v3.0.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
