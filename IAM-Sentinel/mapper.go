package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type ResourceMapper struct {
	logger *logrus.Logger
	dbConn *Neo4jConn
	parser *Parser
}

func NewResourceMapper(logger *logrus.Logger, db *Neo4jConn, parser *Parser) (*ResourceMapper, error) {
	m := ResourceMapper{
		logger: logger,
		dbConn: db,
		parser: parser,
	}
	return &m, nil
}

func (m *ResourceMapper) MapFolder(ctx context.Context, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(dir, fileName)

		logFields := logrus.Fields{"file_name": fileName}
		m.logger.WithFields(logFields).Info("Reading file")
		content, err := os.ReadFile(filePath)
		if err != nil {
			m.logger.WithFields(logFields).Error("Could not read file")
			continue
		}

		if err := m.MapFile(ctx, content); err != nil {
			m.logger.WithFields(logFields).Error("Could not map file")
		}
	}

	return nil
}

func (m *ResourceMapper) MapFile(ctx context.Context, data []byte) error {
	gaad, err := m.parser.parse(data)
	if err != nil {
		return err
	}

	m.mapGroups(ctx, gaad.GroupDetailList)
	m.mapUsers(ctx, gaad.UserDetailList)
	m.mapRoles(ctx, gaad.RoleDetailList)
	m.mapPolicies(ctx, gaad.Policies)
	m.mapRelationships(ctx)

	return nil
}

func (m *ResourceMapper) mapUsers(ctx context.Context, users []UserDetail) {
	for _, user := range users {
		var managedPolicies []Policy
		for _, policy := range user.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		for _, policy := range user.UserPolicyList {
			arn := policy.PolicyArn
			if arn == "" {
				// if no arn, use name as key
				arn = policy.PolicyName
			}
			managedPolicies = append(managedPolicies, Policy{Arn: arn, PolicyName: policy.PolicyName})
		}

		m.mapPolicies(ctx, managedPolicies)

		query := `
		MERGE (user:IAM:User {arn : $arn})
			SET user.UserName = $userName
			SET user.GroupList = $groups 
			SET user.UserID = $userID
			SET user.Path = $path
			SET user.CreateDate = $createDate
		`

		params := map[string]any{
			"arn":        user.Arn,
			"userName":   user.UserName,
			"groups":     user.GroupList,
			"userID":     user.UserID,
			"path":       user.Path,
			"createDate": user.CreateDate,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			m.logger.WithFields(logrus.Fields{"arn": user.Arn}).Error(err)
		}

		m.mapEntityPolicies(ctx, user.Arn, managedPolicies)

		for _, version := range user.UserPolicyList {
			for _, document := range *version.PolicyDocument {
				for _, statement := range document.Statement {
					for _, action := range statement.Action {
						m.mapAction(ctx, "TRUSTS", user.Arn, action, statement.Effect)
					}
				}
			}
		}
	}
}

func (m *ResourceMapper) mapGroups(ctx context.Context, groups []GroupDetail) {
	for _, group := range groups {
		var managedPolicies []Policy
		for _, policy := range group.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		for _, policy := range group.GroupPolicyList {
			arn := policy.PolicyArn
			if arn == "" {
				// if no arn, use name as key
				arn = policy.PolicyName
			}
			managedPolicies = append(managedPolicies, Policy{Arn: arn, PolicyName: policy.PolicyName})
		}

		m.mapPolicies(ctx, managedPolicies)

		query := `
		MERGE (group:IAM:Group {arn : $arn})
			SET group.GroupID = $groupID
			SET group.GroupName = $groupName 
			SET group.Path = $path
			SET group.CreateDate = $createDate
		`

		params := map[string]any{
			"arn":        group.Arn,
			"groupID":    group.GroupID,
			"groupName":  group.GroupName,
			"path":       group.Path,
			"createDate": group.CreateDate,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			m.logger.WithFields(logrus.Fields{"arn": group.Arn}).Error(err)
		}

		m.mapEntityPolicies(ctx, group.Arn, managedPolicies)

		for _, version := range group.GroupPolicyList {
			for _, document := range *version.PolicyDocument {
				for _, statement := range document.Statement {
					for _, action := range statement.Action {
						m.mapAction(ctx, "TRUSTS", group.Arn, action, statement.Effect)
					}
				}
			}
		}
	}
}

func (m *ResourceMapper) mapPolicies(ctx context.Context, policies []Policy) {
	for _, policy := range policies {
		query := `
		MERGE (policy:IAM:Policy {arn : $arn})
			SET policy.PolicyID = $policyID
			SET policy.PolicyName = $policyName 
			SET policy.Path = $path
			SET policy.CreateDate = $createDate
			SET policy.AttachementCount = $attachementCount
			SET policy.IsAttachable = $isAttachable
			SET policy.DefaultVersionId = $defaultVersionID
			SET policy.updatedDate = $updatedDate
		`

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
			m.logger.WithFields(logrus.Fields{"arn": policy.Arn}).Error(err)
		}

		for _, version := range policy.PolicyVersionList {
			for _, document := range version.Document {
				for _, statement := range document.Statement {
					for _, action := range statement.Action {
						m.mapAction(ctx, "TRUSTS", policy.Arn, action, statement.Effect)
					}
				}
			}
		}
	}
}

func (m *ResourceMapper) mapRoles(ctx context.Context, roles []RoleDetail) {
	for _, role := range roles {
		var managedPolicies []Policy
		for _, policy := range role.AttachedManagedPolicies {
			managedPolicies = append(managedPolicies, Policy{Arn: policy.PolicyArn, PolicyName: policy.PolicyName})
		}
		for _, policy := range role.RolePolicyList {
			arn := policy.PolicyArn
			if arn == "" {
				// if no arn, use name as key
				arn = policy.PolicyName
			}
			managedPolicies = append(managedPolicies, Policy{Arn: arn, PolicyName: policy.PolicyName})
		}

		m.mapPolicies(ctx, managedPolicies)

		query := `
		MERGE (role:IAM:Role {arn : $arn})
			SET role.RoleID = $roleID
			SET role.RoleName = $roleName 
			SET role.Path = $path
			SET role.CreateDate = $createDate
		`

		params := map[string]any{
			"arn":        role.Arn,
			"roleID":     role.RoleID,
			"roleName":   role.RoleName,
			"path":       role.Path,
			"createDate": role.CreateDate,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			m.logger.WithFields(logrus.Fields{"arn": role.Arn}).Error(err)
		}

		m.mapEntityPolicies(ctx, role.Arn, managedPolicies)

		for _, version := range role.RolePolicyList {
			for _, document := range *version.PolicyDocument {
				for _, statement := range document.Statement {
					for _, action := range statement.Action {
						m.mapAction(ctx, "TRUSTS", role.Arn, action, statement.Effect)
					}
				}
			}
		}

		for _, statement := range role.AssumeRolePolicyDocument.Statement {
			for _, action := range statement.Action {
				m.mapAction(ctx, "CAN_ASSUME", role.Arn, action, statement.Effect)
			}
		}

		for _, instance := range role.InstanceProfileList {
			for _, role := range instance.Roles {
				for _, statement := range role.AssumeRolePolicyDocument.Statement {
					for _, action := range statement.Action {
						m.mapAction(ctx, "CAN_ASSUME", role.Arn, action, statement.Effect)
					}
				}
			}
		}
	}
}

func (m *ResourceMapper) mapEntityPolicies(ctx context.Context, entityArn string, managedPolicies []Policy) {
	for _, policy := range managedPolicies {
		query := `
		MATCH (entity {arn: $entityArn})
		MATCH (policy:IAM:Policy {arn: $policyArn})
		
		MERGE (entity)-[:MANAGES]->(policy)
		`

		params := map[string]any{
			"entityArn": entityArn,
			"policyArn": policy.Arn,
		}

		if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
			m.logger.WithFields(logrus.Fields{"entityArn": entityArn, "policyArn": policy.Arn}).Error(err)
		}
	}
}

func (m *ResourceMapper) mapRelationships(ctx context.Context) {
	query := `
	MATCH (u:IAM:User)
	UNWIND u.GroupList AS groupName
	MATCH (g:Group {GroupName: groupName})
	MERGE (u)-[:MEMBER_OF]->(g)
	`

	if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, nil); err != nil {
		m.logger.Error(err)
	}
}

func (m *ResourceMapper) mapAction(ctx context.Context, entityActionEdge, arn, action, effect string) {
	query := fmt.Sprintf(`
	MERGE (action:IAM:action {action: $action})
	WITH action

	MATCH (entity {arn : $arn})
	MERGE (entity)-[r:%s]->(action)

	SET r.effect = $effect
	`, entityActionEdge)

	params := map[string]any{
		"arn":    arn,
		"action": action,
		"effect": effect,
	}

	if _, err := m.dbConn.ExecuteQueryWrite(ctx, query, params); err != nil {
		m.logger.WithFields(logrus.Fields{"arn": arn, "action": action}).Error(err)
	}
}
