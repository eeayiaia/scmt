package devices

var global_envs map[string]string

func AddGlobalEnv(k string, v string) {
	if global_envs == nil {
		global_envs = make(map[string]string)
	}

	global_envs[k] = v
}

func GetGlobalEnvs() map[string]string {
	return global_envs
}
