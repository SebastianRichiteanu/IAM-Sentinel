package main

import "time"

type (
	// get-account-authorization-details structure
	// https://docs.aws.amazon.com/cli/latest/reference/iam/get-account-authorization-details.html
	GAAD struct {
		RoleDetailList  []RoleDetail
		GroupDetailList []GroupDetail
		UserDetailList  []UserDetail
		Policies        []Policy
		Marker          string
		IsTruncated     bool // if true, Marker is used to paginate
	}

	RoleDetail struct {
		AssumeRolePolicyDocument any // TODO
		RoleID                   string
		CreateDate               time.Time
		InstanceProfileList      any // TODO
		RoleName                 string
		Path                     string
		AttachedManagedPolicies  []ManagedPolicy
		RoleLastUsed             any   // TODO
		RolePolicyList           []any // TODO
		Arn                      string
	}

	UserDetail struct {
		UserName                string
		GroupList               []string
		CreateDate              time.Time
		UserID                  string
		UserPolicyList          any    // TODO
		Path                    string // ????
		AttachedManagedPolicies []ManagedPolicy
		Arn                     string
	}

	GroupDetail struct {
		GroupID                 string
		AttachedManagedPolicies []ManagedPolicy
		GroupName               string
		Path                    string
		Arn                     string
		CreateDate              time.Time
		GroupPolicyList         []any // TODO: tighten to policy
	}

	Policy struct { // maybe name PolicyDetail?
		// maybe switch to pointer to have default nil?
		PolicyName        string
		CreateDate        time.Time
		AttachmentCount   int
		IsAttachable      bool
		PolicyId          string
		DefaultVersionId  string
		PolicyVersionList any // TODO
		Path              string
		Arn               string
		UpdateDate        time.Time
	}

	ManagedPolicy struct {
		PolicyName string
		PolicyArn  string
	}
)
