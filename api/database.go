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

func (db *Database) GetMqttServers(page, limit int) (servers []models.MqttServer, err error) {
	rows, err := db.Database.Query(`SELECT id_server, server_url, server_status FROM mqtt_server LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var url string
		var status models.JsonNullString
		err = rows.Scan(&id, &url, &status)
		if err != nil {
			return nil, err
		}
		servers = append(servers, models.MqttServer{
			Id:           id,
			ServerUrl:    url,
			ServerStatus: status,
		})
	}

	return servers, nil
}

func (db *Database) GetCategories(page, limit int) (categories []models.Category, err error) {
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

func (db *Database) GetMeasuredUnits(page, limit int) (units []models.MeasuredUnit, err error) {
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

func (db *Database) GetFavorites(page, limit int) (favorites []models.Favorite, err error) {
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

func (db *Database) GetCoordinates(page, limit int) (coordinates []models.Coordinates, err error) {
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

func (db *Database) GetMeasurements(page, limit int, station string) (measurements []models.Measurement, err error) {
	sql := `SELECT 
			ms.id_measurment,
			ms.measurment_time,
			ms.measurment_value,
			st.id_station,
			st.station_name, 
			mu.id_measured_unit,
			mu.title,
			mu.unit,
			mq.unit_message,
			mq.unit_order,
			ct.designation
		FROM measurment as ms
		LEFT JOIN station st ON st.id_station = ms.id_station
		LEFT JOIN measured_unit mu ON mu.id_measured_unit = ms.id_measured_unit
		LEFT JOIN mqtt_unit mq ON mq.id_station = st.id_station AND mq.id_measured_unit = ms.id_measured_unit
		LEFT JOIN optimal_value ov ON ov.id_measured_unit = mu.id_measured_unit
			AND case when ov.bottom_border is null then true else ov.bottom_border <= ms.measurment_value end
			AND case when ov.upper_border is null then true else ov.upper_border >= ms.measurment_value end
		LEFT JOIN category ct ON ct.id_category = ov.id_category`

	if station != "" {
		sql += ` WHERE st.id_station = ` + station
	}

	sql += ` LIMIT $1 OFFSET $2`

	rows, err := db.Database.Query(sql, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var time string
		var value float64
		var idStation int64
		var stationName string
		var idMeasuredUnit int64
		var title string
		var unit string
		var unitMessage models.JsonNullString
		var unitOrder models.JsonNullInt32
		var designation models.JsonNullString

		err = rows.Scan(&id, &time, &value, &idStation, &stationName, &idMeasuredUnit, &title, &unit, &unitMessage, &unitOrder, &designation)
		if err != nil {
			return nil, err
		}
		measurements = append(measurements, models.Measurement{
			Id:             id,
			Time:           time,
			Value:          value,
			IdStation:      idStation,
			StationName:    stationName,
			IdMeasuredUnit: idMeasuredUnit,
			UnitTitle:      title,
			UnitUnit:       unit,
			UnitMessage:    unitMessage,
			UnitOrder:      unitOrder,
			Designation:    designation,
		})
	}

	return measurements, nil
}

func (db *Database) GetMqttUnits(page, limit int) (units []models.MqttUnit, err error) {
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

func (db *Database) GetStations(page, limit int) (stations []models.Station, err error) {
	rows, err := db.Database.Query(`
		SELECT st.id_station, city, 
		   station_name, station_status, 
		   server_url, server_status, 
		   id_saveecobot, 
		   crd.coordinates[0] AS latitude, 
		   crd.coordinates[1] AS longitude
		FROM station as st
		LEFT JOIN mqtt_server as srv ON srv.id_server = st.id_server
		LEFT JOIN coordinates as crd ON crd.id_station = st.id_station
		LIMIT $1 OFFSET $2`, limit, page)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var city string
		var name string
		var stationStatus models.JsonNullString
		var serverUrl models.JsonNullString
		var serverStatus models.JsonNullString
		var idSaveEcoBot string
		var latitude float64
		var longitude float64

		err = rows.Scan(&id, &city, &name, &stationStatus, &serverUrl, &serverStatus, &idSaveEcoBot, &latitude, &longitude)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.Station{
			Id:            id,
			City:          city,
			Name:          name,
			StationStatus: stationStatus,
			ServerUrl:     serverUrl,
			ServerStatus:  serverStatus,
			IdSaveEcoBot:  idSaveEcoBot,
			Coordinates: models.Point{
				Lat: latitude,
				Lng: longitude,
			},
		})
	}

	return stations, nil
}

func (db *Database) GetOptimalValues(page, limit int) (values []models.OptimalValue, err error) {
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

func (db *Database) GetPageCount(tableName string, limit int, specialSql string) (count int, err error) {

	if tableName == "measurement" {
		tableName = "measurment" // TODO: fix typo in database
	}

	rows, err := db.Database.Query(`SELECT COUNT(*) FROM ` + tableName + ` ` + specialSql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count/limit + 1, nil
}

func (db *Database) GetStationsReport() (stations []models.StationReport, err error) {
	rows, err := db.Database.Query(`SELECT
    		st.station_name,
			mt.min AS first_date,
			su.units
		FROM Station st
		LEFT JOIN (SELECT ms.id_station, MIN(ms.measurment_time) 
				   FROM Measurment ms
				   GROUP BY ms.id_station) mt
		ON mt.id_station = st.id_station
		LEFT JOIN (SELECT id_station, STRING_AGG(title||' ('||unit||')', ', ') AS units
				   FROM (SELECT DISTINCT ms.id_station, mu.title, mu.unit
				   FROM Measurment ms
				   LEFT JOIN Measured_Unit mu
				   ON mu.id_measured_unit = ms.id_measured_unit) un
				   GROUP BY id_station) su
		ON su.id_station = st.id_station`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var station string
		var firstDate models.JsonNullString
		var units models.JsonNullString

		err = rows.Scan(&station, &firstDate, &units)
		if err != nil {
			return nil, err
		}
		stations = append(stations, models.StationReport{
			StationName: station,
			FirstDate:   firstDate,
			Units:       units,
		})
	}

	return stations, nil
}

func (db *Database) GetMeasurementsReport(station, dateFrom, dateTo string) (reports []models.MeasurementReport, err error) {
	rows, err := db.Database.Query(`SELECT
			MIN(ms.measurment_value),
			AVG(ms.measurment_value),
			MAX(ms.measurment_value),
			COUNT(ms.measurment_value),
			mu.title||' ('||mu.unit||')' AS unit
		FROM measurment ms
		LEFT JOIN Station st ON st.id_station = ms.id_station
		LEFT JOIN measured_unit mu ON mu.id_measured_unit = ms.id_measured_unit
		WHERE st.station_name = $1
		AND ms.measurment_time >= $2
		AND ms.measurment_time <= $3
		GROUP BY mu.title, mu.unit`, station, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var minV float64
		var avgV float64
		var maxV float64
		var count int
		var unit string

		err = rows.Scan(&minV, &avgV, &maxV, &count, &unit)
		if err != nil {
			return nil, err
		}
		reports = append(reports, models.MeasurementReport{
			Min:   minV,
			Max:   maxV,
			Avg:   avgV,
			Count: count,
			Unit:  unit,
		})
	}

	return reports, nil
}

func (db *Database) GetMaxPMByCity(dateFrom, dateTo string) (reports []models.MaxPMByCity, err error) {
	rows, err := db.Database.Query(`SELECT
    			st.city, 
				CASE WHEN ms.id_measured_unit = 2 THEN 'PM10' ELSE 'PM2.5' END AS unit,
    			MAX(ms.measurment_value) AS max_value
			FROM station st
			LEFT JOIN Measurment ms ON ms.id_station = st.id_station AND (ms.id_measured_unit = 3 OR ms.id_measured_unit = 2)
			WHERE ms.measurment_time >= $1
			AND ms.measurment_time <= $2
			GROUP BY st.city, ms.id_measured_unit`, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city string
		var unit string
		var maxV float64

		err = rows.Scan(&city, &unit, &maxV)
		if err != nil {
			return nil, err
		}
		reports = append(reports, models.MaxPMByCity{
			City: city,
			Unit: unit,
			Max:  maxV,
		})
	}

	return reports, nil
}

func (db *Database) GetBadParticles(station string) (results []models.Particles, err error) {
	rows, err := db.Database.Query(`SELECT 
    		designation, COUNT(measurment_value)
		FROM Station st
		LEFT JOIN Measurment ms ON ms.id_station = st.id_station
		LEFT JOIN measured_unit mu ON mu.id_measured_unit = ms.id_measured_unit
		LEFT JOIN mqtt_unit mq ON mq.id_station = st.id_station AND mq.id_measured_unit = ms.id_measured_unit
		LEFT JOIN optimal_value ov ON ov.id_measured_unit = mu.id_measured_unit
			AND case when ov.bottom_border is null then true else ov.bottom_border <= ms.measurment_value end
			AND case when ov.upper_border is null then true else ov.upper_border >= ms.measurment_value end
		LEFT JOIN category ct ON ct.id_category = ov.id_category
		WHERE mu.id_measured_unit = 3 AND ct.id_category > 3
		AND st.station_name = $1
		GROUP BY designation`, station)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var designation string
		var count int

		err = rows.Scan(&designation, &count)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Particles{
			Designation: designation,
			Count:       count,
		})
	}

	return results, nil
}

func (db *Database) GetSulfurParticles(station string) (results []models.Particles, err error) {
	rows, err := db.Database.Query(`SELECT 
    		designation, COUNT(measurment_value)
		FROM Station st
		LEFT JOIN Measurment ms ON ms.id_station = st.id_station
		LEFT JOIN measured_unit mu ON mu.id_measured_unit = ms.id_measured_unit
		LEFT JOIN mqtt_unit mq ON mq.id_station = st.id_station AND mq.id_measured_unit = ms.id_measured_unit
		LEFT JOIN optimal_value ov ON ov.id_measured_unit = mu.id_measured_unit
		  AND case when ov.bottom_border is null then true else ov.bottom_border <= ms.measurment_value end
		  AND case when ov.upper_border is null then true else ov.upper_border >= ms.measurment_value end
		LEFT JOIN category ct ON ct.id_category = ov.id_category
		WHERE (mu.id_measured_unit = 13 OR mu.id_measured_unit = 15)
		AND st.station_name = $1
		GROUP BY designation`, station)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var designation string
		var count int

		err = rows.Scan(&designation, &count)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Particles{
			Designation: designation,
			Count:       count,
		})
	}

	return results, nil
}

func (db *Database) GetCarbonParticles(station string) (results []models.Particles, err error) {
	rows, err := db.Database.Query(`SELECT 
    		designation, COUNT(measurment_value)
		FROM Station st
		LEFT JOIN Measurment ms ON ms.id_station = st.id_station
		LEFT JOIN measured_unit mu ON mu.id_measured_unit = ms.id_measured_unit
		LEFT JOIN mqtt_unit mq ON mq.id_station = st.id_station AND mq.id_measured_unit = ms.id_measured_unit
		LEFT JOIN optimal_value ov ON ov.id_measured_unit = mu.id_measured_unit
		  AND case when ov.bottom_border is null then true else ov.bottom_border <= ms.measurment_value end
		  AND case when ov.upper_border is null then true else ov.upper_border >= ms.measurment_value end
		LEFT JOIN category ct ON ct.id_category = ov.id_category
		WHERE mu.id_measured_unit = 12
		AND st.station_name = $1
		GROUP BY designation`, station)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var designation string
		var count int

		err = rows.Scan(&designation, &count)
		if err != nil {
			return nil, err
		}
		results = append(results, models.Particles{
			Designation: designation,
			Count:       count,
		})
	}

	return results, nil
}
