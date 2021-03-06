// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages_test

import (
	"testing"

	"github.com/FactomProject/factomd/common/constants"
	. "github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
)

func TestMarshalUnmarshalMissingData(t *testing.T) {
	msg := newMissingData()

	hex, err := msg.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Marshalled - %x", hex)

	msg2, err := UnmarshalMessage(hex)
	if err != nil {
		t.Error(err)
	}
	str := msg2.String()
	t.Logf("str - %v", str)

	if msg2.Type() != constants.MISSING_DATA {
		t.Error("Invalid message type unmarshalled")
	}

	hex2, err := msg2.(*MissingData).MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if len(hex) != len(hex2) {
		t.Error("Hexes aren't of identical length")
	}
	for i := range hex {
		if hex[i] != hex2[i] {
			t.Error("Hexes do not match")
		}
	}

	if msg.IsSameAs(msg2.(*MissingData)) != true {
		t.Errorf("MissingData messages are not identical")
	}
}

func newMissingData() *MissingData {
	msg := new(MissingData)
	msg.Timestamp.SetTimeNow()

	h, err := primitives.NewShaHashFromStr("deadbeef00000000000000000000000000000000000000000000000000000000")
	if err != nil {
		panic(err)
	}
	msg.RequestHash = h

	return msg
}
