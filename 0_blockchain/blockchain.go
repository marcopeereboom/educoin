package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

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
}

func NewBlock(data, previousBlockHash []byte) Block {
	timestamp := time.Now().Unix()
	t := encodeUint64(uint64(timestamp))
	hash := sha256.Sum256(bytes.Join([][]byte{t, data, previousBlockHash},
		[]byte{}))
	return Block{
		Timestamp:         timestamp,
		Data:              data,
		PreviousBlockHash: previousBlockHash,
		Hash:              hash[:],
	}
}

func (b Block) Verify() bool {
	hash := sha256.Sum256(bytes.Join([][]byte{
		encodeUint64(uint64(b.Timestamp)), b.Data, b.PreviousBlockHash},
		[]byte{}))
	return bytes.Equal(hash[:], b.Hash)
}

type Blockchain struct {
	blocks []*Block
}

func (b *Blockchain) Append(data []byte) {
	var previousBlockHash []byte
	if len(b.blocks) == 0 {
		// Genesis
		previousBlockHash = Empty[:]
	} else {
		previousBlockHash = b.blocks[len(b.blocks)-1].Hash
	}
	blk := NewBlock(data, previousBlockHash)
	b.blocks = append(b.blocks, &blk)
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

func NewBlockChain(data []byte) *Blockchain {
	blk := NewBlock(data, Empty[:])
	return &Blockchain{
		blocks: []*Block{&blk},
	}
}
