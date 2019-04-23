package core

import (
	"io"
	"net"
)

type NodesMap map[string]*Node
type Node struct {
	Conn      *net.Conn // Use generic Conn, so that we could use various conn type
	Lastseen  int       // The seconds since last time seen this node
	PublicKey string    // Public key of this node
	Address   string    // TCP-4 Address
}

type Network struct {
	Nodes              NodesMap     // Contacts
	ConnectionsQueue   chan string  // Connections
	Listener           net.Listener // Listener
	Address            string       // TCP-4 address
	ConnectionCallBack chan *Node   // Connections callback
	BroadcastQueue     chan Message // Broadcast queue
	IncommingMessages  chan Message // Incomming messages
}

type Peer struct {
	*KeyPair // Public key and private Key
	*Network // Network
}

// TODO: func NewNode() Node {}

// Use peer as node.
// Though this function is not so commonly used, it's still useful for testing.
func (p Peer) AsNode() Node {
	return Node{nil, 0, string(p.KeyPair.Public), p.Network.Address}
}

// Add node to peer's network.
func (p *Peer) AddNode(n Node) bool {
	pub := n.PublicKey

	if pub != string(p.KeyPair.Public) && p.Network.Nodes[pub] == nil {
		p.Network.Nodes[pub] = &n
		return true
	}

	return false
}

// Test if node n is online.
func (p *Peer) Ping(n Node) bool {
	return false
}

// Send message to specific node.
func (p *Peer) Send(n Node, m *Message) (error, int) {
	mBytes, err := m.MarshalBinary()
	if err != nil {
		return err, 0
	}

	bufLen, err := (*n.Conn).Write(mBytes)
	if err != nil {
		return err, 0
	}

	return nil, bufLen
}

// Receive Message from connected nodes.
func (p *Peer) Recv(n Node) (Message, error) {
	buf := make([]byte, 1024)

	bufLen, err := (*p.Nodes[n.PublicKey].Conn).Read(buf)
	if err != nil {
		if err != io.EOF {
			return Message{}, err
		}
	}

	m := new(Message)

	err = m.UnmarshalBinary(buf[:bufLen])
	if err != nil {
		return Message{}, err
	}

	return *m, nil
}

// Broadcast messages to nodes.
func (p *Peer) BroadcastMessage(m *Message) error {
	return nil
}
