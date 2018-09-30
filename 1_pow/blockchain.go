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

const Difficulty = 16

var Empty [sha256.Size]byte

func encodeUint64(x uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, x)
	return b
}

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             uint64
}

func NewBlock(data, previousBlockHash []byte) Block {
	timestamp := time.Now().Unix()
	return Block{
		Timestamp:         timestamp,
		Data:              data,
		PreviousBlockHash: previousBlockHash,
	}
}

func (b Block) Verify() bool {
	t := encodeUint64(uint64(b.Timestamp))
	n := encodeUint64(b.Nonce)
	hash := sha256.Sum256(bytes.Join([][]byte{t, b.Data,
		b.PreviousBlockHash, n}, []byte{}))
	return bytes.Equal(hash[:], b.Hash)
}

func (b *Block) Mine(difficulty uint) error {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))
	t := encodeUint64(uint64(b.Timestamp))
	n := make([]byte, 8)
	bi := big.Int{}
	for i := uint64(0); i < math.MaxInt64; i++ {
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

type Blockchain struct {
	blocks []*Block
}

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
	b.blocks = append(b.blocks, blk)
	return nil
}

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

func (b Blockchain) Block(block int) (Block, error) {
	if block > len(b.blocks) {
		return Block{}, fmt.Errorf("invalid block: %v", block)
	}
	return *b.blocks[block], nil
}

func (b Blockchain) Len() int {
	return len(b.blocks)
}

func NewBlockChain(data []byte) (*Blockchain, error) {
	b := &Blockchain{}
	blk := b.PrepareBlock(data)
	err := blk.Mine(Difficulty)
	if err != nil {
		return nil, err
	}
	err = b.Append(blk)
	if err != nil {
		return nil, err
	}

	return b, nil
}
