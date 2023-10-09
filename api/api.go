package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func MakeHandler(router *mux.Router) {
	router.HandleFunc("/auth", AuthHandler).Methods("POST")
	router.HandleFunc("/auth", LogoutHandler).Methods("DELETE")

	// NOTE: this handler is developed for testing purposes only and should be removed in production
	router.HandleFunc("/console", ConsoleHandler).Methods("POST")
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"message":"%s","token":"%s"}`, L10n("Already logged in", r), cookie.Value)))
			return
		}
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Username or password is empty", r))))
		return
	}

	// check if we can connect to the database with the given credentials
	db, err := LoginDatabase(username, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
		return
	}

	_ = db.Disconnect()

	// if we can, we return a token which is set as a cookie
	encoded := Encode(username, password)
	cookie := http.Cookie{
		Name:    "token",
		Value:   encoded,
		Expires: time.Now().Add(24 * time.Hour),
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"message":"%s","token":"%s"}`, L10n("Successfully logged in", r), encoded)))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// we delete the token cookie
	cookie := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"message":"%s"}`, L10n("Successfully logged out", r))))
}

func ConsoleHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")

	// get cookie which contains connection credentials ('token')
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			// decode the cookie value
			username, password, err := Decode(cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Invalid token", r))))
				return
			}

			// check if we can connect to the database with the given credentials
			db, err := LoginDatabase(username, password)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
				return
			}

			// if we can, we execute the query and return the result
			rows, err := db.Database.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
				return
			}

			// return rows as JSON (use a library for this)
			columns, err := rows.Columns()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
				return
			}

			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			rowObjects := make([]map[string]interface{}, 0)
			for rows.Next() {
				for i := range columns {
					valuePtrs[i] = &values[i]
				}

				err := rows.Scan(valuePtrs...)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
					_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
					return
				}

				rowObject := make(map[string]interface{})
				for i, col := range columns {
					var v interface{}
					val := values[i]
					b, ok := val.([]byte)
					if ok {
						v = string(b)
					} else {
						v = val
					}
					rowObject[col] = v
				}
				rowObjects = append(rowObjects, rowObject)
			}

			jsonString, err := json.Marshal(rowObjects)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"rows":%s}`, jsonString)))

			_ = db.Disconnect()
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
}

func Decode(encoded string) (username, password string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}
	str := string(decoded)
	parts := strings.Split(str, ":")
	username = parts[0]
	password = parts[1]
	return username, password, nil
}

func Encode(username, password string) string {
	str := fmt.Sprintf("%s:%s", username, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	return encoded
}
