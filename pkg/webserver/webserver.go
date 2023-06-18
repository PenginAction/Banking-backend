package webserver

import (
	"bank/config"
	"bank/pkg/acountData"
	"bank/pkg/transactiondata"
	"bank/pkg/userData"
	"bank/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var user userdata.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = userdata.CreateAccount(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "アカウント作成を失敗しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("app/templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := userdata.GetAccountByEmail(email)
		if err != nil {
			http.Error(w, "メールアドレスかパスワードが一致しません", http.StatusUnauthorized)
			return
		}

		err = userdata.CompareHashAndPassword(user.Password, password)
		if err != nil {
			http.Error(w, "メールアドレスかパスワードが一致しません", http.StatusUnauthorized)
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil {
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
		if err != nil {
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

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "/login", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}

func AccountHandler(w http.ResponseWriter, r *http.Request) {
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

	account, err := acountdata.GetAccountByUserId(userID)
	if err != nil {
		http.Error(w, "ユーザーアカウントの取得に失敗しました", http.StatusInternalServerError)
		return
	}

	transactions, err := transactiondata.GetTransactionByAccountId(account.AccountID)
	if err != nil {
		http.Error(w, "取引履歴を取得することができませんでした", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"User":         user,
		"Balance":      account.Balance,
		"Transactions": transactions,
	}

	err = utils.RenderTemplate(w, "app/templates/acount.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			http.Error(w, "不正な金額です", http.StatusBadRequest)
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "セッションの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		userId := session.Values["user_id"]

		account, err := acountdata.GetAccountByUserId(userId.(int))
		if err != nil {
			http.Error(w, "アカウントが見つかりません", http.StatusInternalServerError)
			return
		}

		err = transactiondata.CreateTransaction(account.AccountID, amount, "入金")
		if err != nil {
			http.Error(w, "入金処理に失敗しました", http.StatusInternalServerError)
			return
		}

		err = acountdata.UpdateBalance(account.AccountID, account.Balance+amount)
		if err != nil {
			http.Error(w, "口座残高の更新に失敗しました", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	} else {
		err := utils.RenderTemplate(w, "app/templates/deposit.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		amountStr := r.FormValue("amount")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			http.Error(w, "無効な金額が入力されました", http.StatusBadRequest)
			return
		}

		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, "セッションの取得に失敗しました", http.StatusInternalServerError)
			return
		}

		userId := session.Values["user_id"]
		account, err := acountdata.GetAccountByUserId(userId.(int))
		if err != nil {
			http.Error(w, "アカウントが見つかりません", http.StatusNotFound)
			return
		}

		if account.Balance < amount {
			http.Error(w, "口座の残高が不足しています．", http.StatusBadRequest)
			return
		}

		err = transactiondata.CreateTransaction(account.AccountID, amount, "出金")
		if err != nil {
			http.Error(w, "取引を処理することができませんでした", http.StatusInternalServerError)
			return
		}

		err = acountdata.UpdateBalance(account.AccountID, account.Balance-amount)
		if err != nil {
			http.Error(w, "口座の残高を更新できませんでした", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/account", http.StatusSeeOther)
	} else {
		err := utils.RenderTemplate(w, "app/templates/withdraw.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func Start() error {
	http.HandleFunc("/create", CreateAccountHandler)
	http.HandleFunc("/", firstHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/account", AuthRequiredHandler(AccountHandler))
	http.HandleFunc("/deposit", DepositHandler)
	http.HandleFunc("/withdraw", WithdrawHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
