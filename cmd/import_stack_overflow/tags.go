package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/lzmadec"
	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func importTagsIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("tags",
		"id",
		"tag_name",
		"count",
		"excerpt_post_id",
		"wiki_post_id",
	))
	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		t := &r.Tag
		_, err = stmt.Exec(
			t.ID,
			t.TagName,
			t.Count,
			t.ExcerptPostID,
			t.WikiPostID,
		)
		if err != nil {
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

func importTags(archive *lzmadec.Archive, db *sql.DB) error {
	name := "Tags.xml"
	entry := getEntryForFile(archive, name)
	if entry == nil {
		return fmt.Errorf("genEntryForFile('%s') returned nil", name)
	}

	reader, err := archive.ExtractReader(entry.Path)
	if err != nil {
		return fmt.Errorf("ExtractReader('%s') failed with %s", entry.Path, err)
	}
	defer reader.Close()
	r, err := stackoverflow.NewTagsReader(reader)
	if err != nil {
		return fmt.Errorf("stackoverflow.NewTagsReader() failed with %s", err)
	}
	defer r.Close()
	n, err := importTagsIntoDB(r, db)
	if err != nil {
		return fmt.Errorf("importTagsIntoDB() failed with %s", err)
	}
	LogVerbosef("processed %d tags records\n", n)
	return nil
}
