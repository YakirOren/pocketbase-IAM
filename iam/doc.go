// Package iam implements AWS IAM-inspired RBAC for PocketBase.
package iam

// IAM collection names used throughout the package.
const (
	colPolicies           = "iam_policies"
	colRolePolicies       = "iam_role_policies"
	colUserPolicies       = "iam_user_policies"
	colUserRoles          = "iam_user_roles"
	colGroups             = "iam_groups"
	colGroupUsers         = "iam_group_users"
	colGroupPolicies      = "iam_group_policies"
	colManagedCollections = "iam_managed_collections"
	colRoles              = "iam_roles"
	colActionRegistry     = "iam_action_registry"
)
