package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.i3wm.org/i3/v4"
)

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func generateMD5Hash(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func getFocusedNode() *i3.Node {
	i3Tree := getI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	return node
}

func markFocusedContainer(mark string) {
	node := getFocusedNode()

	i3.RunCommand(fmt.Sprintf("[con_id=%d] mark --add %s", node.ID, mark))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	randomString := generateRandomString(10)
	md5Hash := generateMD5Hash(randomString)

	//fmt.Printf("Random String: %s\n", randomString)
	//fmt.Printf("MD5 Hash: %s\n", md5Hash)

	markFocusedContainer(md5Hash)
}
