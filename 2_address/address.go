package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/btcsuite/golangcrypto/ripemd160"
)

const AddressVersion = 0 // Version of Address structure

// PrivateKey represent an ECDSA private key.
type PrivateKey struct {
	ecdsa.PrivateKey
}

// NewKey creates a new private key.
func NewKey() (*PrivateKey, error) {
	p, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{*p}, nil
}

// Public returns the corresponding public key.
func (p PrivateKey) Public() []byte {
	return append(p.PublicKey.X.Bytes(), p.PublicKey.Y.Bytes()...)
}

// Sign returns the signature of blob.
func (p PrivateKey) Sign(blob []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, &p.PrivateKey, blob)
	if err != nil {
		return nil, err
	}
	return append(r.Bytes(), s.Bytes()...), nil
}

// PublicKey represents an ECDSA public key.
type PublicKey struct {
	ecdsa.PublicKey
}

// NewPublicKey unpacks pub and creates a corresponding ECDSA public key.
func NewPublicKey(pub []byte) *PublicKey {
	l := len(pub) / 2
	x := new(big.Int).SetBytes(pub[:l])
	y := new(big.Int).SetBytes(pub[l:])
	return &PublicKey{ecdsa.PublicKey{elliptic.P256(), x, y}}
}

// Verify unpacks signature and verifies the integrity of blob.
func (p PublicKey) Verify(blob, signature []byte) bool {
	l := len(signature) / 2
	r := new(big.Int).SetBytes(signature[:l])
	s := new(big.Int).SetBytes(signature[l:])
	return ecdsa.Verify(&p.PublicKey, blob, r, s)
}

// Key return the []byte representation of an ECDSA public key.
func (p PublicKey) Key() []byte {
	return append(p.X.Bytes(), p.Y.Bytes()...)
}

// Address represents all constituent pieces of an address.
type Address struct {
	Version    byte   // Version of the address
	PubKeyHash []byte // Hash of the public key ripemd160(sha256(pk))
	Checksum   []byte // Checksum of PubKeyHash sha256(sha256(pkh))
}

// checksum calculates the checksum of blob by taking the first 4 bytes from
// the double sha256 of blob.  The checksum uses a double sha256 in order to
// prevent length-extension attacks.
func checksum(blob []byte) []byte {
	chk0 := sha256.Sum256(blob)
	chk1 := sha256.Sum256(chk0[:])
	return chk1[0:4]
}

// ripemd160Sum returns the ripemd160 hash of blob.
func ripemd160Sum(blob []byte) []byte {
	r160 := ripemd160.New()
	_, err := r160.Write(blob)
	if err != nil {
		panic(err)
	}
	return r160.Sum(nil)
}

// Address creates an Address structure from a PublicKey.
func (p PublicKey) Address() *Address {
	pksha := sha256.Sum256(p.Key())  // sha256(public key)
	pkhash := ripemd160Sum(pksha[:]) // ripemd160(sha256(public key))
	return &Address{
		Version:    AddressVersion,
		PubKeyHash: pkhash,
		Checksum:   checksum(append([]byte{AddressVersion}, pkhash...)),
	}
}

// String returns the human readable form of an Address. The process is
// base58(Version+PubKeyHash+Checksum).
func (a Address) String() string {
	addr := append([]byte{a.Version}, a.PubKeyHash...)
	return Encode(append(addr, a.Checksum...))
}

// NewAddress decodes a human readable address into an Address structure.
// It recreates the address structure decoding base58 of the provided address
// which results in the following byte array [version][pub key hash][checksum]
func NewAddress(a string) (*Address, error) {
	da := Decode(a)
	l := len(da)
	if l-4 <= 0 {
		return nil, fmt.Errorf("invalid length")
	}
	if da[0] != AddressVersion {
		return nil, fmt.Errorf("invalid address version")
	}
	addr := Address{
		Version:    da[0],
		PubKeyHash: da[1 : l-4],
		Checksum:   da[l-4 : l],
	}
	if !bytes.Equal(checksum(da[0:l-4]), addr.Checksum) {
		return nil, fmt.Errorf("invalid checksum")
	}
	return &addr, nil
}
