package common

const(
	DEFAULT_NAMESPACE = "default"
	DEFAULT_ROLEMAP_NAMESPACE = "default"
	DEFAULT_ROLEMAP_NAME = "role-map"
	DEFAULT_PATHS_SEP = ","
	DEFAULT_PATH_SEGMENT_SEP = "."
)

var(
	USE_PATHS_FOR_ROLES = false
	TOKEN_ROLE_PATH = "realm_access.roles,resource_access.account.roles"
	TOKEN_PATHS_SEP = DEFAULT_PATHS_SEP
	TOKEN_PATH_SEGMENT_SEP = DEFAULT_PATH_SEGMENT_SEP
)