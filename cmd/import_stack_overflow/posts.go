package main

import (
	"database/sql"
	"fmt"

	"github.com/kjk/stackoverflow"
	"github.com/lib/pq"
)

func toTags(tags []string) *string {
	var res string
	for _, tag := range tags {
		res = res + "<" + tag + ">"
	}
	if res == "" {
		return nil
	}
	return &res
}

func importPostsIntoDB(r *stackoverflow.Reader, db *sql.DB) (int, error) {
	txn, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if txn != nil {
			txn.Rollback()
		}
	}()

	stmt, err := txn.Prepare(pq.CopyIn("posts",
		"id",
		"post_type_id",
		"parent_id",
		"accepted_answer_id",
		"creation_date",
		"score",
		"view_count",
		"body",
		"owner_user_id",
		"owner_display_name",
		"last_editor_user_id",
		"last_editor_display_name",
		"last_edit_date",
		"last_activity_date",
		"title",
		"tags",
		"answer_count",
		"comment_count",
		"favorite_count",
		"community_owned_date",
		"closed_date"))
	if err != nil {
		err = fmt.Errorf("txt.Prepare() failed with %s", err)
		return 0, err
	}
	n := 0
	for r.Next() {
		p := &r.Post
		if false && p.ID == 41258 {
			fmt.Printf("p: %#v\n", p)
		}
		_, err = stmt.Exec(
			p.ID,
			p.PostTypeID,
			toIntPtr(p.ParentID),
			toIntPtr(p.AcceptedAnswerID),
			p.CreationDate,
			p.Score,
			p.ViewCount,
			p.Body,
			p.OwnerUserID,
			p.OwnerDisplayName,
			toIntPtr(p.LastEditorUserID),
			toStringPtr(p.LastEditorDisplayName),
			toTimePtr(p.LastEditDate),
			toTimePtr(p.LastActivitityDate),
			p.Title,
			toTags(p.Tags),
			p.AnswerCount,
			p.CommentCount,
			p.FavoriteCount,
			toTimePtr(p.CommunityOwnedDate),
			toTimePtr(p.ClosedDate),
		)
		if err != nil {
			LogVerbosef("p: #+v\n", p)
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

func importPosts(siteName string, db *sql.DB) (int, error) {
	reader, err := getStackOverflowReader(siteName, "Posts")
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	r, err := stackoverflow.NewPostsReader(reader)
	if err != nil {
		return 0, fmt.Errorf("stackoverflow.NewPostsReader() failed with %s", err)
	}
	defer r.Close()
	return importPostsIntoDB(r, db)
}
