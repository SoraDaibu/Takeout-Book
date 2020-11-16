package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"linebot/gurunavi"
)

func TextTakeout(g *gurunavi.GurunaviResponseBody) string {
	var t string
	for _, r := range g.Rest {
		t += r.Name + "\n" + r.URL + "\n"
	}
	return t
}

func FlexTakeout(g *gurunavi.GurunaviResponseBody) *linebot.CarouselContainer {

	var bcs []*linebot.BubbleContainer

	for _, r := range g.Rest {
		b := linebot.BubbleContainer{
			// SDK(linebot)の定数FlexContainerTypeBubbleの値は"bubble"
			Type:   linebot.FlexContainerTypeBubble,
			Hero:   setHero(r),
			Body:   setBody(r),
			Footer: setFooter(r),
		}
		bcs = append(bcs, &b)
	}
}

func setHero(r *gurunavi.Rest) *linebot.ImageComponent {

	if r.ImageURL.ShopImage1 == "" {
		return nil
	}

	return &linebot.ImageComponent{
		Type:        linebot.FlexComponentTypeImage,
		URL:         r.ImageURL.ShopImage1,
		Size:        linebot.FlexImageSizeTypeFull,
		AspectRatio: linebot.FlexImageAspectRatioType20to13,
		AspectMode:  linebot.FlexImageAspectModeTypeCover,
	}
}
