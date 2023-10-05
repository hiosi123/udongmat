package db

import (
	"fmt"
	"math"
	"time"

	"github.com/hiosi123/udongmat/config"
	"github.com/jmoiron/sqlx"
)

func CreateMySqlConnection(env config.EnvVars) *sqlx.DB {
	var counts int64
	var backOff = 1 * time.Second
	var db *sqlx.DB
	var err error

	for {
		db, err = sqlx.Connect("mysql", env.DSN)
		if err != nil {
			counts++
		} else {
			fmt.Println("mysql db connected")
			break
		}

		if counts > 5 {
			fmt.Println(err)
			panic(err)
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
	}

	return db
}
