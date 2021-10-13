package printer

import (
	"context"
	"fmt"

	"github.com/serhiihuberniuk/blog-api/tasks/task4/models"
)

type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) PrintNews(_ context.Context, news models.FootballNews) {
	fmt.Printf("Title: %s\nLink: %s\n\n", news.Title, news.Link)
}
