package main

import (
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================================
// SQLite Order Store
// ============================================================================

const (
	OrderStatusPending = "pending"
	OrderStatusPaid    = "paid"
	OrderStatusExpired = "expired"
)

// Order represents a payment order
type Order struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	EpayOrderNo    string    `gorm:"uniqueIndex;size:128;not null" json:"epay_order_no"`   // fast-api 的 out_trade_no
	AlipayTradeNo  string    `gorm:"size:128" json:"alipay_trade_no"`                       // 支付宝交易号
	PaymentType    string    `gorm:"size:32;not null" json:"payment_type"`                   // 支付方式 (alipay/wxpay)
	ProductName    string    `gorm:"size:256" json:"product_name"`                           // 商品名称
	Amount         string    `gorm:"size:16;not null" json:"amount"`                         // 金额
	Status         string    `gorm:"size:32;not null;index;default:pending" json:"status"`  // pending/paid/expired
	NotifyURL      string    `gorm:"size:1024;not null" json:"notify_url"`                  // epay 异步回调地址
	ReturnURL      string    `gorm:"size:1024" json:"return_url"`                           // epay 同步跳转地址
	MerchantID     string    `gorm:"size:64;not null" json:"merchant_id"`                   // epay 商户号
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Store struct {
	db *gorm.DB
}

func NewStore(dbPath string) (*Store, error) {
	if dbPath == "" {
		dbPath = "epay_gateway.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Order{}); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) CreateOrder(o *Order) error {
	return s.db.Create(o).Error
}

func (s *Store) GetOrderByEpayNo(outTradeNo string) (*Order, error) {
	var o Order
	err := s.db.Where("epay_order_no = ?", outTradeNo).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *Store) UpdateOrderPaid(outTradeNo, alipayTradeNo string) error {
	return s.db.Model(&Order{}).
		Where("epay_order_no = ? AND status = ?", outTradeNo, OrderStatusPending).
		Updates(map[string]any{
			"status":          OrderStatusPaid,
			"alipay_trade_no": alipayTradeNo,
		}).Error
}

func (s *Store) ExpireOldOrders(before time.Time) error {
	return s.db.Model(&Order{}).
		Where("status = ? AND created_at < ?", OrderStatusPending, before).
		Update("status", OrderStatusExpired).Error
}
