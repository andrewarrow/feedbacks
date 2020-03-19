package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"time"

	"github.com/andrewarrow/feedbacks/controllers"
	"github.com/andrewarrow/feedbacks/email"
	"github.com/andrewarrow/feedbacks/models"
	"github.com/andrewarrow/feedbacks/persist"
	"github.com/andrewarrow/feedbacks/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"

	u "net/url"
)

var local = ""
var runners = map[string]*httputil.ReverseProxy{}

func Serve() {
	port := 3001
	controllers.Db = persist.Connection()
	domains, err := models.SelectDomains(controllers.Db, 0)
	hosts := []string{}
	if err != "" {
		fmt.Println(err)
		return
	}
	for i, host := range domains {
		hosts = append(hosts, host.Domain)
		url, _ := u.Parse(fmt.Sprintf("http://localhost:%d", (port + i)))
		runners[host.Domain] = httputil.NewSingleHostReverseProxy(url)
		path := fmt.Sprintf("%s%s", util.AllConfig.Path.Sites, host.Domain)
		go exec.Command("run_feedback", path, fmt.Sprintf("%d", (port+i)), host.Domain).Output()
	}
	local = os.Getenv("LOCAL")
	router := gin.Default()
	prefix := util.AllConfig.Path.Prefix
	router.Static("/feedbacks/assets", prefix+"assets")
	router.GET("/feedbacks", controllers.WelcomeIndex)
	sessions := router.Group("/feedbacks/sessions")
	sessions.GET("/new", controllers.SessionsNew)
	sessions.POST("/", controllers.SessionsCreate)
	sessions.POST("/destroy", controllers.SessionsDestroy)
	fbDomains := router.Group("/feedbacks/domains")
	fbDomains.GET("/", controllers.AdminDomainsIndex)
	fbDomains.GET("/:domain", controllers.AdminDomainsShow)
	fbDomains.POST("/", controllers.AdminDomainsCreate)
	AddTemplates(router, prefix)
	router.NoRoute(handleReq)

	if local == "" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hosts...),
			Cache:      autocert.DirCache("/certs"),
		}

		server := &http.Server{
			Addr:    ":https",
			Handler: router,
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
		}

		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		go email.Run(":25")
		server.ListenAndServeTLS("", "")
	} else {
		go router.Run(":8080")
		go email.Run(":25")
	}

	for {
		time.Sleep(time.Second)
	}

}

func getHost(c *gin.Context) string {
	host := c.Request.Host
	if local != "" {
		host = local
	}
	return host
}
func handleReq(c *gin.Context) {
	defer c.Request.Body.Close()
	c.Writer.Header().Add("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Add("Access-Control-Allow-Methods", "GET,POST")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Filename")
	host := getHost(c)
	controllers.Stats[host]++
	runners[host].ServeHTTP(c.Writer, c.Request)
}
