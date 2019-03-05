package gosnowth

import (
	"encoding/xml"
	"path"

	"github.com/pkg/errors"
)

// Topology values represent IRONdb topology structure.
type Topology struct {
	XMLName     xml.Name       `xml:"nodes" json:"-"`
	NumberNodes int            `xml:"n,attr" json:"-"`
	Hash        string         `xml:"-"`
	Nodes       []TopologyNode `xml:"node"`
}

// TopologyNode represent a node in the IRONdb topology structure.
type TopologyNode struct {
	XMLName     xml.Name `xml:"node" json:"-"`
	ID          string   `xml:"id,attr" json:"id"`
	Address     string   `xml:"address,attr" json:"address"`
	Port        uint16   `xml:"port,attr" json:"port"`
	APIPort     uint16   `xml:"apiport,attr" json:"apiport"`
	Weight      int      `xml:"weight,attr" json:"weight"`
	NumberNodes int      `xml:"-" json:"n"`
}

// GetTopologyInfo retrieves topology information from a node.
func (sc *SnowthClient) GetTopologyInfo(node *SnowthNode) (*Topology, error) {
	t := new(Topology)
	err := sc.do(node, "GET",
		path.Join("/topology/xml", node.GetCurrentTopology()),
		nil, t, decodeXML)
	return t, err
}

// LoadTopology loads a new topology on a node without activating it.
func (sc *SnowthClient) LoadTopology(hash string, t *Topology,
	node *SnowthNode) error {
	b, err := encodeXML(t)
	if err != nil {
		return errors.Wrap(err, "failed to encode request data")
	}

	return sc.do(node, "POST", path.Join("/topology", hash), b, nil, nil)
}

// ActivateTopology activates a new topology on the node. THIS IS DANGEROUS.
func (sc *SnowthClient) ActivateTopology(hash string, node *SnowthNode) error {
	return sc.do(node, "GET", path.Join("/activate", hash), nil, nil, nil)
}
