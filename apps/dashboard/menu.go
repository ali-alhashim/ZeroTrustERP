package dashboard

import (
	"encoding/json"
	"os"

	"zerotrusterp/core"
)

func init() {
	file, err := os.ReadFile("apps/dashboard/menu.json")
	if err != nil {
		panic(err)
	}

	var items []core.MenuItem
	err = json.Unmarshal(file, &items)
	if err != nil {
		panic(err)
	}

	core.RegisterMenus(items)
}
