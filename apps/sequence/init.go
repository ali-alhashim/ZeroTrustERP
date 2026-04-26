package sequence

import(
	"zerotrusterp/core"
	"zerotrusterp/apps/sequence/models"
)

func init() {

	core.Register(SequenceListRoutes)
	core.RegisterModel(models.PrefixSequence{})

}