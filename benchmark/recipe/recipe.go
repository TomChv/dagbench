package recipe

import "quartz/dagbench.io/config"

type Recipe string

func Get(recipe Recipe) []*config.Command {
	switch recipe {
	case SDK:
		return benchSDK
	default:
		return nil
	}
}
