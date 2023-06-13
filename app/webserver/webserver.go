package webserver

import (
	"bank/app/userData"
	"bank/config"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)


func CreateAccountHandler(w http.ResponseWriter, r *http.Request){
	var user userdata.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = userdata.CreateAccount(user.Name, user.Email, user.Password)
	if err != nil{
		http.Error(w, "アカウント作成を失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func homeHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("app/templates/create.html")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Start() error{
	http.HandleFunc("/create", CreateAccountHandler)
	http.HandleFunc("/", homeHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}