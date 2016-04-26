package devices

import (
	log "github.com/Sirupsen/logrus"
)

var global_envs map[string]string

func AddGlobalEnv(k string, v string) {
	_, ok := global_envs[k]
	if !ok {
		Log.WithFields(log.Fields{
			"key": k,
		}).Warning("setting key for envs twice")
	}

	global_envs[k] = v
}

func GetGlobalEnvs() map[string]string {
	return global_envs
}
