package gosnowth

import (
	"context"
	"encoding/xml"
	"path"
)

// TopoRing values represent IRONdb topology ring data.
type TopoRing struct {
	XMLName      xml.Name         `xml:"vnodes" json:"-"`
	VirtualNodes []TopoRingDetail `xml:"vnode"`
	NumberNodes  int              `xml:"n,attr" json:"-"`
}

// TopoRingDetail values represent IRONdb topology ring node details.
type TopoRingDetail struct {
	XMLName  xml.Name `xml:"vnode" json:"-"`
	ID       string   `xml:"id,attr" json:"id"`
	IDX      int      `xml:"idx,attr" json:"idx"`
	Location float64  `xml:"location,attr" json:"location"`
}

// GetTopoRingInfo retrieves topology ring information from a node.
func (sc *SnowthClient) GetTopoRingInfo(hash string,
	node *SnowthNode) (*TopoRing, error) {
	return sc.GetTopoRingInfoContext(context.Background(), hash, node)
}

// GetTopoRingInfoContext is the context aware version of GetTopoRingInfo.
func (sc *SnowthClient) GetTopoRingInfoContext(ctx context.Context,
	hash string, node *SnowthNode) (*TopoRing, error) {
	tr := new(TopoRing)
	err := sc.do(ctx, node, "GET", path.Join("/toporing/xml", hash),
		nil, tr, decodeXML)
	return tr, err
}
