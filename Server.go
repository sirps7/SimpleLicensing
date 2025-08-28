package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	Licensing "github.com/sirps7/SimpleLicensing/SimpleLicensing"
	"github.com/gorilla/mux"
)

var (
	PORT     string
	SSL      bool
	KEY      string
	HOST     string
	DATABASE string
	USERNAME string
	PASSWORD string

	db  *sql.DB
	err error
)

// Helper to generate random license keys
func randomString(n int) string {
	var letterRunes = []rune("1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Load config from environment variables
func loadEnvConfig() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	KEY = os.Getenv("KEY")
	if KEY == "" {
		KEY = randomString(16)
	}
	HOST = os.Getenv("DB_HOST")
	DATABASE = os.Getenv("DB_NAME")
	USERNAME = os.Getenv("DB_USER")
	PASSWORD = os.Getenv("DB_PASS")
	sslEnv := os.Getenv("SSL")
	SSL = sslEnv == "true"
}

// HTTP: Root
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go Simple Licensing Server")
}

// HTTP: Check license
func checkHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	license := r.FormValue("license")
	var tmpexp string
	decrypted := Licensing.Decrypt(KEY, license)

	err := db.QueryRow("SELECT experation FROM licenses WHERE license=?", decrypted).Scan(&tmpexp)
	if err == sql.ErrNoRows {
		fmt.Fprintf(w, "Bad")
		return
	}
	ip := strings.Split(r.RemoteAddr, ":")[0]
	_, _ = db.Exec("UPDATE licenses SET ip=? WHERE license=?", ip, decrypted)

	t, err := time.Parse("2006-01-02", tmpexp)
	if err != nil {
		fmt.Println("ERROR: SQL Table Date not in correct format")
		fmt.Fprintf(w, "Error")
		return
	}
	t2, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if t.After(t2) {
		fmt.Fprintf(w, "Good")
	} else {
		fmt.Fprintf(w, "Expired")
	}
}

// HTTP: Add license (POST JSON: { "email": "", "experation": "YYYY-MM-DD" })
func addLicenseHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email      string `json:"email"`
		Experation string `json:"experation"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	license := randomString(4) + "-" + randomString(4) + "-" + randomString(4)
	var tmpemail string
	err = db.QueryRow("SELECT email FROM licenses WHERE license=?", license).Scan(&tmpemail)
	if err != sql.ErrNoRows {
		http.Error(w, "License collision, try again", 500)
		return
	}

	_, err = db.Exec("INSERT INTO licenses(email, license, experation, ip) VALUES(?, ?, ?, ?)",
		data.Email, license, data.Experation, "none")
	if err != nil {
		http.Error(w, "DB insert error", 500)
		return
	}

	resp := map[string]string{
		"license":    Licensing.Encrypt(KEY, license),
		"email":      data.Email,
		"experation": data.Experation,
	}
	json.NewEncoder(w).Encode(resp)
}

// HTTP: Remove license (POST JSON: { "email": "" })
func removeLicenseHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	res, err := db.Exec("DELETE FROM licenses WHERE email=?", data.Email)
	if err != nil {
		http.Error(w, "DB delete error", 500)
		return
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "License not found", 404)
		return
	}
	fmt.Fprintf(w, "License removed")
}

// HTTP: Count licenses
func countHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM licenses").Scan(&count)
	if err != nil {
		count = 0
	}
	json.NewEncoder(w).Encode(map[string]int{"total_licenses": count})
}

// Setup router
func API() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/check", checkHandler).Methods("POST")
	router.HandleFunc("/add", addLicenseHandler).Methods("POST")
	router.HandleFunc("/remove", removeLicenseHandler).Methods("POST")
	router.HandleFunc("/count", countHandler).Methods("GET")

	if SSL {
		err := http.ListenAndServeTLS(":"+PORT, "server.crt", "server.key", router)
		if err != nil {
			fmt.Println("SSL Server Error: " + err.Error())
			os.Exit(0)
		}
	} else {
		err := http.ListenAndServe(":"+PORT, router)
		if err != nil {
			fmt.Println("Server Error: " + err.Error())
			os.Exit(0)
		}
	}
}

func main() {
	fmt.Println("Go Simple Licensing System (Render-ready)")

	loadEnvConfig()

	db, err = sql.Open("mysql", USERNAME+":"+PASSWORD+"@tcp("+HOST+")/"+DATABASE)
	if err != nil {
		fmt.Println("[!] ERROR: CHECK MYSQL SETTINGS! [!]")
		os.Exit(0)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("[!] ERROR: CHECK IF MYSQL SERVER IS ONLINE! [!]")
		os.Exit(0)
	}

	API() // Start server
}
