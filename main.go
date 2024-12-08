package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"onlineshop/handler"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// create database
	db, err := sql.Open("pgx",os.Getenv("DB_URI")) // open connection

	// condition : if fail create database
	if err != nil{
		fmt.Printf("Gagal membuat koneksi database: %v\n", err)
		os.Exit(1)
	}

	// close program
	defer db.Close()

	if err = db.Ping(); err != nil {
		fmt.Printf("Gagal memverifikasi koneksi database : %v\n", err)
		os.Exit(1)
	}

	if _, err = migrate(db); err != nil {
		fmt.Printf("Gagal melakukan migrasi database: %v\n", err)
		os.Exit(1)
	}

	r := gin.Default()

	//products
	r.GET("/api/v1/products", handler.ListProducts(db)) //show list Product
	r.GET("/api/v1/products/:id", handler.GetProduct(db))
	r.POST("/api/v1/checkout")

	//orders
	r.POST("/api/v1/orders/:id/confirm")
	r.GET("/api/v1/orders/:id")

	//admin
	r.POST("/admin/products", handler.CreateProduct(db))
	r.PUT("/admin/products/:id", handler.UpdateProduct(db))
	r.DELETE("/admin/products/:id", handler.DeleteProduct(db))

	server := &http.Server{
		Addr: ":8080",
		Handler: r,
	}

	if err = server.ListenAndServe(); err != nil{
		fmt.Printf("Gagal menjalankan server: %v\n", err)
		os.Exit(1)
	}
}