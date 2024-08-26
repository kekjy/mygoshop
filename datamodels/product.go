package datamodels

type Product struct {
	ID           int64  `json:"ID" sql:"ID" productshop:"ID" gorm:"column:ID"`
	ProductName  string `json:"productName" sql:"productName" productshop:"productName" gorm:"column:productName" `
	ProductNum   int64  `json:"productNum" sql:"productNum" productshop:"productNum" gorm:"column:productNum"`
	ProductImage string `json:"productImage" sql:"productImage" productshop:"productImage" gorm:"column:productImage"`
	ProductURL   string `json:"productURL" sql:"productURL" productshop:"productURL" gorm:"column:productURL"`
	ProductGHURL string `json:"productGHURL" sql:"productGHURL" productshop:"productGHURL" gorm:"column:productGHURL"`
	ProductInfo  string `json:"productInfo" sql:"productInfo" productshop:"productInfo" gorm:"column:productInfo"`
}
