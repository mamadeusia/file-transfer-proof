package merkletree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashleaf(t *testing.T) {
	l := leaf("5B38Da6a701c568545dCfcB03FcB875f56beddC4")
	leafHash := l.hash().String()
	assert.Equal(t, "5931b4ed56ace4c46b68524cb5bcbf4195f1bbaacbe5228fbd090546c88dd229", leafHash)
}

func TestHashleafHash(t *testing.T) {
	// should not change any thing, supposed the leaf is the hash of the data .
	l := leafHash("c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	leafHash := l.hash().String()
	assert.Equal(t, "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470", leafHash)

}

func TestNewNode(t *testing.T) {
	left := leafHash("c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470")
	right := leaf("5B38Da6a701c568545dCfcB03FcB875f56beddC4")
	node := NewNode(left, right)

	assert.Equal(t, node.right, left)
	assert.Equal(t, node.left, right)
}

func TestEmptyleafHash(t *testing.T) {
	// hash(leafHash)
}

func Test_hashableleafHashFromString(t *testing.T) {
	hashable, err := hashableLeafHashFromString("483095a17331c667da935a465ec026b7100e5d200226ed40da13c15cd2d293b8")
	assert.Nil(t, err)
	assert.Equal(t, hashable.hash().String(), "483095a17331c667da935a465ec026b7100e5d200226ed40da13c15cd2d293b8")

	hashable, err = hashableLeafHashFromString("483095a17331c667da935a465ec026b710adbewqw3lknmcd2d293b8")
	assert.Error(t, err)
	assert.Nil(t, hashable)

}
