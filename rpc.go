package goracoon

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type RPCServer struct{}

func (r *RPCServer) MaintenanceMode(inMaintenance bool, resp *string) error {

	if inMaintenance {
		maintenanceMode = true
		*resp = "Server in maintenance mode."
	} else {
		maintenanceMode = false
		*resp = "Server live."
	}

	return nil
}

func (gr *Goracoon) listenRPC() {
	// if no rpc port is specified, don't start
	rpcPort := os.Getenv("RPC_PORT")

	if rpcPort != "" {
		gr.Log.Info().Msg(fmt.Sprintf("Starting RPC server on port:%s", rpcPort))
		err := rpc.Register(new(RPCServer))
		if err != nil {
			gr.Log.Error().Err(err).Msg("")
			return
		}

		listen, err := net.Listen("tcp", "127.0.0.1:"+rpcPort)
		if err != nil {
			gr.Log.Error().Err(err).Msg("")
			return
		}

		for {
			rpcConn, err := listen.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(rpcConn)
		}
	}
}
