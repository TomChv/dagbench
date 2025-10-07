package recipe

import "quartz/dagbench.io/config"

const SDK Recipe = "sdk"

var benchSDK = []*config.Command{
	// Dagger develop
	{
		Args:      []string{"develop"},
		SpanNames: []string{"develop"},
	},
	// Dagger functions
	{
		Args:      []string{"functions"},
		SpanNames: []string{"load module"},
	},
	// Dagger default call
	{
		Args:      []string{"call", "container-echo", "--string-arg=hello"},
		SpanNames: []string{"load module", "containerEcho"},
	},
}
