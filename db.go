package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connectString string

// Flip - stage of protocol for specific user
type Flip struct {
	userID     int
	stage      int
	userSeed   string
	serverSeed string
	coinSide   string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error when loading env variables")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	fmt.Println(port)
	fmt.Println(host)
	fmt.Println(user)
	fmt.Println(password)
	fmt.Println(dbname)
	connectString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	createTable()
}

func getDB() *sql.DB {
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Println(err)
	}

	return db
}

// create table if not exists
func createTable() {
	db := getDB()
	defer db.Close()

	queryString := `CREATE TABLE IF NOT EXISTS coin_flip (
        user_id integer NOT NULL,
        stage integer NOT NULL, 
        user_seed varchar(200), 
        server_seed varchar(200), 
        coin_side varchar(200), 
        PRIMARY KEY (user_id))`

	_, err := db.Exec(queryString)
	if err != nil {
		fmt.Println(err)
	}
}

// create row in table with protocol stage
func createUpdateFlip(flip Flip) {
	db := getDB()
	defer db.Close()

	queryString := `INSERT INTO coin_flip (user_id, stage)
        VALUES
        (
            $1,
            $2
        ) 
        ON CONFLICT (user_id) 
        DO
            UPDATE
            SET stage = $2, user_seed = $3, server_seed = $4, coin_side = $5`

	_, err := db.Exec(queryString, flip.userID, flip.stage, flip.userSeed, flip.serverSeed, flip.coinSide)
	if err != nil {
		fmt.Println(err)
	}
}

// fetch current protocol stage for specific user
func getFlipStage(userID int) Flip {
	db := getDB()
	defer db.Close()

	queryString := `SELECT user_id, stage, user_seed, server_seed, coin_side FROM coin_flip WHERE user_id = $1`

	var flip = Flip{}
	row := db.QueryRow(queryString, userID)

	err := row.Scan(&flip.userID, &flip.stage, &flip.userSeed, &flip.serverSeed, &flip.coinSide)
	if err != nil {
		fmt.Println(err)
	}

	return flip
}
