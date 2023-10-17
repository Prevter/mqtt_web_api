package api

import (
	"database/sql"
	_ "github.com/lib/pq"
	"mqtt_web_api/api/models"
	"os"
)

type Database struct {
	ConnectionString string
	Database         *sql.DB
}

func LoginDatabase(login, password string) (db *Database, err error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connStr := "postgres://" + login + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable"

	db, e := NewDatabase(connStr)
	if e != nil {
		return nil, e
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewDatabase(connectionString string) (db *Database, err error) {
	database, e := sql.Open("postgres", connectionString)
	if e != nil {
		return nil, e
	}
	return &Database{
		ConnectionString: connectionString,
		Database:         database,
	}, nil
}

func (db *Database) Disconnect() (err error) {
	return db.Database.Close()
}

func (db *Database) Ping() (err error) {
	return db.Database.Ping()
}

func (db *Database) GetMqttServers(page, limit string) (servers []models.MqttServer, err error) {
	rows, err := db.Database.Query(`SELECT id_server, server_url, server_status FROM mqtt_server LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var url string
		var status string
		err = rows.Scan(&id, &url, &status)
		if err != nil {
			return nil, err
		}
		servers = append(servers, models.MqttServer{
			Id:           id,
			ServerUrl:    url,
			ServerStatus: models.ServerStatus(status),
		})
	}

	return servers, nil
}

func (db *Database) GetCategories(page, limit string) (categories []models.Category, err error) {
	rows, err := db.Database.Query(`SELECT id_category, designation FROM category LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var designation string
		err = rows.Scan(&id, &designation)
		if err != nil {
			return nil, err
		}
		categories = append(categories, models.Category{
			Id:          id,
			Designation: designation,
		})
	}

	return categories, nil
}

func (db *Database) GetMeasuredUnits(page, limit string) (units []models.MeasuredUnit, err error) {
	rows, err := db.Database.Query(`SELECT id_measured_unit, title, unit FROM measured_unit LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var title string
		var unit string
		err = rows.Scan(&id, &title, &unit)
		if err != nil {
			return nil, err
		}
		units = append(units, models.MeasuredUnit{
			Id:    id,
			Title: title,
			Unit:  unit,
		})
	}

	return units, nil
}

func (db *Database) GetFavorites(page, limit string) (favorites []models.Favorite, err error) {
	rows, err := db.Database.Query(`SELECT user_name, id_station FROM favorite LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var idStation int64

		err = rows.Scan(&username, &idStation)
		if err != nil {
			return nil, err
		}

		favorites = append(favorites, models.Favorite{
			Username:  username,
			IdStation: idStation,
		})
	}

	return favorites, nil
}

func (db *Database) GetCoordinates(page, limit string) (coordinates []models.Coordinates, err error) {
	rows, err := db.Database.Query(`SELECT id_station, coordinates[0] AS latitude, coordinates[1] AS longitude FROM coordinates LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idStation int64
		var latitude float64
		var longitude float64

		err = rows.Scan(&idStation, &latitude, &longitude)
		if err != nil {
			return nil, err
		}
		coordinates = append(coordinates, models.Coordinates{
			IdStation: idStation,
			Coordinates: models.Point{
				Lat: latitude,
				Lng: longitude,
			},
		})
	}

	return coordinates, nil
}

func (db *Database) GetMeasurements(page, limit string) (measurements []models.Measurement, err error) {
	rows, err := db.Database.Query(`SELECT id_measurment, measurment_time, measurment_value, id_station, id_measured_unit FROM measurment LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var time string
		var value float64
		var idStation int64
		var idMeasuredUnit int64

		err = rows.Scan(&id, &time, &value, &idStation, &idMeasuredUnit)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, models.Measurement{
			Id:             id,
			Time:           time,
			Value:          value,
			IdStation:      idStation,
			IdMeasuredUnit: idMeasuredUnit,
		})
	}

	return measurements, nil
}

func (db *Database) GetMqttUnits(page, limit string) (units []models.MqttUnit, err error) {
	rows, err := db.Database.Query(`SELECT id_station, id_measured_unit, unit_message, unit_order FROM mqtt_unit LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idStation int64
		var idMeasuredUnit int64
		var unitMessage string
		var unitOrder int

		err = rows.Scan(&idStation, &idMeasuredUnit, &unitMessage, &unitOrder)
		if err != nil {
			return nil, err
		}
		units = append(units, models.MqttUnit{
			IdStation:      idStation,
			IdMeasuredUnit: idMeasuredUnit,
			UnitMessage:    unitMessage,
			UnitOrder:      unitOrder,
		})
	}

	return units, nil
}

func (db *Database) GetStations(page, limit string) (stations []models.Station, err error) {
	rows, err := db.Database.Query(`SELECT id_station, city, station_name, station_status, id_server, id_saveecobot FROM station LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var city string
		var name string
		var status string
		var idServer models.JsonNullInt64
		var idSaveEcoBot string

		err = rows.Scan(&id, &city, &name, &status, &idServer, &idSaveEcoBot)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.Station{
			Id:            id,
			City:          city,
			Name:          name,
			StationStatus: models.ServerStatus(status),
			IdServer:      idServer,
			IdSaveEcoBot:  idSaveEcoBot,
		})
	}

	return stations, nil
}

func (db *Database) GetOptimalValues(page, limit string) (values []models.OptimalValue, err error) {
	rows, err := db.Database.Query(`SELECT id_category, id_measured_unit, bottom_border, upper_border FROM optimal_value LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idCategory int64
		var idMeasuredUnit int64
		var bottomBorder int
		var upperBorder models.JsonNullInt32

		err = rows.Scan(&idCategory, &idMeasuredUnit, &bottomBorder, &upperBorder)
		if err != nil {
			return nil, err
		}
		values = append(values, models.OptimalValue{
			IdCategory:     idCategory,
			IdMeasuredUnit: idMeasuredUnit,
			BottomBorder:   bottomBorder,
			UpperBorder:    upperBorder,
		})
	}

	return values, nil
}
