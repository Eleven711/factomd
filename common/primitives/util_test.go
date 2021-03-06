// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package primitives_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/FactomProject/factomd/common/primitives"
	"github.com/FactomProject/factomd/testHelper"
)

func TestConversions(test *testing.T) {
	v, err := ConvertFixedPoint(".999")
	if err != nil || v != "99900000" {
		fmt.Println("1", v, err)
		test.Fail()
	}
	v, err = ConvertFixedPoint("0.999")
	if err != nil || v != "99900000" {
		fmt.Println("2", v, err)
		test.Fail()
	}
	v, err = ConvertFixedPoint("10.999")
	if err != nil || v != "1099900000" {
		fmt.Println("3", v, err)
		test.Fail()
	}
	v, err = ConvertFixedPoint(".99999999999999")
	if err != nil || v != "99999999" {
		fmt.Println("4", v, err)
		test.Fail()
	}
}

func TestWriteNumber(t *testing.T) {
	out := new(Buffer)

	WriteNumber8(out, 0x01)
	WriteNumber16(out, 0x0203)
	WriteNumber32(out, 0x04050607)
	WriteNumber64(out, 0x0809101112131415)

	answer := "010203040506070809101112131415"
	if out.String() != answer {
		t.Errorf("Failed WriteNumbers. Expected %v, got %v", out.String())
	}
}

func TestConvertion(t *testing.T) {
	var num uint64 = 123456789
	if ConvertDecimalToString(num) != "1.23456789" {
		t.Error("Failed ConvertDecimalToString")
	}
	if ConvertDecimalToPaddedString(num) != "            1.23456789" {
		t.Errorf("Failed ConvertDecimalToPaddedString - '%v'", ConvertDecimalToPaddedString(num))
	}
}

// func DecodeVarInt(data []byte)                   (uint64, []byte)
// func EncodeVarInt(out *bytes.Buffer, v uint64)   error

func TestVariable_Integers(test *testing.T) {
	for i := 0; i < 1000; i++ {
		var out Buffer

		v := make([]uint64, 10)

		for j := 0; j < len(v); j++ {
			var m uint64           // 64 bit mask
			sw := rand.Int63() % 4 // Pick a random choice
			switch sw {
			case 0:
				m = 0xFF // Random byte
			case 1:
				m = 0xFFFF // Random 16 bit integer
			case 2:
				m = 0xFFFFFFFF // Random 32 bit integer
			case 3:
				m = 0xFFFFFFFFFFFFFFFF // Random 64 bit integer
			}
			n := uint64(rand.Int63() + (rand.Int63() << 32))
			v[j] = n & m
		}

		for j := 0; j < len(v); j++ { // Encode our entire array of numbers
			err := EncodeVarInt(&out, v[j])
			if err != nil {
				fmt.Println(err)
				test.Fail()
				return
			}
			//              fmt.Printf("%x ",v[j])
		}
		//          fmt.Println( "Length: ",out.Len())

		data := out.Bytes()

		//          PrtData(data)
		//          fmt.Println()
		sdata := data // Decode our entire array of numbers, and
		var dv uint64 // check we got them back correctly.
		for k := 0; k < 1000; k++ {
			data = sdata
			for j := 0; j < len(v); j++ {
				dv, data = DecodeVarInt(data)
				if dv != v[j] {
					fmt.Printf("Values don't match: decode:%x expected:%x (%d)\n", dv, v[j], j)
					test.Fail()
					return
				}
			}
		}
	}
}

func TestValidateUserStr(t *testing.T) {
	fctAdd := "FA2jK2HcLnRdS94dEcU27rF3meoJfpUcZPSinpb7AwQvPRY6RL1Q"
	fctAddSecret := "Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"
	ecAdd := "EC2DKSYyRcNWf7RS963VFYgMExoHRYLHVeCfQ9PGPmNzwrcmgm2r"

	ok := ValidateFUserStr(fctAdd)
	if ok == false {
		t.Errorf("Valid address not validating - %v", fctAdd)
	}

	ok = ValidateECUserStr(ecAdd)
	if ok == false {
		t.Errorf("Valid address not validating - %v", fctAdd)
	}

	ok = ValidateFPrivateUserStr(fctAddSecret)
	if ok == false {
		t.Errorf("Valid address not validating - %v", fctAdd)
	}

	factoidAddresses := []string{}
	//ecAddresses:=[]string{}

	max := 1000

	for i := 0; i < max; i++ {
		_, _, add := testHelper.NewFactoidAddressStrings(uint64(i))
		factoidAddresses = append(factoidAddresses, add)

		//ecAddresses = append(ecAddresses, add)
	}

	for _, v := range factoidAddresses {
		ok := ValidateFUserStr(v)
		if ok == false {
			t.Errorf("Valid address not validating - %v", v)
		}
	}
}
