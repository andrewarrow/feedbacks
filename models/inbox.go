package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"html/template"

	"github.com/emersion/go-message"
	"github.com/jmoiron/sqlx"
)

var SELECT_INBOX = "SELECT id, sent_to as sentto, sent_from as sentfrom, body, subject, UNIX_TIMESTAMP(created_at) as createdat from inboxes order by created_at desc limit 1000"

type Inbox struct {
	Id          int    `json:"id"`
	SentTo      string `json:"sent_to"`
	SentFrom    string `json:"sent_from"`
	Body        string `json:"body"`
	Subject     string `json:"subject"`
	CreatedAt   int64  `json:"created_at"`
	MessageText template.HTML
	MessageHTML template.HTML
}

func SelectInboxes(db *sqlx.DB) ([]Inbox, string) {
	items := []Inbox{}
	err := db.Select(&items, SELECT_INBOX)
	s := ""
	if err != nil {
		s = err.Error()
	}

	return items, s
}
func SelectInboxByDomain(db *sqlx.DB, domain string) ([]Inbox, string) {
	items := []Inbox{}
	sql := fmt.Sprintf("select sent_to as sentto, sent_from as sentfrom, body, subject, UNIX_TIMESTAMP(created_at) as createdat from inbox where sent_to like :domain order by created_at desc")
	rows, err := db.NamedQuery(sql, map[string]interface{}{"domain": "%" + domain})
	if err != nil {
		return items, err.Error()
	}
	for rows.Next() {
		item := Inbox{}
		rows.StructScan(&item)
		text, html := ReadEmailBody(item.Body)
		item.MessageText = template.HTML(text)
		item.MessageHTML = template.HTML(html)
		items = append(items, item)
	}

	return items, ""
}
func ReadEmailBody(s string) (string, string) {
	var r io.Reader = strings.NewReader(s)

	m, err := message.Read(r)
	if message.IsUnknownCharset(err) {
		log.Println("Unknown encoding:", err)
	} else if err != nil {
		log.Fatal(err)
	}

	plain := ""
	html := ""
	if mr := m.MultipartReader(); mr != nil {
		log.Println("This is a multipart message containing:")
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err)
			}

			t, _, _ := p.Header.ContentType()
			log.Println("A part with type", t)
			all, _ := ioutil.ReadAll(p.Body)
			if t == "text/plain" {
				plain = string(all)
			} else if t == "text/html" {
				html = string(all)
			}
		}
	} else {
		t, _, _ := m.Header.ContentType()
		log.Println("This is a non-multipart message with type", t)
		all, _ := ioutil.ReadAll(m.Body)
		if t == "text/plain" {
			plain = string(all)
		} else if t == "text/html" {
			html = string(all)
		}
	}

	return plain, html
}
