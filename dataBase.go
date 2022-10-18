package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

//Todo добавить таблицы для userName и других данных о пользователе

type DataBase struct {
	pool *pgxpool.Pool
}

func (db *DataBase) Connect() {
	var err error

	db.pool, err = pgxpool.Connect(context.Background(), DataBaseUrl)
	if err != nil {
		log.Printf("Unable to connect to database\n%s\n", err)
		panic(err)
	}

	db.CreateTablesIfNotExists()
}

func (db *DataBase) Close() {
	db.pool.Close()
}

func (db *DataBase) Add(message Message) {
	_, err := db.pool.Exec(context.Background(),
		`insert into users ("userId", "phoneNumber") values ($1, $2) on conflict ("userId", "phoneNumber") do nothing`,
		message.Contact.UserId, message.Contact.PhoneNumber)
	if err != nil {
		log.Printf("Error DataBase Exec insert users (userId, phoneNumber)\n%s\n", err)
		return
	}
}

func (db *DataBase) CreateTablesIfNotExists() {
	//users
	_, err := db.pool.Exec(context.Background(),
		`CREATE TABLE if not exists public.users
(
    id bigserial NOT NULL,
    "userId" bigint NOT NULL,
    "phoneNumber" character varying(24) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE ("userId", "phoneNumber")
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;`)
	if err != nil {
		log.Printf("Error DataBase create table users\n%s\n", err)
		return
	}
}
