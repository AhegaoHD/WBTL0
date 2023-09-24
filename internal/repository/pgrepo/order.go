package pgrepo

import (
	"context"
	"errors"
	"github.com/AhegaoHD/WBTL0/internal/model"
	"github.com/AhegaoHD/WBTL0/pkg/postgres"
)

type OrderRepo struct {
	*postgres.Postgres
}

func NewOrderRepo(db *postgres.Postgres) *OrderRepo {
	return &OrderRepo{db}
}

func (o OrderRepo) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	// Выполнение SQL-запроса для получения всех заказов
	rows, err := o.Pool.Query(ctx, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		// Сканирование данных из результатов запроса в структуры
		err := rows.Scan(
			&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id,
			&order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard,
		)
		if err != nil {
			return nil, err
		}
		models = append(models, order)
	}

	return models, nil
}

func (o OrderRepo) GetDelivery(ctx context.Context, orderId string) (model.Delivery, error) {
	var delivery model.Delivery
	err := o.Pool.QueryRow(ctx, "SELECT name,phone,zip,city,address,region,email FROM delivery WHERE order_uid = $1", orderId).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
	)
	if err != nil {
		return delivery, err
	}

	return delivery, nil
}

func (o OrderRepo) GetPayment(ctx context.Context, orderId string) (model.Payment, error) {
	var payment model.Payment
	err := o.Pool.QueryRow(ctx, "SELECT transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee FROM payment WHERE order_uid = $1", orderId).Scan(
		&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount,
		&payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (o OrderRepo) GetItems(ctx context.Context, orderId string) ([]model.Item, error) {
	items := make([]model.Item, 0)
	itemRows, err := o.Pool.Query(ctx, "SELECT chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status FROM items WHERE order_uid = $1", orderId)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item model.Item
		// Сканирование данных из результатов запроса в структуру
		err := itemRows.Scan(
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (o OrderRepo) GetOrdersWithDetails(ctx context.Context) ([]model.Order, error) {
	rows, err := o.Pool.Query(ctx, `
        SELECT 
            o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
            o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
            d.name as delivery_name, d.phone as delivery_phone, d.zip as delivery_zip, d.city as delivery_city, d.address as delivery_address, d.region as delivery_region, d.email as delivery_email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
            i.chrt_id, i.track_number as item_track_number, i.price as item_price, i.rid as item_rid, i.name as item_name, i.sale as item_sale, i.size as item_size, i.total_price as item_total_price, i.nm_id as item_nm_id, i.brand as item_brand, i.status as item_status
        FROM orders o
        JOIN delivery d ON o.order_uid = d.order_uid
        JOIN payment p ON o.order_uid = p.order_uid
        JOIN items i ON o.order_uid = i.order_uid
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersMap = make(map[string]model.Order)

	for rows.Next() {
		var orderUid string
		var order model.Order
		var item model.Item

		err := rows.Scan(
			&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id,
			&order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}

		if existingOrder, ok := ordersMap[orderUid]; ok {
			existingOrder.Items = append(existingOrder.Items, item)
			ordersMap[orderUid] = existingOrder
		} else {
			order.Items = append(order.Items, item)
			ordersMap[orderUid] = order
		}
	}

	var orders []model.Order
	for _, order := range ordersMap {
		orders = append(orders, order)
	}

	return orders, nil
}

func (o OrderRepo) GetOrderWithDetailsByUid(ctx context.Context, orderUID string) (model.Order, error) {
	rows, err := o.Pool.Query(ctx, `
        SELECT 
            o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
            o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
            d.name as delivery_name, d.phone as delivery_phone, d.zip as delivery_zip, d.city as delivery_city, d.address as delivery_address, d.region as delivery_region, d.email as delivery_email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
            i.chrt_id, i.track_number as item_track_number, i.price as item_price, i.rid as item_rid, i.name as item_name, i.sale as item_sale, i.size as item_size, i.total_price as item_total_price, i.nm_id as item_nm_id, i.brand as item_brand, i.status as item_status
        FROM orders o
        JOIN delivery d ON o.order_uid = d.order_uid
        JOIN payment p ON o.order_uid = p.order_uid
        JOIN items i ON o.order_uid = i.order_uid
        WHERE o.order_uid = $1
    `, orderUID)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var ordersMap = make(map[string]model.Order)

	for rows.Next() {
		var orderUid string
		var order model.Order
		var item model.Item

		err := rows.Scan(
			&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id,
			&order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)
		if err != nil {
			return model.Order{}, err
		}

		if existingOrder, ok := ordersMap[orderUid]; ok {
			existingOrder.Items = append(existingOrder.Items, item)
			ordersMap[orderUid] = existingOrder
		} else {
			order.Items = append(order.Items, item)
			ordersMap[orderUid] = order
		}
	}

	var orders model.Order
	if len(ordersMap) == 0 {
		return model.Order{}, errors.New("zero")
	}
	for _, order := range ordersMap {
		orders = order
	}
	return orders, nil
}

func (o OrderRepo) SetOrdersWithDetails(ctx context.Context, order *model.Order) error {
	tx, err := o.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	_, err = o.Pool.Exec(ctx, `INSERT INTO orders (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		return err
	}
	_, err = o.Pool.Exec(ctx, `INSERT INTO delivery (order_uid, name, phone, zip, city,address,region,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}
	_, err = o.Pool.Exec(ctx, `INSERT INTO payment (order_uid,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Order_uid, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		_, err = o.Pool.Exec(ctx, `INSERT INTO items (order_uid,chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			order.Order_uid, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	// Подтверждение транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
