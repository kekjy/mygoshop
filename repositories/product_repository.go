package repositories

import (
	"mygoshop/datamodels"
	"mygoshop/db"

	"gorm.io/gorm"
)

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProduct(int64, int64) error
}

type ProductManager struct {
	sqlConn *gorm.DB
}

// 新建商品管理接口
func NewProductManager(db *gorm.DB) IProduct {
	return &ProductManager{
		sqlConn: db,
	}
}

func (p *ProductManager) Conn() error {
	if p.sqlConn == nil {
		__sqlConn, err := db.NewDbConn()
		if err != nil {
			return err
		}
		p.sqlConn = __sqlConn
	}
	return nil
}

func (p *ProductManager) Insert(product *datamodels.Product) (int64, error) {
	if err := p.Conn(); err != nil {
		return -1, err
	}
	err := p.sqlConn.Create(product).Error
	if err != nil {
		return -1, err
	}
	return product.ID, nil
}

func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	err := p.sqlConn.Where("ID=?", productId).Delete(&datamodels.Product{}).Error
	return err == nil
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	err := p.sqlConn.Save(product).Error
	return err
}

func (p *ProductManager) SelectByKey(productId int64) (*datamodels.Product, error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	productResult := &datamodels.Product{}
	err := p.sqlConn.Where("ID=?", productId).First(productResult).Error
	if err != nil {
		return nil, err
	}
	return productResult, err
}

func (p *ProductManager) SelectAll() (products []*datamodels.Product, err error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	var productResult []datamodels.Product
	if err := p.sqlConn.Find(&productResult).Error; err != nil {
		return nil, err
	}
	for index := range productResult {
		products = append(products, &productResult[index])
	}
	return products, nil
}

func (p *ProductManager) SubProduct(productId int64, subNum int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	productResult := &datamodels.Product{}
	if err := p.sqlConn.Where("ID=?", productId).Take(productResult).Error; err != nil {
		return err
	}

	if err := p.sqlConn.Model(productResult).Update("productNum", productResult.ProductNum-subNum).Error; err != nil {
		return err
	}
	return nil
}
