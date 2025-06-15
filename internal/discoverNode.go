package internal

import (
	"context"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	libp2p "github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p/core/host"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type Node struct {
	Host  host.Host
	Peers map[string]peer.ID
	Lock  sync.Mutex
	Role  string
}

func NewNode() (*Node, error) {
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0")) // create new libp2p host
	if err != nil {
		return nil, err
	}
	role := strings.ToUpper(os.Getenv("ROLE")) // reading role from env
	return &Node{
		Host:  h,
		Peers: make(map[string]peer.ID),
		Role:  role,
	}, nil
}

func (n *Node) SetupDiscovery() error {
	notif := &discoveryNotifee{node: n}
	service := mdns.NewMdnsService(n.Host, DiscoveryTag, notif) //mDNS discovery (libp2p)
	return service.Start()
}

type discoveryNotifee struct {
	node *Node
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.node.Host.ID() {
		return
	}

	n.node.Lock.Lock()
	n.node.Host.Connect(context.Background(), pi) //connect to discovered peer (libp2p)
	n.node.Lock.Unlock()

	go n.node.SendMyRole(pi.ID)
}

func (n *Node) SendMyRole(peerID peer.ID) {
	s, err := n.Host.NewStream(context.Background(), peerID, RoleProtocolID) //open stream (libp2p)
	if err != nil {
		log.Printf("Error opening role stream to %s: %v", peerID, err)
		return
	}
	defer s.Close()

	s.Write([]byte(n.Role))
	log.Printf("Sent my role [%s] to peer %s", n.Role, peerID)
}

func (n *Node) HandleRoleExchange(s network.Stream) {
	defer s.Close()
	data, err := io.ReadAll(s)
	if err != nil {
		return
	}
	role := strings.ToUpper(string(data)) // store peer ID for role

	n.Lock.Lock()
	n.Peers[role] = s.Conn().RemotePeer()
	n.Lock.Unlock()

	log.Printf("Received role [%s] from peer %s", role, s.Conn().RemotePeer())
}
