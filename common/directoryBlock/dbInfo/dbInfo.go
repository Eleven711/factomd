package dbInfo

import (
	"bytes"
	"encoding/gob"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

type DirBlockInfo struct {
	// Serial hash for the directory block
	DBHash    interfaces.IHash
	DBHeight  uint32 //directory block height
	Timestamp int64  // time of this dir block info being created
	// BTCTxHash is the Tx hash returned from rpcclient.SendRawTransaction
	BTCTxHash interfaces.IHash // use string or *btcwire.ShaHash ???
	// BTCTxOffset is the index of the TX in this BTC block
	BTCTxOffset int32
	// BTCBlockHeight is the height of the block where this TX is stored in BTC
	BTCBlockHeight int32
	//BTCBlockHash is the hash of the block where this TX is stored in BTC
	BTCBlockHash interfaces.IHash // use string or *btcwire.ShaHash ???
	// DBMerkleRoot is the merkle root of the Directory Block
	// and is written into BTC as OP_RETURN data
	DBMerkleRoot interfaces.IHash
	// A flag to to show BTC anchor confirmation
	BTCConfirmed bool
}

var _ interfaces.Printable = (*DirBlockInfo)(nil)
var _ interfaces.BinaryMarshallableAndCopyable = (*DirBlockInfo)(nil)
var _ interfaces.DatabaseBatchable = (*DirBlockInfo)(nil)
var _ interfaces.IDirBlockInfo = (*DirBlockInfo)(nil)

func NewDirBlockInfo() *DirBlockInfo {
	dbi := new(DirBlockInfo)
	dbi.DBHash = primitives.NewZeroHash()
	dbi.BTCTxHash = primitives.NewZeroHash()
	dbi.BTCBlockHash = primitives.NewZeroHash()
	dbi.DBMerkleRoot = primitives.NewZeroHash()
	return dbi
}

func (e *DirBlockInfo) String() string {
	str, _ := e.JSONString()
	return str
}

func (e *DirBlockInfo) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *DirBlockInfo) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *DirBlockInfo) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (c *DirBlockInfo) New() interfaces.BinaryMarshallableAndCopyable {
	return NewDirBlockInfo()
}

func (c *DirBlockInfo) GetDatabaseHeight() uint32 {
	return c.DBHeight
}

func (c *DirBlockInfo) GetDBHeight() uint32 {
	return c.DBHeight
}

func (c *DirBlockInfo) GetBTCConfirmed() bool {
	return c.BTCConfirmed
}

func (c *DirBlockInfo) GetChainID() interfaces.IHash {
	id := make([]byte, 32)
	copy(id, []byte("DirBlockInfo"))
	return primitives.NewHash(id)
}

func (c *DirBlockInfo) DatabasePrimaryIndex() interfaces.IHash {
	return c.DBMerkleRoot
}

func (c *DirBlockInfo) DatabaseSecondaryIndex() interfaces.IHash {
	return c.DBHash
}

func (e *DirBlockInfo) GetDBMerkleRoot() interfaces.IHash {
	return e.DBMerkleRoot
}

func (e *DirBlockInfo) GetBTCTxHash() interfaces.IHash {
	return e.BTCTxHash
}

func (e *DirBlockInfo) GetTimestamp() int64 {
	return e.Timestamp
}

func (e *DirBlockInfo) GetBTCBlockHeight() int32 {
	return e.BTCBlockHeight
}

func (e *DirBlockInfo) MarshalBinary() ([]byte, error) {
	var data primitives.Buffer

	enc := gob.NewEncoder(&data)

	err := enc.Encode(newDirBlockInfoCopyFromDBI(e))
	if err != nil {
		return nil, err
	}
	return data.DeepCopyBytes(), nil
}

func (e *DirBlockInfo) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	dec := gob.NewDecoder(primitives.NewBuffer(data))
	dbic := newDirBlockInfoCopy()
	err = dec.Decode(dbic)
	if err != nil {
		return nil, err
	}
	e.parseDirBlockInfoCopy(dbic)
	return nil, nil
}

func (e *DirBlockInfo) UnmarshalBinary(data []byte) (err error) {
	_, err = e.UnmarshalBinaryData(data)
	return
}

type dirBlockInfoCopy struct {
	// Serial hash for the directory block
	DBHash    interfaces.IHash
	DBHeight  uint32 //directory block height
	Timestamp int64  // time of this dir block info being created
	// BTCTxHash is the Tx hash returned from rpcclient.SendRawTransaction
	BTCTxHash interfaces.IHash // use string or *btcwire.ShaHash ???
	// BTCTxOffset is the index of the TX in this BTC block
	BTCTxOffset int32
	// BTCBlockHeight is the height of the block where this TX is stored in BTC
	BTCBlockHeight int32
	//BTCBlockHash is the hash of the block where this TX is stored in BTC
	BTCBlockHash interfaces.IHash // use string or *btcwire.ShaHash ???
	// DBMerkleRoot is the merkle root of the Directory Block
	// and is written into BTC as OP_RETURN data
	DBMerkleRoot interfaces.IHash
	// A flag to to show BTC anchor confirmation
	BTCConfirmed bool
}

func newDirBlockInfoCopyFromDBI(dbi *DirBlockInfo) *dirBlockInfoCopy {
	dbic := new(dirBlockInfoCopy)
	dbic.DBHash = dbi.DBHash
	dbic.DBHeight = dbi.DBHeight
	dbic.Timestamp = dbi.Timestamp
	dbic.BTCTxHash = dbi.BTCTxHash
	dbic.BTCTxOffset = dbi.BTCTxOffset
	dbic.BTCBlockHeight = dbi.BTCBlockHeight
	dbic.BTCBlockHash = dbi.BTCBlockHash
	dbic.DBMerkleRoot = dbi.DBMerkleRoot
	dbic.BTCConfirmed = dbi.BTCConfirmed
	return dbic
}

func newDirBlockInfoCopy() *dirBlockInfoCopy {
	dbi := new(dirBlockInfoCopy)
	dbi.DBHash = primitives.NewZeroHash()
	dbi.BTCTxHash = primitives.NewZeroHash()
	dbi.BTCBlockHash = primitives.NewZeroHash()
	dbi.DBMerkleRoot = primitives.NewZeroHash()
	return dbi
}

func (dbic *DirBlockInfo) parseDirBlockInfoCopy(dbi *dirBlockInfoCopy) {
	dbic.DBHash = dbi.DBHash
	dbic.DBHeight = dbi.DBHeight
	dbic.Timestamp = dbi.Timestamp
	dbic.BTCTxHash = dbi.BTCTxHash
	dbic.BTCTxOffset = dbi.BTCTxOffset
	dbic.BTCBlockHeight = dbi.BTCBlockHeight
	dbic.BTCBlockHash = dbi.BTCBlockHash
	dbic.DBMerkleRoot = dbi.DBMerkleRoot
	dbic.BTCConfirmed = dbi.BTCConfirmed
}

// NewDirBlockInfoFromDirBlock creates a DirDirBlockInfo from DirectoryBlock
func NewDirBlockInfoFromDirBlock(dirBlock interfaces.IDirectoryBlock) *DirBlockInfo {
	dbic := new(DirBlockInfo)
	dbic.DBHash = dirBlock.GetHash()
	dbic.DBHeight = dirBlock.GetDatabaseHeight()
	dbic.DBMerkleRoot = dirBlock.GetKeyMR()
	dbic.Timestamp = int64(dirBlock.GetHeader().GetTimestamp()) // * 60 ???
	dbic.BTCTxHash = primitives.NewZeroHash()
	dbic.BTCBlockHash = primitives.NewZeroHash()
	dbic.BTCConfirmed = false
	return dbic
}
