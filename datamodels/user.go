package datamodels

type User struct {
	ID           int64  `json:"ID" form:"ID" sql:"ID"  gorm:"column:ID"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName" gorm:"column:nickName"`
	UserName     string `json:"userName" form:"userName" sql:"userName" gorm:"column:userName"`
	HashPassword string `json:"-" form:"password" sql:"password" gorm:"column:password"`
}
