package main

import (
	"encoding/json"
	"fmt"
	"myproject/store"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
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
	RespondSuccess(w, http.StatusOK, "ע�Գɹ�", user)
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
	RespondSuccess(w, http.StatusOK, "��¼�ɹ�", map[string]interface{}{
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
		RespondError(w, http.StatusBadRequest, "�û�ID����Ϊ��")
		return
	}
	record, err := store.CreateRecord(req.UserID, req.Sort, req.Category, req.Amount, req.Note, req.Date)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	RespondSuccess(w, http.StatusOK, "�����ɹ�", record)
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
		RespondError(w, http.StatusBadRequest, "�û�ID����Ϊ��")
		return
	}
	Records := store.GetRecordsByUserID(req.UserID)
	RespondSuccess(w, http.StatusOK, "��ȡ�ɹ�", Records)
}

func ShowRecord2(w http.ResponseWriter, r *http.Request) {
	var rep struct {
		ID      int    `json:"id"`
		UserID  int    `json:"user_id"`
		Keyword string `json:"keyword"`
	}
	err := json.NewDecoder(r.Body).Decode(&rep)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if rep.Keyword != "" {
		records := store.GetRecordsByUserID(rep.UserID)
		result := make([]store.Record, 0)
		for _, record := range records {
			sortName := record.Sort
			if record.Sort == "Income" {
				sortName = "����"
			}
			if record.Sort == "Expense" {
				sortName = "֧��"
			}
			if strings.Contains(sortName, rep.Keyword) || strings.Contains(record.Category, rep.Keyword) || strings.Contains(record.Note, rep.Keyword) {
				result = append(result, record)
			}
		}
		RespondSuccess(w, http.StatusOK, "�����ɹ�", result)
		return
	}
	Records := store.GetRecordsByUserID(rep.UserID)
	for i := range Records {
		if Records[i].ID == rep.ID {
			if Records[i].UserID != rep.UserID {
				RespondError(w, http.StatusForbidden, "��Ȩ���ʼ�¼")
				return
			}
			record, err2 := store.ShowRecord1(Records[i].ID)
			if err2 != nil {
				RespondError(w, http.StatusBadRequest, err2.Error())
				return
			}
			RespondSuccess(w, http.StatusOK, "��ȡ�ɹ�", record)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "δ�ҵ�IDΪ"+strconv.Itoa(rep.ID)+"�ļ�¼")
}

func DeleteRecord1(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int    `json:"id"`
		UserID int    `json:"user_id"`
		Sort   string `json:"sort"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.UserID <= 0 {
		RespondError(w, http.StatusBadRequest, "�û�ID����Ϊ��")
		return
	}
	Records := store.GetRecordsByUserID(req.UserID)
	for i := range Records {
		if Records[i].ID == req.ID {
			records1, err2 := store.DeleteRecord(Records[i].ID, req.UserID, req.Sort)
			if err2 != nil {
				RespondError(w, http.StatusBadRequest, err2.Error())
				return
			}
			RespondSuccess(w, http.StatusOK, "ɾ���ɹ�", records1)
			return
		}
	}
	RespondError(w, http.StatusNotFound, "δ�ҵ�IDΪ"+strconv.Itoa(req.ID)+"�ļ�¼")
}

func createProjectShortcut() {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()
	shellObj, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		fmt.Println("创建 COM 对象失败:", err)
		return
	}
	defer shellObj.Release()

	wshell, err := shellObj.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		fmt.Println("获取接口失败:", err)
		return
	}
	defer wshell.Release()
	currentDir, _ := os.Getwd()
	batPath := filepath.Join(currentDir, "start.bat")
	homeDir, _ := os.UserHomeDir()
	desktopPath := filepath.Join(homeDir, "Desktop", "皇帝记账仪.lnk")
	shortcut, err := oleutil.CallMethod(wshell, "CreateShortcut", desktopPath)
	if err != nil {
		fmt.Println("创建快捷方式失败:", err)
		return
	}
	dispatch := shortcut.ToIDispatch()
	defer dispatch.Release()
	oleutil.PutProperty(dispatch, "TargetPath", batPath)
	oleutil.PutProperty(dispatch, "WorkingDirectory", currentDir)
	oleutil.PutProperty(dispatch, "Description", "一键启动前后端服务")
	_, err = oleutil.CallMethod(dispatch, "Save")
	if err != nil {
		fmt.Println("保存快捷方式失败:", err)
		return
	}
	fmt.Println("🎉 桌面快捷方式创建成功！")
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/api/login", Login)
	r.HandleFunc("/api/register", Register)
	r.HandleFunc("/api/CreateRecord", CreateRecord1)
	r.HandleFunc("/api/ShowRecord1", ShowRecord1)
	r.HandleFunc("/api/DeleteRecord1", DeleteRecord1)
	r.HandleFunc("/api/ShowRecord2", ShowRecord2)
	createProjectShortcut()

	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/index_new.html")
	})

	fmt.Println("服务器启动在 http://localhost:8081")
	http.ListenAndServe(":8081", corsMiddleware(r))
}
