package ast

import (
	"fmt"
)

func (n *Node) String() string {
	if !n.IsParent() {
		if s, ok := n.Value.(string); ok {
			return fmt.Sprintf("[%q,%q]", n.TypeString(), s)
		}
		return fmt.Sprintf("[%q,\"%v\"]", n.TypeString(), n.Value)
	}
	jsonString := fmt.Sprintf("[%q,", n.TypeString())
	for idx, child := range n.Children() {
		if idx != 0 {
			jsonString += ","
		}
		jsonString += child.String()
	}
	jsonString += "]"
	return jsonString
}

func (n *Node) MarshalJSON() ([]byte, error) {
	jsonString, err := n.MarshalJSONString()
	return []byte(jsonString), err
}

func (n *Node) MarshalJSONString() (string, error) {
	if !n.IsParent() {
		if s, ok := n.Value.(string); ok {
			return fmt.Sprintf("[%d,%q]", n.Type, s), nil
		}
		return fmt.Sprintf("[%d,\"%v\"]", n.Type, n.Value), nil
	}
	jsonString := "["
	for idx, child := range n.Children() {
		if idx != 0 {
			jsonString += ","
		}
		childString, err := child.MarshalJSONString()
		if err != nil {
			return "", err
		}
		jsonString += childString
	}
	jsonString += "]"
	return jsonString, nil
}
