package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrDuplicateRimanMainOrdersPK = errors.New("duplicate main_orders_pk")
)

// RimanOrder persists the entire riman.RimanOrder struct along with metadata.
// JSON-like fields in riman.RimanOrder are stored as jsonb in the database.
// We marshal/unmarshal those interface{} fields when writing/reading.
//
// We embed riman.RimanOrder to keep the same shape for consumers.
// The unique natural key is MainOrdersPK.
type RimanOrder struct {
	ID                      int64       `json:"id"`
	CreatedAt               time.Time   `json:"created_at"`
	OrderDate               string      `json:"orderDate"`
	MainOrdersPK            int         `json:"mainOrdersPK"`
	OrderType               string      `json:"orderType"`
	FinalOrderType          interface{} `json:"finalOrderType"`
	SiteURL                 string      `json:"siteURL"`
	EncOrderNumber          string      `json:"encOrderNumber"`
	CurrencySymbol          string      `json:"currencySymbol"`
	CurrencyCode            string      `json:"currencyCode"`
	PaidStatus              bool        `json:"paidStatus"`
	HasTaxInvoice           bool        `json:"hasTaxInvoice"`
	HasCommercialInvoice    bool        `json:"hasCommercialInvoice"`
	HasCreditNote           bool        `json:"hasCreditNote"`
	IsShippingPending       bool        `json:"isShippingPending"`
	IsPB                    bool        `json:"isPB"`
	IsPA                    bool        `json:"isPA"`
	IsCC                    bool        `json:"isCC"`
	MainFK                  int         `json:"mainFK"`
	MainOrderTypeFK         int         `json:"mainOrderTypeFK"`
	VoucherURL              interface{} `json:"voucherURL"`
	ShipCountry             string      `json:"shipCountry"`
	ShippingStatus          string      `json:"shippingStatus"`
	OrderShippingStatus     string      `json:"orderShippingStatus"`
	OrderTypeValue          interface{} `json:"orderTypeValue"`
	PaidStatusValue         string      `json:"paidStatusValue"`
	Quantity                int         `json:"quantity"`
	Email                   interface{} `json:"email"`
	Phone                   interface{} `json:"phone"`
	ShipFirstName           interface{} `json:"shipFirstName"`
	ShipLastName            interface{} `json:"shipLastName"`
	MarkedPaidDate          string      `json:"markedPaidDate"`
	Total                   float64     `json:"total"`
	ConvTotal               float64     `json:"convTotal"`
	ConvTotalFormat         string      `json:"convTotalFormat"`
	SubTotal                int         `json:"subTotal"`
	ConvSubtotal            int         `json:"convSubtotal"`
	ShipTotal               float64     `json:"shipTotal"`
	ConvShipTotal           float64     `json:"convShipTotal"`
	Taxes                   float64     `json:"taxes"`
	TaxLabel                string      `json:"taxLabel"`
	ProductTax              float64     `json:"productTax"`
	ShippingTax             float64     `json:"shippingTax"`
	TotalProductTax         float64     `json:"totalProductTax"`
	AdditionalTaxLabel      string      `json:"additionalTaxLabel"`
	AdditionalTax           interface{} `json:"additionalTax"`
	ConvTaxes               float64     `json:"convTaxes"`
	OrderProcessingFees     interface{} `json:"orderProcessingFees"`
	ConvOrderProcessingFees interface{} `json:"convOrderProcessingFees"`
	Discount                int         `json:"discount"`
	ConvDiscount            int         `json:"convDiscount"`
	RefundAmount            int         `json:"refundAmount"`
	ConvRefund              int         `json:"convRefund"`
	SalesCampaignFK         interface{} `json:"salesCampaignFK"`
	Paidstatusfk            int         `json:"paidstatusfk"`
	DeliveryDate            interface{} `json:"deliveryDate"`
	ShippingDetails         interface{} `json:"shippingDetails"`
	OrderItems              interface{} `json:"orderItems"`
	Payments                interface{} `json:"payments"`
	IsPrepaidOrder          interface{} `json:"isPrepaidOrder"`
	ShowInvoice             bool        `json:"showInvoice"`
	ShowOrderInvoice        bool        `json:"showOrderInvoice"`
	KrGuaranteeNo           string      `json:"krGuaranteeNo"`
	WeChatOrderNumber       interface{} `json:"weChatOrderNumber"`
	MemberID                interface{} `json:"memberID"`
	Version                 int         `json:"-"`
}

type RimanOrderModel struct {
	DB *sql.DB
}

func toJSONB(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	return json.Marshal(v)
}

func fromJSONB(data []byte, dest *interface{}) error {
	if dest == nil {
		return nil
	}
	if data == nil {
		*dest = nil
		return nil
	}
	var any interface{}
	if err := json.Unmarshal(data, &any); err != nil {
		return err
	}
	*dest = any
	return nil
}

// Insert creates a new riman_orders row, storing all fields of riman.RimanOrder.
func (m RimanOrderModel) Insert(o *RimanOrder) error {
	query := `
        INSERT INTO riman_orders (
            main_orders_pk, order_date, order_type, final_order_type, site_url, enc_order_number,
            currency_symbol, currency_code, paid_status, has_tax_invoice, has_commercial_invoice,
            has_credit_note, is_shipping_pending, is_pb, is_pa, is_cc, main_fk, main_order_type_fk,
            voucher_url, ship_country, shipping_status, order_shipping_status, order_type_value,
            paid_status_value, quantity, email, phone, ship_first_name, ship_last_name,
            marked_paid_date, total, conv_total, conv_total_format, sub_total, conv_subtotal,
            ship_total, conv_ship_total, taxes, tax_label, product_tax, shipping_tax,
            total_product_tax, additional_tax_label, additional_tax, conv_taxes,
            order_processing_fees, conv_order_processing_fees, discount, conv_discount,
            refund_amount, conv_refund, sales_campaign_fk, paidstatusfk, delivery_date,
            shipping_details, order_items, payments, is_prepaid_order, show_invoice,
            show_order_invoice, kr_guarantee_no, we_chat_order_number, member_id
        ) VALUES (
            $1,$2,$3,$4,$5,$6,
            $7,$8,$9,$10,$11,
            $12,$13,$14,$15,$16,$17,$18,
            $19,$20,$21,$22,$23,
            $24,$25,$26,$27,$28,$29,
            $30,$31,$32,$33,$34,$35,
            $36,$37,$38,$39,$40,$41,
            $42,$43,$44,$45,
            $46,$47,$48,$49,
            $50,$51,$52,$53,$54,
            $55,$56,$57,$58,$59,
            $60,$61,$62,$63
        )
        RETURNING id, created_at, version`

	finalOrderType, err := toJSONB(o.FinalOrderType)
	if err != nil {
		return err
	}
	voucherURL, err := toJSONB(o.VoucherURL)
	if err != nil {
		return err
	}
	orderTypeValue, err := toJSONB(o.OrderTypeValue)
	if err != nil {
		return err
	}
	email, err := toJSONB(o.Email)
	if err != nil {
		return err
	}
	phone, err := toJSONB(o.Phone)
	if err != nil {
		return err
	}
	shipFirstName, err := toJSONB(o.ShipFirstName)
	if err != nil {
		return err
	}
	shipLastName, err := toJSONB(o.ShipLastName)
	if err != nil {
		return err
	}
	additionalTax, err := toJSONB(o.AdditionalTax)
	if err != nil {
		return err
	}
	orderProcessingFees, err := toJSONB(o.OrderProcessingFees)
	if err != nil {
		return err
	}
	convOrderProcessingFees, err := toJSONB(o.ConvOrderProcessingFees)
	if err != nil {
		return err
	}
	salesCampaignFK, err := toJSONB(o.SalesCampaignFK)
	if err != nil {
		return err
	}
	deliveryDate, err := toJSONB(o.DeliveryDate)
	if err != nil {
		return err
	}
	shippingDetails, err := toJSONB(o.ShippingDetails)
	if err != nil {
		return err
	}
	orderItems, err := toJSONB(o.OrderItems)
	if err != nil {
		return err
	}
	payments, err := toJSONB(o.Payments)
	if err != nil {
		return err
	}
	isPrepaidOrder, err := toJSONB(o.IsPrepaidOrder)
	if err != nil {
		return err
	}
	weChatOrderNumber, err := toJSONB(o.WeChatOrderNumber)
	if err != nil {
		return err
	}
	memberID, err := toJSONB(o.MemberID)
	if err != nil {
		return err
	}

	args := []any{
		o.MainOrdersPK, o.OrderDate, o.OrderType, finalOrderType, o.SiteURL, o.EncOrderNumber,
		o.CurrencySymbol, o.CurrencyCode, o.PaidStatus, o.HasTaxInvoice, o.HasCommercialInvoice,
		o.HasCreditNote, o.IsShippingPending, o.IsPB, o.IsPA, o.IsCC, o.MainFK, o.MainOrderTypeFK,
		voucherURL, o.ShipCountry, o.ShippingStatus, o.OrderShippingStatus, orderTypeValue,
		o.PaidStatusValue, o.Quantity, email, phone, shipFirstName, shipLastName,
		o.MarkedPaidDate, o.Total, o.ConvTotal, o.ConvTotalFormat, o.SubTotal, o.ConvSubtotal,
		o.ShipTotal, o.ConvShipTotal, o.Taxes, o.TaxLabel, o.ProductTax, o.ShippingTax,
		o.TotalProductTax, o.AdditionalTaxLabel, additionalTax, o.ConvTaxes,
		orderProcessingFees, convOrderProcessingFees, o.Discount, o.ConvDiscount,
		o.RefundAmount, o.ConvRefund, salesCampaignFK, o.Paidstatusfk, deliveryDate,
		shippingDetails, orderItems, payments, isPrepaidOrder, o.ShowInvoice,
		o.ShowOrderInvoice, o.KrGuaranteeNo, weChatOrderNumber, memberID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&o.ID, &o.CreatedAt, &o.Version); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "riman_orders_main_orders_pk_key"`:
			return ErrDuplicateRimanMainOrdersPK
		default:
			return err
		}
	}
	return nil
}

// GetByID fetches a RimanOrder by internal ID.
func (m RimanOrderModel) GetByID(id int64) (*RimanOrder, error) {
	query := `
        SELECT 
            id, created_at,
            main_orders_pk, order_date, order_type, final_order_type, site_url, enc_order_number,
            currency_symbol, currency_code, paid_status, has_tax_invoice, has_commercial_invoice,
            has_credit_note, is_shipping_pending, is_pb, is_pa, is_cc, main_fk, main_order_type_fk,
            voucher_url, ship_country, shipping_status, order_shipping_status, order_type_value,
            paid_status_value, quantity, email, phone, ship_first_name, ship_last_name,
            marked_paid_date, total, conv_total, conv_total_format, sub_total, conv_subtotal,
            ship_total, conv_ship_total, taxes, tax_label, product_tax, shipping_tax,
            total_product_tax, additional_tax_label, additional_tax, conv_taxes,
            order_processing_fees, conv_order_processing_fees, discount, conv_discount,
            refund_amount, conv_refund, sales_campaign_fk, paidstatusfk, delivery_date,
            shipping_details, order_items, payments, is_prepaid_order, show_invoice,
            show_order_invoice, kr_guarantee_no, we_chat_order_number, member_id,
            version
        FROM riman_orders
        WHERE id = $1`

	var o RimanOrder

	var finalOrderType, voucherURL, orderTypeValue, email, phone, shipFirstName, shipLastName []byte
	var additionalTax, orderProcessingFees, convOrderProcessingFees []byte
	var salesCampaignFK, deliveryDate, shippingDetails, orderItems, payments, isPrepaidOrder []byte
	var weChatOrderNumber, memberID []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.CreatedAt,
		&o.MainOrdersPK, &o.OrderDate, &o.OrderType, &finalOrderType, &o.SiteURL, &o.EncOrderNumber,
		&o.CurrencySymbol, &o.CurrencyCode, &o.PaidStatus, &o.HasTaxInvoice, &o.HasCommercialInvoice,
		&o.HasCreditNote, &o.IsShippingPending, &o.IsPB, &o.IsPA, &o.IsCC, &o.MainFK, &o.MainOrderTypeFK,
		&voucherURL, &o.ShipCountry, &o.ShippingStatus, &o.OrderShippingStatus, &orderTypeValue,
		&o.PaidStatusValue, &o.Quantity, &email, &phone, &shipFirstName, &shipLastName,
		&o.MarkedPaidDate, &o.Total, &o.ConvTotal, &o.ConvTotalFormat, &o.SubTotal, &o.ConvSubtotal,
		&o.ShipTotal, &o.ConvShipTotal, &o.Taxes, &o.TaxLabel, &o.ProductTax, &o.ShippingTax,
		&o.TotalProductTax, &o.AdditionalTaxLabel, &additionalTax, &o.ConvTaxes,
		&orderProcessingFees, &convOrderProcessingFees, &o.Discount, &o.ConvDiscount,
		&o.RefundAmount, &o.ConvRefund, &salesCampaignFK, &o.Paidstatusfk, &deliveryDate,
		&shippingDetails, &orderItems, &payments, &isPrepaidOrder, &o.ShowInvoice,
		&o.ShowOrderInvoice, &o.KrGuaranteeNo, &weChatOrderNumber, &memberID,
		&o.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Unmarshal JSONB fields back into interface{} fields
	if err := fromJSONB(finalOrderType, &o.FinalOrderType); err != nil {
		return nil, err
	}
	if err := fromJSONB(voucherURL, &o.VoucherURL); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderTypeValue, &o.OrderTypeValue); err != nil {
		return nil, err
	}
	if err := fromJSONB(email, &o.Email); err != nil {
		return nil, err
	}
	if err := fromJSONB(phone, &o.Phone); err != nil {
		return nil, err
	}
	if err := fromJSONB(shipFirstName, &o.ShipFirstName); err != nil {
		return nil, err
	}
	if err := fromJSONB(shipLastName, &o.ShipLastName); err != nil {
		return nil, err
	}
	if err := fromJSONB(additionalTax, &o.AdditionalTax); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderProcessingFees, &o.OrderProcessingFees); err != nil {
		return nil, err
	}
	if err := fromJSONB(convOrderProcessingFees, &o.ConvOrderProcessingFees); err != nil {
		return nil, err
	}
	if err := fromJSONB(salesCampaignFK, &o.SalesCampaignFK); err != nil {
		return nil, err
	}
	if err := fromJSONB(deliveryDate, &o.DeliveryDate); err != nil {
		return nil, err
	}
	if err := fromJSONB(shippingDetails, &o.ShippingDetails); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderItems, &o.OrderItems); err != nil {
		return nil, err
	}
	if err := fromJSONB(payments, &o.Payments); err != nil {
		return nil, err
	}
	if err := fromJSONB(isPrepaidOrder, &o.IsPrepaidOrder); err != nil {
		return nil, err
	}
	if err := fromJSONB(weChatOrderNumber, &o.WeChatOrderNumber); err != nil {
		return nil, err
	}
	if err := fromJSONB(memberID, &o.MemberID); err != nil {
		return nil, err
	}

	return &o, nil
}

// GetByMainOrdersPK fetches a RimanOrder by the unique external key.
func (m RimanOrderModel) GetByMainOrdersPK(mainOrdersPK int) (*RimanOrder, error) {
	query := `
        SELECT 
            id, created_at,
            main_orders_pk, order_date, order_type, final_order_type, site_url, enc_order_number,
            currency_symbol, currency_code, paid_status, has_tax_invoice, has_commercial_invoice,
            has_credit_note, is_shipping_pending, is_pb, is_pa, is_cc, main_fk, main_order_type_fk,
            voucher_url, ship_country, shipping_status, order_shipping_status, order_type_value,
            paid_status_value, quantity, email, phone, ship_first_name, ship_last_name,
            marked_paid_date, total, conv_total, conv_total_format, sub_total, conv_subtotal,
            ship_total, conv_ship_total, taxes, tax_label, product_tax, shipping_tax,
            total_product_tax, additional_tax_label, additional_tax, conv_taxes,
            order_processing_fees, conv_order_processing_fees, discount, conv_discount,
            refund_amount, conv_refund, sales_campaign_fk, paidstatusfk, delivery_date,
            shipping_details, order_items, payments, is_prepaid_order, show_invoice,
            show_order_invoice, kr_guarantee_no, we_chat_order_number, member_id,
            version
        FROM riman_orders
        WHERE main_orders_pk = $1`

	var o RimanOrder

	var finalOrderType, voucherURL, orderTypeValue, email, phone, shipFirstName, shipLastName []byte
	var additionalTax, orderProcessingFees, convOrderProcessingFees []byte
	var salesCampaignFK, deliveryDate, shippingDetails, orderItems, payments, isPrepaidOrder []byte
	var weChatOrderNumber, memberID []byte

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, mainOrdersPK).Scan(
		&o.ID, &o.CreatedAt,
		&o.MainOrdersPK, &o.OrderDate, &o.OrderType, &finalOrderType, &o.SiteURL, &o.EncOrderNumber,
		&o.CurrencySymbol, &o.CurrencyCode, &o.PaidStatus, &o.HasTaxInvoice, &o.HasCommercialInvoice,
		&o.HasCreditNote, &o.IsShippingPending, &o.IsPB, &o.IsPA, &o.IsCC, &o.MainFK, &o.MainOrderTypeFK,
		&voucherURL, &o.ShipCountry, &o.ShippingStatus, &o.OrderShippingStatus, &orderTypeValue,
		&o.PaidStatusValue, &o.Quantity, &email, &phone, &shipFirstName, &shipLastName,
		&o.MarkedPaidDate, &o.Total, &o.ConvTotal, &o.ConvTotalFormat, &o.SubTotal, &o.ConvSubtotal,
		&o.ShipTotal, &o.ConvShipTotal, &o.Taxes, &o.TaxLabel, &o.ProductTax, &o.ShippingTax,
		&o.TotalProductTax, &o.AdditionalTaxLabel, &additionalTax, &o.ConvTaxes,
		&orderProcessingFees, &convOrderProcessingFees, &o.Discount, &o.ConvDiscount,
		&o.RefundAmount, &o.ConvRefund, &salesCampaignFK, &o.Paidstatusfk, &deliveryDate,
		&shippingDetails, &orderItems, &payments, &isPrepaidOrder, &o.ShowInvoice,
		&o.ShowOrderInvoice, &o.KrGuaranteeNo, &weChatOrderNumber, &memberID,
		&o.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if err := fromJSONB(finalOrderType, &o.FinalOrderType); err != nil {
		return nil, err
	}
	if err := fromJSONB(voucherURL, &o.VoucherURL); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderTypeValue, &o.OrderTypeValue); err != nil {
		return nil, err
	}
	if err := fromJSONB(email, &o.Email); err != nil {
		return nil, err
	}
	if err := fromJSONB(phone, &o.Phone); err != nil {
		return nil, err
	}
	if err := fromJSONB(shipFirstName, &o.ShipFirstName); err != nil {
		return nil, err
	}
	if err := fromJSONB(shipLastName, &o.ShipLastName); err != nil {
		return nil, err
	}
	if err := fromJSONB(additionalTax, &o.AdditionalTax); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderProcessingFees, &o.OrderProcessingFees); err != nil {
		return nil, err
	}
	if err := fromJSONB(convOrderProcessingFees, &o.ConvOrderProcessingFees); err != nil {
		return nil, err
	}
	if err := fromJSONB(salesCampaignFK, &o.SalesCampaignFK); err != nil {
		return nil, err
	}
	if err := fromJSONB(deliveryDate, &o.DeliveryDate); err != nil {
		return nil, err
	}
	if err := fromJSONB(shippingDetails, &o.ShippingDetails); err != nil {
		return nil, err
	}
	if err := fromJSONB(orderItems, &o.OrderItems); err != nil {
		return nil, err
	}
	if err := fromJSONB(payments, &o.Payments); err != nil {
		return nil, err
	}
	if err := fromJSONB(isPrepaidOrder, &o.IsPrepaidOrder); err != nil {
		return nil, err
	}
	if err := fromJSONB(weChatOrderNumber, &o.WeChatOrderNumber); err != nil {
		return nil, err
	}
	if err := fromJSONB(memberID, &o.MemberID); err != nil {
		return nil, err
	}

	return &o, nil
}

// Update modifies an existing row using optimistic locking on version.
func (m RimanOrderModel) Update(o *RimanOrder) error {
	query := `
        UPDATE riman_orders SET
            order_date = $1, order_type = $2, final_order_type = $3, site_url = $4, enc_order_number = $5,
            currency_symbol = $6, currency_code = $7, paid_status = $8, has_tax_invoice = $9, has_commercial_invoice = $10,
            has_credit_note = $11, is_shipping_pending = $12, is_pb = $13, is_pa = $14, is_cc = $15, main_fk = $16, main_order_type_fk = $17,
            voucher_url = $18, ship_country = $19, shipping_status = $20, order_shipping_status = $21, order_type_value = $22,
            paid_status_value = $23, quantity = $24, email = $25, phone = $26, ship_first_name = $27, ship_last_name = $28,
            marked_paid_date = $29, total = $30, conv_total = $31, conv_total_format = $32, sub_total = $33, conv_subtotal = $34,
            ship_total = $35, conv_ship_total = $36, taxes = $37, tax_label = $38, product_tax = $39, shipping_tax = $40,
            total_product_tax = $41, additional_tax_label = $42, additional_tax = $43, conv_taxes = $44,
            order_processing_fees = $45, conv_order_processing_fees = $46, discount = $47, conv_discount = $48,
            refund_amount = $49, conv_refund = $50, sales_campaign_fk = $51, paidstatusfk = $52, delivery_date = $53,
            shipping_details = $54, order_items = $55, payments = $56, is_prepaid_order = $57, show_invoice = $58,
            show_order_invoice = $59, kr_guarantee_no = $60, we_chat_order_number = $61, member_id = $62,
            version = version + 1
        WHERE id = $63 AND version = $64
        RETURNING version`

	finalOrderType, err := toJSONB(o.FinalOrderType)
	if err != nil {
		return err
	}
	voucherURL, err := toJSONB(o.VoucherURL)
	if err != nil {
		return err
	}
	orderTypeValue, err := toJSONB(o.OrderTypeValue)
	if err != nil {
		return err
	}
	email, err := toJSONB(o.Email)
	if err != nil {
		return err
	}
	phone, err := toJSONB(o.Phone)
	if err != nil {
		return err
	}
	shipFirstName, err := toJSONB(o.ShipFirstName)
	if err != nil {
		return err
	}
	shipLastName, err := toJSONB(o.ShipLastName)
	if err != nil {
		return err
	}
	additionalTax, err := toJSONB(o.AdditionalTax)
	if err != nil {
		return err
	}
	orderProcessingFees, err := toJSONB(o.OrderProcessingFees)
	if err != nil {
		return err
	}
	convOrderProcessingFees, err := toJSONB(o.ConvOrderProcessingFees)
	if err != nil {
		return err
	}
	salesCampaignFK, err := toJSONB(o.SalesCampaignFK)
	if err != nil {
		return err
	}
	deliveryDate, err := toJSONB(o.DeliveryDate)
	if err != nil {
		return err
	}
	shippingDetails, err := toJSONB(o.ShippingDetails)
	if err != nil {
		return err
	}
	orderItems, err := toJSONB(o.OrderItems)
	if err != nil {
		return err
	}
	payments, err := toJSONB(o.Payments)
	if err != nil {
		return err
	}
	isPrepaidOrder, err := toJSONB(o.IsPrepaidOrder)
	if err != nil {
		return err
	}
	weChatOrderNumber, err := toJSONB(o.WeChatOrderNumber)
	if err != nil {
		return err
	}
	memberID, err := toJSONB(o.MemberID)
	if err != nil {
		return err
	}

	args := []any{
		o.OrderDate, o.OrderType, finalOrderType, o.SiteURL, o.EncOrderNumber,
		o.CurrencySymbol, o.CurrencyCode, o.PaidStatus, o.HasTaxInvoice, o.HasCommercialInvoice,
		o.HasCreditNote, o.IsShippingPending, o.IsPB, o.IsPA, o.IsCC, o.MainFK, o.MainOrderTypeFK,
		voucherURL, o.ShipCountry, o.ShippingStatus, o.OrderShippingStatus, orderTypeValue,
		o.PaidStatusValue, o.Quantity, email, phone, shipFirstName, shipLastName,
		o.MarkedPaidDate, o.Total, o.ConvTotal, o.ConvTotalFormat, o.SubTotal, o.ConvSubtotal,
		o.ShipTotal, o.ConvShipTotal, o.Taxes, o.TaxLabel, o.ProductTax, o.ShippingTax,
		o.TotalProductTax, o.AdditionalTaxLabel, additionalTax, o.ConvTaxes,
		orderProcessingFees, convOrderProcessingFees, o.Discount, o.ConvDiscount,
		o.RefundAmount, o.ConvRefund, salesCampaignFK, o.Paidstatusfk, deliveryDate,
		shippingDetails, orderItems, payments, isPrepaidOrder, o.ShowInvoice,
		o.ShowOrderInvoice, o.KrGuaranteeNo, weChatOrderNumber, memberID,
		o.ID, o.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&o.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete removes a row by ID.
func (m RimanOrderModel) Delete(id int64) error {
	query := `DELETE FROM riman_orders WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}
