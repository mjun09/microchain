package core

import (
	"net"
	"strconv"
)

// RemoteNode ... Represent other nodes.
type RemoteNode struct {
	PublicKey  []byte        // Public key
	Address    *net.TCPAddr  // Address
	Lastseen   int           // The unix time of seeing this node last time
	VerifiedBy []*RemoteNode // Nodes that verify this node
}

// Packet ... Received packect.
type Packet struct {
	Content []byte       // Raw bytes
	Conn    *net.TCPConn // TCP connection
}

// IncommingMessage ...
type IncommingMessage struct {
	Content Message      // Message
	Conn    *net.TCPConn // TCP connection
}

// Node ... Represent ourselves.
type Node struct {
	Keypair        *KeyPair               // Key pair
	IP             string                 // IP address
	Port           int                    // Port
	RoutingTable   map[string]*RemoteNode // Routing table (public key, node)
	Listerner      *net.TCPListener       // TCP listener
	MessageChannel chan IncommingMessage  // Incomming message
}

// NewNode ... Generate new node.
func NewNode(ip string, port int) (*Node, error) {
	kp, err := NewECDSAKeyPair()
	if err != nil {
		return nil, err
	}

	return &Node{
		Keypair:        kp,
		IP:             ip,
		Port:           port,
		RoutingTable:   make(map[string]*RemoteNode),
		Listerner:      new(net.TCPListener),
		MessageChannel: make(chan IncommingMessage),
	}, nil
}

// Run ... Run a simple TCP server.
func (n *Node) Run() error {
	// TODO: Error handling
	addr, err := net.ResolveTCPAddr("tcp", n.IP+":"+strconv.Itoa(n.Port))
	if err != nil {
		return err
	}

	// TODO: Error handling
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	n.Listerner = listener

	incommingPacket := make(chan Packet)

	go n.receivePacket(incommingPacket)
	go n.processPacket(incommingPacket)

	return nil
}

// receivePacket ... Listen on binding address.
func (n *Node) receivePacket(packetch chan Packet) {
	for {
		// TODO: Error handling
		conn, _ := n.Listerner.AcceptTCP()

		buf := make([]byte, 4096)

		// TODO: Error handling
		bufLen, _ := conn.Read(buf)

		p := Packet{
			Content: buf[:bufLen],
			Conn:    conn,
		}

		// send packet to channel
		packetch <- p
	}
}

// processPacket ... Process packet.
func (n *Node) processPacket(packetch chan Packet) {
	for p := range packetch {
		var m Message
		// TODO: Error handling
		err := m.UnmarshalJson(p.Content)
		if err != nil {
			// We just drop the malformed message
			continue
		}

		n.MessageChannel <- IncommingMessage{Content: m, Conn: p.Conn}
	}
}
