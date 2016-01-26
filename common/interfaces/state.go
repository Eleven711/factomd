// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

import ()

// Holds the state information for factomd.  This does imply that we will be
// using accessors to access state information in the consensus algorithm.
// This is a bit tedious, but does provide single choke points where information
// can be logged about the execution of Factom.  Also ensures that we do not
// accidentally
type IState interface {

	// Server
	
	GetServerIndex() int			// Returns this server's index, if a federated server
	GetCfg() IFactomConfig
	Init(string)
	String() string
	Sign([]byte) IFullSignature
	GetProcessListLen(list int) int

	GetServer() IServer
	SetServer(IServer)

	// Channels
	//==========

	// Network Processor
	NetworkInMsgQueue() chan IMsg // Not sure that IMsg is the right type... TBD
	NetworkOutMsgQueue() chan IMsg
	NetworkInvalidMsgQueue() chan IMsg

	// Consensus
	InMsgQueue() chan IMsg         // Read by Validate
	LeaderInMsgQueue() chan IMsg   // Processed by the Leader
	FollowerInMsgQueue() chan IMsg // Processed by the Follower

	// Lists and Maps
	// =====
	// The leader CANNOT touch these lists!  Only the FollowerExecution
	// methods can touch them safely.
	GetAuditServers() []IServer   // List of Audit Servers
	GetFedServers() []IServer     // List of Federated Servers
	GetServerOrder() [][]IServer  // 10 lists for Server Order for each minute
	GetAuditHeartBeats() []IMsg   // The checklist of HeartBeats for this period
	GetFedServerFaults() [][]IMsg // Keep a fault list for every server

	GetNewEBlks([32]byte) IEntryBlock
	PutNewEBlks([32]byte, IEntryBlock)

	GetCommits(IHash) IMsg
	PutCommits(IHash, IMsg)
	// Server Configuration
	// ====================

	//Network MAIN = 0, TEST = 1, LOCAL = 2, CUSTOM = 3
	GetNetworkNumber() int  // Encoded into Directory Blocks
	GetNetworkName() string // Some networks have defined names

	// Number of Servers acknowledged by Factom
	GetTotalServers() int
	GetServerState() int    // (0 if client, 1 if server, 2 if audit server
	GetMatryoshka() []IHash // Reverse Hash

	LeaderFor([]byte) bool // Tests if this server is the leader for this key

	// Database
	// ========
	GetDB() DBOverlay
	SetDB(DBOverlay)

	// Directory Block State
	// =====================
	GetPreviousDirectoryBlock() IDirectoryBlock // The previous directory block
	GetCurrentDirectoryBlock() IDirectoryBlock  // The directory block under construction
	SetCurrentDirectoryBlock(IDirectoryBlock)

	GetCurrentEntryCreditBlock() IEntryCreditBlock
	SetCurrentEntryCreditBlock(IEntryCreditBlock)

	GetCurrentAdminBlock() IAdminBlock
	SetCurrentAdminBlock(IAdminBlock)

	GetDBHeight() uint32 // The index of the directory block under construction.

	// Message State
	GetLastAck() IMsg // Return the last Acknowledgement set by this server
	SetLastAck(IMsg)

	// Server Methods
	// ==============
	UpdateProcessLists()
	ProcessEndOfBlock()

	// Web Services
	// ============
	SetPort(int)
	GetPort() int

	// Factoid State
	// =============
	GetFactoidState() IFactoidState
	GetPrevFactoidKeyMR() IHash
	SetPrevFactoidKeyMR(IHash)

	// MISC
	// ====

	// Returns true if it found a match
	MatchAckFollowerExecute(m IMsg) (bool, error)
	FollowerExecuteAck(m IMsg) error
	GetTimestamp() Timestamp
	GetNewHash() IHash // Return a new Hash object
	CreateDBlock() (b IDirectoryBlock, err error)
	PrintType(int) bool // Debugging

	RecalculateBalances() error
	LogInfo(args ...interface{})
}
