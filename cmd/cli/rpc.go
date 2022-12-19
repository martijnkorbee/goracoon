package main

import (
	"net/rpc"
	"os"

	"github.com/fatih/color"
)

func rpcMaintenanceMode(inMaintenanceMode bool) {
	rpcPort := os.Getenv("RPC_PORT")

	c, err := rpc.Dial("tcp", "127.0.0.1:"+rpcPort)
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("\tConnected to RPC client...")
	var result string
	err = c.Call("RPCServer.MaintenanceMode", inMaintenanceMode, &result)
	if err != nil {
		exitGracefully(err)
	}

	color.HiWhite(result)
}
