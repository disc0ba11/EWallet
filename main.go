package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db, err = sql.Open("sqlite3", "file:wallets.db")
	mutex   sync.Mutex
)

type wallet struct {
	address string
	balance *big.Float
}
type transaction struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func newWallet() wallet {
	var wallet wallet
	wallet.address = randToken()
	wallet.balance, _ = new(big.Float).SetString("100")
	return wallet
}
func dbCreate_wallets(db *sql.DB) {
	table := `CREATE TABLE wallets (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "address" TEXT,
        "balance" TEXT);`
	query, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'wallets' created successfully!\n--------------")
}
func dbCreate_transactions(db *sql.DB) {
	table := `CREATE TABLE transactions (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "from" TEXT,
        "to" TEXT,
		"amount" TEXT);`
	query, err := db.Prepare(table)
	if err != nil {
		log.Fatal(err)
	}
	_, err = query.Exec()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'transactions' created successfully!\n--------------")
}
func dbSelect(db *sql.DB) int64 {
	query := `SELECT * FROM wallets`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var wallet wallet
	var id int64
	var amount int64
	var balance string
	for rows.Next() {
		rows.Scan(&id, &wallet.address, &balance)
		wallet.balance, _ = new(big.Float).SetString(balance)
		fmt.Println(id, wallet.address, wallet.balance)
		amount++
	}
	fmt.Println("--------------")
	return amount
}
func dbUpdate(db *sql.DB, amount *big.Float, address string) {
	query := `UPDATE wallets SET balance = ? WHERE address = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(amount.String(), address)
	if err != nil {
		log.Fatal(err)
	}
}
func dbTransaction(db *sql.DB, from string, to string, amount string) (int, error) {
	var fromWallet wallet
	var toWallet wallet
	var fromId int64
	var toId int64
	var balance string
	var hasDigit bool = false
	if strings.Contains(amount, `.`) {
		if len(strings.Split(amount, ".")) >= 2 {
			tail := strings.Split(amount, ".")[1]
			if len(tail) > 2 {
				return 1, errors.New("wrong amount format")
			}
			if len(strings.Split(amount, ".")) > 2 {
				return 1, errors.New("wrong amount format")
			}
		}
	}
	if amount == "0" {
		return 1, errors.New("wrong amount format")
	}
	for i := 0; i < len(amount); i++ {
		if (amount[i] < '0' || amount[i] > '9') && amount[i] != '.' {
			hasDigit = true
		}
	}
	if hasDigit {
		return 1, errors.New("wrong amount format")
	}
	amountFloat, _ := new(big.Float).SetString(amount)
	rows, err := db.Query(`SELECT * FROM wallets WHERE address = ?`, from)
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	err = rows.Scan(&fromId, &fromWallet.address, &balance)
	if err != nil {
		return 4, errors.New("no such 'from' address")
	}
	rows.Close()
	fromWallet.balance, _ = new(big.Float).SetString(balance)
	rows, err = db.Query(`SELECT * FROM wallets WHERE address = ?`, to)
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	err = rows.Scan(&toId, &toWallet.address, &balance)
	if err != nil {
		return 5, errors.New("no such 'to' address")
	}
	rows.Close()
	toWallet.balance, _ = new(big.Float).SetString(balance)
	zero, _ := new(big.Float).SetString("0")
	switch new(big.Float).Sub(fromWallet.balance, amountFloat).Cmp(zero) {
	case 1, 0:
		dbUpdate(db, new(big.Float).Sub(fromWallet.balance, amountFloat), from)
		dbUpdate(db, new(big.Float).Add(toWallet.balance, amountFloat), to)
		var transaction transaction
		transaction.From = fromWallet.address
		transaction.To = toWallet.address
		transaction.Amount = amount
		dbInsert_transaction(db, transaction)
		fmt.Println("Success!")
		return 0, nil
	case -1:
		return 2, errors.New("not enough money")
	default:
		return 0, nil
	}
}
func dbInsert_wallet(db *sql.DB, wallet wallet) {
	query := `INSERT INTO wallets(address, balance) VALUES (?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(wallet.address, wallet.balance.String())
	if err != nil {
		log.Fatal(err)
	}
}
func dbInsert_transaction(db *sql.DB, transaction transaction) {
	query := "INSERT INTO transactions(`from`, `to`, `amount`) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal("prepare", err)
	}
	_, err = stmt.Exec(transaction.From, transaction.To, transaction.Amount)
	if err != nil {
		log.Fatal("exec", err)
	}
}
func dbDelete(db *sql.DB, id int64) {
	query := `DELETE FROM wallets WHERE id = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
	}
}
func dbFill(db *sql.DB) {
	query := `SELECT * FROM wallets`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var amount int64
	var id int64
	var wallet wallet
	var ids []int64
	var balance string
	for rows.Next() {
		amount++
	}
	if amount < 10 {
		for i := amount; i < 10; i++ {
			dbInsert_wallet(db, newWallet())
		}
		fmt.Println("Succesfully filled to 10 rows")
	} else if amount > 10 {
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			rows.Next()
		}
		for i := amount - 10; i > 0; i-- {
			rows.Next()
			rows.Scan(&id, &wallet.address, &balance)
			wallet.balance, _ = new(big.Float).SetString(balance)
			ids = append(ids, id)
		}
		rows.Close()
		for _, value := range ids {
			dbDelete(db, value)
		}
		log.Println("Succesfully decreased to 10 rows")
	}
	fmt.Println("--------------")
}
func firstStart(db *sql.DB) {
	for i := 0; i < 10; i++ {
		dbInsert_wallet(db, newWallet())
	}
}
func getBalance(w http.ResponseWriter, r *http.Request) {
	const lenPath = len("/api/wallet/")
	split := r.URL.Path[lenPath:]
	request := strings.Split(split, "/")
	var wallet wallet
	var id int64
	var balance string
	if len(request) < 2 {
		fmt.Fprint(w, "Not enought params. Specify address and method for wallet request.")
		return
	} else if len(request) > 2 {
		fmt.Fprint(w, "Too many params. Wrong request. Try http://localhost/api/wallet/"+request[0]+"/balance")
		return
	} else if request[1] != "balance" {
		fmt.Fprint(w, "Wrong method.")
		return
	}
	rows, err := db.Query(`SELECT * FROM wallets WHERE address = ?`, request[0])
	if err != nil {
		log.Fatal(err)
	}
	rows.Next()
	err = rows.Scan(&id, &wallet.address, &balance)
	if err != nil {
		fmt.Fprint(w, "No such address in database.")
		return
	}
	rows.Close()
	wallet.balance, _ = new(big.Float).SetString(balance)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"Id": `+strconv.FormatInt(id, 10)+`, "Address": "`+wallet.address+`", "Balance": `+wallet.balance.String()+`}`)
}
func getLast(w http.ResponseWriter, r *http.Request) {
	const lenPath = len("/api/")
	request := r.URL.Path[lenPath:]
	if request != "transactions" {
		fmt.Fprint(w, "Wrong method.")
		return
	}
	if string(r.URL.RawQuery[0:6]) != "count=" || strings.Contains(string(r.URL.RawQuery[6:]), "count") {
		fmt.Fprint(w, "Wrong query.")
		return
	}
	count := r.FormValue("count")
	var hasDigit bool = false
	for i := 0; i < len(count); i++ {
		if count[i] < '0' || count[i] > '9' {
			hasDigit = true
		}
	}
	if hasDigit {
		fmt.Fprint(w, "Count of transactions has to be an integer.")
		return
	}
	var transaction transaction
	var id int64
	rows, err := db.Query(`SELECT * FROM transactions ORDER BY id DESC LIMIT ?`, count)
	if err != nil {
		log.Fatal(err)
	}
	var json string = `[`
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&id, &transaction.From, &transaction.To, &transaction.Amount)
		json += `{"Id": ` + strconv.FormatInt(id, 10) + `, "From": "` + transaction.From + `", "To": "` + transaction.To + `", "Amount": ` + transaction.Amount + `},`
	}
	json = strings.TrimSuffix(json, ",") + "]"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, json)
}
func send(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if strings.Count(strings.ToLower(string(body)), `"from"`) != 1 ||
		strings.Count(strings.ToLower(string(body)), `"to"`) != 1 ||
		strings.Count(strings.ToLower(string(body)), `"amount"`) != 1 ||
		!strings.Contains(strings.ToLower(string(body)), `"from"`) ||
		!strings.Contains(strings.ToLower(string(body)), `"to"`) ||
		!strings.Contains(strings.ToLower(string(body)), `"amount"`) {
		fmt.Fprint(w, "Wrong POST body.")
		log.Print(errors.New("wrong post body"))
		return
	}
	jsonBody := &transaction{}
	err = json.Unmarshal(body, jsonBody)
	if err != nil {
		fmt.Fprint(w, "Wrong POST body.")
		log.Print(errors.New("wrong post body"))
		return
	}
	fmt.Println("Transaction From: " + jsonBody.From + " To: " + jsonBody.To + " Amount: " + jsonBody.Amount)
	mutex.Lock()
	errCode, err := dbTransaction(db, jsonBody.From, jsonBody.To, jsonBody.Amount)
	mutex.Unlock()
	if err != nil {
		switch errCode {
		case 1:
			fmt.Fprint(w, "Wrong amount format. Example: 51.49")
			log.Print(err)
			return
		case 2:
			fmt.Fprint(w, "Not enough money.")
			log.Print(err)
			return
		case 3:
			fmt.Fprint(w, "Amount cannot be zero.")
			log.Print(err)
			return
		case 4:
			fmt.Fprint(w, "No such 'from' address.")
			return
		case 5:
			fmt.Fprint(w, "No such 'to' address.")
			return
		default:
			fmt.Fprint(w, "Something gone wrong.")
			log.Print(err)
			return
		}
	}
}

func main() {
	db.SetMaxOpenConns(1)
	dbExists, _ := exists("wallets.db")
	if !dbExists {
		dbCreate_wallets(db)
		dbCreate_transactions(db)
		firstStart(db)
	}
	if dbSelect(db) != 10 {
		dbFill(db)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", getLast)
	mux.HandleFunc("/api/send", send)
	mux.HandleFunc("/api/wallet/", getBalance)
	err = http.ListenAndServe(":8080", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
