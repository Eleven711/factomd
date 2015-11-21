// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package databaseOverlay_test

import (
	"bytes"
	"encoding/gob"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	. "github.com/FactomProject/factomd/database/databaseOverlay"
	"github.com/FactomProject/factomd/database/mapdb"
	"testing"
)

func TestInsertFetch(t *testing.T) {
	dbo := createOverlay()
	defer dbo.Close()
	b := NewDBTestObject()
	b.Data = []byte{0x00, 0x01, 0x02, 0x03}

	err := dbo.Insert(TestBucket, b)
	if err != nil {
		t.Error(err)
	}

	index := b.DatabasePrimaryIndex()

	b2 := NewDBTestObject()
	resp, err := dbo.FetchBlock(TestBucket, index, b2)
	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error("Response is nil while it shouldn't be.")
	}

	bResp := resp.(*DBTestObject)

	bytes1 := b.Data
	bytes2 := bResp.Data

	if primitives.AreBytesEqual(bytes1, bytes2) == false {
		t.Errorf("Bytes are not equal - %x vs %x", bytes1, bytes2)
	}
}

func TestFetchBy(t *testing.T) {
	dbo := createOverlay()
	defer dbo.Close()

	blocks := []*DBTestObject{}

	max := 10
	for i := 0; i < max; i++ {
		b := NewDBTestObject()
		b.ChainID = primitives.NewHash(CopyZeroHash())
		b.Data = []byte{byte(i), byte(i + 1), byte(i + 2), byte(i + 3)}

		primaryIndex := CopyZeroHash()
		primaryIndex[len(primaryIndex)-1] = byte(i)

		secondaryIndex := CopyZeroHash()
		secondaryIndex[0] = byte(i)

		b.PrimaryIndex = primitives.NewHash(primaryIndex)
		b.SecondaryIndex = primitives.NewHash(secondaryIndex)
		b.DatabaseHeight = uint32(i)
		blocks = append(blocks, b)

		err := dbo.ProcessBlockBatch(TestBucket, TestNumberBucket, TestSecondaryIndexBucket, b)
		if err != nil {
			t.Error(err)
		}
	}
	headIndex, err := dbo.FetchHeadIndexByChainID(primitives.NewHash(CopyZeroHash()))
	if err != nil {
		t.Error(err)
	}
	if headIndex.IsSameAs(blocks[max-1].PrimaryIndex) == false {
		t.Error("Wrong chain head")
	}

	head, err := dbo.FetchChainHeadByChainID(TestBucket, primitives.NewHash(CopyZeroHash()), new(DBTestObject))
	if err != nil {
		t.Error(err)
	}
	if blocks[max-1].IsEqual(head.(*DBTestObject)) == false {
		t.Error("Heads are not equal")
	}

	for i := 0; i < max; i++ {
		primaryIndex := CopyZeroHash()
		primaryIndex[len(primaryIndex)-1] = byte(i)

		secondaryIndex := CopyZeroHash()
		secondaryIndex[0] = byte(i)

		dbHeight := uint32(i)

		block, err := dbo.FetchBlockByHeight(TestNumberBucket, TestBucket, dbHeight, new(DBTestObject))
		if err != nil {
			t.Error(err)
		}
		if blocks[i].IsEqual(block.(*DBTestObject)) == false {
			t.Error("Blocks are not equal")
		}

		index, err := dbo.FetchBlockIndexByHeight(TestNumberBucket, dbHeight)
		if err != nil {
			t.Error(err)
		}
		if primitives.AreBytesEqual(index.Bytes(), primaryIndex) == false {
			t.Error("Wrong primary index returned")
		}

		index, err = dbo.FetchPrimaryIndexBySecondaryIndex(TestSecondaryIndexBucket, primitives.NewHash(secondaryIndex))
		if err != nil {
			t.Error(err)
		}
		if primitives.AreBytesEqual(index.Bytes(), primaryIndex) == false {
			t.Error("Wrong primary index returned")
		}

		block, err = dbo.FetchBlockBySecondaryIndex(TestSecondaryIndexBucket, TestBucket, primitives.NewHash(secondaryIndex), new(DBTestObject))
		if err != nil {
			t.Error(err)
		}
		if blocks[i].IsEqual(block.(*DBTestObject)) == false {
			t.Error("Blocks are not equal")
		}
	}

	fetchedBlocks, err := dbo.FetchAllBlocksFromBucket(TestBucket, new(DBTestObject))
	if err != nil {
		t.Error(err)
	}
	if len(fetchedBlocks) != len(blocks) {
		t.Error("Invalid amount of blocks returned")
	}
	for i := 0; i < max; i++ {
		if blocks[i].IsEqual(fetchedBlocks[i].(*DBTestObject)) == false {
			t.Error("Block from batch is not equal")
		}
	}

	startIndex := 3
	indexCount := 4
	fetchedIndexes, err := dbo.FetchBlockIndexesInHeightRange(TestNumberBucket, int64(startIndex), int64(startIndex+indexCount))
	if len(fetchedIndexes) != indexCount {
		t.Error("Invalid amount of indexes returned")
	}
	for i := 0; i < indexCount; i++ {
		primaryIndex := CopyZeroHash()
		primaryIndex[len(primaryIndex)-1] = byte(i + startIndex)

		if primitives.AreBytesEqual(primaryIndex, fetchedIndexes[i].Bytes()) == false {
			t.Error("Index from batch is not equal")
		}
	}
}

func createOverlay() *Overlay {
	return NewOverlay(new(mapdb.MapDB))
}

var TestBucket []byte = []byte{0x01}
var TestNumberBucket []byte = []byte{0x02}
var TestSecondaryIndexBucket []byte = []byte{0x03}

type bareDBTestObject struct {
	Data           []byte
	DatabaseHeight uint32
	PrimaryIndex   interfaces.IHash
	SecondaryIndex interfaces.IHash
	ChainID        interfaces.IHash
}

func NewBareDBTestObject() *bareDBTestObject {
	d := new(bareDBTestObject)
	d.Data = []byte{}
	d.DatabaseHeight = 0
	d.PrimaryIndex = new(primitives.Hash)
	d.SecondaryIndex = new(primitives.Hash)
	d.ChainID = new(primitives.Hash)
	return d
}

type DBTestObject struct {
	Data           []byte
	DatabaseHeight uint32
	PrimaryIndex   interfaces.IHash
	SecondaryIndex interfaces.IHash
	ChainID        interfaces.IHash
}

func NewDBTestObject() *DBTestObject {
	d := new(DBTestObject)
	d.Data = []byte{}
	d.DatabaseHeight = 0
	d.PrimaryIndex = new(primitives.Hash)
	d.SecondaryIndex = new(primitives.Hash)
	d.ChainID = new(primitives.Hash)
	return d
}

var _ interfaces.DatabaseBatchable = (*DBTestObject)(nil)

func (d *DBTestObject) GetDatabaseHeight() uint32 {
	return d.DatabaseHeight
}

func (d *DBTestObject) DatabasePrimaryIndex() interfaces.IHash {
	return d.PrimaryIndex
}

func (d *DBTestObject) DatabaseSecondaryIndex() interfaces.IHash {
	return d.SecondaryIndex
}

func (d *DBTestObject) GetChainID() []byte {
	return d.ChainID.Bytes()
}

func (d *DBTestObject) New() interfaces.BinaryMarshallableAndCopyable {
	return NewDBTestObject()
}

func (d *DBTestObject) UnmarshalBinaryData(data []byte) ([]byte, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(data))

	tmp := NewBareDBTestObject()

	err := dec.Decode(tmp)
	if err != nil {
		return nil, err
	}

	d.Data = tmp.Data
	d.DatabaseHeight = tmp.DatabaseHeight
	d.PrimaryIndex = tmp.PrimaryIndex
	d.SecondaryIndex = tmp.SecondaryIndex
	d.ChainID = tmp.ChainID

	return nil, nil
}

func (d *DBTestObject) UnmarshalBinary(data []byte) error {
	_, err := d.UnmarshalBinaryData(data)
	return err
}

func (d *DBTestObject) MarshalBinary() ([]byte, error) {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data)

	tmp := new(bareDBTestObject)

	tmp.Data = d.Data
	tmp.DatabaseHeight = d.DatabaseHeight
	tmp.PrimaryIndex = d.PrimaryIndex
	tmp.SecondaryIndex = d.SecondaryIndex
	tmp.ChainID = d.ChainID

	err := enc.Encode(tmp)
	if err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}

func (d1 *DBTestObject) IsEqual(d2 *DBTestObject) bool {
	if d1.DatabaseHeight != d2.DatabaseHeight {
		return false
	}

	if primitives.AreBytesEqual(d1.Data, d2.Data) == false {
		return false
	}

	if primitives.AreBytesEqual(d1.PrimaryIndex.Bytes(), d2.PrimaryIndex.Bytes()) == false {
		return false
	}

	if primitives.AreBytesEqual(d1.SecondaryIndex.Bytes(), d2.SecondaryIndex.Bytes()) == false {
		return false
	}

	if primitives.AreBytesEqual(d1.ChainID.Bytes(), d2.ChainID.Bytes()) == false {
		return false
	}

	return true
}

func CopyZeroHash() []byte {
	answer := make([]byte, len(constants.ZERO_HASH))
	return answer
}