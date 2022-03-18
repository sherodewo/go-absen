package dto

type MenuDto struct {
	ParentID    string   `json:"parent_id" form:"parent_id"`
	MenuTitle   string   `json:"menu_title" form:"menu_title"`
	Slug        string   `json:"slug" form:"slug"`
	Url         string   `json:"url" form:"url"`
	Icon        string   `json:"icon" form:"icon"`
	MenuOrder   int      `json:"menu_order" form:"menu_order"`
	Description string   `json:"description" form:"description"`
	IsActive    bool     `json:"is_active" form:"is_active"`
	Role        []string `json:"roles" form:"roles[]"`
}

type MenuUpdateDto struct {
	ParentID    string   `json:"parent_id" form:"parent_id"`
	MenuTitle   string   `json:"menu_title" form:"menu_title"`
	Slug        string   `json:"slug" form:"slug"`
	Url         string   `json:"url" form:"url"`
	Icon        string   `json:"icon" form:"icon"`
	MenuOrder   int      `json:"menu_order" form:"menu_order"`
	Description string   `json:"description" form:"description"`
	IsActive    bool     `json:"is_active" form:"is_active"`
	Role        []string `json:"roles" form:"roles[]"`
}

