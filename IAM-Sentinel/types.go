package main

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	DefaultGraphProjectionName = "fullGraph"
)

var (
	DefaultNodeProjectionNames = []string{"IAM"}
	DefaultRelProjectionNames  = []string{"*"}
)

type (
	CentralityNode struct {
		Score float64     `json:"score"`
		Node  *neo4j.Node `json:"node"`
	}

	CommunityMap map[int64][]*neo4j.Node

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
		AssumeRolePolicyDocument PolicyDocument
		RoleID                   string
		CreateDate               time.Time
		InstanceProfileList      []Profile
		RoleName                 string
		Path                     string
		AttachedManagedPolicies  []ManagedPolicy
		RoleLastUsed             LastUsed
		RolePolicyList           []ManagedPolicy
		Arn                      string
	}

	Profile struct {
		InstanceProfileId   string
		Roles               []RoleDetail
		CreateDate          time.Time
		InstanceProfilename string
		Path                string
		Arn                 string
	}

	LastUsed struct {
		Region       string
		LastUsedDate time.Time
	}

	UserDetail struct {
		UserName                string
		GroupList               []string
		CreateDate              time.Time
		UserID                  string
		UserPolicyList          []ManagedPolicy
		Path                    string
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
		GroupPolicyList         []ManagedPolicy
	}

	Policy struct {
		PolicyName        string
		CreateDate        time.Time
		AttachmentCount   int
		IsAttachable      bool
		PolicyId          string
		DefaultVersionId  string
		PolicyVersionList []PolicyVersion
		Path              string
		Arn               string
		UpdateDate        time.Time
	}

	ManagedPolicy struct {
		PolicyName     string
		PolicyArn      string
		PolicyDocument *[]PolicyDocument
	}

	PolicyVersion struct {
		CreateDate       time.Time
		VersionID        string
		IsDefaultVersion bool
		Document         []PolicyDocument
	}

	PolicyDocument struct {
		Version   string
		Statement []PolicyStatement
	}

	PolicyStatement struct {
		Sid       string
		Effect    string
		Principal any // Prinicpal or string :(
		Action    []string
	}

	// It can be a star also :(
	Principal struct {
		Service       string
		AWS           string
		Federated     string
		CanonicalUser string
	}
)
