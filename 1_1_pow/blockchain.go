package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"time"
)

const Difficulty = 16 // Static difficulty for PoW calculation

var Empty [sha256.Size]byte // All zero sha256 value

// encodeUint64 encodes a uint64 to big endian notation. This code uses big
// endian in order to make the resulting values more readable for humans.
func encodeUint64(x uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, x)
	return b
}

// Block represents a single block in the blockchain. It is linked to the prior
// block via the PreviousBlockHash.
type Block struct {
	Timestamp         int64  // Timestamp block was mined
	Data              []byte // Blockchain data
	PreviousBlockHash []byte // Previous block hash in order link blocks
	Hash              []byte // PoW hash of this block
	Nonce             uint64 // Nonce used to calculate Hash
}

// NewBlock returns a block that is linked to previousBlockHash.
func NewBlock(data, previousBlockHash []byte) Block {
	timestamp := time.Now().Unix()
	return Block{
		Timestamp:         timestamp,
		Data:              data,
		PreviousBlockHash: previousBlockHash,
	}
}

// Verify ensures that the block is valid by hashing timestamp, data and nonce.
func (b Block) Verify() bool {
	t := encodeUint64(uint64(b.Timestamp))
	n := encodeUint64(b.Nonce)
	hash := sha256.Sum256(bytes.Join([][]byte{t, b.Data,
		b.PreviousBlockHash, n}, []byte{}))
	return bytes.Equal(hash[:], b.Hash)
}

// Mine attempts to mine the block within the provided range.
func (b *Block) Mine(difficulty uint, start, end uint64) error {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	t := encodeUint64(uint64(b.Timestamp))
	n := make([]byte, 8)
	bi := big.Int{}
	for i := start; i < end; i++ {
		binary.BigEndian.PutUint64(n, i)
		hash := sha256.Sum256(bytes.Join([][]byte{t, b.Data,
			b.PreviousBlockHash, n}, []byte{}))
		bi.SetBytes(hash[:])
		if bi.Cmp(target) == -1 {
			b.Hash = hash[:]
			b.Nonce = i
			return nil
		}
	}
	return fmt.Errorf("no solution for block")
}

// Blockchain is the blockchain context that houses an array of blocks.
type Blockchain struct {
	blocks []*Block
}

// Append adds a block, if valid, to the end of the blockchain.
func (b *Blockchain) Append(blk *Block) error {
	var previousBlockHash []byte
	if len(b.blocks) == 0 {
		// Genesis
		previousBlockHash = Empty[:]
	} else {
		previousBlockHash = b.blocks[len(b.blocks)-1].Hash
	}
	if !bytes.Equal(previousBlockHash, blk.PreviousBlockHash) {
		return fmt.Errorf("block does not link to previous block %x %x",
			previousBlockHash, blk.PreviousBlockHash)
	}
	if !blk.Verify() {
		return fmt.Errorf("can't append invalid block")
	}
	b.blocks = append(b.blocks, blk)
	return nil
}

// PrepareBlock returns a block template based on the current height of the
// blockchain.
func (b *Blockchain) PrepareBlock(data []byte) *Block {
	var previousBlockHash []byte
	if len(b.blocks) == 0 {
		// Genesis
		previousBlockHash = Empty[:]
	} else {
		previousBlockHash = b.blocks[len(b.blocks)-1].Hash
	}
	blk := NewBlock(data, previousBlockHash)
	return &blk
}

// Block returns a copy of the block at the specified block height.
func (b Blockchain) Block(block int) (Block, error) {
	if block > len(b.blocks) {
		return Block{}, fmt.Errorf("invalid block: %v", block)
	}
	return *b.blocks[block], nil
}

// Len returns the current blockchain height.
func (b Blockchain) Len() int {
	return len(b.blocks)
}

// NewBlockChain returns a blockchain context that has a genesis block.
func NewBlockChain(data []byte) (*Blockchain, error) {
	b := &Blockchain{}
	blk := b.PrepareBlock(data)
	err := blk.Mine(Difficulty, 0, math.MaxUint64)
	if err != nil {
		return nil, err
	}
	err = b.Append(blk)
	if err != nil {
		return nil, err
	}

	return b, nil
}
