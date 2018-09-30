package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func (b Block) dump(t *testing.T) {
	t.Logf("Timestamp        : %v\n", time.Unix(b.Timestamp, 0))
	t.Logf("PreviousBlockHash: %x\n", b.PreviousBlockHash)
	t.Logf("Hash             : %x\n", b.Hash)
	t.Logf("Data             : %s\n", string(b.Data))
	t.Logf("Nonce            : %v\n", b.Nonce)
}

func (b *Blockchain) corrupt(block int, data []byte) error {
	if block > len(b.blocks) {
		return fmt.Errorf("invalid block: %v", block)
	}
	b.blocks[block].Data = data
	return nil
}

func TestBlockChain(t *testing.T) {
	b, err := NewBlockChain([]byte("Decred is money!"))
	if err != nil {
		t.Fatal(err)
	}

	blk := b.PrepareBlock([]byte("Send 1 Decred to Alice"))
	err = blk.Mine(Difficulty)
	if err != nil {
		t.Fatal(err)
	}
	err = b.Append(blk)
	if err != nil {
		t.Fatal(err)
	}

	blk = b.PrepareBlock([]byte("Send 2 Decred to Bob"))
	err = blk.Mine(Difficulty)
	if err != nil {
		t.Fatal(err)
	}
	b.Append(blk)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < b.Len(); i++ {
		t.Log(strings.Repeat("=", 80))
		blk, err := b.Block(i)
		if err != nil {
			t.Fatal(err)
		}
		blk.dump(t)
		if !blk.Verify() {
			t.Fatalf("corrupt")
		}
	}

	// Corrupt data and try again
	if err := b.corrupt(1, []byte("Send 2 Decred to Alice")); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < b.Len(); i++ {
		t.Log(strings.Repeat("=", 80))
		blk, err := b.Block(i)
		if err != nil {
			t.Fatal(err)
		}
		blk.dump(t)
		if !blk.Verify() && i == 1 {
			t.Logf("Block 1 corrupt")
			return
		}
	}
}
