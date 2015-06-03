package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func importPostLinksIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("postlinks",
		"id",
		"creation_date",
		"post_id",
		"related_post_id",
		"link_type_id",
	))
	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		l := &r.PostLink
		_, err = stmt.Exec(
			l.ID,
			l.CreationDate,
			l.PostID,
			l.RelatedPostID,
			l.LinkTypeID,
		)
		if err != nil {
			LogVerbosef("l: %+v\n", l)
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

func importPostLinks(siteName string, db *sql.DB) (int, error) {
	reader, err := getStackOverflowReader(siteName, "PostLinks")
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	r, err := stackoverflow.NewPostLinksReader(reader)
	if err != nil {
		return 0, fmt.Errorf("stackoverflow.NewPostLinksReader() failed with %s", err)
	}
	defer r.Close()
	return importPostLinksIntoDB(r, db)
}
