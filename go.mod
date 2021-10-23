module github.com/rkass/ava-sim

go 1.16

replace github.com/ava-labs/avalanchego => ../avalanchego

require (
	github.com/ava-labs/avalanchego v1.6.3
	github.com/fatih/color v1.13.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

require github.com/hashicorp/go-plugin v1.4.3 // indirect
