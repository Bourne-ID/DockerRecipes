package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"sync"

	"golang.org/x/time/rate"
)
import "time"
import "flag"
import "strconv"

var (
	difficulty    = 3
	length        = 1024
	output_method = "activemq"
	count         = 0 //unlimited
	mps           = 20
	destination   = "address.foo"
	server        = "127.0.0.1:61616"
	username      = ""
	password      = ""
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
	flag.IntVar(&difficulty, "difficulty", LookupEnvOrInt("DIFFICULTY", difficulty), "Sets the difficulty of the calculation")
	flag.IntVar(&length, "length", LookupEnvOrInt("LENGTH", length), "Sets the length of string to produce (IE message size)")
	flag.IntVar(&count, "count", LookupEnvOrInt("COUNT", count), "How many messages to produce (0=unlimited)")
	flag.IntVar(&mps, "mps", LookupEnvOrInt("MPS", mps), "How Many Messages Per Second (default 20)")
	flag.StringVar(&output_method, "output_method", LookupEnvOrString("OUTPUT", output_method), "Where to output messages")
	flag.StringVar(&destination, "destination", LookupEnvOrString("DESTINATION", destination), "Queue name (ActiveMQ)")
	flag.StringVar(&server, "server", LookupEnvOrString("SERVER", server), "server of the service to connect to")
	flag.StringVar(&username, "username", LookupEnvOrString("USERNAME", username), "Username of the service to connect to")
	flag.StringVar(&password, "password", LookupEnvOrString("PASSWORD", password), "Password of the service to connect to")

	flag.Parse()
	rlim := rate.NewLimiter(rate.Every(time.Second/time.Duration(mps)), 10)

	var output output
	if output_method == "activemq" {
		output = activeMq{destination: destination, server: server, username: username, password: password}
		output.init()
	} else {
		output = stdout{}
	}

	var wg sync.WaitGroup

	for i := 0; count == 0 || i < count; i++ {
		test := RandStringRunes(length)
		data := strconv.Itoa(difficulty) + ":" + test
		err := rlim.Wait(context.Background())
		if err != nil {
			log.Println("Error: Rate Wait:", err)
			return
		}
		wg.Add(1)

		go func(dataToSend string) {
			defer wg.Done()
			output.out(data)
		}(data)
	}
	wg.Wait()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
