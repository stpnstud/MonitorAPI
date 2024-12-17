package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func main() {
	if len(os.Args) <= 1 {
		printHelp()
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "--help":
		printHelp()

	case "--createdb":
		if _, err := os.Stat("monitors.txt"); os.IsNotExist(err) {
			fmt.Println("ERROR! File \"monitors.txt\" does not exist!")
			return
		}

		if _, err := os.Stat("./products.db"); err == nil {
			err = os.Remove("./products.db")

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		CreateDB()
		AddMonitorsFromFile("./monitors.txt")

		fmt.Println("OK. File product.db is created!")
		return

	case "--start":
		http.HandleFunc("/category/monitors", GetMonitors)
		http.HandleFunc("/category/monitor/", GetStatForMonitor)
		http.HandleFunc("/category/monitor_click/", AddClickForMonitor)

		fmt.Println("The server is running!")
		fmt.Println("Looking forward to requests...")

		if err := http.ListenAndServe(":8030", nil); err != nil {
			log.Fatal("Failed to start server!", err)
		}

	default:
		printHelp()
	}
}

func CreateDB() {
	OpenDB()

	_, err := DB.Exec("create table monitors (id integer, name varchar(255) not null, count integer)")
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	DB.Close()
}

func OpenDB() {
	db, err := sql.Open("sqlite3", "products.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	DB = db
}

func printHelp() {
	fmt.Println()
	fmt.Println("Help:                    ./counter --help")
	fmt.Println("Create products database: ./counter --createdb")
	fmt.Println("Start server:            ./counter --start")
	fmt.Println()
}

func AddMonitorsFromFile(filename string) {
	var file *os.File
	var err error

	if file, err = os.Open(filename); err != nil {
		log.Fatal("Failed to open the file:", err)
		os.Exit(2)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	OpenDB()

	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), ",")
		id := arr[0]
		monitorName := arr[1]

		_, err = DB.Exec("insert into monitors(id, name, count) values ($1, $2, 0)", id, monitorName)
	}
}

func AddClickForMonitor(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{%s}", err)
		return
	}

	OpenDB()

	countValue := 0

	monitorId := strings.TrimPrefix(request.URL.Path, "/category/monitor_click/")
	rows, _ := DB.Query("select count from monitors where id=" + monitorId)

	for rows.Next() {
		rows.Scan(&countValue)
	}

	countValue++
	_, err = DB.Exec("update monitors set count=" + strconv.Itoa(countValue) + " where id=" + monitorId)
}

func GetFromDBNameModel(tblName string) [][]string {
	var monitors [][]string
	var name string
	var id int

	rows, err := DB.Query("select id, name from " + tblName)

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		rows.Scan(&id, &name)
		monitors = append(monitors, []string{strconv.Itoa(id), name})
	}

	return monitors
}

func GetMonitors(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{%s}", err)
		return
	}

	OpenDB()

	monitors := GetFromDBNameModel("monitors")

	strOut := "{ \"monitors\": ["

	for i, arr := range monitors {
		strOut += "[" + arr[0] + "," + arr[1] + "]"

		if i != len(monitors)-1 {
			strOut += ", "
		}
	}

	strOut += "] }"

	fmt.Fprintf(w, strOut)
}

func GetStatForMonitor(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()

	if err != nil {
		fmt.Fprintf(w, "{%s}", err)
		return
	}

	OpenDB()

	countValue := 0

	monitorId := strings.TrimPrefix(request.URL.Path, "/category/monitor/")
	rows, _ := DB.Query("select count from monitors where id=" + monitorId)

	for rows.Next() {
		rows.Scan(&countValue)
	}

	strOut := "{ \"id\": \"" + monitorId + "\", \"count\": \"" + strconv.Itoa(countValue) + "\"}"

	fmt.Fprintf(w, strOut)
}
