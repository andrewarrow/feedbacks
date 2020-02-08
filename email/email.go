package email

import "fmt"
import "strings"
import "net"
import "net/smtp"

func Send(to, from, subj, content string) {
	tokens := strings.Split(to, "@")
	domain := tokens[1]
	fmt.Println("|", domain, "|")
	mxrecords, err := net.LookupMX(domain)
	fmt.Println(err)
	for _, mx := range mxrecords {
		fmt.Println(mx.Host, mx.Pref)
	}
	server := mxrecords[0].Host
	recipients := []string{to}
	headers := []string{"From: " + from,
		"To: " + to,
		"Subject: " + subj}
	body := []string{content}

	ports := []int{25, 465, 587, 2525}
	for _, port := range ports {
		err := smtp.SendMail(fmt.Sprintf("%s:%d", server, port), nil, from, recipients,
			[]byte(strings.Join(append(headers, body...), "\r\n")))
		if err == nil {
			break
		} else {
			fmt.Println(err)
		}
	}
}
