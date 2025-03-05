package lock_container

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"math/rand"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

const (
	randomStringLength = 10
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

func markFocusedContainer(node *i3.Node, mark string) error {
	if _, err := i3.RunCommand(fmt.Sprintf("[con_id=%d] mark --add %s", node.ID, mark)); err != nil {
		return fmt.Errorf("failed to mark focused container: %w", err)
	}

	return nil
}

func Execute() error {
	node := common.GetFocusedNode()
	md5Hash := generateMD5Hash(generateRandomString(randomStringLength))

	if err := markFocusedContainer(node, md5Hash); err != nil {
		return fmt.Errorf("failed to mark focused container: %w", err)
	}

	return nil
}
