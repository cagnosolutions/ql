package main

import (
	"fmt"
	"lab/ql"
	"log"
	"strings"
)

func main() {

	qry := "select * from users where id = 23, name ^ greg"
	r := strings.NewReader(qry)
	parser := ql.NewParser(r)
	stmt, err := parser.Parse()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%+v\n", stmt)
}
