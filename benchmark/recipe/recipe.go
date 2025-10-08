package recipe

import "quartz/dagbench.io/config"

type Recipe string

// Get returns the commands for a given recipe or nil if not found.
func Get(recipe Recipe) []*config.Command {
	switch recipe {
	case SDK:
		return benchSDK
	default:
		return nil
	}
}
