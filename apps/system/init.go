package system

import(
	  "zerotrusterp/core"
)

func init() {

    // Register routes
    core.Register(SystemRoutes)
}