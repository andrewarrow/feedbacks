package server

import "github.com/gin-gonic/gin"
import "github.com/andrewarrow/feedbacks/util"
import "github.com/andrewarrow/feedbacks/email"
import "time"
import "fmt"
import "net/http/httputil"
import "os"
import "os/exec"
import "net/http"
import "crypto/tls"
import "golang.org/x/crypto/acme/autocert"
import u "net/url"

var local = ""
var runners = map[string]*httputil.ReverseProxy{}

func Serve() {
	port := 3001
	for i, host := range util.AllConfig.Http.Hosts {
		url, _ := u.Parse(fmt.Sprintf("http://localhost:%d", (port + i)))
		runners[host] = httputil.NewSingleHostReverseProxy(url)
		path := fmt.Sprintf("%s%s", util.AllConfig.Path.Sites, host)
		go exec.Command("run_feedback", path, fmt.Sprintf("%d", (port+i)), host).Output()
	}
	local = os.Getenv("LOCAL")
	router := gin.Default()
	router.GET("/*name", handleReq)
	router.POST("/*name", handleReq)

	if local == "" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(util.AllConfig.Http.Hosts...),
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
		server.ListenAndServeTLS("", "")
		go email.Run(":25")
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
	runners[host].ServeHTTP(c.Writer, c.Request)
}
