package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mqtt_web_api/api/models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func MakeHandler(router *mux.Router) {
	router.HandleFunc("/auth", AuthHandler).Methods("POST")
	router.HandleFunc("/auth", LogoutHandler).Methods("DELETE")

	router.HandleFunc("/console", ConsoleHandler).Methods("POST")

	router.HandleFunc("/select/{table}", GetHandler).Methods("GET")
	router.HandleFunc("/report/{name}", ReportHandler).Methods("GET")
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
	cookie, err := GetCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
		return
	}

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

	// save to a file
	f, err := os.Create("console.json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
		return
	}

	// convert to JSON
	jsonString, err := json.Marshal(rowObjects)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
		return
	}

	_, _ = f.Write(jsonString)
	_ = f.Close()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"rows":%s}`, jsonString)))

	_ = db.Disconnect()
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET %s\n", r.URL.Path)

	// get page and limit from get parameters
	pageStr := r.FormValue("page")
	limitStr := r.FormValue("limit")

	// set default values if not provided
	if pageStr == "" {
		pageStr = "0"
	}
	if limitStr == "" {
		limitStr = "10"
	}

	// convert to int
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Invalid page number", r))))
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Invalid limit number", r))))
		return
	}

	// page is actually offset, so we need to multiply it by limit
	page *= limit

	// get table name from url
	vars := mux.Vars(r)
	table := vars["table"]

	fmt.Printf("table: %s, (page: %d, limit: %d)", table, page, limit)

	// get cookie which contains connection credentials ('token')
	cookie, err := GetCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
		return
	}

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

	var rows interface{}

	switch table {
	case "category":
		var categories []models.Category
		categories, err = db.GetCategories(page, limit)
		rows = categories
	case "coordinates":
		var coordinates []models.Coordinates
		coordinates, err = db.GetCoordinates(page, limit)
		rows = coordinates
	case "favorite":
		var favorites []models.Favorite
		favorites, err = db.GetFavorites(page, limit)
		rows = favorites
	case "measured_unit":
		var units []models.MeasuredUnit
		units, err = db.GetMeasuredUnits(page, limit)
		rows = units
	case "measurement":
		stationParam := r.FormValue("station")
		var measurments []models.Measurement
		measurments, err = db.GetMeasurements(page, limit, stationParam)
		rows = measurments
	case "mqtt_server":
		var servers []models.MqttServer
		servers, err = db.GetMqttServers(page, limit)
		rows = servers
	case "mqtt_unit":
		var units []models.MqttUnit
		units, err = db.GetMqttUnits(page, limit)
		rows = units
	case "optimal_value":
		var values []models.OptimalValue
		values, err = db.GetOptimalValues(page, limit)
		rows = values
	case "station":
		var stations []models.Station
		stations, err = db.GetStations(page, limit)
		rows = stations
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Invalid table name", r))))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
		return
	}

	jsonString, err := json.Marshal(rows)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
		return
	}

	specialSql := ""
	if table == "measurement" {
		stationParam := r.FormValue("station")
		if stationParam != "" {
			specialSql = fmt.Sprintf(` WHERE id_station = %s`, stationParam)
		}
	}
	pages, err := db.GetPageCount(table, limit, specialSql)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
		return
	}

	w.Header().Set("X-Page-Count", strconv.Itoa(pages))

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"total_pages":%d,"rows":%s}`, pages, jsonString)))

	_ = db.Disconnect()
}

func ReportHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	// get cookie which contains connection credentials ('token')
	cookie, err := GetCookie(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":401,"error":"%s"}`, L10n("Invalid credentials", r))))
		return
	}

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

	if name == "station" {
		var stations []models.StationReport
		stations, err = db.GetStationsReport()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(stations)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else if name == "measurement" {
		station := r.FormValue("station")
		start := r.FormValue("start")
		end := r.FormValue("end")

		var measurments []models.MeasurementReport
		measurments, err = db.GetMeasurementsReport(station, start, end)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(measurments)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else if name == "maxparticles" {
		start := r.FormValue("start")
		end := r.FormValue("end")

		var maxes []models.MaxPMByCity
		maxes, err = db.GetMaxPMByCity(start, end)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(maxes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else if name == "badparticles" {
		station := r.FormValue("station")

		var particles []models.Particles
		particles, err = db.GetBadParticles(station)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(particles)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else if name == "sulfur" {
		station := r.FormValue("station")

		var particles []models.Particles
		particles, err = db.GetSulfurParticles(station)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(particles)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else if name == "carbon" {
		station := r.FormValue("station")

		var particles []models.Particles
		particles, err = db.GetCarbonParticles(station)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, escaped)))
			return
		}

		jsonString, err := json.Marshal(particles)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			escaped := strings.ReplaceAll(err.Error(), `"`, `\"`)
			_, _ = w.Write([]byte(fmt.Sprintf(`{"status":500,"error":"%s"}`, escaped)))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":200,"data":%s}`, jsonString)))

		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"status":400,"error":"%s"}`, L10n("Invalid report name", r))))
		return
	}
}

func GetCookie(r *http.Request) (cookie *http.Cookie, err error) {
	cookies := r.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			return cookie, nil
		}
	}
	return nil, fmt.Errorf("no token cookie found")
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
