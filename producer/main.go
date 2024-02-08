package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
	// "encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/kafka-go"
)

const (
	testTopic     = "test"
	userTopic     = "user"
	productTopic  = "product"
	broker0Address = "kafka-0:9092"
	broker1Address = "kafka-1:9092"
	broker2Address = "kafka-2:9092"
)

type myJSON struct {
    Array []string
}

func produce(ctx context.Context, timeoutMs int, limit int)  {
	i := 0

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker0Address, broker1Address, broker2Address},
		Topic:   testTopic,
	})

	for {
		if (i > limit) {
			return
		}
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			Value: []byte("this is message" + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		fmt.Println("writes:", i)
		i++
		time.Sleep(time.Duration(timeoutMs))
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
		return  c.HTML(http.StatusOK, "Topics initialized!")
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

		return c.HTML(http.StatusOK, "Test started!")
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
