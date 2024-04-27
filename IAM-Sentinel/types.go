package main

import "time"

type (
	// get-account-authorization-details structure
	GAAD struct {
		RoleDetailList  any          // Tighten to aws internal stuff but arrray???
		GroupDetailList any          // same as above
		UserDetailList  []UserDetail // same
		Policies        any
		Marker          string
		IsTruncated     bool // if true, Marker is used to paginate
	}

	UserDetail struct {
		UserName                string
		GroupList               []string
		CreateDate              time.Time
		UserID                  string
		UserPolicyList          any      // TODO
		Path                    string   // ????
		AttachedManagedPolicies []string // ???
		Arn                     string
	}
)
