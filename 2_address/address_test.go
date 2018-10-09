package main

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestPK(t *testing.T) {
	key, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}

	hash := sha256.Sum256([]byte("Hello world!"))

	signature, err := key.Sign(hash[:])
	if err != nil {
		t.Fatal(err)
	}

	pk := NewPublicKey(key.Public())
	if !pk.Verify(hash[:], signature) {
		t.Fatalf("verify failed")
	}

	// Corrupt hash
	hash = sha256.Sum256([]byte("hello world!"))
	if pk.Verify(hash[:], signature) {
		t.Fatalf("verify succeeded")
	}
}

func TestAddress(t *testing.T) {
	key, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}
	pk := NewPublicKey(key.Public())
	a := pk.Address()
	t.Logf("Version   : %v", a.Version)
	t.Logf("PubKeyHash: %x", a.PubKeyHash)
	t.Logf("Checksum  : %x", a.Checksum)
	t.Logf("Address  : %v", a)

	aa, err := NewAddress(a.String())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Version   : %v", aa.Version)
	t.Logf("PubKeyHash: %x", aa.PubKeyHash)
	t.Logf("Checksum  : %x", aa.Checksum)
	t.Logf("Address  : %v", aa)

	if aa.String() != a.String() {
		t.Fatalf("stringers don't match")
	}
}

func TestAddressCorrupt(t *testing.T) {
	key, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}
	pk := NewPublicKey(key.Public())
	a := pk.Address()
	aa, err := NewAddress(a.String())
	if err != nil {
		t.Fatal(err)
	}

	assertPanic(t, func() {
		// Corrupt Address
		aa.PubKeyHash[0] = ^aa.PubKeyHash[0]
		if aa.String() == a.String() {
			t.Fatalf("stringers match")
		}
	})

	// Fix checksum and try again
	aa.Checksum = checksum(append([]byte{aa.Version}, aa.PubKeyHash...))
	if bytes.Equal(aa.Checksum, a.Checksum) {
		t.Fatal("checksums should not be equal")
	}

	if aa.String() == a.String() {
		t.Fatalf("stringers match")
	}
}
