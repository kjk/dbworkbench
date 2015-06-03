package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func importPostHistoryIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("posthistory",
		"id",
		"post_history_type_id",
		"post_id",
		"revision_guid",
		"creation_date",
		"user_id",
		"user_display_name",
		"text",
		"comment",
	))
	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		p := &r.PostHistory
		_, err = stmt.Exec(
			p.ID,
			p.PostHistoryTypeID,
			p.PostID,
			p.RevisionGUID,
			p.CreationDate,
			p.UserID,
			p.UserDisplayName,
			p.Text,
			p.Comment,
		)
		if err != nil {
			LogVerbosef("p: %+v\n", p)
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

func importPostHistory(siteName string, db *sql.DB) (int, error) {
	reader, err := getStackOverflowReader(siteName, "PostHistory")
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	r, err := stackoverflow.NewPostHistoryReader(reader)
	if err != nil {
		return 0, fmt.Errorf("stackoverflow.NewPostHistoryReader() failed with %s", err)
	}
	defer r.Close()
	return importPostHistoryIntoDB(r, db)
}
