package ast

import (
	"fmt"
)

func (n *Node) String() string {
	if !n.IsParent() {
		return fmt.Sprintf("[%q,%q]", n.TypeString(), n.Value)
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
		return fmt.Sprintf("[%d,%q]", n.Type, n.Value), nil
	}
	jsonString := fmt.Sprintf("[%d,", n.Type)
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
