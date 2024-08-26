package datamodels

type Order struct {
	ID          int64
	UserId      int64 `gorm:"column:userId"`
	ProductId   int64 `gorm:"column:productId"`
	OrderStatus int64 `gorm:"column:orderStatus"`
}

const (
	OrderWait    = iota
	OrderSuccess // 1
	OrderFailed  // 2
)
