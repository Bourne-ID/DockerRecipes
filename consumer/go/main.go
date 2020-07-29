package main

import (
	"flag"
	"github.com/go-stomp/stomp"
	"log"
	"os"
)
import "strings"
import "crypto/sha256"
import "encoding/hex"
import "strconv"

var (
	count    = 0 //unlimited
	queue    = "address.foo"
	server   = "127.0.0.1:61616"
	username = ""
	password = ""
)

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("LookupEnvOrInt[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func main() {
	flag.StringVar(&server, "server", LookupEnvOrString("SERVER", server), "server of the service to connect to")
	flag.StringVar(&username, "username", LookupEnvOrString("USERNAME", username), "Username of the service to connect to")
	flag.StringVar(&password, "password", LookupEnvOrString("PASSWORD", password), "Password of the service to connect to")
	flag.StringVar(&queue, "queue", LookupEnvOrString("QUEUE", queue), "Queue Name")
	flag.IntVar(&count, "count", LookupEnvOrInt("COUNT", count), "How many messages to produce (0=unlimited)")
	flag.Parse()

	var options = []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.Host("/"),
	}

	conn, err := stomp.Dial("tcp", server, options...)

	if err != nil {
		println("cannot connect to server", err.Error())
		return
	}

	sub, err := conn.Subscribe(queue, stomp.AckAuto,
		stomp.SubscribeOpt.Header("subscription-type", "ANYCAST"),
	)

	if err != nil {
		println("cannot subscribe to", queue, err.Error())
		panic(err)
	}

	for i := 1; count == 0 || i <= count; i++ {
		msg := <-sub.C
		input := string(msg.Body)
		splitLoc := strings.Index(input, ":")
		if splitLoc < 0 {
			panic("Input " + input + " does not contain a colon")
		}
		difficulty := input[0:splitLoc]
		content := input[splitLoc+1:]

		diffInt, err := strconv.Atoi(difficulty)
		if err != nil {
			panic(err)
		}

		target := strings.Repeat("0", diffInt)

		match := false
		nounce := 0
		var sha string
		for !match {

			hasher := sha256.New()
			hasher.Write([]byte(content + strconv.Itoa(nounce)))
			sha = hex.EncodeToString(hasher.Sum(nil))

			if sha[0:diffInt] == target {
				match = true
			} else {
				nounce++
			}
		}

		log.Println(input + "|" + strconv.Itoa(nounce))
		log.Println(sha)
	}
	println("receiver finished")
}
