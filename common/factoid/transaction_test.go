// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package factoid_test

import (
	"fmt"
	"github.com/FactomProject/ed25519"
	"github.com/FactomProject/factomd/common/constants"
	. "github.com/FactomProject/factomd/common/factoid"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	"math/rand"
	"testing"
)

// Random first "address".  It isn't a real one, but one we are using for now.
var adr1 = [constants.ADDRESS_LENGTH]byte{
	0x61, 0xe3, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72, 0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
	0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c, 0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
}

type zeroReader struct{}

var zero zeroReader

var r *rand.Rand

func (zeroReader) Read(buf []byte) (int, error) {
	//if r==nil { r = rand.New(rand.NewSource(time.Now().Unix())) }
	if r == nil {
		r = rand.New(rand.NewSource(1))
	}
	for i := range buf {
		buf[i] = byte(r.Int())
	}
	return len(buf), nil
}

func nextAddress() interfaces.IAddress {

	public, _, _ := ed25519.GenerateKey(zero)

	addr := new(Address)
	addr.SetBytes(public[:])
	return addr
}

func nextSig() []byte {
	// Get me a private key.
	public, _, _ := ed25519.GenerateKey(zero)

	return public[:]
}

func nextAuth2() interfaces.IRCD {
	if r == nil {
		r = rand.New(rand.NewSource(1))
	}
	n := r.Int()%4 + 1
	m := r.Int()%4 + n
	addresses := make([]interfaces.IAddress, m, m)
	for j := 0; j < m; j++ {
		addresses[j] = nextAddress()
	}

	rcd, _ := NewRCD_2(n, m, addresses)
	return rcd
}

var nb interfaces.IBlock

func getSignedTrans() interfaces.IBlock {

	if nb != nil {
		return nb
	}

	nb = new(Transaction)
	t := nb.(*Transaction)

	for i := 0; i < 5; i++ {
		t.AddInput(nextAddress(), uint64(rand.Int63n(10000000000)))
	}

	for i := 0; i < 3; i++ {
		t.AddOutput(nextAddress(), uint64(rand.Int63n(10000000000)))
	}

	for i := 0; i < 3; i++ {
		t.AddECOutput(nextAddress(), uint64(rand.Int63n(10000000)))
	}

	for i := 0; i < 3; i++ {
		sig := NewRCD_1(nextSig())
		t.AddAuthorization(sig)
	}

	for i := 0; i < 2; i++ {

		t.AddAuthorization(nextAuth2())
	}

	return nb
}

// This test prints bunches of stuff that must be visually checked.
// Mostly we keep it commented out.
func TestTransaction(test *testing.T) {
	nb = getSignedTrans()
	bytes, _ := nb.CustomMarshalText()
	fmt.Printf("Transaction:\n%slen: %d\n", string(bytes), len(bytes))
	fmt.Println("\n---------------------------------------------------------------------")
}

func Test_Address_MarshalUnMarshal(test *testing.T) {
	a := nextAddress()
	adr, err := a.MarshalBinary()
	if err != nil {
		primitives.Prtln(err)
		test.Fail()
	}
	_, err = a.UnmarshalBinaryData(adr)
	if err != nil {
		primitives.Prtln(err)
		test.Fail()
	}
}

func Test_Multisig_MarshalUnMarshal(test *testing.T) {
	rcd := nextAuth2()
	auth2, err := rcd.MarshalBinary()
	if err != nil {
		primitives.Prtln(err)
		test.Fail()
	}

	_, err = rcd.UnmarshalBinaryData(auth2)

	if err != nil {
		primitives.Prtln(err)
		test.Fail()
	}
}

func Test_Transaction_MarshalUnMarshal(test *testing.T) {

	getSignedTrans()                // Make sure we have a signed transaction
	data, err := nb.MarshalBinary() // Marshal our signed transaction
	if err != nil {                 // If we have an error, print our stack
		primitives.Prtln(err) //   and fail our test
		test.Fail()
	}

	xb := new(Transaction)

	err = xb.UnmarshalBinary(data) // Now Unmarshal

	if err != nil {
		primitives.Prtln(err)
		test.Fail()
		return
	}

	//     txt1,_ := xb.CustomMarshalText()
	//     txt2,_ := nb.CustomMarshalText()
	//     primitives.Prtln(string(txt1))
	//     primitives.Prtln(string(txt2))

	if xb.IsEqual(nb) != nil {
		primitives.Prtln("Trans\n", nb, "Unmarshal Trans\n", xb)
		test.Fail()
	}

}

func Test_ValidateAmounts(test *testing.T) {
	var zero uint64
	_, err := ValidateAmounts(zero - 1)
	if err != nil {
		test.Failed()
	}
	_, err = ValidateAmounts(1, 2, 3, 4, 5, zero-1)
	if err != nil {
		test.Failed()
	}
	_, err = ValidateAmounts(0x6FFFFFFFFFFFFFFF, 1)
	if err != nil {
		test.Failed()
	}
	_, err = ValidateAmounts(1, 0x6FFFFFFFFFFFFFFF, 1)
	if err != nil {
		test.Failed()
	}
}
