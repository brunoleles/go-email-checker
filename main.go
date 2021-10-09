package main

import (
	"flag"
	"fmt"
	"syscall"
)

func init() {

	flag.BoolVar(&runOptions.runAsServer, "serve", runOptions.runAsServer, "Run as server API")
	flag.StringVar(&runOptions.serverPort, "serve-port", runOptions.serverPort, "API server listening port")

	flag.StringVar(&runOptions.email, "email", runOptions.email, "Check one e-mail")
	flag.StringVar(&runOptions.testEmail, "email-from-test", runOptions.testEmail, "The e-mail that will be used as sender for the test (use some valid email under your controll)")

	//TODO: add verbosity flag
}

func main() {
	flag.Parse()

	if runOptions.runAsServer {
		//TODO: validate server related flags

		runRestServer()
		return
	}

	if runOptions.email != "" {
		//TODO: validate email related flags

		runCheckEmail()
		return
	}

	message(VERBOSITY_ERROR, "-serve or -email options must be defined")
	flag.Usage()
	syscall.Exit(1)
}

func runCheckEmail() {
	valid, err := CheckEmailAddress(runOptions.email)

	if err != nil {
		message(VERBOSITY_INFO, fmt.Sprintf("%s nok, err: \"%s\"", runOptions.email, err.Error()))
		syscall.Exit(1)
	}

	if valid {
		message(VERBOSITY_INFO, fmt.Sprintf("%s ok", runOptions.email))
		syscall.Exit(0)
	}

	message(VERBOSITY_INFO, fmt.Sprintf("%s nok", runOptions.email))
	syscall.Exit(1)
}

func runRestServer() {
	gin := setupGin()
	gin.Run(runOptions.serverPort)
}
