package pg

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func SendOrder(order Order, db *sql.DB) error {
	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO orders (order_uid,track_number,entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Order_uid, order.Track_number, order.Entry, order.Locale, order.Internal_signature, order.Customer_id, order.Delivery_service, order.Shardkey, order.Sm_id, order.Date_created, order.Oof_shard)
	if err != nil {
		tx.Rollback() // Откат транзакции в случае ошибки
		return err
	}
	_, err = db.Exec(`INSERT INTO delivery (order_uid, name, phone, zip, city,address,region,email) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		order.Order_uid, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback() // Откат транзакции в случае ошибки
		return err
	}
	_, err = db.Exec(`INSERT INTO payment (order_uid,transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		order.Order_uid, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback() // Откат транзакции в случае ошибки
		return err
	}
	for _, item := range order.Items {
		_, err = db.Exec(`INSERT INTO items (order_uid,chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
			order.Order_uid, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
		if err != nil {
			tx.Rollback() // Откат транзакции в случае ошибки
			return err
		}
	}

	// Подтверждение транзакции
	err = tx.Commit()
	if err != nil {
		return err
	}
	log.Println("Data stored successfully")
	return nil
}

func GetOrder(orderID string) (*Order, error) {
	//connect
	connStr := "user=postgres password=qweQWE dbname=WB_Tech_level_0 sslmode=disable port=5555"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	//ping
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Выполнение запроса SELECT для получения информации о заказе
	row := db.QueryRow(`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1`, orderID)
	order := &Order{}
	err = row.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id, &order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
	if err != nil {
		return nil, err
	}

	// Выполнение запроса SELECT для получения информации о доставке заказа
	row = db.QueryRow(`SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1`, orderID)
	delivery := &Delivery{}
	err = row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	// Выполнение запроса SELECT для получения информации о платеже заказа
	row = db.QueryRow(`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1`, orderID)
	payment := &Payment{}
	err = row.Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	// Выполнение запроса SELECT для получения информации о товарах заказа
	rows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale, size, total_prise, nm_id, brand, status FROM items WHERE order_uid = $1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := Item{}
		err = rows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}
