package server

import "github.com/gin-gonic/gin"
import "github.com/andrewarrow/feedbacks/util"
import "time"
import "net/http/httputil"
import "os"
import "net/http"
import "crypto/tls"
import "golang.org/x/crypto/acme/autocert"
import u "net/url"

var local = ""
var runners = map[string]*httputil.ReverseProxy{}

func Serve() {
	for host, port := range util.HostToPort {
		url, _ := u.Parse("http://localhost:" + port)
		runners[host] = httputil.NewSingleHostReverseProxy(url)
	}
	local = os.Getenv("MANY_LOCAL")
	router := gin.Default()
	router.GET("/*name", handleReq)
	router.POST("/*name", handleReq)

	if local == "" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(util.Hosts...),
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
		//go email.Run(":2525")
		server.ListenAndServeTLS("", "")
	} else {
		go router.Run(":8080")
		//go email.Run(":2525")
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
	c.Writer.Header().Add("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Filename")
	host := getHost(c)
	runners[host].ServeHTTP(c.Writer, c.Request)
}
