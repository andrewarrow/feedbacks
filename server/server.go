package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/foomo/simplecert"
	"github.com/foomo/tlsconfig"

	"github.com/andrewarrow/feedbacks/controllers"
	"github.com/andrewarrow/feedbacks/email"
	"github.com/andrewarrow/feedbacks/models"
	"github.com/andrewarrow/feedbacks/persist"
	"github.com/andrewarrow/feedbacks/util"
	"github.com/gin-gonic/gin"

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
		fmt.Println("2", host, path)
		go exec.Command("/root/ice/feedbacks/run_feedback", path, fmt.Sprintf("%d", (port+i)), host.Domain).Output()
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
		cfg := simplecert.Default
		cfg.Domains = hosts
		cfg.CacheDir = "/certs"
		cfg.SSLEmail = "oneone@gmail.com"
		certReloader, err := simplecert.Init(cfg, nil)
		fmt.Println("err", err)

		go http.ListenAndServe(":80", http.HandlerFunc(simplecert.Redirect))
		tlsconf := tlsconfig.NewServerTLSConfig(tlsconfig.TLSModeServerStrict)
		tlsconf.GetCertificate = certReloader.GetCertificateFunc()

		s := &http.Server{
			Addr:      ":443",
			Handler:   router,
			TLSConfig: tlsconf,
		}

		go email.Run(":25")
		s.ListenAndServeTLS("", "")

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
	controllers.Mutex.Lock()
	fmt.Println(c.ClientIP())
	if controllers.Stats[host] == nil {
		controllers.Stats[host] = map[string]int{}
	}
	controllers.Stats[host][c.ClientIP()]++
	if controllers.RefererStats[host] == nil {
		controllers.RefererStats[host] = map[string]int{}
	}
	controllers.RefererStats[host][strings.Join(c.Request.Header["Referer"], ",")]++
	controllers.Mutex.Unlock()
	if runners[host] != nil {
		runners[host].ServeHTTP(c.Writer, c.Request)
	}
}
