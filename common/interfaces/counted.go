// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

import ()

/**************************
 * IRCD  Interface for Redeem Condition Datastructures (RCD)
 *
 * https://github.com/FactomProject/FactomDocs/blob/master/factomDataStructureDetails.md#factoid-transaction
 **************************/

type ICounted interface {
	GetCount() int
	IncCount()
	SetCount(int)
}
