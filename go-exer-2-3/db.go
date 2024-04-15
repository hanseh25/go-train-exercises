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

	log.Printf("row username %s", rows)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var credentials []credential

	for rows.Next() {
		var c credential
		err = rows.Scan(&c.Id, &c.Url, &c.Username, &c.Password, &c.CreatedAt)
		log.Printf("row username %s", c.Username)
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

	log.Printf("row username %s", rows)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []user

	for rows.Next() {
		var u user
		err = rows.Scan(&u.Id, &u.Name, &u.Username, &u.Password, &u.CreatedAt)
		log.Printf("row username %s", u.Username)
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

func callMe() {
	log.Printf("hello")
}
