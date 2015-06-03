package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func importVotesIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("votes",
		"id",
		"post_id",
		"vote_type_id",
		"user_id",
		"bounty_amount",
		"creation_date",
	))
	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		v := &r.Vote
		_, err = stmt.Exec(
			v.ID,
			v.PostID,
			v.VoteTypeID,
			v.UserID,
			v.BountyAmount,
			v.CreationDate,
		)
		if err != nil {
			err = fmt.Errorf("stmt.Exec() 1 failed with %s", err)
			return 0, err
		}
		n++
	}
	if err = r.Err(); err != nil {
		return 0, err
	}
	_, err = stmt.Exec()
	if err != nil {
		err = fmt.Errorf("stmt.Exec() 2 failed with %s", err)
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

func importVotes(siteName string, db *sql.DB) error {
	reader, err := getStackOverflowReader(siteName, "Votes")
	if err != nil {
		return err
	}
	defer reader.Close()

	r, err := stackoverflow.NewVotesReader(reader)
	if err != nil {
		return fmt.Errorf("stackoverflow.NewVotesReader() failed with %s", err)
	}
	defer r.Close()
	n, err := importVotesIntoDB(r, db)
	if err != nil {
		return fmt.Errorf("importVotesIntoDB() failed with %s", err)
	}
	LogVerbosef("processed %d votes records\n", n)
	return nil
}
