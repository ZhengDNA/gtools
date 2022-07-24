package test

import (
	"encoding/json"
	"fmt"
	"gtools/contains"
	"testing"
)

func TestMarshal(t *testing.T) {
	tree, _ := contains.NewHuffTree([]byte("7u8irt3e4vuyb8irtedvuybserwvuyb"))
	bytes, err := json.Marshal(tree)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(string(bytes))
	fmt.Println("-----------------------------------")
	marshalHuffTree := contains.MarshalHuffTree(tree)
	huffTree := contains.UnmarshalHuffTree(marshalHuffTree)
	bytes, err = json.Marshal(huffTree)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(string(bytes))
}
