package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/fatih/color"
	"github.com/zapalabs/ava-sim/constants"
	"github.com/zapalabs/ava-sim/manager"
	"github.com/zapalabs/ava-sim/runner"
	"golang.org/x/sync/errgroup"
)

func main() {
	os.Setenv("AVASIM", "true")
	var vm, vmGenesis string
	switch len(os.Args) {
	case 1: // normal network
	case 2: // new node for existing network
		vm = path.Clean(os.Args[1])
		if _, err := os.Stat(vm); os.IsNotExist(err) {
			panic(fmt.Sprintf("%s does not exist", vm))
		}
		color.Yellow("vm set to: %s", vm)
	case 3:
		vm = path.Clean(os.Args[1])
		if _, err := os.Stat(vm); os.IsNotExist(err) {
			panic(fmt.Sprintf("%s does not exist", vm))
		}
		color.Yellow("vm set to: %s", vm)

		vmGenesis = path.Clean(os.Args[2])
		if _, err := os.Stat(vmGenesis); os.IsNotExist(err) {
			panic(fmt.Sprintf("%s does not exist", vmGenesis))
		}
		color.Yellow("vm-genesis set to: %s", vmGenesis)
	case 4:
		// staking.InitNodeStakingKeyPair("/Users/rkass/repos/zapa/ava-sim/manager/certs/keys6/staker.key", "/Users/rkass/repos/zapa/ava-sim/manager/certs/keys6/staker.crt")
		p, e := tls.LoadX509KeyPair("/Users/rkass/repos/zapa/ava-sim/manager/certs/keys6/staker.crt", "/Users/rkass/repos/zapa/ava-sim/manager/certs/keys6/staker.key")
		color.Yellow("initialized staking pair for node 6 %s %s", p, e)
		os.Exit(0)
	default:
		panic("invalid arguments (expecting no arguments or [vm] [vm-genesis])")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// register signals to kill the application
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT)
		signal.Notify(signals, syscall.SIGTERM)
		defer func() {
			// shut down the signal go routine
			signal.Stop(signals)
			close(signals)
		}()

		select {
		case sig := <-signals:
			color.Red("signal received: %v", sig)
			cancel()
		case <-gctx.Done():
		}
		return nil
	})

	if vmGenesis != "" {
		// Start local network
		bootstrapped := make(chan struct{})

		nodeNums := make([]int, constants.NumNodes)
		for i := 0; i < constants.NumNodes; i++ {
			nodeNums[i] = i
		}

		g.Go(func() error {
			
			return manager.StartNetwork(gctx, vm, nodeNums, bootstrapped)
		})

		// Only setup network if a custom VM is provided and the network has finished
		// bootstrapping
		select {
		case <-bootstrapped:
			if len(vm) > 0 && gctx.Err() == nil {
				g.Go(func() error {
					return runner.SetupSubnet(gctx, nodeNums, vmGenesis, true)
				})
			}
		case <-gctx.Done():
		}
	} else {
		bootstrapped := make(chan struct{})
		nodeNums := []int{5}
		color.Yellow("setting up 1 node")
		g.Go(func() error {
			
			return manager.StartNetwork(gctx, vm, nodeNums, bootstrapped)
		})

		// Only setup network if a custom VM is provided and the network has finished
		// bootstrapping
		select {
		case <-bootstrapped:
			if len(vm) > 0 && gctx.Err() == nil {
				g.Go(func() error {
					return runner.SetupSubnet(gctx, nodeNums, vmGenesis, false)
				})
			}
		case <-gctx.Done():
		}
	}

	color.Red("ava-sim exited with error: %s", g.Wait())
	os.Exit(1)
}
