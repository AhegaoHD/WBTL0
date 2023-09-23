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
