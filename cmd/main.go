package main

import (
	"log"
	"os"
	"os/signal"
	"p2papp/internal"
	"syscall"
	"time"
)

func main() {
	node, err := internal.NewNode() // initialize new p2p node(libp2p)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Node ID:", node.Host.ID().String()) //node's Peer ID (libp2p)
	for _, addr := range node.Host.Addrs() {
		log.Printf("Listening on: %s/p2p/%s", addr, node.Host.ID().String())
	}

	err = node.SetupDiscovery() // start peer discovery (libp2p-mDNS)
	if err != nil {
		log.Fatal(err)
	}

	// register protocol handlers for task and role protocols
	node.Host.SetStreamHandler(internal.TaskProtocolID, node.HandleTask)
	node.Host.SetStreamHandler(internal.RoleProtocolID, node.HandleRoleExchange)

	go func() {
		time.Sleep(5 * time.Second)
		if os.Getenv("IS_DISPATCHER") == "1" {
			node.RunDispatcher()
		}
	}()

	waitForShutdown(node)
}

func waitForShutdown(n *internal.Node) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("Shutting down...")
	n.Host.Close() // Close libp2p host (libp2p)
}
