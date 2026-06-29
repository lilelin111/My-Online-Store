import os

content = r"""package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"unicode"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var usersFilePath = "store/users.json"
var recordsFilePath = "store/records.json"

var (
	users  []User
	records []Record
	db     *sql.DB
	useDB  = false
)

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		var err error
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
			if err == nil {
				useDB = true
				createTables()
				return
			}
			fmt.Println("数据库连接失败:", err)
		}
	}
	loadUsers()
	loadRecords()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS records (
			id SERIAL PRIMARY KEY,
			sort VARCHAR(10) NOT NULL,
			category VARCHAR(255) NOT NULL,
			amount DOUBLE PRECISION NOT NULL,
			note TEXT DEFAULT '',
			"date" TIMESTAMP NOT NULL,
			total DOUBLE PRECISION NOT NULL DEFAULT 0
		)`,
	}
	for _, q := range queries {
		db.Exec(q)
	}
}

func loadUsers() {
	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("加载用户数据失败:", err)
		}
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println("解析用户数据失败:", err)
	}
}

func saveUsers() error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(usersFilePath, data, 0644)
}

func loadRecords() {
	data, err := os.ReadFile(recordsFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("加载记录数据失败:", err)
		}
		return
	}
	err = json.Unmarshal(data, &records)
	if err != nil {
		fmt.Println("解析记录数据失败:", err)
	}
}

func saveRecords() error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(recordsFilePath, data, 0644)
}
"""  # truncated for brevity

with open("C:\\Users\\27128\\Desktop\\git\\store\\store.go", "w", encoding="utf-8") as f:
    f.write(content)
print("written")
