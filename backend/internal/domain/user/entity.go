package user

const (
	RoleAdmin  = "admin"
	RoleParent = "parent"
	RoleKid    = "kid"
)

type User struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Avatar        string `json:"avatar"`
	ParentID      *uint  `json:"parentId"`
	PreferredLang string `json:"preferredLang"`
}
