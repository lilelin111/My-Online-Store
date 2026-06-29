package main

import (
	"encoding/json"
	"fmt"
	"myproject/store"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LogRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func RespondError(w http.ResponseWriter, status int, message string) {
	ResponseJSON(w, status, Response{Success: false, Message: message})
}
func RespondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	ResponseJSON(w, status, Response{Success: true, Message: message, Data: data})
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})

}
func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := store.CreateUser(req.Username, req.Password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondSuccess(w, http.StatusOK, "注册成功", user)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var rep LogRequest
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := store.LoginService(rep.Username, rep.Password)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	ResponseJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "登录成功",
		Data: map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		},
	})
}
func CreateRecord1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req store.Record
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	record, err := store.CreateRecord(req.Sort, req.Category, req.Amount, req.Note, req.Date, req.Total)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondSuccess(w, http.StatusOK, "创建成功", record)
}
func ShowRecord1(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	Records := store.GetAllRecords()
	if len(Records) == 0 {
		RespondError(w, http.StatusBadRequest, "不存在账单！")
		return
	}
	RespondSuccess(w, http.StatusOK, "获取成功", Records)

}
func ShowRecord2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var rep struct {
		ID int `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	Records := store.GetAllRecords()
	for i := range Records {
		if Records[i].ID == rep.ID {
			record, err2 := store.ShowRecord(Records[i].ID)
			if err2 != nil {
				RespondError(w, http.StatusBadRequest, err2.Error())
				return
			}
			ResponseJSON(w, http.StatusOK, record)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "未找到ID为"+strconv.Itoa(rep.ID)+"的记录")
}
func DeleteRecord1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		ID int `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	Records := store.GetAllRecords()
	for i := range Records {
		if Records[i].ID == req.ID {
			records1, err2 := store.DeleteRecord(Records[i].ID)
			if err2 != nil {
				RespondError(w, http.StatusBadRequest, err2.Error())
				return
			}
			RespondSuccess(w, http.StatusOK, "删除成功", records1)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "未找到ID为"+strconv.Itoa(req.ID)+"的记录")
}
func main() {
	r := http.NewServeMux()
	r.HandleFunc("/api/login", Login)
	r.HandleFunc("/api/register", Register)
	r.HandleFunc("/api/CreateRecord", CreateRecord1)
	r.HandleFunc("/api/ShowRecord1", ShowRecord1)
	r.HandleFunc("/api/DeleteRecord1", DeleteRecord1)
	r.HandleFunc("/api/ShowRecord2", ShowRecord2)
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index.html")
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("服务器启动在 :" + port)
	http.ListenAndServe(":"+port, corsMiddleware(r))
}
