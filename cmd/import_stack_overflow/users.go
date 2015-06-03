package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func importUsersIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := txn.Prepare(pq.CopyIn("users",
		"id",
		"reputation",
		"creation_date",
		"display_name",
		"last_access_date",
		"website_url",
		"location",
		"about_me",
		"views",
		"up_votes",
		"down_votes",
		"age",
		"account_id",
		"profile_image_url",
	))

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		u := &r.User
		_, err = stmt.Exec(
			u.ID,
			u.Reputation,
			u.CreationDate,
			toStringPtr(u.DisplayName),
			toTimePtr(u.LastAccessDate),
			toStringPtr(u.WebsiteURL),
			toStringPtr(u.Location),
			toStringPtr(u.AboutMe),
			u.Views,
			u.UpVotes,
			u.DownVotes,
			u.Age,
			u.AccountID,
			toStringPtr(u.ProfileImageURL),
		)
		if err != nil {
			LogVerbosef("n: %+v\n", u)
			err = fmt.Errorf("stmt.Exec() failed with %s", err)
			return 0, err
		}
		n++
	}
	if err = r.Err(); err != nil {
		return 0, err
	}
	_, err = stmt.Exec()
	if err != nil {
		err = fmt.Errorf("stmt.Exec() failed with %s", err)
		return 0, err
	}
	err = stmt.Close()
	stmt = nil
	if err != nil {
		err = fmt.Errorf("stmt.Close() failed with %s", err)
		return 0, err
	}
	err = txn.Commit()
	txn = nil
	if err != nil {
		err = fmt.Errorf("txn.Commit() failed with %s", err)
		return 0, err
	}
	return n, nil
}

func importUsers(siteName string, db *sql.DB) error {
	reader, err := getStackOverflowReader(siteName, "Users")
	if err != nil {
		return err
	}
	defer reader.Close()
	r, err := stackoverflow.NewUsersReader(reader)
	if err != nil {
		return fmt.Errorf("stackoverflow.NewUsersReader() failed with %s", err)
	}
	defer r.Close()
	n, err := importUsersIntoDB(r, db)
	if err != nil {
		return fmt.Errorf("importUsersIntoDB() failed with %s", err)
	}
	LogVerbosef("processed %d user records\n", n)
	return nil
}
