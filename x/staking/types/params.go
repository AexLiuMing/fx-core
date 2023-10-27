package types

import (
	"time"
)

// GrantExpirationTime period ends at Max Time supported by Amino
// (Dec 31, 9999 - 23:59:59 GMT).
var GrantExpirationTime = time.Unix(253402300799, 0)
