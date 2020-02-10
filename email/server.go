package email

import (
	"fmt"
	"github.com/emersion/go-smtp"
	"github.com/saintienn/go-spamc"
	"io"
	"io/ioutil"
	"strings"
	"time"
)
import "github.com/jmoiron/sqlx"
import "github.com/andrewarrow/feedbacks/persist"

var db *sqlx.DB
var spam *spamc.Client

var insertChannel chan map[string]interface{} = make(chan map[string]interface{}, 1024)

type Backend struct{}
type Session struct {
	SentFrom string
	SentTo   string
	Subject  string
	Body     string
	Host     string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	fmt.Println("Mail from:", from)
	s.SentFrom = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	s.SentTo = to
	tokens := strings.Split(to, "@")
	if len(tokens) > 0 {
		s.Host = tokens[1]
	}
	//fmt.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	tokens := strings.Split(string(b), "\n")
	fmt.Println("MMMMM", tokens)
	subjectDone := false
	for _, token := range tokens {
		if subjectDone {
			s.Body += token + "\n"
		}
		if strings.HasPrefix(token, "Subject:") {
			s.Subject = token[9:]
			subjectDone = true
		}

	}

	isSpam := 0
	spamScore := 0.0
	reply, err := spam.Check(string(b))

	if err == nil {
		i1, ok := reply.Vars["isSpam"]
		if ok {

			if i1.(bool) {
				isSpam = 1
			}
			i2, ok := reply.Vars["spamScore"]
			if ok {
				spamScore = i2.(float64)
			}
		}
	}
	m := map[string]interface{}{"host": s.Host, "body": s.Body,
		"is_spam":    isSpam,
		"spam_score": spamScore,
		"sent_from":  s.SentFrom, "sent_to": s.SentTo,
		"subject": s.Subject}
	fmt.Println("33333", m)
	insertChannel <- m
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{}, nil
}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{}, nil
}

func Run(ss string) {
	db = persist.Connection()
	spam = spamc.New("127.0.0.1:783", 10)

	go func() {
		for m := range insertChannel {
			db.NamedExec("INSERT INTO inbox (host, is_spam, spam_score, body, sent_from, sent_to, subject) values (:host, :is_spam, :spam_score, :body, :sent_from, :sent_to, :subject)", m)
		}
	}()

	be := &Backend{}
	s := smtp.NewServer(be)
	s.Addr = ":25"
	s.Domain = "many.pw"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = false
	s.AuthDisabled = true
	s.ListenAndServe()
}
