package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/lzmadec"
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

func importPostHistory(archive *lzmadec.Archive, db *sql.DB) error {
	name := "PostHistory.xml"
	entry := getEntryForFile(archive, name)
	if entry == nil {
		return fmt.Errorf("genEntryForFile('%s') returned nil", name)
	}

	reader, err := archive.ExtractReader(entry.Path)
	if err != nil {
		return fmt.Errorf("ExtractReader('%s') failed with %s", entry.Path, err)
	}
	defer reader.Close()
	r, err := stackoverflow.NewPostHistoryReader(reader)
	if err != nil {
		return fmt.Errorf("stackoverflow.NewPostHistoryReader() failed with %s", err)
	}
	defer r.Close()
	n, err := importPostHistoryIntoDB(r, db)
	if err != nil {
		return fmt.Errorf("importPostHistoryIntoDB() failed with %s", err)
	}
	LogVerbosef("processed %d postshistory records\n", n)
	return nil
}
