<div align="center">
  <img src="resources/AvalancheLogoRed.png?raw=true">
</div>

# ava-sim
`ava-sim` makes it easy for anyone to spin up a local instance of an Avalanche network
to interact with the [standard APIs](https://docs.avax.network/build/avalanchego-apis)
or to test a [custom
VM](https://docs.avax.network/build/tutorials/platform/create-custom-blockchain). This version has been adapted for:
- Compatiblity with the latest version of avalanchego
- Specific needs for running zapavm
- Ability to spin up a 6th node that bootstraps from an already running network.

### Running

- You must have [Golang](https://golang.org/doc/install) >= `1.16` and a configured
[`$GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH).

1. Build Zapavm or use the pre-built binary. See https://github.com/zapalabs/zapavm#building for instructions on how to build the plugin.
2. Write out the following files. These files are queried by the nodes starting up to know which node number they are, which lets them infer the port number of their corresponding zcash.

```
=>echo "NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg" > ~/node-ids/0
=>echo "NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ" > ~/node-ids/1
=>echo "NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN" > ~/node-ids/2
=>echo "NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu" > ~/node-ids/3
=>echo "NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5" > ~/node-ids/4
```

3. From the `ava-sim` project root, run `./scripts/prepare-system-plugins.sh`

4. Run it.
```
go run main/main.go ../zapavm/builds/zapavm ../zapavm/builds/emptygenesis.txt
```

5. Shutting down. After terminating the process, run `pkill -f ava-sim`. I've found this step necessary as terminating the process sometimes leaves orphaned processes.

## Debugging with VS Code

If you'd prefer to debug with VS Code, follow steps 1-3 above, and then in VS Code navigate to `Run and Debug`, and select `Zapavm` which is defined in the [launch.json](./.vscode/launch.json) to run. You will have to update this spec to be compatible with your own directory structure.

## Logs

After initial network startup,
you'll see the following logs when all validators in the network are validating
your custom blockchain (the actual blockchain ID will vary every time):
```txt
NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg validating blockchain 28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ validating blockchain 28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN validating blockchain 28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu validating blockchain 28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5 validating blockchain 28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
```

When the VM is ready to interact with, the URLs it is accessible on will be
printed out:
```txt
Custom VM endpoints now accessible at:
NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg: http://127.0.0.1:9652/ext/bc/28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-MFrZFVCXPv5iCn6M9K6XduxGTYp891xXZ: http://127.0.0.1:9654/ext/bc/28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-NFBbbJ4qCmNaCzeW7sxErhvWqvEQMnYcN: http://127.0.0.1:9656/ext/bc/28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-GWPcbFJZFfZreETSoWjPimr846mXEKCtu: http://127.0.0.1:9658/ext/bc/28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
NodeID-P7oB2McjBGgW2NXXWVYjV8JEDFoW9xDE5: http://127.0.0.1:9660/ext/bc/28TtJ7sdYvdgfj1CcXo5o3yXFMhKLrv4FQC9WhgSHgY6YNYRs2
```

## Interacting with the Chain

Use the per-node endpoints defined above to issue API requests. See [available API methods](https://github.com/zapalabs/zapavm#api)

## Bootstrapping a 6th Node

Once you have 5 nodes running, you might want to observe bootstrapping behavior. To spin up a 6th node, run:

```
go run main/main.go ../zapavm/builds/zapavm
```

or see the `New Node for Existing Network` in [launch.json](./.vscode/launch.json) .
