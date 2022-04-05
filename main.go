package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type Poem []struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

type Author []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AuthorList []struct {
	Id       int `json:"id"`
	PoemId   int `json:"poem_id"`
	AuthorId int `json:"author_id"`
}

func main() {
	// ! Poems.json
	jsonFile, err := os.Open("poems.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var poem Poem

	if err := json.Unmarshal(byteValue, &poem); err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println("Successfully Write users.json")

	// ! Authors.json
	jsonFile, err = os.Open("authors.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened authors.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	var author Author

	if err := json.Unmarshal(byteValue, &author); err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println("Successfully Write authors.json")

	// ! Authors.json
	jsonFile, err = os.Open("authorslist.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened authorslist.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ = ioutil.ReadAll(jsonFile)

	var authorList AuthorList

	if err := json.Unmarshal(byteValue, &authorList); err != nil {
		fmt.Println(err.Error())
		return
	}
	
	fmt.Println("Successfully Write authorslist.json")

	db, err := sqlx.Open("postgres", "host='localhost' port='5432' user='postgres' dbname='postgres' password='qwerty123' sslmode='disable'")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Starting write poems")
	query := "INSERT INTO poems (title, text) VALUES ($1, $2) RETURNING id"
	var id int
	for _, p := range poem {
		row := tx.QueryRow(query, p.Title, p.Text)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
		// fmt.Println("Write Complete: " + strconv.Itoa(id))
	}
	fmt.Println("Finish write poems")

	fmt.Println("Starting write authors")
	query = "INSERT INTO authors (name) VALUES ($1) RETURNING id"
	for _, a := range author {
		row := tx.QueryRow(query, a.Name)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			fmt.Println(err)

			return
		}
		// fmt.Println("Write Complete: " + strconv.Itoa(id))
	}
	fmt.Println("Finish write authors")

	fmt.Println("Starting write authors_list")
	query = "INSERT INTO authors_list (author_id, poem_id) VALUES ($1, $2) RETURNING id"
	for _, a := range authorList {
		row := tx.QueryRow(query, a.AuthorId, a.PoemId)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			fmt.Println(err)

			return
		}
		// fmt.Println("Write Complete: " + strconv.Itoa(id))
	}
	fmt.Println("Finish write authors_list")

	if err := tx.Commit(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Write Complete!!!")

	// var authorList AuthorList
	// i := 0
	// for _, p := range poem {
	// 	for _, a := range author {
	// 		if p.Author == a.Name {
	// 			authorList = append(authorList, struct {
	// 				Id       int "json:\"id\""
	// 				PoemId   int "json:\"poem_id\""
	// 				AuthorId int "json:\"author_id\""
	// 			}{Id: i, AuthorId: a.Id, PoemId: p.Id})
	// 			i++
	// 		}
	// 	}
	// }

	// file, _ := json.MarshalIndent(authorList, "", " ")

	// _ = ioutil.WriteFile("authorslist.json", file, 0644)
}

// fmt.Println("Id: " + strconv.Itoa(p.Id))
// fmt.Println("Title: " + p.Title)
// fmt.Println("Text: " + p.Text)
// fmt.Println("Author: " + p.Author)
