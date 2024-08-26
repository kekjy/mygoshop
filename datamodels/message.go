package datamodels

type Message struct {
	UserID     int64
	ProductID  int64
	ProductNum int64
}

func NewMessage(userID int64, productID int64, productNum int64) *Message {
	return &Message{
		UserID:     userID,
		ProductID:  productID,
		ProductNum: productNum,
	}
}
