package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/smelton01/strearning-bot/scrape"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type factory struct {
	bot reddit.Bot
}

func Start() {
	bot, err := reddit.NewBotFromAgentFile("bot.agent", 10*time.Second)
	if err != nil {
		log.Fatalf("failed to initialize bot: %v", err)
	}

	cfg := graw.Config{Subreddits: []string{"bottesting", "LivestreamFail"}, SubredditComments: []string{"LivestreamFail", "bottesting"}}

	handler := &factory{bot: bot}
	fmt.Println("starting run")

	_, wait, err := graw.Run(handler, bot, cfg)
	if err != nil {
		log.Fatalf("failed to run: %v", err)
	}
	fmt.Println("graw run failed", wait())
}

func (f *factory) Post(p *reddit.Post) error {
	numbers := []int{}
	if p.Author == "strugglingstrimerbot" || p.Author == "LSFmoderator" {
		return nil
	}
	if strings.Contains(p.SelfText, "#") || strings.Contains(p.Title, "#") {
		whole := strings.Split(p.Title, " ")
		whole = append(whole, strings.Split(p.SelfText, " ")...)
		for _, word := range whole {
			if word[0] == '#' {
				rank := strings.TrimSpace(word[1:])
				r, err := strconv.Atoi(rank)
				if err != nil {
					log.Printf("failed to convert %v: %v", rank, err)
					continue
				}
				if r > 5000 || r < 1 {
					log.Printf("number out of bounds: %v", r)
					continue
				}
				numbers = append(numbers, r)
			}
		}
		if len(numbers) == 0 {
			return nil
		}
		message := reply(numbers)
		err := f.bot.Reply(p.Name, message)
		if err != nil {
			log.Fatalf("failed to reply: %v", err)
		}
		log.Printf("replied to u/%v", p.Author)
	}
	return nil
}

func (f *factory) Comment(c *reddit.Comment) error {
	numbers := []int{}
	if c.Author == "strugglingstrimerbot" || c.Author == "LSFmoderator" {
		return nil
	}
	if strings.Contains(c.Body, "#") {
		whole := strings.Split(c.Body, " ")
		for _, word := range whole {
			if word[0] == '#' {
				rank := strings.TrimSpace(word[1:])
				r, err := strconv.Atoi(rank)
				if err != nil {
					log.Printf("failed to convert %v: %v", rank, err)
					continue
				}
				if r > 5000 || r < 1 {
					log.Printf("number out of bounds: %v", r)
					continue
				}
				numbers = append(numbers, r)
			}
		}
		if len(numbers) == 0 {
			return nil
		}
		message := reply(numbers)
		err := f.bot.Reply(c.Name, message)
		if err != nil {
			log.Fatalf("failed to reply: %v", err)
		}
		log.Printf("replied to u/%v", c.Author)
	}
	return nil
}

func reply(nums []int) string {
	replyMessage := ""
	base := "Streamer #%v is [%v](%v) earning *only* ^(%v) over the last 2 years.\n  "
	for _, rank := range nums {
		data := scrape.GetData()
		streamer := data[rank-1]
		replyMessage += fmt.Sprintf(base, rank, streamer.Name, streamer.Url, streamer.Money)
	}
	replyMessage += "\n  ^(please PM [u/Smelton09](https://www.reddit.com/user/Smelton09/) with issues or feedback! [Code](https://github.com/Smelton01/struggling-streamer-bot))"
	return replyMessage
}
