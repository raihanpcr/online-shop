package handler

import (
	"database/sql"
	"log"
	"math/rand"
	"onlineshop/model"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckoutOrder (db *sql.DB) gin.HandlerFunc{
	return func (c *gin.Context)  {

		// todo : buat kata sandi
		
		// todo : buat order dan detail

		// todo : ambil pesanan dari request body
		var CheckoutOrder model.Checkout
		if err := c.BindJSON(&CheckoutOrder); err != nil {
			log.Printf("Terjadi Kesalahan saat membaca request body: %v\n", err)
			c.JSON(400, gin.H{"error":"Data Product Tidak Valid"})
			return
		}

		ids := []string{}
		orderQty := make(map[string]int32)
		for _, o := range CheckoutOrder.Product {
			ids = append(ids, o.ID)
			orderQty[o.ID] = int32(o.Quantity)
		}

		log.Println("ids : ",ids)
		// todo : ambil product dari db
		products, err := model.SelectProductIn(db, ids)
		if err != nil {
			_, file, line, _ := runtime.Caller(1)
			log.Printf("Terjadi Kesalahan saat Mengambil Product: %v (file: %s, line: %d)", err, file, line)
			c.JSON(500, gin.H{"error":"Terjadi Kesalahan Pada Server"})
			return
		}

		c.JSON(200, products)

		// todo : hast kata sandi
		passcode := GeneratePasscode(5)

		hashcode, err := bcrypt.GenerateFromPassword([]byte(passcode), 10)

		if err != nil {
			log.Printf("Terjadi Kesalahan saat Mengambil Hash: %v (file: %s, line: %d)", err)
			c.JSON(500, gin.H{"error":"Terjadi Kesalahan Pada Server"})
			return
		}

		hashcodeString := string(hashcode)

		// todo : buat order & detail
		order := model.Order{
			ID : uuid.New().String(),
			Email : CheckoutOrder.Email,
			Address: CheckoutOrder.Address,
			Passcode: &hashcodeString,
			GrandTotal: 0,
		}

		details := []model.OrderDetail{}

		for _, p := range products{
			total := p.Price * int64(orderQty[p.ID])

			detail := model.OrderDetail{
				ID : uuid.New().String(),
				OrderID: order.ID,
				ProductID: p.ID,
				Quantity: orderQty[p.ID],
				Price: p.Price,
				Total: total,
			}

			details = append(details, detail)
			order.GrandTotal += total
		}

		model.CreateOrder(db, order, details)

		// orderWithDetail := model.OrderWithDetail{
		// 	Order: order,
		// 	Detail: details,
		// }

		

		// response := map[string]interface{}{
		// 	"status":  "success",
		// 	"message": "Order created successfully",
		// 	"data": orderWithDetail,
		// }

		c.JSON(200, gin.H{
			"status": "success",
			"message": "Order created successfully",
			"data": model.OrderWithDetail{
				Order: order,
				Detail: details,
			},
		})
	}
}

func GeneratePasscode(length int ) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	randomGenerator := rand.New(rand.NewSource(time.Now().UnixNano())) //untuk menghasilkan angka angka acak

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[randomGenerator.Intn(len(charset))]
	}

	return string(code)
}

func ConfirmOrder (db *sql.DB) gin.HandlerFunc{
	return func (ctx *gin.Context)  {
		
	}
}
func GetOrder (db *sql.DB) gin.HandlerFunc{
	return func (ctx *gin.Context)  {
		
	}
}