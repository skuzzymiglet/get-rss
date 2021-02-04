package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// GetYoutubeKey retrieves the YouTube API key from the YOUTUBE_API_KEY env variable
func GetYoutubeKey() string {
	k, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		log.Warn("YOUTUBE_API_KEY not set - YouTube API requests will fail")
	}
	return k
}

// YoutubeChannelFinder finds feeds for YouTube channels
type YoutubeChannelFinder struct {
	// TODO: cache in client (so cache-control works)
	Service *youtube.Service
}

func (y *YoutubeChannelFinder) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var err error
	y.Service, err = youtube.NewService(ctx, option.WithAPIKey(GetYoutubeKey()))
	return err
}

func (y *YoutubeChannelFinder) Find(query string) (string, error) {
	// BUG: channel must appear within first 5 results
	// TODO: limit to channels
	call := y.Service.Search.List([]string{"snippet"}).Type("channel").Q(query).MaxResults(int64(1))
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#channel":
			return "https://www.youtube.com/feeds/videos.xml?channel_id=" + item.Id.ChannelId, nil
		}
	}
	return "", nil
}
