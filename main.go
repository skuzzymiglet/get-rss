package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/cheggaaa/pb"
	"github.com/peterbourgon/diskv"
)

// FeedFinder finds RSS feeds for a given search query
type FeedFinder interface {
	Find(query string) (feed string, err error)
	Init() error
}

func main() {
	handlers := map[string]FeedFinder{}
	handlers["yt_nc"] = &YoutubeChannelFinder{} // No cache
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Println("Error getting cache dir:", err)
		handlers["yt"] = handlers["yt_nc"]
		delete(handlers, "yt_nc")

	} else {
		handlers["yt"] = &DiskvCacheFinder{
			Child: &YoutubeChannelFinder{},
			Opts:  diskv.Options{BasePath: filepath.Join(cacheDir, "get-rss")},
		}
	}
	for k, v := range handlers {
		err := v.Init()
		if err != nil {
			log.Fatalf("Error initializing service %s: %s", k, err)
		}
	}
	lines := make([]string, 0)
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	bar := pb.New(len(lines))
	bar.Output = os.Stderr // Just an URL
	bar.Start()
	for _, l := range lines {
		bar.Increment()
		if strings.HasPrefix(l, "#") {
			continue
		}
		if !strings.Contains(l, " ") { // Just an URL
			fmt.Println(l)
		} else {
			q := strings.SplitN(l, " ", 2)
			if len(q) != 2 {
				log.Fatalf("Invalid query: %s", l)
			}
			handler, found := handlers[q[0]]
			if !found {
				log.Fatalf("Handler not found: %s", q[0])
			}
			s, err := handler.Find(q[1])
			if err != nil {
				log.Fatalf("Error on %s finding %s: %s", q[0], q[1], err)
			}
			if s == "" {
				log.Warnf("No result for %s on %s", q[1], q[0])
			}
			fmt.Println(s)
		}

	}
}
