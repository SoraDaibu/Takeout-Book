package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/gurunavi"
	"log"
	"net/http"
	"os"
)

// for line channel's request header, body
type Webhook struct {
	Destication string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// botにはClient型の変数が入ったメモリの場所(ポインタ)が入ってる
	// line sdkのNew関数は引数としてチャネルシークレットと
	// アクセストークンを引数として渡す必要があるから、渡してる
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)

	// 環境変数が正しく取得できず、New関数がエラーとなった際の処理
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"message":"%s"`+"\n", http.StatusText(http.StatusInternalServerError)),
		}, nil
	}

	log.Print(request.Headers)
	log.Print(request.Body)

	var webhook Webhook

	if err := json.Unmarshal([]byte(request.Body), &webhook); err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
		}, nil
	}

	// LINE channelから来たメッセージをオウム返しする
	for _, event := range webhook.Events {
		// LINE channelからEventTypeMessageのリクエストが来たときの処理
		// EventTypeMessage以外でも以下に同じ様にcase文書くと良い
		switch event.Type {
		case linebot.EventTypeMessage:
			// m.TextはLINE channelから送られてきたメッセージが入ってる
			switch m := event.Message.(type) {
			case *linebot.TextMessage:
				switch request.Path {
				case "/parrot":
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(m.Text)).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
						}, nil
					}
				case "/takeout":
					g, err := gurunavi.SearchTakeoutRestaurants(m.Text)
					if err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body:       fmt.Sprintf(`{"message":"%s"`+"\n", http.StatusText(http.StatusInternalServerError)),
						}, nil
					}

					var sm linebot.SendingMessage

					switch {
					case g.Error != nil:
						t := g.Error[0].Message
						sm = linebot.NewTextMessage(t)
					default:
						f := FlexTakeout(g)
						sm = linebot.NewFlexMessage("テイクアウト可能なお店の検索結果", f)
					}

					if _, err = bot.ReplyMessage(event.ReplyToken, sm).Do(); err != nil {
						log.Print(err)
						return events.APIGatewayProxyResponse{
							StatusCode: http.StatusInternalServerError,
							Body: fmt.Sprint(`{"message":"%s"}`+"\n", http.StatusText(http.StatusText(http.StatusInternalServerError)),
						}, err
					}
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(handler)
}
