package main

import (
	"WB_Tech_level_0/pg"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"sync"
)

var db *sql.DB
var orderCache *OrderCache
var nc *nats.Conn

type OrderCache struct {
	mu    sync.RWMutex
	cache map[string]*pg.Order
}

func main() {
	// Инициализируем соединение с БД и кэш заказов
	initDB()
	initCache()
	initNats()
	nc.Subscribe("subject", func(msg *nats.Msg) {
		// Print message data
		var newOrder pg.Order
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Fatal(err)
		}
		err = pg.SendOrder(newOrder, db)
		if err != nil {
			log.Fatal(err)
		}
		saveOrderToCache(&newOrder)
	})
	// Запускаем HTTP-сервер
	http.HandleFunc("/getOrder", getOrderHandler)
	http.ListenAndServe(":3333", nil)

}

func initDB() {
	// Инициализируем соединение с БД
	connStr := "user=postgres password=qweQWE dbname=WB_Tech_level_0 sslmode=disable port=5555"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	// Проверяем соединение с БД
	err = db.Ping()
	if err != nil {
		panic(err)
	}
}

func initNats() {
	var err error
	nc, err = nats.Connect("127.0.0.1:4222")
	if err != nil {
		panic(err)
	}
	//defer nc.Close()
}

func initCache() {
	orderCache = &OrderCache{
		cache: make(map[string]*pg.Order),
	}
	// Загружаем данные из БД в кэш
	err := loadOrdersIntoCache()
	if err != nil {
		fmt.Println(err)
	}

}

func loadOrdersIntoCache() error {
	// Выполнение SQL-запроса для получения всех заказов
	rows, err := db.Query("SELECT * FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order pg.Order
		var delivery pg.Delivery
		var payment pg.Payment
		var items []pg.Item

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
			var item pg.Item
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

		saveOrderToCache(&order)
	}

	fmt.Println("загружено в кеш!")
	return nil
}
func saveOrderToCache(order *pg.Order) {
	orderCache.mu.Lock()
	orderCache.cache[order.Order_uid] = order
	orderCache.mu.Unlock()
}
func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("id")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	orderCache.mu.RLock()
	order, found := orderCache.cache[orderID]
	orderCache.mu.RUnlock()

	if !found {
		// Если заказ не найден в кэше, можно вернуть сообщение об ошибке
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Возвращаем заказ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
