package goracoon

import (
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
		gr.InfoLog.Println("Starting RPC server on port:", rpcPort)
		err := rpc.Register(new(RPCServer))
		if err != nil {
			gr.ErrorLog.Println(err)
			return
		}

		listen, err := net.Listen("tcp", "127.0.0.1:"+rpcPort)
		if err != nil {
			gr.ErrorLog.Println(err)
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
