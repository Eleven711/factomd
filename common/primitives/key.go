// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/interfaces"
)

// PrivateKey contains Public/Private key pair
type PrivateKey struct {
	Key *[ed25519.PrivateKeySize]byte
	Pub *PublicKey
}

var _ interfaces.Signer = (*PrivateKey)(nil)

func (pk *PrivateKey) CustomMarshalText2(string) ([]byte, error) {
	return ([]byte)(hex.EncodeToString(pk.Key[:]) + pk.Pub.String()), nil
}

func (pk *PrivateKey) Public() []byte {
	return pk.Pub[:]
}

func (pk *PrivateKey) AllocateNew() {
	pk.Key = new([ed25519.PrivateKeySize]byte)
	pk.Pub = new(PublicKey)
}

// Create a new private key from a hex string
func NewPrivateKeyFromHex(s string) (pk PrivateKey, err error) {
	privKeybytes, err := hex.DecodeString(s)
	if err != nil {
		return
	}
	if privKeybytes == nil {
		return pk, errors.New("Invalid private key input string!")
	}
	if len(privKeybytes) == ed25519.PrivateKeySize-ed25519.PublicKeySize {
		_, privKeybytes, err = GenerateKeyFromPrivateKey(privKeybytes)
		if err != nil {
			return
		}
	}
	if len(privKeybytes) != ed25519.PrivateKeySize {
		return pk, errors.New("Invalid private key input string!")
	}
	pk.AllocateNew()
	copy(pk.Key[:], privKeybytes)
	err = pk.Pub.UnmarshalBinary(privKeybytes[len(privKeybytes)-ed25519.PublicKeySize:])
	return
}

func NewPrivateKeyFromHexBytes(privKeybytes []byte) *PrivateKey {
	pk := new(PrivateKey)
	pk.AllocateNew()
	copy(pk.Key[:], privKeybytes)
	pk.Pub.UnmarshalBinary(privKeybytes)
	return pk
}

// Sign signs msg with PrivateKey and return Signature
func (pk *PrivateKey) Sign(msg []byte) (sig interfaces.IFullSignature) {
	sig = new(Signature)
	sig.SetPub(pk.Pub[:])
	s := ed25519.Sign(pk.Key, msg)
	sig.SetSignature(s[:])
	return
}

// Sign signs msg with PrivateKey and return Signature
func (pk *PrivateKey) MarshalSign(msg interfaces.BinaryMarshallable) (sig interfaces.IFullSignature) {
	data, _ := msg.MarshalBinary()
	return pk.Sign(data)
}

//Generate creates new PrivateKey / PublciKey pair or returns error
func (pk *PrivateKey) GenerateKey() error {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	pk.Key = priv
	pk.Pub = new(PublicKey)
	err = pk.Pub.UnmarshalBinary(pub[:])
	return err
}

// Returns hex-encoded string of first 32 bytes of key (private key portion)
func (pk *PrivateKey) PrivateKeyString() string {
	return hex.EncodeToString(pk.Key[:32])
}

/******************PublicKey*******************************/

// PublicKey contains only Public part of Public/Private key pair
type PublicKey [ed25519.PublicKeySize]byte

var _ interfaces.Verifier = (*PublicKey)(nil)

func (a *PublicKey) IsSameAs(b *PublicKey) bool {
	if b == nil {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (pk *PublicKey) MarshalText() ([]byte, error) {
	return []byte(pk.String()), nil
}

func (pk *PublicKey) UnmarshalText(b []byte) error {
	p, err := hex.DecodeString(string(b))
	if err != nil {
		return err
	}
	copy(pk[:], p)
	return nil
}

func (pk *PublicKey) String() string {
	return hex.EncodeToString(pk[:])
}

func PubKeyFromString(instr string) (pk PublicKey) {
	p, _ := hex.DecodeString(instr)
	copy(pk[:], p)
	return
}

func (k *PublicKey) Verify(msg []byte, sig *[ed25519.SignatureSize]byte) bool {
	return ed25519.VerifyCanonical((*[32]byte)(k), msg, sig)
}

func (k *PublicKey) MarshalBinary() ([]byte, error) {
	var buf Buffer
	buf.Write(k[:])
	return buf.DeepCopyBytes(), nil
}

func (k *PublicKey) UnmarshalBinaryData(p []byte) ([]byte, error) {
	if len(p) < ed25519.PublicKeySize {
		return nil, fmt.Errorf("Invalid data passed")
	}
	copy(k[:], p)
	return p[ed25519.PublicKeySize:], nil
}

func (k *PublicKey) UnmarshalBinary(p []byte) (err error) {
	_, err = k.UnmarshalBinaryData(p)
	return
}

// Verify returns true iff sig is a valid signature of message by publicKey.
func Verify(publicKey *[ed25519.PublicKeySize]byte, message []byte, sig *[ed25519.SignatureSize]byte) bool {
	return ed25519.VerifyCanonical(publicKey, message, sig)
}

// Verify returns true iff sig is a valid signature of message by publicKey.
func VerifySlice(p []byte, message []byte, s []byte) bool {
	sig := new([ed25519.PrivateKeySize]byte)
	pub := new([ed25519.PublicKeySize]byte)
	copy(sig[:], s)
	copy(pub[:], p)
	return Verify(pub, message, sig)
}
