package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	// "encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/kafka-go"
)

const (
	testTopic      = "test"
	userTopic      = "user"
	productTopic   = "product"
	broker0Address = "kafka-0:9092"
	broker1Address = "kafka-1:9092"
	broker2Address = "kafka-2:9092"
)

func produce(ctx context.Context, timeoutMs int, limit int) {
	var wg sync.WaitGroup

	i := 0

	w := &kafka.Writer{
		Addr:  kafka.TCP(broker0Address, broker1Address, broker2Address),
		Topic: testTopic,
	}

	for {
		if i >= limit {
			wg.Wait()
			fmt.Println(fmt.Sprintf("%d items written.", limit))
			return
		}

		wg.Add(1)
		go func(key int) {
			defer wg.Done()
			err := w.WriteMessages(ctx, kafka.Message{
				Key:   []byte(strconv.Itoa(key)),
				Value: []byte(fmt.Sprintf("this is message %d. Time: %d", key, time.Now().Unix())),
			})
			if err != nil {
				panic("could not write message " + err.Error())
			}
			fmt.Println("writes:", key)
		}(i)

		i++
	}
}

func initTopics() {
	conn, err := kafka.Dial("tcp", broker0Address)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: testTopic, NumPartitions: 3, ReplicationFactor: 2}, {Topic: userTopic, NumPartitions: 2, ReplicationFactor: 2}, {Topic: productTopic, NumPartitions: 2, ReplicationFactor: 2}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello, Docker! <3")
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.GET("/initTopics", func(c echo.Context) error {
		initTopics()
		return c.HTML(http.StatusOK, "Topics initialized!")
	})

	e.GET("/test", func(c echo.Context) error {
		var realTimeout int
		var realLimit int

		if t, err := strconv.Atoi(c.QueryParam("timeout")); err == nil {
			realTimeout = t
		} else {
			realTimeout = 1000
		}

		if l, err := strconv.Atoi(c.QueryParam("limit")); err == nil {
			realLimit = l
		} else {
			realLimit = 1000
		}

		ctx := context.Background()
		go produce(ctx, realTimeout, realLimit)

		res := fmt.Sprintf("Test started at %s!\nNumber of items: %d", time.Now().Format("02/01/2006 15:04:05"), realLimit)
		return c.HTML(http.StatusOK, res)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
