package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Permissions   PermissionModel
	Tokens        TokenModel
	Users         UserModel
	ShopifyOrders ShopifyOrderModel
	RimanOrders   RimanOrderModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Permissions:   PermissionModel{DB: db},
		Tokens:        TokenModel{DB: db},
		Users:         UserModel{DB: db},
		ShopifyOrders: ShopifyOrderModel{DB: db},
		RimanOrders:   RimanOrderModel{DB: db},
	}
}
