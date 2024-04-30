package wrapper

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/oracle"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/trie/trienode"
	"github.com/ethereum/go-ethereum/triedb"
)

type StateDB struct {
	blockNumber *big.Int
	ethdb       ethdb.Database
	inner       state.Database
	triedb      *triedb.Database
}

func NewStateDB(header types.Header, inner state.Database, triedb *triedb.Database, ethdb ethdb.Database) state.Database {
	return &StateDB{
		blockNumber: header.Number,
		inner:       inner,
		ethdb:       ethdb,
		triedb:      triedb,
	}
}

func (db *StateDB) ContractCode(addr common.Address, codeHash common.Hash) ([]byte, error) {
	addrHash := crypto.Keccak256Hash(addr.Bytes())
	fmt.Println("want contract code of ", addrHash)
	oracle.PrefetchCode(db.blockNumber, addrHash)
	code := oracle.Preimage(codeHash)
	return code, nil
}

func (db *StateDB) ContractCodeSize(addr common.Address, codeHash common.Hash) (int, error) {
	addrHash := crypto.Keccak256Hash(addr.Bytes())
	oracle.PrefetchCode(db.blockNumber, addrHash)
	code := oracle.Preimage(codeHash)
	return len(code), nil
}

func (db *StateDB) CopyTrie(t state.Trie) state.Trie {
	return db.inner.CopyTrie(t)
}

func (db *StateDB) DiskDB() ethdb.KeyValueStore {
	return db.inner.DiskDB()
}

// it seems this function will be called
func (db *StateDB) OpenStorageTrie(stateRoot common.Hash, address common.Address, root common.Hash, trie state.Trie) (state.Trie, error) {
	return db.inner.OpenStorageTrie(stateRoot, address, root, trie)
}

func (db *StateDB) OpenTrie(root common.Hash) (state.Trie, error) {
	// fmt.Println("open trie of root:", root)
	trie, err := db.inner.OpenTrie(root)
	if err != nil {
		return trie, err
	}
	return NewTrie(root, db.blockNumber, trie, db.ethdb, db, db.triedb), nil
}

func (db *StateDB) TrieDB() *triedb.Database {
	return db.inner.TrieDB()
}

type Trie struct {
	root        common.Hash
	blockNumber *big.Int
	inner       state.Trie
	ethdb       ethdb.Database
	statedb     *StateDB
	triedb      *triedb.Database
}

func NewTrie(root common.Hash, blockNumber *big.Int, inner state.Trie, ethdb ethdb.Database, statedb *StateDB, triedb *triedb.Database) state.Trie {
	return &Trie{
		root:        root,
		blockNumber: blockNumber,
		inner:       inner,
		ethdb:       ethdb,
		statedb:     statedb,
		triedb:      triedb,
	}
}

func (t *Trie) Commit(collectLeaf bool) (common.Hash, *trienode.NodeSet, error) {
	return t.inner.Commit(collectLeaf)
}

func (t *Trie) DeleteAccount(address common.Address) error {
	panic("delete account")
	// PrefetchAccount(t.blockNumber, address, nil)
	return t.inner.DeleteAccount(address)
}

func (t *Trie) DeleteStorage(addr common.Address, key []byte) error {
	panic("delete storage")
	return t.inner.DeleteStorage(addr, key)
}

func (t *Trie) GetAccount(address common.Address) (*types.StateAccount, error) {
	// fmt.Println("get account:", address)
	// panic(t.blockNumber.Uint64())
	oracle.PrefetchAccount(t.blockNumber, address, nil)
	// inner, err := t.statedb.inner.OpenTrie(t.root)
	// if err != nil {
	// 	panic(err)
	// }
	// t.inner = inner

	return t.inner.GetAccount(address)
}

func (t *Trie) GetKey(key []byte) []byte {
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>here get key", key)
	return t.inner.GetKey(key)
}

func (t *Trie) GetStorage(addr common.Address, key []byte) ([]byte, error) {
	oracle.PrefetchStorage(t.blockNumber, addr, crypto.Keccak256Hash(key[:]), nil)
	return t.inner.GetStorage(addr, key)
}

func (t *Trie) Hash() common.Hash {
	return t.inner.Hash()
}

func (t *Trie) NodeIterator(startKey []byte) (trie.NodeIterator, error) {
	return t.inner.NodeIterator(startKey)
}

func (t *Trie) Prove(key []byte, proofDb ethdb.KeyValueWriter) error {
	return t.inner.Prove(key, proofDb)
}

// seems ok because 15537372 is ok
func (t *Trie) UpdateAccount(address common.Address, account *types.StateAccount) error {
	oracle.UpdateAccount(t.blockNumber, address)
	return t.inner.UpdateAccount(address, account)
}

func (t *Trie) UpdateContractCode(address common.Address, codeHash common.Hash, code []byte) error {
	// panic("udpate contract code")
	return t.inner.UpdateContractCode(address, codeHash, code)
}

func (t *Trie) UpdateStorage(addr common.Address, key []byte, value []byte) error {
	// panic("update storage")
	return t.inner.UpdateStorage(addr, key, value)
}
