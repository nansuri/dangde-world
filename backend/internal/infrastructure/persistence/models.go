package persistence

type userModel struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Role          string `gorm:"index"`
	Avatar        string
	ParentID      *uint
	PreferredLang string
}

type categoryModel struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Slug     string `gorm:"uniqueIndex"`
	Type     string
	ParentID *uint
}

type activityModel struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	Language    string
	Difficulty  string
	AgeGroup    string
	CategoryID  uint
	Category    categoryModel `gorm:"foreignKey:CategoryID"`
	Prompt      string
	Icon        string
}

type assignmentModel struct {
	ID         uint `gorm:"primaryKey"`
	ParentID   uint
	Parent     userModel `gorm:"foreignKey:ParentID"`
	KidID      uint
	Kid        userModel `gorm:"foreignKey:KidID"`
	ActivityID uint
	Activity   activityModel `gorm:"foreignKey:ActivityID"`
	Status     string
	Progress   int
}

type activityDataModel struct {
	ID           uint   `gorm:"primaryKey"`
	ActivityID   uint   `gorm:"index"`
	AssignmentID uint   `gorm:"index"`
	KidID        uint   `gorm:"index"`
	Key          string `gorm:"index"`
	Value        string `gorm:"type:text"`
}
