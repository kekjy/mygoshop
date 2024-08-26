package repositories

import (
	"log"
	"mygoshop/datamodels"
	"mygoshop/db"

	"gorm.io/gorm"
)

type IOrder interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithMap() (map[int]map[string]string, error)
}

type OrderManager struct {
	sqlConn *gorm.DB
}

func NewOrderManager(db *gorm.DB) *OrderManager {
	return &OrderManager{
		sqlConn: db,
	}
}

func (o *OrderManager) Conn() error {
	if o.sqlConn == nil {
		__sqlConn, err := db.NewDbConn()
		if err != nil {
			return err
		}
		o.sqlConn = __sqlConn
	}
	return nil
}

func (o *OrderManager) Insert(order *datamodels.Order) (orderId int64, err error) {
	if err = o.Conn(); err != nil {
		return -1, err
	}

	err = o.sqlConn.Create(order).Error
	if err != nil {
		log.Println("insert fail : ", err)
	}

	return orderId, nil
}

func (o *OrderManager) Delete(orderId int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	err := o.sqlConn.Where("id=?", orderId).Delete(&datamodels.Order{})
	if err != nil {
		log.Println("delete order fail : ", err)
		return false
	}

	return true
}
func (o *OrderManager) Update(order *datamodels.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}

	err := o.sqlConn.Save(order).Error
	if err != nil {
		log.Println("update order fail", err)
		return err
	}

	return nil
}

func (o *OrderManager) SelectByKey(id int64) (*datamodels.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	order := &datamodels.Order{}
	err := o.sqlConn.First(&order, id).Error
	if err != nil {
		log.Println("Select By Key fail", err)
		return nil, err
	}
	return order, nil
}

func (o *OrderManager) SelectAll() (order []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}

	err = o.sqlConn.Find(order).Error
	if err != nil {
		log.Println("Select All fail", err)
		return nil, err
	}
	return order, nil
}

func (o *OrderManager) SelectAllWithMap() (OrderMap map[int]map[string]string, err error) {
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	rows, err := o.sqlConn.Model(&datamodels.Order{}).Select("orders.ID, users.userName, products.productName, orders.orderStatus").Joins("left join products on products.ID = orders.productID").Joins("left join users on orders.userID = users.ID").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := db.GetAllResult(rows)
	return result, nil
}
