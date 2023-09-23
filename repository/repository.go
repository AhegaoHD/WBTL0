package repository

import (
	"WB_Tech_level_0/infrastructure"
	"WB_Tech_level_0/model"
	"database/sql"
	"fmt"
	"log"
)

func LoadOrdersIntoCache(db *sql.DB, orderCache *infrastructure.OrderCache) error {
	// Выполнение SQL-запроса для получения всех заказов
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		var delivery model.Delivery
		var payment model.Payment
		var items []model.Item

		// Сканирование данных из результатов запроса в структуры
		err := rows.Scan(
			&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature, &order.Customer_id,
			&order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard,
		)
		if err != nil {
			return err
		}

		// Выполнение SQL-запроса для получения доставки по каждому заказу
		err = db.QueryRow("SELECT * FROM delivery WHERE order_uid = $1", order.Order_uid).Scan(
			&order.Order_uid, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
		)
		if err != nil {
			return err
		}

		// Выполнение SQL-запроса для получения оплаты по каждому заказу
		err = db.QueryRow("SELECT * FROM payment WHERE order_uid = $1", order.Order_uid).Scan(&order.Order_uid,
			&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount,
			&payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		)
		if err != nil {
			return err
		}

		// Выполнение SQL-запроса для получения всех товаров по каждому заказу
		itemRows, err := db.Query("SELECT * FROM items WHERE order_uid = $1", order.Order_uid)
		if err != nil {
			return err
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item model.Item
			// Сканирование данных из результатов запроса в структуру
			err := itemRows.Scan(&order.Order_uid,
				&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
			)
			if err != nil {
				return err
			}
			items = append(items, item)
		}

		// Наполняем структуру заказа полученными данными
		order.Delivery = delivery
		order.Payment = payment
		order.Items = items

		SaveOrderToCache(orderCache, &order)
	}

	fmt.Println("загружено в кеш!")
	return nil
}
func SaveOrderToCache(orderCache *infrastructure.OrderCache, order *model.Order) {
	orderCache.Set(order.Order_uid, order)
}
func SendOrder(order model.Order, db *sql.DB) error {
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
