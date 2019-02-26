// Poonam Phowakande
package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "Port to connect NSQD")

	AppCtx = App{}
)

func main() {
	flag.Parse()
	Welcome()

	AppCtx.Initialize()

	log.WithFields(log.Fields{
		"Port":   *port,
		"Host":   AppCtx.Hostname,
		"Module": AppCtx.Module,
	}).Info("Execution environment")

	addr := net.JoinHostPort("0.0.0.0", strconv.Itoa(*port))

	AppCtx.Run(addr)
}

// Welcome banner http://patorjk.com/software/taag - font Stop
func Welcome() {
	var str string = `

  ______         _                                ______   _______     _     _______            ______   _____
 / _____)       | |                              (_____ \ (_______)   | |   (_______)     /\   (_____ \ (_____)
| /  ___   ___  | |  ____  ____    ____    ___    _____) ) _____       \ \   _           /  \   _____) )   _
| | (___) / _ \ | | / _  ||  _ \  / _  |  (___)  (_____ ( |  ___)       \ \ | |         / /\ \ |  ____/   | |
| \____/|| |_| || |( ( | || | | |( ( | |               | || |_____  _____) )| |_____   | |__| || |       _| |_
 \_____/  \___/ |_| \_||_||_| |_| \_|| |               |_||_______)(______/  \______)  |______||_|      (_____)
                                 (_____|

	`
	fmt.Println(str)
}
