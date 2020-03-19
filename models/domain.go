package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Domain struct {
	Id        int    `json:"id"`
	Domain    string `json:"domain"`
	CreatedAt int64  `json:"created_at"`
	Hits int
}

const DOMAIN_SELECT = "SELECT id, domain, UNIX_TIMESTAMP(created_at) as createdat from domains"

func SelectDomains(db *sqlx.DB, userId int) ([]*Domain, string) {
	items := []*Domain{}
	sql := fmt.Sprintf("%s where user_id=:id order by created_at desc", DOMAIN_SELECT)
	if userId == 0 {
		sql = fmt.Sprintf("%s order by created_at desc", DOMAIN_SELECT)
	}
	rows, err := db.NamedQuery(sql, map[string]interface{}{"id": userId})
	if err != nil {
		return items, err.Error()
	}
	for rows.Next() {
		item := Domain{}
		rows.StructScan(&item)
		items = append(items, &item)
	}

	return items, ""
}

func InsertDomain(db *sqlx.DB, domain string, userId int) string {
	_, err := db.NamedExec("INSERT INTO domains (domain, user_id) values (:domain, :id)",
		map[string]interface{}{"domain": domain, "id": userId})
	if err != nil {
		return err.Error()
	}
	return ""
}
