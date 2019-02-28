package gosnowth

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
)

const topologyTestData = `[
	{
		"id": "1f846f26-0cfd-4df5-b4f1-e0930604e577",
		"address": "10.8.20.1",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "765ac4cc-1929-4642-9ef1-d194d08f9538",
		"address": "10.8.20.2",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "8c2fc7b8-c569-402d-a393-db433fb267aa",
		"address": "10.8.20.3",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	},
	{
		"id": "07fa2237-5744-4c28-a622-a99cfc1ac87e",
		"address": "10.8.20.4",
		"port": 8112,
		"apiport": 8112,
		"weight": 32,
		"n": 2
	}
]`

const topologyXMLTestData = `<nodes n="2">
	<node id="1f846f26-0cfd-4df5-b4f1-e0930604e577"
		address="10.8.20.1"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="765ac4cc-1929-4642-9ef1-d194d08f9538"
		address="10.8.20.2"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="8c2fc7b8-c569-402d-a393-db433fb267aa"
		address="10.8.20.3"
		port="8112"
		apiport="8112"
		weight="32"/>
	<node id="07fa2237-5744-4c28-a622-a99cfc1ac87e"
		address="10.8.20.4"
		port="8112"
		apiport="8112"
		weight="32"/>
</nodes>`

func TestTopologyJSONDeserialization(t *testing.T) {
	dec := json.NewDecoder(bytes.NewBufferString(topologyTestData))
	dec.UseNumber()
	topo := []TopologyNode{}
	err := dec.Decode(&topo)
	if err != nil {
		t.Errorf("failed to decode topology, %s\n", err.Error())
	}

	if len(topo) != 4 {
		t.Error("should be 4 elements")
	}
}

func TestTopologyXMLDeserialization(t *testing.T) {
	dec := xml.NewDecoder(bytes.NewBufferString(topologyXMLTestData))
	topo := new(Topology)
	err := dec.Decode(topo)
	if err != nil {
		t.Errorf("failed to decode topology, %s\n", err.Error())
	}

	if len(topo.Nodes) != 4 {
		t.Error("should be 4 elements")
	}
}

func TestTopologyXMLSerialization(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	enc := xml.NewEncoder(buf)
	topo := Topology{
		NumberNodes: 2,
		Nodes: []TopologyNode{
			TopologyNode{
				ID:      "1f846f26-0cfd-4df5-b4f1-e0930604e577",
				Address: "10.8.20.1",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			TopologyNode{
				ID:      "765ac4cc-1929-4642-9ef1-d194d08f9538",
				Address: "10.8.20.2",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			TopologyNode{
				ID:      "8c2fc7b8-c569-402d-a393-db433fb267aa",
				Address: "10.8.20.3",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
			TopologyNode{
				ID:      "07fa2237-5744-4c28-a622-a99cfc1ac87e",
				Address: "10.8.20.4",
				Port:    8112,
				APIPort: 8112,
				Weight:  32,
			},
		},
	}

	err := enc.Encode(topo)
	if err != nil {
		t.Errorf("failed to decode node stats, %s\n", err.Error())
	}

	if strings.Count(buf.String(), "id=") != 4 {
		t.Error("should have 4 nodes")
	}
}
