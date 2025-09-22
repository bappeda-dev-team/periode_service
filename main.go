package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// var db *sql.DB

// func initDB() {
// 	dsn := os.Getenv("PERENCANAAN_DB_URL")
// 	if dsn == "" {
// 		log.Fatal("PERENCANAAN_DB_URL env tidak terdefinisi")
// 	}

// 	var err error
// 	db, err = sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("[FATAL] Error connecting to db: %v", err)
// 	}

// 	log.Printf("koneksi ke database berhasil")
// 	db.SetMaxOpenConns(70)
// 	db.SetMaxIdleConns(300)
// 	db.SetConnMaxIdleTime(5 * time.Minute)
// 	db.SetConnMaxLifetime(60 * time.Minute)

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
// 		log.Printf("Mencoba koneksi ulang...")

// 		// Coba lagi dengan timeout yang lebih lama
// 		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
// 		defer cancel()

// 		err = db.PingContext(ctx)
// 		if err != nil {
// 			db.Close()
// 			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
// 		}
// 	}

// 	log.Print("Berhasil terhubung ke database")
// 	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
// 	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
// 	log.Printf("In Use Connections: %d", db.Stats().InUse)
// 	log.Printf("Idle Connections: %d", db.Stats().Idle)
// }

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"PERIODE SERVICE UP"}`))
}

func main() {
	log.Print("PERIODE SERVICE")

	// initDB()

	http.HandleFunc("/health", healthCheckHandler)
	// http.HandleFunc("/periode", periodeHandler)

	handler := corsMiddleware(http.DefaultServeMux)
	log.Println("Server running di :8080")

	http.ListenAndServe(":8080", handler)
}

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Untuk development, bisa pakai "*"
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
