package accesscontrol

import (
	"start-feishubot/initialization"
	"sync"
)

var accessCountMap = sync.Map{}

/*
AllowAccess
If user has accessed more than 100 times according to accessCountMap, return false.
Otherwise, return true and increase the access count by 1
*/
func AllowAccess(userId *string) bool {

	// Get the access count for the user
	accessCount, ok := accessCountMap.Load(*userId)

	if !ok {
		// If the user has not accessed before, set the access count to 1
		accessCountMap.Store(*userId, 1)
		return true
	}

	// If the user has accessed more than 100 times, return false
	if accessCount.(int) >= initialization.GetConfig().AccessControlMaxCountPerUserPerDay {
		return false
	}

	// If the user has accessed before, increase the access count by 1
	accessCountMap.Store(*userId, accessCount.(int)+1)

	// Otherwise, return true
	return true
}
