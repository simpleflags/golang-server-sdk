# Simple Flags Server-side SDK for Go

---

this is experimental repo!!!

## Supported GO versions
This version of SDK has been tested with Go 1.14

## Install
`go get github.com/simpleflags/golang-server-sdk`

## Usage
First we need to import lib
```go
import sdk "github.com/simpleflags/golang-server-sdk"
```

Next we initialize client instance for interaction with api
```go
err := sdk.Initialize(sdkKey)
```

Target definition can be user, device, app etc.
```go
target := map[string]interface{}{
    "identifier": "enver",
}
```

Evaluating Feature Flag
```go
showFeature, err := sdk.Evaluate(featureFlagKey, &target).Bool(false)
```

Flush any changes and close the SDK
```go
sdk.close()
```

## Interface

very simple and small interfaces:
```go
type Client interface {
    WaitForInitialization()
    Evaluate(feature string, target evaluation.Target) evaluation.Evaluation
    Close() error
}

type Logger interface {
    Debug(args ...interface{})
    Debugf(template string, args ...interface{})
    Info(args ...interface{})
    Infof(template string, args ...interface{})
    Error(args ...interface{})
    Errorf(template string, args ...interface{})
}
```

## Logger

It is very simple to set logger from your current app configuration:
```go
sdk.SetLogger(your_logger)
```