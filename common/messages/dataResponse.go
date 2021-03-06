// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/entryBlock"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

// Communicate a Directory Block State

type DataResponse struct {
	MessageBase
	Timestamp interfaces.Timestamp

	DataType   int // 0 = Entry, 1 = EntryBlock
	DataHash   interfaces.IHash
	DataObject interfaces.BinaryMarshallable //Entry or EntryBlock

	//Not signed!
}

var _ interfaces.IMsg = (*DataResponse)(nil)

func (a *DataResponse) IsSameAs(b *DataResponse) bool {
	if b == nil {
		return false
	}
	if a.Timestamp != b.Timestamp {
		return false
	}
	if a.DataType != b.DataType {
		return false
	}

	if a.DataHash == nil && b.DataHash != nil {
		return false
	}
	if a.DataHash != nil {
		if a.DataHash.IsSameAs(b.DataHash) == false {
			return false
		}
	}

	if a.DataObject == nil && b.DataObject != nil {
		return false
	}
	if a.DataObject != nil {
		hex1, err := a.DataObject.MarshalBinary()
		if err != nil {
			return false
		}
		hex2, err := b.DataObject.MarshalBinary()
		if err != nil {
			return false
		}
		if primitives.AreBytesEqual(hex1, hex2) == false {
			return false
		}
	}
	return true
}

func (m *DataResponse) GetHash() interfaces.IHash {
	return m.GetMsgHash()
}

func (m *DataResponse) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *DataResponse) Type() byte {
	return constants.DATA_RESPONSE
}

func (m *DataResponse) Int() int {
	return -1
}

func (m *DataResponse) Bytes() []byte {
	return nil
}

func (m *DataResponse) GetTimestamp() interfaces.Timestamp {
	return m.Timestamp
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *DataResponse) Validate(state interfaces.IState) int {
	var dataHash interfaces.IHash
	var err error
	switch m.DataType {
	case 0: // DataType = entry
		dataObject, ok := m.DataObject.(interfaces.IEBEntry)
		if !ok {
			return -1
		}
		dataHash = dataObject.GetHash()
	case 1: // DataType = eblock
		dataObject, ok := m.DataObject.(interfaces.IEntryBlock)
		if !ok {
			return -1
		}
		dataHash, err = dataObject.KeyMR()
		if err != nil {
			return -1
		}
	default:
		// DataType currently not supported, treat as invalid
		return -1
	}

	if dataHash.IsSameAs(m.DataHash) {
		return 1
	}

	return -1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *DataResponse) Leader(state interfaces.IState) bool {
	return false
}

// Execute the leader functions of the given message
func (m *DataResponse) LeaderExecute(state interfaces.IState) error {
	return fmt.Errorf("Should never execute a DataResponse in the Leader")
}

// Returns true if this is a message for this server to execute as a follower
func (m *DataResponse) Follower(interfaces.IState) bool {
	return true
}

func (m *DataResponse) FollowerExecute(state interfaces.IState) error {
	if state.HasDataRequest(m.DataHash) {
		switch m.DataType {
		case 1: // Data is an entryBlock
			eblock, ok := m.DataObject.(interfaces.IEntryBlock)
			if !ok {
				return fmt.Errorf("Wrong DataType -- not IEntryBlock")
			}
			ebKeyMR, err := eblock.KeyMR()

			if err == nil {
				if ebKeyMR.IsSameAs(m.DataHash) {
					if !state.DatabaseContains(ebKeyMR) {
						err := state.FollowerExecuteAddData(m) // Save EBlock

						for _, hashMatchAttempt := range state.GetDirectoryBlockByHeight(state.GetEBDBHeightComplete()).GetEntryHashes() {
							if hashMatchAttempt.IsSameAs(ebKeyMR) {
								if state.GetAllEntries(ebKeyMR) {
									state.SetEBDBHeightComplete(state.GetEBDBHeightComplete() + 1)
								}
							}
						}

						if err != nil { // If there was an error saving the data, return err
							return err
						}
					}
				}
			}
		case 0: // Data is an entry
			if !state.DatabaseContains(m.DataHash) {
				entry, ok := m.DataObject.(interfaces.IEBEntry)
				if !ok {
					return fmt.Errorf("Wrong DataType -- not IEBEntry")
				}
				err := state.FollowerExecuteAddData(m) // Save entry

				ebKeyMR := state.GetEBlockKeyMRFromEntryHash(entry.GetHash())

				if ebKeyMR != nil {
					if state.DatabaseContains(ebKeyMR) { // Node already has eBlock in database
						if state.GetAllEntries(ebKeyMR) {
							state.SetEBDBHeightComplete(state.GetEBDBHeightComplete() + 1)
						}
					} else {
						if !state.HasDataRequest(ebKeyMR) {
							// Need to get eblock itself
							eBlockRequest := NewMissingData(state, ebKeyMR)
							state.NetworkOutMsgQueue() <- eBlockRequest
						}
					}
				}

				if err != nil { // If there was an error saving the data, return err
					return err
				}

			}
		}
	}
	return nil
}

// Acknowledgements do not go into the process list.
func (e *DataResponse) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Should never have its Process() method called")
}

func (e *DataResponse) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *DataResponse) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *DataResponse) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (m *DataResponse) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()
	newData = data
	if newData[0] != m.Type() {
		return nil, fmt.Errorf("Invalid Message type")
	}
	newData = newData[1:]

	newData, err = m.Timestamp.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	m.DataType = int(newData[0])
	newData = newData[1:]

	m.DataHash = primitives.NewHash(constants.ZERO_HASH)
	newData, err = m.DataHash.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}
	switch m.DataType {
	case 0:
		entryAttempt, err := attemptEntryUnmarshal(newData)
		if err != nil {
			return nil, err
		} else {
			m.DataObject = entryAttempt
		}
	case 1:
		eblockAttempt, err := attemptEBlockUnmarshal(newData)
		if err != nil {
			return nil, err
		} else {
			m.DataObject = eblockAttempt
		}
	default:
		return nil, fmt.Errorf("DataResponse's DataType not supported for unmarshalling yet")
	}

	m.Peer2Peer = true // Always a peer2peer request.
	return data, nil
}

func (m *DataResponse) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func attemptEntryUnmarshal(data []byte) (entry interfaces.IEBEntry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Bytes do not represent an entry")
		}
	}()

	entry, err = entryBlock.UnmarshalEntry(data)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func attemptEBlockUnmarshal(data []byte) (eblock interfaces.IEntryBlock, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Bytes do not represent an eblock: %v\n", r)
		}
	}()

	eblock, err = entryBlock.UnmarshalEBlock(data)
	if err != nil {
		return nil, err
	}

	return eblock, nil
}

func (m *DataResponse) MarshalBinary() ([]byte, error) {
	var buf primitives.Buffer
	buf.Write([]byte{m.Type()})
	if d, err := m.Timestamp.MarshalBinary(); err != nil {
		return nil, err
	} else {
		buf.Write(d)
	}

	binary.Write(&buf, binary.BigEndian, uint8(m.DataType))

	if d, err := m.DataHash.MarshalBinary(); err != nil {
		return nil, err
	} else {
		buf.Write(d)
	}

	if m.DataObject != nil {
		d, err := m.DataObject.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(d)
	}

	return buf.DeepCopyBytes(), nil
}

func (m *DataResponse) String() string {
	return fmt.Sprintf("DataResponse Type: %v\n Hash: %x\n Object: %v\n",
		m.DataType,
		m.DataHash.Bytes()[:5],
		m.DataObject)
}

func NewDataResponse(state interfaces.IState, dataObject interfaces.BinaryMarshallable,
	dataType int,
	dataHash interfaces.IHash) interfaces.IMsg {

	msg := new(DataResponse)

	msg.Peer2Peer = true
	msg.Timestamp = state.GetTimestamp()

	msg.DataHash = dataHash
	msg.DataType = dataType
	msg.DataObject = dataObject

	//fmt.Println("DATARESPONSE: ", msg.DataObject)

	return msg
}
