package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// connect to a databased server
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=postgres password=password")
	if err != nil {
		log.Fatalf("Failed to connect to: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to database!")

	// test my connection
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping database!")
	}

	log.Println("Pinged database!")

	// get rows from table
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("Failed to get rows from database: %v", err)
	}

	// insert a row
	query := `insert into users (first_name, last_name) values ($1, $2)`
	//_, err = conn.Exec(insQuery, "Helga", "Hinderberhg")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted a row!")

	//get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("Failed to get rows from database: %v", err)
	}

	// update a row
	stmt := `update users set first_name = $1 where first_name = $2 and id = $3`
	_, err = conn.Exec(stmt, "Therese", "Helga", 7)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Updated one or more rows!")

	// get rows from table again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("Failed to get rows from database: %v", err)
	}

	// get one row by id
	query = `select id, first_name, last_name from users where id = $1`
	var first_name, last_name string
	var id int
	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &first_name, &last_name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Record is:", id, first_name, last_name)
	fmt.Println("-------------------------------------------------------------")

	// delete a row
	query = `delete from users where id = $1`
	_, err = conn.Exec(query, 9)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted a row!")

	// get rows again
	err = getAllRows(conn)
	if err != nil {
		log.Fatalf("Failed to get rows from database: %v", err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users order by id asc")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	var first_name, last_name string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &first_name, &last_name)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is:", id, first_name, last_name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error scanning rows: ", err)
	}

	fmt.Println("-------------------------------------------------------------")

	return nil
}
