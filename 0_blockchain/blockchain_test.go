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
}

func newBlockChain() *Blockchain {
	b := NewBlockChain([]byte("Decred is money!"))
	b.Append([]byte("Send 1 Decred to Alice"))
	b.Append([]byte("Send 2 Decred to Bob"))
	return b
}

func (b Blockchain) dumpBlockchain(t *testing.T) error {
	for i := 0; i < b.Len(); i++ {
		t.Logf("%v", strings.Repeat("=", 80))
		t.Logf("Height           : %v", i)
		block, err := b.Block(i)
		if err != nil {
			return err
		}
		block.dump(t)
		t.Logf("Block valid      : %v", block.Verify())
		if !block.Verify() {
			return fmt.Errorf("block corrupt: %v", i)
		}
	}
	return nil
}

func (b *Blockchain) corrupt(block int, data []byte) error {
	if block > len(b.blocks) {
		return fmt.Errorf("invalid block: %v", block)
	}
	b.blocks[block].Data = data
	return nil
}

func TestSuccess(t *testing.T) {
	b := newBlockChain()
	if err := b.dumpBlockchain(t); err != nil {
		t.Fatal(err)
	}
}

func TestFailure(t *testing.T) {
	b := newBlockChain()
	if err := b.corrupt(1, []byte("Send 2 Decred to Alice")); err != nil {
		t.Fatal(err)
	}
	if err := b.dumpBlockchain(t); err == nil {
		t.Fatalf("Unexpected success")
	}
}
