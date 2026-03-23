package employees

import (
	"encoding/json"
	"os"

	"zerotrusterp/core"
)

func init() {
	file, err := os.ReadFile("apps/employees/menu.json")
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