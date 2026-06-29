import sys
with open("store/store.go", "w", encoding="utf-8") as f:
    # Write Part 1
    f.write('''package store

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
	users   []User
	records []Record
	db      *sql.DB
	useDB   = false
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
	db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name VARCHAR(255) UNIQUE NOT NULL, password VARCHAR(255) NOT NULL)")
	db.Exec(
		"CREATE TABLE IF NOT EXISTS records (" +
		"id SERIAL PRIMARY KEY, sort VARCHAR(10) NOT NULL, " +
		"category VARCHAR(255) NOT NULL, amount DOUBLE PRECISION NOT NULL, " +
		"note TEXT DEFAULT '', \"date\" TIMESTAMP NOT NULL, total DOUBLE PRECISION NOT NULL DEFAULT 0" +
		")",
	)
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
''')
    
    # Write Part 2 - functions and types
    # Use explicit string with backtick handling
    bt = chr(96)  # backtick character
    
    f.write(f'''
func FindUserByName(name string) *User {{
	if useDB {{
		row := db.QueryRow("SELECT id, name, password FROM users WHERE name = $1", name)
		var u User
		err := row.Scan(&u.ID, &u.Name, &u.Password)
		if err != nil {{
			return nil
		}}
		return &u
	}}
	for i := range users {{
		if users[i].Name == name {{
			return &users[i]
		}}
	}}
	return nil
}}

func CreateUser(Name string, password string) (*User, error) {{
	if Name == "" || password == "" {{
		return nil, errors.New("用户名和密码不能为空")
	}}
	if FindUserByName(Name) != nil {{
		return nil, errors.New("用户名已存在")
	}}
	if len(password) < 8 {{
		return nil, errors.New("密码长度必须大于8")
	}}
	count := 0
	count1 := 0
	for _, i := range password {{
		if unicode.IsLetter(i) {{
			count++
		}} else if unicode.IsDigit(i) {{
			count1++
		}}
	}}
	if count == 0 {{
		return nil, errors.New("密码必须包含字母")
	}}
	if count1 == 0 {{
		return nil, errors.New("密码必须包含数字")
	}}
	if (len(password) - count - count1) <= 0 {{
		return nil, errors.New("密码中必须包含特殊字符")
	}}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {{
		fmt.Println("密码加密失败")
		return nil, err
	}}

	if useDB {{
		var newID int
		err := db.QueryRow(
			"INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id",
			Name, string(newPassword),
		).Scan(&newID)
		if err != nil {{
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}}
		return &User{{ID: newID, Name: Name, Password: string(newPassword)}}, nil
	}}

	newID := 1
	if len(users) > 0 {{
		newID = users[len(users)-1].ID + 1
	}}
	user1 := &User{{
		ID:       newID,
		Name:     Name,
		Password: string(newPassword),
	}}
	users = append(users, *user1)
	err = saveUsers()
	if err != nil {{
		fmt.Println("保存用户失败:", err)
		return nil, err
	}}
	return user1, nil
}}

func CheckUser(name string, password string, u *User) (*User, error) {{
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {{
		return nil, fmt.Errorf("密码错误")
	}}
	if name != u.Name {{
		return nil, fmt.Errorf("用户名错误")
	}}
	return u, nil
}}

func LoginService(username string, password string) (*User, error) {{
	user := FindUserByName(username)
	if user == nil {{
		return nil, fmt.Errorf("用户不存在")
	}}
	user1, err := CheckUser(username, password, user)
	if err != nil {{
		return nil, err
	}}
	return user1, nil
}}

type User struct {{
	ID       int    {bt}json:"id"{bt}
	Name     string {bt}json:"name"{bt}
	Password string {bt}json:"-"{bt}
}}

var IncomeCategories = []string{{
	"工资",
	"奖金",
	"投资收益",
	"兼职",
	"其他收入",
}}

var ExpenseCategories = []string{{
	"餐饮",
	"交通",
	"购物",
	"娱乐",
	"医疗",
	"教育",
	"住房",
	"其他支出",
}}

type Record struct {{
	ID       int       {bt}json:"id"{bt}
	Sort     string    {bt}json:"sort"{bt}
	Category string    {bt}json:"category"{bt}
	Amount   float64   {bt}json:"amount"{bt}
	Note     string    {bt}json:"note"{bt}
	Date     time.Time {bt}json:"date"{bt}
	Total    float64   {bt}json:"total"{bt}
}}

func BoolSort(a string) bool {{
	return a == "Income" || a == "Expense"
}}

func contains(list []string, category string) bool {{
	for _, item := range list {{
		if item == category {{
			return true
		}}
	}}
	return false
}}

var (
	ErrInvalidAmount = errors.New("金额必须大于0")
	ErrInvalidSort   = errors.New("无效收支类型，必须为'Income'或'Expense'")
)

func CreateRecord(sort string, category string, amount float64, note string, date time.Time, total float64) (Record, error) {{
	if amount <= 0 {{
		return Record{{}}, ErrInvalidAmount
	}}
	if !BoolSort(sort) {{
		return Record{{}}, ErrInvalidSort
	}}
	if sort == "Income" {{
		if !contains(IncomeCategories, category) {{
			return Record{{}}, fmt.Errorf("无效收入类型'%s',必须是：%v", category, IncomeCategories)
		}}
	}} else if sort == "Expense" {{
		if !contains(ExpenseCategories, category) {{
			return Record{{}}, fmt.Errorf("无效支出类型'%s',必须是：%v", category, ExpenseCategories)
		}}
	}}

	if useDB {{
		var newID int
		err := db.QueryRow(
			"INSERT INTO records (sort, category, amount, note, date, total) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			sort, category, amount, note, date, total,
		).Scan(&newID)
		if err != nil {{
			return Record{{}}, fmt.Errorf("创建记录失败: %v", err)
		}}
		return Record{{ID: newID, Sort: sort, Category: category, Amount: amount, Note: note, Date: date, Total: total}}, nil
	}}

	newID := 1
	if len(records) > 0 {{
		newID = records[len(records)-1].ID + 1
	}}
	record := Record{{
		ID:       newID,
		Sort:     sort,
		Category: category,
		Amount:   amount,
		Note:     note,
		Date:     date,
		Total:    total,
	}}
	records = append(records, record)
	if err := saveRecords(); err != nil {{
		return Record{{}}, err
	}}
	return record, nil
}}

func ShowRecord(id int) (*Record, error) {{
	if useDB {{
		row := db.QueryRow("SELECT id, sort, category, amount, note, date, total FROM records WHERE id = $1", id)
		var r Record
		err := row.Scan(&r.ID, &r.Sort, &r.Category, &r.Amount, &r.Note, &r.Date, &r.Total)
		if err != nil {{
			return nil, fmt.Errorf("未找到ID为%d的记录", id)
		}}
		return &r, nil
	}}
	for i := range records {{
		if records[i].ID == id {{
			return &records[i], nil
		}}
	}}
	return nil, fmt.Errorf("未找到ID为%d的记录", id)
}}

func DeleteRecord(id int) ([]Record, error) {{
	if useDB {{
		_, err := db.Exec("DELETE FROM records WHERE id = $1", id)
		if err != nil {{
			return nil, fmt.Errorf("删除记录失败: %v", err)
		}}
		return GetAllRecords(), nil
	}}
	for i, record := range records {{
		if record.ID == id {{
			records = append(records[:i], records[i+1:]...)
			if err := saveRecords(); err != nil {{
				return nil, err
			}}
			return records, nil
		}}
	}}
	return nil, fmt.Errorf("未找到 ID 为 %d 的账单记录", id)
}}

func GetAllRecords() []Record {{
	if useDB {{
		rows, err := db.Query("SELECT id, sort, category, amount, note, date, total FROM records ORDER BY id")
		if err != nil {{
			return nil
		}}
		defer rows.Close()
		var result []Record
		for rows.Next() {{
			var r Record
			err := rows.Scan(&r.ID, &r.Sort, &r.Category, &r.Amount, &r.Note, &r.Date, &r.Total)
			if err != nil {{
				continue
			}}
			result = append(result, r)
		}}
		return result
	}}
	return records
}}
''')

print("store.go written successfully")

