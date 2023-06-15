package webserver

import (
	"bank/app/userData"
	"bank/config"
	"bank/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))


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

func firstHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("app/templates/create.html")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := userdata.GetAccountByEmail(email)
		if err != nil{
			http.Error(w, "メールアドレスかパスワードが一致しません", http.StatusUnauthorized)
			return
		}

		err = userdata.CompareHashAndPassword(user.Password, password)
		if err != nil{
			http.Error(w, "メールアドレスかパスワードが一致しません", http.StatusUnauthorized)
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil{
			http.Error(w, "Sessionの取得に失敗しました", http.StatusInternalServerError)
		}
		session.Values["authenticated"] = true
		session.Values["user_id"] = user.UserId
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Sessionの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	} else {
		// tmpl, err := template.ParseFiles("app/templates/login.html")
		err := utils.RenderTemplate(w, "app/templates/login.html", nil)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AuthRequiredHandler(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "Sessionの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth{
			http.Error(w, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request){
    session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Sessionの取得に失敗しました", http.StatusInternalServerError)
		return
	}
    userID, ok := session.Values["user_id"].(int)
	if !ok {
		http.Error(w, "Session内のユーザーIDの取得に失敗しました", http.StatusInternalServerError)
	}

    user, err := userdata.GetAccountById(userID)
    if err != nil {
        http.Error(w, "ユーザー情報の取得に失敗しました", http.StatusInternalServerError)
        return
    }
	err = utils.RenderTemplate(w, "app/templates/ballanceAcount.html", user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func Start() error{
	http.HandleFunc("/create", CreateAccountHandler)
	http.HandleFunc("/", firstHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/account", AuthRequiredHandler(AccountHandler))
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}