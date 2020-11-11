package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// botにはClient型の変数が入ったメモリの場所(ポインタ)が入ってる
	// line sdkのNew関数は引数としてチャネルシークレットと
	// アクセストークンを引数として渡す必要があるから、渡してる
	_, err := linebot.New(
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

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
