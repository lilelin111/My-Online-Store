package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var usersFilePath = "store/users.json"
var recordsFilePath = "store/records.json"

var users []User
var records []Record

func init() {
	loadUsers()
	loadRecords()
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

func FindUserByName(name string) *User {
	for i := range users {
		if users[i].Name == name {
			return &users[i]
		}
	}
	return nil
}
func CreateUser(Name string, password string) (*User, error) {
	if Name == "" || password == "" {
		return nil, errors.New("用户名和密码不能为空！")
	}
	if FindUserByName(Name) != nil {
		return nil, errors.New("用户名已存在！")
	}
	if len(password) < 8 {
		return nil, errors.New("密码长度必须大于8！")
	}
	count := 0
	count1 := 0
	for _, i := range password {
		if unicode.IsLetter(i) {
			count++
		} else if unicode.IsDigit(i) {
			count1++
		}
	}
	if count == 0 {
		return nil, errors.New("密码必须包含字母！")
	}
	if count1 == 0 {
		return nil, errors.New("密码必须包含数字！")
	}
	if (len(password) - count - count1) <= 0 {
		return nil, errors.New("密码中必须包含特殊符号！")
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败")
		return nil, err
	}
	newID := 1
	if len(users) > 0 {
		newID = users[len(users)-1].ID + 1
	}
	user1 := &User{
		ID:       newID,
		Name:     Name,
		Password: string(newPassword),
	}
	users = append(users, *user1)
	err = saveUsers()
	if err != nil {
		fmt.Println("保存用户失败:", err)
		return nil, err
	}
	return user1, nil
}

func CheckUser(name string, password string, u *User) (*User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("密码错误！")
	}
	if name != u.Name {
		return nil, fmt.Errorf("用户名错误！")
	}
	return u, nil
}
func LoginService(username string, password string) (*User, error) {
	user := FindUserByName(username)
	if user == nil {
		return nil, fmt.Errorf("用户不存在！")
	}
	user1, err := CheckUser(username, password, user)
	if err != nil {
		return nil, err
	}
	return user1, nil
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

var IncomeCategories = []string{
	"工资",
	"奖金",
	"投资收益",
	"兼职",
	"其他收入",
}

var ExpenseCategories = []string{
	"餐饮",
	"交通",
	"购物",
	"娱乐",
	"医疗",
	"教育",
	"住房",
	"其他支出",
}

type Record struct {
	ID       int       `json:"id"`
	Sort     string    `json:"sort"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
	Note     string    `json:"note"`
	Date     time.Time `json:"date"`
	Total    float64   `json:"total"`
}

func BoolSort(a string) bool {
	if a == "Income" || a == "Expense" {
		return true
	}
	return false
}

func contains(list []string, category string) bool {
	for _, item := range list {
		if item == category {
			return true
		}
	}
	return false
}

var (
	ErrInvalidAmount = errors.New("金额必须大于0")
	ErrInvalidSort   = errors.New("无效收支类型，必须为'Income'或'Expense'")
)

func CreateRecord(sort string, category string, amount float64, note string, date time.Time, total float64) (Record, error) {
	if amount <= 0 {
		fmt.Println("金额必须大于0")
		return Record{}, ErrInvalidAmount
	}
	if !BoolSort(sort) {
		return Record{}, ErrInvalidSort
	}
	if sort == "Income" {
		if !contains(IncomeCategories, category) {
			return Record{}, fmt.Errorf("无效收入类型'%s',必须是：%v", category, IncomeCategories)
		}
	} else if sort == "Expense" {
		if !contains(ExpenseCategories, category) {
			return Record{}, fmt.Errorf("无效支出类型'%s',必须是：%v", category, ExpenseCategories)
		}
	} else {
		return Record{}, ErrInvalidSort
	}
	newID := 1
	if len(records) > 0 {
		newID = records[len(records)-1].ID + 1
	}
	record := Record{
		ID:       newID,
		Sort:     sort,
		Category: category,
		Amount:   amount,
		Note:     note,
		Date:     date,
		Total:    total,
	}
	records = append(records, record)
	if err := saveRecords(); err != nil {
		fmt.Println("保存记录失败:", err)
		return Record{}, err
	}
	return record, nil
}

func ShowRecord(id int) (*Record, error) {
	if len(records) <= 0 {
		fmt.Println("现在没有任何记录")
		return nil, errors.New("现在没有任何数据")
	}
	fmt.Println("账单如下：")
	for i := range records {
		if records[i].ID == id {
			return &records[i], nil
		}
	}
	return nil, fmt.Errorf("未找到'ID'为%d的记录", id)
}

func DeleteRecord(id int) ([]Record, error) {
	if len(records) == 0 {
		return nil, errors.New("没有任何账单记录！")
	}
	for i, record := range records {
		if record.ID == id {
			fmt.Println("删除后账单如下：")
			records = append(records[:i], records[i+1:]...)
			if err := saveRecords(); err != nil {
				fmt.Println("保存记录失败:", err)
				return nil, err
			}
			return records, nil
		}
	}
	return nil, fmt.Errorf("未找到 ID 为 %d 的账单记录，请检查后重试", id)
}

func GetAllRecords() []Record {
	if len(records) == 0 {
		return nil
	}
	return records
}
