package controllers

import (
	"github.com/andrewarrow/feedbacks/models"
	"github.com/gin-gonic/gin"
	"sync"
	"net/http"
)
var Mutex = sync.Mutex{}
var Stats = map[string]int{}
var EmailStats = map[string]int{}

func AdminDomainsIndex(c *gin.Context) {
	if !BeforeAll("admin", c) {
		return
	}
	domains, err := models.SelectDomains(Db, user.Id)
	for _, domain := range domains {
		domain.Hits = Stats[domain.Domain]
		domain.Emails = EmailStats[domain.Domain]
	}

	c.HTML(http.StatusOK, "admin__domains__index.tmpl", gin.H{
		"flash":   err,
		"user":    user,
		"domains": domains,
	})

}
func AdminDomainsCreate(c *gin.Context) {
	if !BeforeAll("admin", c) {
		return
	}

	domain := c.PostForm("domain")
	err := models.InsertDomain(Db, domain, user.Id)
	if err != "" {
		SetFlash(err, c)
	}
	c.Redirect(http.StatusFound, "/feedbacks/domains")
	c.Abort()
}
func AdminDomainsShow(c *gin.Context) {
	if !BeforeAll("admin", c) {
		return
	}
	domain := c.Param("domain")
	items, err := models.SelectInboxByDomain(Db, domain)
	if err != "" {
		SetFlash(err, c)
		c.Redirect(http.StatusFound, "/feedbacks/domains")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "admin__domains__show.tmpl", gin.H{
		"user":   user,
		"flash":  "",
		"items":  items,
		"domain": domain,
	})
}
