package task

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

type Task interface {
	Run(ctx context.Context) error
	Stop() error
}

type ExampleTask struct {
	Name string
}

func (t *ExampleTask) Run(ctx context.Context) error {
	log.Printf("Task %s started", t.Name)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var title string
	var err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err = chromedp.Navigate("https://example.com").Do(ctx)
		if err != nil {
			return err
		}
		err = chromedp.Title(&title).Do(ctx)
		return err
	}))
	if err != nil {
		log.Printf("Task %s error: %v", t.Name, err)
		return err
	}
	log.Printf("Page title: %s", title)
	return nil
}

func (t *ExampleTask) Stop() error {
	log.Printf("Task %s stopped", t.Name)
	return nil
}
