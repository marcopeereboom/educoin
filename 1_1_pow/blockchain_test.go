package main

import (
	"strings"
	"sync"
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

func TestMiningPool(t *testing.T) {
	// May have to play with the increment value on a fast machine.
	mp, err := NewMiningPool(100000)
	if err != nil {
		t.Fatal(err)
	}

	// Start racing miners.
	var wg sync.WaitGroup
	maxWorkers := 10 // increse for more racing
	for x := 0; x < maxWorkers; x++ {
		minerID := x
		wg.Add(1)
		go func() {
			defer wg.Done()
			start, end, blk := mp.GetWork(minerID) // Obtain work
			err := blk.Mine(Difficulty, start, end)
			if err != nil {
				t.Logf("%v %v %v: %v", minerID, start, end, err)
				return
			}

			err = mp.CommitWork(blk) // Send to pool
			if err != nil {
				t.Logf("commit: %v %v %v: %v",
					minerID, start, end, err)
				return
			}
			t.Logf("%v %v %v: nonce %v", minerID, start, end,
				blk.Nonce)
		}()
	}

	wg.Wait()

	// Dump blockchain
	for i := 0; i < mp.blockchain.Len(); i++ {
		t.Log(strings.Repeat("=", 80))
		blk, err := mp.blockchain.Block(i)
		if err != nil {
			t.Fatal(err)
		}
		blk.dump(t)
	}
}
