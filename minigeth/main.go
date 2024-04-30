package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime/pprof"
	"strconv"

	"minigeth/oracle"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/triedb"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func shouldPreserveFn(header *types.Header) bool {
	return false
}

// preimageKey = PreimagePrefix + hash
func preimageKey(hash common.Hash) []byte {
	PreimagePrefix := []byte("secure-key-")
	return append(PreimagePrefix, hash.Bytes()...)
}

func main() {
	if len(os.Args) > 2 {
		f, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// non mips
	if len(os.Args) > 1 {
		newNodeUrl, setNewNodeUrl := os.LookupEnv("NODE")
		if setNewNodeUrl {
			fmt.Println("override node url", newNodeUrl)
			oracle.SetNodeUrl(newNodeUrl)
		}
		basedir := os.Getenv("BASEDIR")
		if len(basedir) == 0 {
			basedir = "/tmp/cannon"
		}

		pkw := &oracle.PreimageKeyValueWriter{}
		// todo
		pkwtrie := trie.NewStackTrie(pkw.WriteFn)

		blockNumber, _ := strconv.Atoi(os.Args[1])
		// TODO: get the chainid
		oracle.SetRoot(fmt.Sprintf("%s/0_%d", basedir, blockNumber))
		oracle.PrefetchBlock(big.NewInt(int64(blockNumber)), true, nil)
		oracle.PrefetchBlock(big.NewInt(int64(blockNumber)+1), false, pkwtrie)
		hash := pkwtrie.Hash()
		fmt.Println("committed transactions", hash)
	}

	// init secp256k1BytePoints
	crypto.S256()

	// get inputs
	inputBytes := oracle.Preimage(oracle.InputHash())
	var inputs [6]common.Hash
	for i := 0; i < len(inputs); i++ {
		inputs[i] = common.BytesToHash(inputBytes[i*0x20 : i*0x20+0x20])
	}

	// read start block header
	var parent types.Header
	check(rlp.DecodeBytes(oracle.Preimage(inputs[0]), &parent))

	// read header
	var newheader types.Header
	// from parent
	newheader.ParentHash = parent.Hash()
	newheader.Number = big.NewInt(0).Add(parent.Number, big.NewInt(1))
	// todo: basefee may be different
	newheader.BaseFee = eip1559.CalcBaseFee(params.MainnetChainConfig, &parent)

	// from input oracle
	newheader.TxHash = inputs[1]
	newheader.Coinbase = common.BigToAddress(inputs[2].Big())
	newheader.UncleHash = inputs[3]
	newheader.GasLimit = inputs[4].Big().Uint64()
	newheader.Time = inputs[5].Big().Uint64()
	fmt.Printf("new header: %+v\n", newheader)

	vmconfig := vm.Config{}

	db := oracle.NewDB()

	// testAddr := common.HexToAddress("0xE01B7Ac63fF9f87C2BFbe4E28015B80bbf1ACD02")
	// oracle.PrefetchAccount(parent.Number, testAddr, nil)
	// for hash, value := range oracle.Preimages() {
	// 	db.Put(hash.Bytes(), value)
	// }

	// trieDb.Update(newheader.Root, parent.Root, newheader.Number.Uint64(), nil, nil)
	// fmt.Println("udpate done")
	// todo: genesis block is for what?
	bc, err := core.NewBlockChain(db, nil, nil, nil, &ethash.Ethash{}, vmconfig, nil, nil)
	if err != nil {
		panic(fmt.Sprint("failed to new blockchain:", err))
	}
	// bc.SetHead(parent.Number.Uint64())

	processor := core.NewStateProcessor(params.MainnetChainConfig, bc, bc.Engine())
	fmt.Println("processing state:", parent.Number, "->", newheader.Number)

	newheader.Difficulty = bc.Engine().CalcDifficulty(bc, newheader.Time, &parent)

	// read txs
	//traverseStackTrie(newheader.TxHash)

	//fmt.Println(fn)
	//fmt.Println(txTrieRoot)
	var txs []*types.Transaction

	// todo: no transactions
	oracle.PrefetchAccount(parent.Number, common.Address{}, nil)
	oracle.PrefetchStorage(parent.Number, common.HexToAddress("0x0000000000004946c0e9F43F4Dee607b0eF1fA1c"), common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"), nil)
	oracle.PrefetchStorage(parent.Number, common.HexToAddress("0x0000000000004946c0e9F43F4Dee607b0eF1fA1c"), common.HexToHash("0x15f4db769e203d3965138078e54ad9d833f17c87a721a8e75e4c8bfa1c083a4a"), nil)
	trieDb := triedb.NewDatabase(db, nil)
	// trieDb.Initialized(newheader.TxHash)
	// trieDb.InsertPreimage(oracle.Preimages())
	// fmt.Println("preimages size is ", len(oracle.Preimages()))
	// for key := range oracle.Preimages() {
	// 	fmt.Println(key)
	// }
	// for hash := range oracle.Preimages() {
	// 	_, err := db.Get(preimageKey(hash))
	// 	fmt.Println("err is ", err)
	// }
	// fmt.Println("write preimages..")
	// trieDb.WritePreimages()
	// for hash := range oracle.Preimages() {
	// 	_, err := db.Get(preimageKey(hash))
	// 	fmt.Println("err is ", err)
	// }

	// trieDb.Commit(newheader.TxHash, true)
	// tt, _ := trie.New(newheader.TxHash, &triedb)
	fmt.Println("tx hash is ", newheader.TxHash)
	fmt.Println("root is ", newheader.Root)
	t, err := trie.NewStateTrie(trie.StateTrieID(newheader.TxHash), trieDb)
	if err != nil {
		panic(fmt.Sprint("new state trie failed:", err))
	}
	tni, _ := t.NodeIterator([]byte{})
	for tni.Next(true) {
		fmt.Println(tni.Hash(), tni.Leaf(), tni.Path(), tni.Error())
		if tni.Leaf() {
			tx := types.Transaction{}
			var rlpKey uint64
			check(rlp.DecodeBytes(tni.LeafKey(), &rlpKey))
			check(tx.UnmarshalBinary(tni.LeafBlob()))
			// TODO: resize an array in go?
			for uint64(len(txs)) <= rlpKey {
				txs = append(txs, nil)
			}
			txs[rlpKey] = &tx
		}
	}
	fmt.Println("read", len(txs), "transactions")
	// // TODO: OMG the transaction ordering isn't fixed

	var uncles []*types.Header
	check(rlp.DecodeBytes(oracle.Preimage(newheader.UncleHash), &uncles))

	var receipts []*types.Receipt
	block := types.NewBlock(&newheader, txs, uncles, receipts, trie.NewStackTrie(nil))
	fmt.Println("made block, parent:", newheader.ParentHash)

	// if this is correct, the trie is working
	// TODO: it's the previous block now
	if newheader.TxHash != block.Header().TxHash {
		panic("wrong transactions for block")
	}
	if newheader.UncleHash != block.Header().UncleHash {
		panic("wrong uncles for block " + newheader.UncleHash.String() + " " + block.Header().UncleHash.String())
	}

	statedb, err := state.New(parent.Root, oracle.NewStateDB(parent, state.NewDatabaseWithNodeDB(db, trieDb), trieDb, db), nil)
	if err != nil {
		panic(fmt.Sprint("failed to new state:", err))
	}

	// validateState is more complete, gas used + bloom also
	receipts, _, _, err = processor.Process(block, statedb, vmconfig)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", receipts[0])

	receiptSha := types.DeriveSha(types.Receipts(receipts), trie.NewStackTrie(nil))

	newroot := statedb.IntermediateRoot(bc.Config().IsEIP158(newheader.Number))

	fmt.Println("receipt count", len(receipts), "hash", receiptSha)
	fmt.Println("process done with hash", parent.Root, "->", newroot)
	oracle.Output(newroot, receiptSha)
}
