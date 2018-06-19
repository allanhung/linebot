package main

import (
    "log"
    "os"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    bot, err := linebot.New(
        os.Getenv("CHANNEL_SECRET"),
        os.Getenv("CHANNEL_TOKEN"),
    )
    if err != nil {
        log.Fatal(err)
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.POST("/callback", func(c *gin.Context) {
        events, err := bot.ParseRequest(c.Request)
        if err != nil {
            if err == linebot.ErrInvalidSignature {
                log.Print(err)
            }
            return
        }
        for _, event := range events {
            if event.Type == linebot.EventTypeMessage {
                switch message := event.Message.(type) {
                case *linebot.TextMessage:
                    reply_msg := fmt.Sprintf("Your Token: %s, you type: %s", event.ReplyToken, message.Text)
                    if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply_msg)).Do(); err != nil {
                        log.Print(err)
                    }
                }
            }
        }
    })

    router.POST("/mysql", func(c *gin.Context) {
//        buf  := make([]byte, 1024)  
//        n, _ := c.Request.Body.Read(buf) 
//        body := string(buf[0:n])

        if _, err := bot.PushMessage("96911779b0fc4755aa72608b53173986", linebot.NewTextMessage("hello")).Do(); err != nil {
            log.Print(err)
        }
    })

    router.Run(":" + port)
}
