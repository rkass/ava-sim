module github.com/zapalabs/ava-sim

go 1.16

replace github.com/ava-labs/avalanchego => ../avalanchego

require (
	github.com/ava-labs/avalanchego v1.7.1
	github.com/fatih/color v1.9.0
	github.com/spf13/viper v1.9.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)
