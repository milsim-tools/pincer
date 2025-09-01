package authz

const (
	// Full access to the unit, including managing members and settings.
	PermissionAdministrator = 1 << 0

	// No permissions, the member is banned.
	PermissionBanned = 1 << 1

	// Can view the member list.
	PermissionViewMembers = 1 << 2

	// Can manage members, including inviting, removing, and changing permissions.
	PermissionManageMembers = 1 << 3

	// Can view and manage applications to join the unit.
	PermissionViewApplications = 1 << 4

	// Can manage applications to join the unit.
	PermissionManageApplications = 1 << 5

	// Can view events.
	PermissionViewEvents = 1 << 6

	// Can respond to event invitations.
	PermissionRespondEvents = 1 << 7

	// Can manage events, including creating, updating, and deleting events.
	PermissionManageEvents = 1 << 8

	// Can view sections.
	PermissionViewSections = 1 << 9

	// Can manage sections, including creating, updating, and deleting sections.
	PermissionManageSections = 1 << 10
)

// Can checks if the given member permissions include the required permission(s).
func Can(member int32, check int32) bool {
	return (member & check) == check
}

// With adds the given permission(s) to the member permissions.
func With(member int32, add int32) int32 {
	return member | add
}

// Without removes the given permission(s) from the member permissions.
func Without(member int32, add int32) int32 {
	return member &^ add
}
