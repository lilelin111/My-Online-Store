package main

import (
	"encoding/json"
	"fmt"
	"myproject/store"
	"net/http"
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
	RespondSuccess(w, http.StatusOK, "注冊成功", user)
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
	RespondSuccess(w, http.StatusOK, "登录成功", map[string]interface{}{
		"id":   user.ID,
		"name": user.Name,
	})
}

func CreateRecord1(w http.ResponseWriter, r *http.Request) {
	var req store.Record
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.UserID <= 0 {
		RespondError(w, http.StatusBadRequest, "用户ID不能为空")
		return
	}
	record, err := store.CreateRecord(req.UserID, req.Sort, req.Category, req.Amount, req.Note, req.Date, req.Total)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondSuccess(w, http.StatusOK, "创建成功", record)
}

func ShowRecord1(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.UserID <= 0 {
		RespondError(w, http.StatusBadRequest, "用户ID不能为空")
		return
	}
	Records := store.GetRecordsByUserID(req.UserID)
	RespondSuccess(w, http.StatusOK, "获取成功", Records)
}

func ShowRecord2(w http.ResponseWriter, r *http.Request) {
	var rep struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	Records := store.GetRecordsByUserID(rep.UserID)
	for i := range Records {
		if Records[i].ID == rep.ID {
			if Records[i].UserID != rep.UserID {
				RespondError(w, http.StatusForbidden, "无权访问记录")
				return
			}
			record, err2 := store.ShowRecord(Records[i].ID)
			if err2 != nil {
				RespondError(w, http.StatusBadRequest, err2.Error())
				return
			}
			RespondSuccess(w, http.StatusOK, "获取成功", record)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "未找到ID为"+strconv.Itoa(rep.ID)+"的记录")
}

func DeleteRecord1(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.UserID <= 0 {
		RespondError(w, http.StatusBadRequest, "用户ID不能为空")
		return
	}
	Records := store.GetRecordsByUserID(req.UserID)
	for i := range Records {
		if Records[i].ID == req.ID {
			records1, err2 := store.DeleteRecord(Records[i].ID, req.UserID)
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

	// API 路由
	r.HandleFunc("/api/login", Login)
	r.HandleFunc("/api/register", Register)
	r.HandleFunc("/api/CreateRecord", CreateRecord1)
	r.HandleFunc("/api/ShowRecord1", ShowRecord1)
	r.HandleFunc("/api/DeleteRecord1", DeleteRecord1)
	r.HandleFunc("/api/ShowRecord2", ShowRecord2)

	// 静态文件与前端页面托管
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index_new.html")
	})

	fmt.Println("服务器启动在 http://localhost:8081")
	http.ListenAndServe(":8081", corsMiddleware(r))
}
