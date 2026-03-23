package core

type MenuItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
}

var Menus []MenuItem

func RegisterMenus(items []MenuItem) {
	Menus = append(Menus, items...)
}
