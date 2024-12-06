package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type TextEmail struct {
	ID          int
	DESCRIPTION string
}

func GetTextEmail(db *sql.DB) ([]TextEmail, error) {
	rows, err := db.Query("select * from form_email")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var textemail []TextEmail
	for rows.Next() {
		var text TextEmail
		if err := rows.Scan(&text.ID, &text.DESCRIPTION); err != nil {
			return nil, err
		}
		textemail = append(textemail, text)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return textemail, nil
}
