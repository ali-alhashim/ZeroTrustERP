package core

type MenuItem struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
}

type MenuGroup struct {
    GroupName string     `json:"group_name"`
    Items     []MenuItem `json:"items"`
}

var MenuGroups []MenuGroup

func RegisterAppMenu(groupName string, items []MenuItem) {
    newGroup := MenuGroup{
        GroupName: groupName,
        Items:     items,
    }
    MenuGroups = append(MenuGroups, newGroup)
}
