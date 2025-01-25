package handler

import (
	"database/sql"
	"errors"
	"log"
	"onlineshop/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListProducts(db *sql.DB) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		//TODO : ambil dari database
		products, err := model.SelectProduct(db)
		if err != nil {
			log.Printf("Terjadi Kesalahan saat mengambil data product : %v\n", err)
			ctx.JSON(500, gin.H{"error": "Terjadi Kesalahan Pada Server"})
			return
		}

		// struct
		response := gin.H{
			"status":  "success",
			"message": "Produk berhasil diambil",
			"data":    products,
		}
	
		ctx.JSON(200, response)
	}
}

func GetProduct(db *sql.DB) gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		
		// Mengambil parameter id
		id := ctx.Param("id")

		product, err := model.SelectProductByID(db, id)

		if err != nil {

			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("data in id not found %v\n", err)
				ctx.JSON(404, gin.H{"error" : "Product Tidak ditemukan"})
				return
			}

			log.Printf("Terjadi Kesalahan saat mengambil data product %v\n", err)
			ctx.JSON(500, gin.H{"error" : "Terjadi Kesalahan Pada Server"})
			return
		}

		response := gin.H{
			"status" : "success",
			"message" : "Success Get Products",
			"data" : product,
		}

		ctx.JSON(200, response)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc{
	return func (ctx *gin.Context)  {

		var product model.Product

		//handle error request
		if err := ctx.Bind(&product); err != nil {
			log.Printf("Terjadi Kesalahan saat membaca request body: %v\n", err)
			ctx.JSON(400, gin.H{"error" : "Data product tidak valid"})
			return
		}

		product.ID = uuid.New().String()

		if err := model.InsertProduct(db, product); err != nil {
			log.Printf("Terjadi Kesalahan saat membaca request body: %v\n", err)
			ctx.JSON(500, gin.H{"error" : "Data product tidak valid"})
			return
		}

		ctx.JSON(201, product)
	}
}

// todo : update product
func UpdateProduct(db *sql.DB) gin.HandlerFunc{
	return func (ctx *gin.Context)  {

		id := ctx.Param("id")

		var product model.Product

		//handle error request
		if err := ctx.Bind(&product); err != nil {
			log.Printf("Terjadi Kesalahan saat membaca request body: %v\n", err)
			ctx.JSON(400, gin.H{"error" : "Data product tidak valid"})
			return
		}

		productExisting, err := model.SelectProductByID(db, id)

		if err != nil {
			log.Printf("Terjadi Kesalahan saat mengambil product: %v\n", err)
			ctx.JSON(500, gin.H{"error" : "Terjadi Kesalahan Server"})
			return
		}

		if product.Name != "" {
			productExisting.Name = product.Name
		}

		if product.Name != "" {
			productExisting.Price = product.Price
		}

		if err := model.UpdateProduct(db, productExisting); err != nil {
			log.Printf("Terjadi Kesalahan saat memperbarui product: %v\n", err)
			ctx.JSON(500, gin.H{"error" : "Terjadi Kesalahan Server"})
			return
		}

		ctx.JSON(201, productExisting)
	}
}

// todo : delete product
func DeleteProduct(db *sql.DB) gin.HandlerFunc{
	return func (ctx *gin.Context)  {
		id := ctx.Param("id")

		if err := model.DeleteProduct(db, id); err != nil {
			log.Printf("Terjadi Kesalahan saat memperbarui product: %v\n", err)
			ctx.JSON(500, gin.H{"error" : "Terjadi Kesalahan Server"})
			return
		}

		ctx.JSON(201, gin.H{"message" : "Product berhasil dihapus"})
	}
}