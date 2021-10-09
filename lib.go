package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

// App verbosity type
type Verbosity int

// App verbosity levels available
const (
	VERBOSITY_SILENT Verbosity = iota
	VERBOSITY_ERROR
	VERBOSITY_INFO
	VERBOSITY_DEBUG
)

// Application defined errors
var (
	ErrInvalidEmailAddress = errors.New("invalid email address")
	ErrNoDomainFound       = errors.New("unable to string email from email address")
	ErrMXLookupError       = errors.New("mx lookup error")
)

// Application run/init options struct
type RunOptions struct {
	runAsServer bool
	serverPort  string
	email       string
	testEmail   string
	verbosity   Verbosity
}

// Default application options, will be updateded on func main
var runOptions = RunOptions{
	runAsServer: false,
	serverPort:  "8000",
	email:       "",
	testEmail:   "",
	verbosity:   VERBOSITY_INFO,
}

// Send a debug message, respects app verbosity
func message(level Verbosity, mesage string) {
	if level > runOptions.verbosity {
		return
	}
	println(mesage)
}

// Creates the web server instance
//
// TODO: create a response struct, and return more information about the test to check e-mail request via API
func setupGin() *gin.Engine {
	// gin.DisableConsoleColor()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/", func(c *gin.Context) {
		email := c.Query("email")

		valid, err := CheckEmailAddress(email)
		if err != nil {
			c.String(http.StatusOK, "ERROR 2")
			return
		}

		if valid {
			c.String(http.StatusOK, "ok")
			return
		}

		c.String(http.StatusOK, "nok")
	})

	return r
}

// Execs the mx lookup
//
// TOOD: cache for domain lookup, to avoid multiple lookups to same domain
func WrappedMxLookup(domain string) ([]*net.MX, error) {
	mxs, err := net.LookupMX(domain)
	if err != nil {
		return nil, err
	}
	return mxs, nil
}

// Checks email address
func CheckEmailAddress(email string) (bool, error) {
	message(VERBOSITY_DEBUG, fmt.Sprintf("Checking email address:\"%s\"...", email))

	// Validating email address
	address, err := mail.ParseAddress(email)

	if err != nil {
		message(VERBOSITY_DEBUG, fmt.Sprintf("err, invalid email address, err: \"%s\"", err.Error()))
		return false, ErrInvalidEmailAddress
	}

	// getting the email domain, for mx lookup
	domain := regexp.MustCompile(`.*@(.*)`).ReplaceAllString(address.Address, "$1")
	if err != nil {
		message(VERBOSITY_DEBUG, fmt.Sprintf("err, unable to extract domain from, email: \"%s\", err: \"%s\"", email, err.Error()))
		return false, ErrNoDomainFound
	}

	// MX lookup
	mxs, err := WrappedMxLookup(domain)
	if err != nil {
		message(VERBOSITY_DEBUG, fmt.Sprintf("err, mxlookup error, domain: \"%s\", err: \"%s\"", domain, err.Error()))
		return false, ErrMXLookupError
	}

	sucess := false
	for _, mx := range mxs {
		message(VERBOSITY_DEBUG, fmt.Sprintf("Testing email on MX host: %s:%d", mx.Host, 25))

		tcp_conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", mx.Host, 25), time.Second*5)
		if err != nil {
			message(VERBOSITY_DEBUG, fmt.Sprintf("err, unabe to dial host, err: \"%s\" ", err.Error()))
			continue
		}

		tcp_conn_timeout := time.AfterFunc(time.Second*5, func() { tcp_conn.Close() })
		defer tcp_conn_timeout.Stop()

		client, _ := smtp.NewClient(tcp_conn, mx.Host)
		defer client.Close()

		err = client.Hello(mx.Host)
		if err != nil {
			message(VERBOSITY_DEBUG, fmt.Sprintf("err, unabe greet host, err: \"%s\" ", err.Error()))
			continue
		}

		err = client.Mail(runOptions.testEmail)
		if err != nil {
			message(VERBOSITY_DEBUG, fmt.Sprintf("err, unable to create mail, err: \"%s\" ", err.Error()))
			continue
		}

		err = client.Rcpt(address.Address)
		if err != nil {
			message(VERBOSITY_DEBUG, fmt.Sprintf("err, unable to create rcpt, err: \"%s\" ", err.Error()))
			continue
		}

		message(VERBOSITY_DEBUG, "Email adress found on host")
		sucess = true
		break
	}

	return sucess, nil
}
