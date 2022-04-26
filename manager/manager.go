package manager

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/zapalabs/ava-sim/constants"
	"github.com/zapalabs/ava-sim/utils"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/app/process"
	"github.com/ava-labs/avalanchego/node"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

const (
	bootstrapID = "NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg"
	bootstrapIP = "127.0.0.1:9653"
	waitDiff    = 10 * time.Second
)

// Embed certs in binary and write to tmp file on startup (full binary)
var (
	//go:embed certs/keys1/staker.crt
	keys1StakerCrt []byte
	//go:embed certs/keys1/staker.key
	keys1StakerKey []byte

	//go:embed certs/keys2/staker.crt
	keys2StakerCrt []byte
	//go:embed certs/keys2/staker.key
	keys2StakerKey []byte

	//go:embed certs/keys3/staker.crt
	keys3StakerCrt []byte
	//go:embed certs/keys3/staker.key
	keys3StakerKey []byte

	//go:embed certs/keys4/staker.crt
	keys4StakerCrt []byte
	//go:embed certs/keys4/staker.key
	keys4StakerKey []byte

	//go:embed certs/keys5/staker.crt
	keys5StakerCrt []byte
	//go:embed certs/keys5/staker.key
	keys5StakerKey []byte

	//go:embed certs/keys6/staker.crt
	keys6StakerCrt []byte
	//go:embed certs/keys6/staker.key
	keys6StakerKey []byte


	nodeCerts = [][]byte{keys1StakerCrt, keys2StakerCrt, keys3StakerCrt, keys4StakerCrt, keys5StakerCrt, keys6StakerCrt}
	nodeKeys  = [][]byte{keys1StakerKey, keys2StakerKey, keys3StakerKey, keys4StakerKey, keys5StakerKey, keys6StakerKey}
)

func NodeIDs(nodeNums[] int) []string {
	nodeCerts := [][]byte{keys1StakerCrt, keys2StakerCrt, keys3StakerCrt, keys4StakerCrt, keys5StakerCrt, keys6StakerCrt}
	nodeIDs := make([]string, len(nodeNums))
	for i, nodeNum := range nodeNums {
		cert := nodeCerts[nodeNum]
		id, err := utils.LoadNodeID(cert)
		if err != nil {
			panic(err)
		}
		nodeIDs[i] = id
	}
	return nodeIDs
}

func NodeURLs(nodeNums[] int) []string {
	urls := make([]string, len(nodeNums))
	for i, nodeNum := range nodeNums {
		urls[i] = fmt.Sprintf("http://127.0.0.1:%d", constants.BaseHTTPPort+nodeNum*2)
	}
	return urls
}

func StartNetwork(ctx context.Context, vmPath string, nodeNums []int, bootstrapped chan struct{}) error {
	dir, err := ioutil.TempDir("", "ava-sim")
	if err != nil {
		panic(err)
	}
	color.Cyan("tmp dir located at: %s", dir)
	defer func() {
		color.Cyan("tmp dir located at: %s", dir)
	}()

	// Copy files into custom plugins
	pluginsDir := fmt.Sprintf("%s/plugins", dir)
	if err := os.MkdirAll(pluginsDir, os.FileMode(constants.FilePerms)); err != nil {
		panic(err)
	}
	if err := utils.CopyFile("build/system-plugins/evm", fmt.Sprintf("%s/evm", pluginsDir)); err != nil {
		panic(err)
	}
	if len(vmPath) > 0 {
		if err := utils.CopyFile(vmPath, fmt.Sprintf("%s/%s", pluginsDir, constants.VMID)); err != nil {
			panic(err)
		}
	}



	nodeConfigs := make([]node.Config, len(nodeNums))
	for i, nodeNum := range nodeNums {
		nodeConfigs[i] = getNodeConfig(nodeNum, dir, vmPath, pluginsDir)
	}

	// Start all nodes and check if bootstrapped
	g, gctx := errgroup.WithContext(ctx)
	for i, config := range nodeConfigs {
		c := config
		j := i

		g.Go(func() error {
			return runApp(g, gctx, j, c)
		})
	}
	g.Go(func() error {
		return checkBootstrapped(gctx, nodeNums, bootstrapped)
	})
	return g.Wait()
}

// nodeNum is 0-indexed
func getNodeConfig(nodeNum int, dir string, vmPath string, pluginsDir string) node.Config {
	nodeDir := fmt.Sprintf("%s/node%d", dir, nodeNum+1)
	if err := os.MkdirAll(nodeDir, os.FileMode(constants.FilePerms)); err != nil {
		panic(err)
	}
	certFile := fmt.Sprintf("%s/staker.crt", nodeDir)
	if err := ioutil.WriteFile(certFile, nodeCerts[nodeNum], os.FileMode(constants.FilePerms)); err != nil {
		panic(err)
	}
	keyFile := fmt.Sprintf("%s/staker.key", nodeDir)
	if err := ioutil.WriteFile(keyFile, nodeKeys[nodeNum], os.FileMode(constants.FilePerms)); err != nil {
		panic(err)
	}

	df := defaultFlags()
	df.LogLevel = "info"
	df.LogDir = fmt.Sprintf("%s/logs", nodeDir)
	df.DBDir = fmt.Sprintf("%s/db", nodeDir)
	df.StakingEnabled = true
	df.HTTPPort = uint(constants.BaseHTTPPort + 2*nodeNum)
	df.StakingPort = uint(constants.BaseHTTPPort + 2*nodeNum + 1)
	if nodeNum != 0 {
		df.BootstrapIPs = bootstrapIP
		df.BootstrapIDs = bootstrapID
	} else {
		df.BootstrapIPs = ""
		df.BootstrapIDs = ""
	}
	if len(vmPath) > 0 {
		df.WhitelistedSubnets = constants.WhitelistedSubnets
	}
	df.StakingTLSCertFile = certFile
	df.StakingTLSKeyFile = keyFile
	nodeConfig, err := createNodeConfig(pluginsDir, flagsToArgs(df))
	if err != nil {
		panic(err)
	}
	nodeConfig.PluginDir = pluginsDir

	// write node id of node
		h, _ := os.LookupEnv("HOME")
		f, _ := os.Create(h + "/node-ids/" + strconv.Itoa(nodeNum))
	nodeNums := []int{nodeNum}
	nodeIds := NodeIDs(nodeNums)
	f.Write([]byte(nodeIds[0]))
	f.Close()
	return nodeConfig
}

func checkBootstrapped(ctx context.Context, nodeNums []int, bootstrapped chan struct{}) error {
	if bootstrapped == nil {
		return nil
	}

	var (
		nodeURLs = NodeURLs(nodeNums)
		nodeIDs  = NodeIDs(nodeNums)
	)

	for i, url := range nodeURLs {
		client := info.NewClient(url)

		for {
			if ctx.Err() != nil {
				color.Red("stopping bootstrapped check: %v", ctx.Err())
				return ctx.Err()
			}
			bootstrapped := true
			for _, chain := range constants.Chains {
				chainBootstrapped, _ := client.IsBootstrapped(ctx, chain)
				if !chainBootstrapped {
					color.Yellow("waiting for %s to bootstrap %s-chain", nodeIDs[i], chain)
					bootstrapped = false
					break
				}
			}
			if !bootstrapped {
				time.Sleep(waitDiff)
				continue
			}
			if peers, _ := client.Peers(ctx); len(peers) < constants.NumNodes-1 {
				color.Yellow("waiting for %s to connect to all peers (%d/%d)", nodeIDs[i], len(peers), constants.NumNodes-1)
				time.Sleep(waitDiff)
				continue
			}

			us, err := client.Uptime(ctx)
			if err != nil {
				color.Yellow("waiting for %s to begin staking", nodeIDs[i])
				time.Sleep(waitDiff)
				continue
			}

			color.Yellow("Uptime %s", us)

			color.Cyan("%s is bootstrapped and connected", nodeIDs[i])
			break
		}

	}

	color.Cyan("all nodes bootstrapped")
	close(bootstrapped)

	// Print endpoints where VM is accessible
	color.Green("standard VM endpoints now accessible at:")
	for i, url := range nodeURLs {
		color.Green("%s: %s", nodeIDs[i], url)
	}

	return nil
}

func runApp(g *errgroup.Group, ctx context.Context, nodeNum int, config node.Config) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("couldn't marshal config: %w", err)
	}
	fmt.Printf("node config %s", configJSON)
	app := process.NewApp(config)

	// Start running the AvalancheGo application
	if err := app.Start(); err != nil {
		return fmt.Errorf("node%d failed to start: %w", nodeNum+1, err)
	}

	g.Go(func() error {
		<-ctx.Done()
		_ = app.Stop()
		return ctx.Err()
	})

	exitCode, err := app.ExitCode()
	if (exitCode > 0 || err != nil) && ctx.Err() == nil {
		color.Red("node%d exited with code %d: %v", nodeNum+1, exitCode, err)
	}
	return err
}
