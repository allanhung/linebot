package main

import (
    "log"
    "net/http"
    "encoding/json"
    "os"
    "fmt"

    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/linebot"
)

func parsejson(raw []byte) map[string]interface{} {
    result := make(map[string]interface{})
    err := json.Unmarshal(raw, &result)
    if err != nil {
        log.Fatal("json decode failed")
    }
    return result
}

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
                    reply_msg := fmt.Sprintf("Your Token: %s\nyou type: %s\nmessage id: %s\nuser_id:%s\ngroup_id:%s\nroom_id:%s", event.ReplyToken, message.Text, message.ID, event.Source.UserID, event.Source.GroupID, event.Source.RoomID)
                    if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply_msg)).Do(); err != nil {
                        log.Print(err)
                    }
                }
            }
        }
    })

    router.POST("/tt", func(c *gin.Context) {
        buf  := make([]byte, 1024)  
        n, _ := c.Request.Body.Read(buf) 
        jsonbody := parsejson(buf[0:n])
        var return_body string
        for k, v := range jsonbody {
            return_body = fmt.Sprintf("%s_%s_%s",return_body,k,v)
        }
        c.String(http.StatusOK, "Hello%s", return_body)
    })

    router.POST("/json", func(c *gin.Context) {
        buf  := make([]byte, 1024)  
        n, _ := c.Request.Body.Read(buf) 
        jsonbody := parsejson(buf[0:n])

        if _, err := bot.PushMessage(jsonbody["roomid"].(string), linebot.NewTextMessage(jsonbody["msg"].(string))).Do(); err != nil {
            log.Print(err)
        }
    })

    router.POST("/raw", func(c *gin.Context) {
        buf  := make([]byte, 1024)  
        n, _ := c.Request.Body.Read(buf) 
        body := string(buf[0:n])

        if _, err := bot.PushMessage("R6dcb63709978fed802b24764686c3ea8", linebot.NewTextMessage(body)).Do(); err != nil {
            log.Print(err)
        }
    })

    router.Run(":" + port)
}
