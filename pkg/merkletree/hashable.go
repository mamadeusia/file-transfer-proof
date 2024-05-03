package merkletree

import (
	"encoding/hex"
	"io"

	"golang.org/x/crypto/sha3"
)

func HashFromReader(r io.Reader) (string, error) {
	hash := sha3.NewLegacyKeccak256()
	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}
	val := hash.Sum(nil)
	return hex.EncodeToString(val[:]), nil

}

func hash(data []byte) Hash {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	val := hash.Sum(nil)
	var output Hash
	copy(output[:], val[0:32])
	return output
}

type Hashable interface {
	hash() Hash
}

type Hash [32]byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}
