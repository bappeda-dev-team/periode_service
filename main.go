package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	dsn := os.Getenv("PERENCANAAN_DB_URL")
	if dsn == "" {
		log.Fatal("PERENCANAAN_DB_URL env tidak terdefinisi")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[FATAL] Error connecting to db: %v", err)
	}

	log.Printf("koneksi ke database berhasil")
	db.SetMaxOpenConns(70)
	db.SetMaxIdleConns(300)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Gagal terhubung ke database dalam 10 detik: %v", err)
		log.Printf("Mencoba koneksi ulang...")

		// Coba lagi dengan timeout yang lebih lama
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err = db.PingContext(ctx)
		if err != nil {
			db.Close()
			log.Fatalf("Koneksi database gagal setelah percobaan ulang: %v", err)
		}
	}

	log.Print("Berhasil terhubung ke database")
	log.Printf("Max Open Connections: %d", db.Stats().MaxOpenConnections)
	log.Printf("Open Connections: %d", db.Stats().OpenConnections)
	log.Printf("In Use Connections: %d", db.Stats().InUse)
	log.Printf("Idle Connections: %d", db.Stats().Idle)
}

func periodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed, pakai GET", http.StatusMethodNotAllowed)
		return
	}

	// query pohon tematik
	rows, err := db.Query(`SELECT id, tahun_awal, tahun_akhir
                           FROM tb_periode`)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var list []Periode
	for rows.Next() {
		var periode Periode
		if err := rows.Scan(&periode.Id, &periode.TahunAwal, &periode.TahunAkhir); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, periode)
	}

	msg := "Daftar Periode"
	response := Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    list}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listOpdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed, pakai GET", http.StatusMethodNotAllowed)
		return
	}

	// query opd
	rows, err := db.Query(`SELECT kode_opd, nama_opd
                           FROM tb_operasional_daerah`)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var list []Opd
	for rows.Next() {
		var opd Opd
		if err := rows.Scan(&opd.KodeOpd, &opd.NamaOpd); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		list = append(list, opd)
	}

	msg := "Daftar OPD"
	response := Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    list}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func listUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed, pakai GET", http.StatusMethodNotAllowed)
		return
	}

	kodeOpd := r.URL.Query().Get("kode_opd")
	if kodeOpd == " " || len(kodeOpd) != 22 {
		http.Error(w, "KODE OPD TIDAK DITEMUKAN", http.StatusBadRequest)
		return
	}

	// query opd
	rows, err := db.Query(`SELECT DISTINCT u.nip, p.nama, u.email, u.is_active, r.role
                           FROM tb_users u
                           JOIN tb_pegawai p ON u.nip = p.nip
                           JOIN tb_user_role ur ON u.id = ur.user_id
                           JOIN tb_role r ON ur.role_id = r.id
                           WHERE p.kode_opd = ?`, kodeOpd)
	if err != nil {
		http.Error(w, "query error: "+err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var list []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Nip, &user.Nama, &user.Email, &user.IsActive, &user.Roles); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		statusText := "Tidak Aktif"
		if user.IsActive == 1 {
			statusText = "Aktif"
		}
		user.Status = statusText

		list = append(list, user)
	}

	msg := "Daftar User OPD"
	response := Response{
		Status:  http.StatusOK,
		Message: msg,
		Data:    list}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"PERIODE SERVICE UP"}`))
}

func main() {
	log.Print("BUKAN PERIODE SERVICE")

	initDB()

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/periode", periodeHandler)
	http.HandleFunc("/list_opd", listOpdHandler)
	http.HandleFunc("/list_user", listUserHandler)

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
