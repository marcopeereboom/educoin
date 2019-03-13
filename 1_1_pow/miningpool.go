package main

import (
	"fmt"
	"sync"
)

// MiningPool is the context that encapsulates the blockchain.
type MiningPool struct {
	sync.Mutex // write mutex to synchronize MiningPool access

	at        uint64 // current nonce
	increment uint64 // nonce increment

	blockchain *Blockchain
}

// getWork returns a mining range and a block to mine.
func (p *MiningPool) GetWork(minerID int) (uint64, uint64, *Block) {
	p.Lock()
	defer p.Unlock()

	// Update where we are at
	start := p.at
	p.at += p.increment

	txpool := fmt.Sprintf("Send 1 Decred to miner %v", minerID)
	blk := p.blockchain.PrepareBlock([]byte(txpool))

	return start, p.at, blk
}

// commitWork attempts to commit a block to the pool.
func (p *MiningPool) CommitWork(blk *Block) error {
	// Verify block
	if !blk.Verify() {
		return fmt.Errorf("invalid block")
	}

	// Add block
	p.Lock()
	defer p.Unlock()
	err := p.blockchain.Append(blk)
	if err != nil {
		return err
	}

	return nil
}

// NewMiningPool returns a miningpool context.
func NewMiningPool(increment uint64) (*MiningPool, error) {
	b, err := NewBlockChain([]byte("Decred is money!"))
	if err != nil {
		return nil, err
	}

	return &MiningPool{
		blockchain: b,
		increment:  increment,
	}, nil
}
