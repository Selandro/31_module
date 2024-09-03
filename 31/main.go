package main

import (
	"fmt"
	"log"

	"main.go/pkg/storage/postgres"
)

func main() {
	dbInfo := "host=localhost user=postgres password=admin dbname=postgres port=5432 sslmode=disable"
	db, err := postgres.Connect(dbInfo)
	if err != nil {
		log.Fatal(err)
	}

	//вывод всех имеющихся тасков по id, автору или тегу
	tasks, err := db.Tasks(0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tasks)

	//создание таски
	a := postgres.Task{
		ID:         1,
		AuthorID:   1,
		AssignedID: 1,
		Title:      "title",
		Content:    "content",
	}
	db.NewTask(a)
	tasks, err = db.Tasks(0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tasks)

	//изменение таски передается объект Task, из него извлекается id, по которому
	//апдейтится таска
	b := postgres.Task{
		ID:         1,
		AuthorID:   1,
		AssignedID: 1,
		Title:      "обновлено",
		Content:    "обновлено",
	}
	err = db.UpdateTask(b)
	if err != nil {
		fmt.Println(err)
	}
	tasks, err = db.Tasks(0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tasks)

	//удаление таски по id
	err = db.DeleteTask(1)
	if err != nil {
		fmt.Println(err)
	}
	tasks, err = db.Tasks(0, 0, 0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tasks)

}
