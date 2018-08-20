package main

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gocolly/colly"
	"google.golang.org/api/option"
)

type Ticket struct {
	Date  string `firebase:"string"`
	Main  []int  `firebase:"main"`
	Extra []int  `firebase:"extra"`
}

const historyURL = "https://www.national-lottery.co.uk/results/euromillions/draw-history"

var (
	ctx    context.Context
	app    *firebase.App
	client *firestore.Client
)

func init() {
	var err error
	ctx = context.Background()
	sa := option.WithCredentialsFile("./database.json")

	app, err = firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

func ticketToNumbers(t string) []int {
	var result []int

	for _, value := range strings.Split(t, " - ") {
		n, _ := strconv.Atoi(value)
		result = append(result, n)
	}

	return result
}

func queryElementText(e *colly.HTMLElement, query string) string {
	return strings.TrimSpace(e.ChildText(query))
}

func addTicket(ticket Ticket) {
	_, _, err := client.Collection("tickets").Add(ctx, ticket)
	if err != nil {
		log.Fatalln(err)
	}
}

func crawl(draws map[string]bool, wg *sync.WaitGroup) func(*colly.HTMLElement) {
	return func(e *colly.HTMLElement) {
		date := queryElementText(e, ".table_cell.table_cell_1 .table_cell_block")
		main := queryElementText(e, ".table_cell.table_cell_3 .table_cell_block")
		extra := queryElementText(e, ".table_cell.table_cell_4 .table_cell_block")

		if !draws[date] {
			log.Printf("Adding %s draw to the list", date)

			t := Ticket{
				Date:  date,
				Main:  ticketToNumbers(main),
				Extra: ticketToNumbers(extra),
			}

			wg.Add(1)

			go func() {
				defer wg.Done()
				addTicket(t)
			}()
		}
	}
}

func main() {
	var wg sync.WaitGroup

	entries, err := client.Collection("tickets").Documents(ctx).GetAll()
	if err != nil {
		log.Fatalln(err)
	}

	var draws = make(map[string]bool)

	for _, entry := range entries {
		var t Ticket
		entry.DataTo(&t)
		draws[t.Date] = true
	}

	c := colly.NewCollector()
	c.OnHTML(".list_table.list_table_presentation.table_row_odd", crawl(draws, &wg))
	c.OnHTML(".list_table.list_table_presentation.table_row_even", crawl(draws, &wg))
	c.Visit(historyURL)

	wg.Wait()

	log.Println("All done.")
}
