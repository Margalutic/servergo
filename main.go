package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
    "golang.org/x/crypto/bcrypt"
)

type DisplayInfo struct { //Структура
	ID_display int     `json:"ID_Display"`
	Diagonal   float32 `json:"diagonal"`
	Resolution string  `json:"resolution"`
	Type       string  `json:"type"`
	GSync      bool    `json:"gsync"`
}

type MonitorInfo struct { //Структура
	ID_monitor   int         `json:"ID_Monitor"`
	PowerVoltage int         `json:"powerVoltage"`
	Display      DisplayInfo `json:"display"`
	GSyncPremium bool        `json:"gSyncPremium"`
	IsCurved     bool        `json:"isCurved"`
}
type User struct {
    Username string `json:"username"`
    password string `json:"password"`
    Email    string `json:"email"`
    IsAdmin  bool   `json:"isAdmin"`
}


var db *sql.DB

func main() {
	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Запуск сервера...")
	startServer()
}

func startServer() {
	http.HandleFunc("/addDisplay", addDisplayHandler)
	http.HandleFunc("/addMonitor", addMonitorHandler)
	http.HandleFunc("/removeDisplay", removeDisplayHandler)
	http.HandleFunc("/removeMonitor", removeMonitorHandler)
	http.HandleFunc("/allDisplays", allDisplaysHandler)
	http.HandleFunc("/allMonitors", allMonitorsHandler)
	http.HandleFunc("/getMonitor", getMonitorHandler)
    http.HandleFunc("/registerUser", registerUserHandler)
   

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		panic(err)
	}
    
}
func registerUserHandler(w http.ResponseWriter, r *http.Request) {
    newUser := User{}
    body, _ := io.ReadAll(r.Body)
    err := json.Unmarshal(body, &newUser)
    if err != nil {
        http.Error(w, "Ошибка при разборе JSON", http.StatusBadRequest)
        return
    }
  var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec(`INSERT INTO users (username, password_users, email, isadmin)
        VALUES ($1, $2, $3, $4)`, newUser.Username, string(hashedPassword), newUser.Email, newUser.IsAdmin)
    if err != nil {
        http.Error(w, "Ошибка при вставке данных в БД", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Пользователь успешно зарегистрирован"))
}

func addMonitorHandler(w http.ResponseWriter, r *http.Request) {
	tempMonitor := MonitorInfo{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempMonitor)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var displayID int
	err = db.QueryRow("insert into displays (diagonal, resolution, type, gsync) values ($1, $2, $3, $4) returning ID",
		tempMonitor.Display.Diagonal, tempMonitor.Display.Resolution, tempMonitor.Display.Type, tempMonitor.Display.GSync).Scan(&displayID)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into Monitors (powervoltage, display_id, GSyncPremium, curved) values ($1, $2, $3, $4)",
		tempMonitor.PowerVoltage, displayID, tempMonitor.GSyncPremium, tempMonitor.IsCurved)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w.Write([]byte("Новый Монитор добавлен."))
}

func addDisplayHandler(w http.ResponseWriter, r *http.Request) {
	tempDisplay := DisplayInfo{}
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &tempDisplay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("insert into Displays (Diagonal, Resolution, type, Gsync) values ($1, $2, $3, $4)", tempDisplay.Diagonal, tempDisplay.Resolution, tempDisplay.Type, tempDisplay.GSync)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	w.Write([]byte("Новый Дисплей добавлен."))
}

func removeDisplayHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	displayID := string(body)

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("delete from Displays where id_display = $1", displayID)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w.Write([]byte("Вы удалили дисплей."))
}

func removeMonitorHandler(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	monitorID := string(body)

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("delete from Monitors where id_monitor = $1", monitorID)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	w.Write([]byte("Вы удалили монитор."))
}

func allDisplaysHandler(w http.ResponseWriter, r *http.Request) {

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Query("select * from Displays")
	if err != nil {
		log.Fatal(err)
	}

	var displays []DisplayInfo
	for res.Next() {
		var display DisplayInfo
		err := res.Scan(&display.ID_display, &display.Diagonal, &display.Resolution, &display.Type, &display.GSync)
		if err != nil {
			log.Fatal(err)
		}
		displays = append(displays, display)
	}
	defer db.Close()

	out, err := json.Marshal(displays)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func allMonitorsHandler(w http.ResponseWriter, r *http.Request) {

	var connStr = "user=postgres password=Pa$$vv0rd dbname=gomonitor sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	res, err := db.Query("select m.ID_Monitor, m.PowerVoltage, d.ID_Display, d.Diagonal, d.Resolution, d.type, d.Gsync, m.Gsync_premium, m.Curved from Monitors m join Displays d on m.Display_id = d.ID_Display")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var monitors []MonitorInfo
	for res.Next() {
		var monitor MonitorInfo
		err := res.Scan(&monitor.ID_monitor, &monitor.PowerVoltage, &monitor.Display.ID_display, &monitor.Display.Diagonal, &monitor.Display.Resolution, &monitor.Display.Type, &monitor.Display.GSync, &monitor.GSyncPremium, &monitor.IsCurved)
		if err != nil {
			log.Fatal(err)
		}
		monitors = append(monitors, monitor)
	}
	defer db.Close()

	out, err := json.Marshal(monitors)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func getMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitorID := r.URL.Query().Get("id")

	var connStr = "user=postgres password=Pa$$w0rd dbname=go_on_wed sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "select m.ID_Monitor, m.PowerVoltage, d.ID_Display, d.Diagonal, d.Resolution, d.type, d.Gsync, m.Gsync_premium, m.Curved from Monitors m join Displays d on m.Display_id = d.ID_Display where m.ID_Monitor = $1"
	var monitor MonitorInfo
	err = db.QueryRow(query, monitorID).Scan(&monitor.ID_monitor, &monitor.PowerVoltage, &monitor.Display.ID_display, &monitor.Display.Diagonal, &monitor.Display.Resolution, &monitor.Display.Type, &monitor.Display.GSync, &monitor.GSyncPremium, &monitor.IsCurved)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	out, err := json.Marshal(monitor)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
