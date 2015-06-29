// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"github.com/FactomProject/FactomCode/common"
	"io"
	"encoding/binary"	
)

// Acknowledgement Type
const (
	ACK_FACTOID_TX uint8 = iota
	END_MINUTE_1
	END_MINUTE_2
	END_MINUTE_3
	END_MINUTE_4
	END_MINUTE_5
	END_MINUTE_6
	END_MINUTE_7
	END_MINUTE_8
	END_MINUTE_9
	END_MINUTE_10
	ACK_REVEAL_ENTRY
	ACK_COMMIT_CHAIN
	ACK_REVEAL_CHAIN
	ACK_COMMIT_ENTRY

	FORCE_FACTOID_GENESIS_REBUILD
	//	FORCE_FACTOID_VALIDATION // at a specific block height; everything higher will be discarded by btcd-side (possibly creating orphaned blocks)
	INFO_CURRENT_HEIGHT // info message to the wire-side to indicate the current known block height; a duplicate of FORCE_FACTOID_VALIDATION (???)
)

type MsgAcknowledgement struct {
	Height      uint32
	ChainID     *common.Hash
	Index       uint32
	Type        byte
	Affirmation *ShaHash // affirmation value -- hash of the message/object in question
	SerialHash  [32]byte
	Signature   [64]byte
}

// Write out the MsgAcknowledgement (excluding Signature) to binary.
func (msg *MsgAcknowledgement) GetBinaryForSignature() (data []byte, err error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, msg.Height)
	
	data, err = msg.ChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	
	binary.Write(&buf, binary.BigEndian, msg.Index)
	
	buf.Write([]byte{msg.Type})

	buf.Write(msg.Affirmation.Bytes())

	buf.Write(msg.SerialHash[:])

	return buf.Bytes(), err
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgAcknowledgement) BtcDecode(r io.Reader, pver uint32) error {
	err := readElements(r, &msg.Height, &msg.ChainID, &msg.Index, &msg.Affirmation, &msg.SerialHash, &msg.Signature)
	if err != nil {
		return err
	}

	return nil
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgAcknowledgement) BtcEncode(w io.Writer, pver uint32) error {
	err := writeElements(w, &msg.Height, &msg.ChainID, &msg.Index, &msg.Affirmation, &msg.SerialHash, &msg.Signature)
	if err != nil {
		return err
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgAcknowledgement) Command() string {
	return CmdAcknowledgement
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgAcknowledgement) MaxPayloadLength(pver uint32) uint32 {

	// 10K is too big of course, TODO: adjust
	return MaxAppMsgPayload
}

// NewMsgAcknowledgement returns a new bitcoin ping message that conforms to the Message
// interface.  See MsgAcknowledgement for details.
func NewMsgAcknowledgement(height uint32, index uint32, affirm *ShaHash, ackType byte) *MsgAcknowledgement {
	
	if affirm == nil {
		affirm = new (ShaHash)
	}
	return &MsgAcknowledgement{
		Height:      height,
		Index:       index,
		Affirmation: affirm,
		Type:        ackType,
	}
}

// Create a sha hash from the message binary (output of BtcEncode)
func (msg *MsgAcknowledgement) Sha() (ShaHash, error) {

	buf := bytes.NewBuffer(nil)
	msg.BtcEncode(buf, ProtocolVersion)
	var sha ShaHash
	_ = sha.SetBytes(Sha256(buf.Bytes()))

	return sha, nil
}

