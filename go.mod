module github.com/simpleflags/golang-server-sdk

go 1.14

require (
	github.com/hashicorp/go-retryablehttp v0.7.1
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kr/text v0.2.0 // indirect
	github.com/looplab/fsm v0.3.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/r3labs/sse/v2 v2.8.1
	github.com/simpleflags/evaluation v0.2.1
	go.uber.org/atomic v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20220706163947-c90051bbdb60 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/r3labs/sse/v2 => github.com/simpleflags/sse/v2 v2.8.1
