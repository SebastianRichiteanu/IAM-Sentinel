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
	m.mapUsers(ctx, gaad.UserDetailList)
	// m.mapGroups(gaad.GroupDetailList)
	// m.mapRoles(gaad.RoleDetailList)
	// m.mapPolicies(gaad.Policies)

	return nil
}

func (m *ResourceMapper) mapUsers(ctx context.Context, users []UserDetail) {
	for _, user := range users {
		// TODO: idk if I should add groups (and policies) here, probably just ->[:IS_PART_OF] or whatever
		// probably create them before and in create user just link :D

		query := `
		MERGE (user:IAMUser {key : $arn})
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
