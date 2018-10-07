package main

import (
	"crypto/sha256"
	"testing"
)

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
	aaa, err := NewAddress("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Version   : %v", aaa.Version)
	t.Logf("PubKeyHash: %x", aaa.PubKeyHash)
	t.Logf("Checksum  : %x", aaa.Checksum)
	t.Logf("Address  : %v", aaa)

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
