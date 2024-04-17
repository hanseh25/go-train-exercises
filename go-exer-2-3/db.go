package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

type user struct {
	Id        int64
	Name      string
	Username  string
	Password  string
	CreatedAt time.Time
}

type credential struct {
	Id        int64
	Url       string
	Username  string
	Password  string
	CreatedAt time.Time
}

func dbAllCredentialsForUser(ctx context.Context, conn *pgx.Conn, userId int64) ([]credential, error) {
	rows, err := conn.Query(ctx, `
        select id, url, username, password, created_at
        from credentials
		inner join user_credential on credentials.id = user_credential.credential_id
		where credentials.id = $1`, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var credentials []credential

	for rows.Next() {
		var c credential
		err = rows.Scan(&c.Id, &c.Url, &c.Username, &c.Password, &c.CreatedAt)

		if err != nil {
			return nil, err
		}
		credentials = append(credentials, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return credentials, nil
}

func dbGetUserByUsername(ctx context.Context, conn *pgx.Conn, username string) ([]user, error) {
	rows, err := conn.Query(ctx, `
	select *
	from users
	where users.username = $1`, username)

	//log.Printf("row username %s", rows)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []user

	for rows.Next() {
		var u user
		err = rows.Scan(&u.Id, &u.Name, &u.Username, &u.Password, &u.CreatedAt)

		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func saveAllCredentialsForUser(ctx context.Context, conn *pgx.Conn, url string, username string, password string, userID int64) {
	sql := "INSERT INTO credentials (url, username, password)  VALUES ($1, $2, $3)  RETURNING id"

	// Execute the insert statement
	var lastID int64
	err := conn.QueryRow(ctx, sql, url, username, password).Scan(&lastID)

	log.Printf("Last inserted ID : %v", lastID)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("userID %v", userID)

	sql2 := "INSERT INTO user_credential VALUES ($2, $1)"

	res, err := conn.Exec(ctx, sql2, lastID, userID)

	log.Printf("results %s", res)

	if err != nil {
		log.Fatal(err)
		log.Fatal(res)
	}
}

func callMe() {
	log.Printf("hello")
}
