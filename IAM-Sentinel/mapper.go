package main

import "context"

type ResourceMapper struct {
	dbConn *Neo4jConn
	parser *Parser
}

func NewResourceMapper(db *Neo4jConn, parser *Parser) (*ResourceMapper, error) {
	m := ResourceMapper{
		dbConn: db,
		parser: parser,
	}
	return &m, nil
}

func (m *ResourceMapper) MapFile(ctx context.Context, data []byte) error {
	gaad, err := m.parser.parse(data)
	if err != nil {
		return err
	}

	// TODO: This should be go func() but after it's stable enough :( or create var for PARALLEL?
	m.mapGroups(ctx, gaad.GroupDetailList)
	m.mapUsers(ctx, gaad.UserDetailList)
	m.mapRoles(ctx, gaad.RoleDetailList)
	m.mapPolicies(ctx, gaad.Policies)

	return nil
}

func (m *ResourceMapper) mapUsers(ctx context.Context, users []UserDetail) {
	for _, user := range users {
		var managedPolicies []Policy
		for _, policy := range user.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		m.mapPolicies(ctx, managedPolicies)
		// TODO: idk if I should add groups (and policies) here, probably just ->[:IS_PART_OF] or whatever
		// probably create them before and in create user just link :D

		query := `
		MERGE (user:User {key : $arn})
			SET user.UserName = $userName
			SET user.GroupList = $groups 
			SET user.UserID = $userID
			SET user.Path = $path
			SET user.CreateDate = $createDate
		`
		// SET user.ManagedPolicies = $managedPolicies
		// SET user.PolicyList = $policyList

		params := map[string]any{
			"arn":        user.Arn,
			"userName":   user.UserName,
			"groups":     user.GroupList,
			"userID":     user.UserID,
			"path":       user.Path,
			"createDate": user.CreateDate,
			//"managedPolicies": user.AttachedManagedPolicies,
			//"policyList":      user.UserPolicyList,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			// change this to a logger
			panic(err)
		}
	}
}

func (m *ResourceMapper) mapGroups(ctx context.Context, groups []GroupDetail) {
	for _, group := range groups {
		var managedPolicies []Policy
		for _, policy := range group.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		m.mapPolicies(ctx, managedPolicies)
		// TODO: idk if I should add policies here, probably just ->[:IS_PART_OF] or whatever
		// probably create them before and in create user just link :D

		query := `
		MERGE (user:Group {key : $arn})
			SET user.GroupID = $groupID
			SET user.GroupName = $groupName 
			SET user.Path = $path
			SET user.CreateDate = $createDate
		`
		// SET user.ManagedPolicies = $managedPolicies
		// SET user.PolicyList = $policyList

		params := map[string]any{
			"arn":        group.Arn,
			"groupID":    group.GroupID,
			"groupName":  group.GroupName,
			"path":       group.Path,
			"createDate": group.CreateDate,
			//"managedPolicies": user.AttachedManagedPolicies,
			//"policyList":      user.UserPolicyList,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			// change this to a logger
			panic(err)
		}
	}
}

func (m *ResourceMapper) mapPolicies(ctx context.Context, policies []Policy) {
	for _, policy := range policies {
		// TODO: idk if I should add policies here, probably just ->[:IS_PART_OF] or whatever
		// probably create them before and in create user just link :D

		query := `
		MERGE (policy:Policy {key : $arn})
			SET policy.PolicyID = $policyID
			SET policy.PolicyName = $policyName 
			SET policy.Path = $path
			SET policy.CreateDate = $createDate
			SET policy.AttachementCount = $attachementCount
			SET policy.IsAttachable = $isAttachable
			SET policy.DefaultVersionId = $defaultVersionID
			SET policy.updatedDate = $updatedDate
		`
		// TODO: PolicyVersionList

		params := map[string]any{
			"arn":              policy.Arn,
			"policyID":         policy.PolicyId,
			"policyName":       policy.PolicyName,
			"path":             policy.Path,
			"createDate":       policy.CreateDate,
			"attachementCount": policy.AttachmentCount,
			"isAttachable":     policy.IsAttachable,
			"defaultVersionID": policy.DefaultVersionId,
			"updatedDate":      policy.UpdateDate,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			// change this to a logger
			panic(err)
		}
	}
}

func (m *ResourceMapper) mapRoles(ctx context.Context, roles []RoleDetail) {
	for _, role := range roles {
		var managedPolicies []Policy
		for _, policy := range role.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		m.mapPolicies(ctx, managedPolicies)

		// TODO: idk if I should add policies here, probably just ->[:IS_PART_OF] or whatever
		// probably create them before and in create user just link :D

		query := `
		MERGE (role:Role {key : $arn})
			SET role.RoleID = $roleID
			SET role.RoleName = $roleName 
			SET role.Path = $path
			SET role.CreateDate = $createDate
		`
		// TODO: AssumeRolePolicyDocument
		// InstanceProfileList
		// RoleLastUsed
		// RolePolicyList

		params := map[string]any{
			"arn":        role.Arn,
			"roleID":     role.RoleID,
			"roleName":   role.RoleName,
			"path":       role.Path,
			"createDate": role.CreateDate,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			// change this to a logger
			panic(err)
		}
	}
}
