package main 

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)

// based on http://mindbowser.com/golang-go-database-sql/

func main() {
	db, err := sql.Open("mysql", "<user>:<password>@tcp(<ip>:3306)/UserDirectory")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connection established")
	}
	defer db.Close()
}


// SQL used:
// CREATE DATABASE UserDirectory;
// USE UserDirectory;
// CREATE TABLE User (
//   id  INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
//   email     VARCHAR(255)    NOT NULL,
//   create_date   DATETIME    NOT NULL,
//   password      VARCHAR(255)   NOT NULL,
//   last_name     VARCHAR(255),
//   first_name    VARCHAR(255),
// );
