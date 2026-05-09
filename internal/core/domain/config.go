package domain

const (
	ConfigKeyJwtTtl             = "jwt_ttl"
	ConfigKeyNetworkWhitelist   = "network_whitelist"
	ConfigKeyRebuildPolicyCache = "rebuild_policy_cache"
)

var ValidConfigKeys = map[string]bool{
	ConfigKeyJwtTtl:             true,
	ConfigKeyNetworkWhitelist:   true,
	ConfigKeyRebuildPolicyCache: true,
}

func IsValidConfigKey(key string) bool {
	return ValidConfigKeys[key]
}
