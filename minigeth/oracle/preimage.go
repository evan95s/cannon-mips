//go:build !mips
// +build !mips

package oracle

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
)

var memoryDB = memorydb.New()
var preimages = make(map[common.Hash][]byte)
var root = "/tmp/cannon"
var parentBlockNumber = big.NewInt(0)

func SetRoot(newRoot string) {
	root = newRoot
	err := os.MkdirAll(root, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func GetMemoryDB() *memorydb.Database {
	return memoryDB
}

func SetBlockNumber(num int64) {
	parentBlockNumber = big.NewInt(num)
}

func GetBlockNumber() *big.Int {
	return parentBlockNumber
}

func Preimage(hash common.Hash) []byte {
	val, ok := preimages[hash]
	key := fmt.Sprintf("%s/%s", root, hash)
	// We write the preimage even if its value is nil (will result in an empty file).
	// This can happen if the hash represents a full node that is the child of another full node
	// that collapses due to a key deletion. See fetching-preimages.md for more details.
	err := ioutil.WriteFile(key, val, 0644)
	check(err)
	comphash := crypto.Keccak256Hash(val)
	if ok && hash != comphash {
		panic("corruption in hash " + hash.String())
	}
	return val
}

func Preimages() map[common.Hash][]byte {
	return preimages
}

// PreimageKeyValueWriter wraps the Put method of a backing data store.
type PreimageKeyValueWriter struct{}

// Put inserts the given value into the key-value data store.
func (kw *PreimageKeyValueWriter) Put(key []byte, value []byte) error {
	hash := crypto.Keccak256Hash(value)
	if hash != common.BytesToHash(key) {
		panic("bad preimage value write")
	}
	preimages[hash] = common.CopyBytes(value)
	memoryDB.Put(hash.Bytes(), common.CopyBytes(value))
	return nil
}

// WriteFn is a replacer of Put
// However, path is still unknown
// hash is common.BytesToHash(key)
// blob is value
func (kw *PreimageKeyValueWriter) WriteFn(path []byte, hash common.Hash, blob []byte) {
	if hash != crypto.Keccak256Hash(blob) {
		panic("bad preimage value write")
	}
	preimages[hash] = common.CopyBytes(blob)
	memoryDB.Put(hash.Bytes(), common.CopyBytes(blob))
}
