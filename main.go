package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gocolly/colly"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
)

var rowz = map[string]string{
	"Fri 17 Aug 2018": "02 - 15 - 18 - 24 - 43",
	"Fri 10 Aug 2018": "18 - 20 - 36 - 43 - 44",
	"Fri 03 Aug 2018": "07 - 26 - 36 - 38 - 43",
	"Tue 31 Jul 2018": "20 - 25 - 34 - 42 - 45",
	"Tue 24 Jul 2018": "02 - 04 - 23 - 39 - 40",
	"Tue 17 Jul 2018": "04 - 06 - 27 - 48 - 50",
	"Tue 10 Jul 2018": "03 - 08 - 26 - 33 - 45",
	"Tue 03 Jul 2018": "01 - 12 - 15 - 29 - 48",
	"Fri 29 Jun 2018": "15 - 21 - 23 - 40 - 48",
	"Fri 22 Jun 2018": "14 - 25 - 39 - 41 - 44",
	"Fri 15 Jun 2018": "23 - 26 - 33 - 38 - 49",
	"Fri 08 Jun 2018": "08 - 19 - 32 - 43 - 46",
	"Fri 01 Jun 2018": "17 - 18 - 24 - 29 - 40",
	"Tue 29 May 2018": "06 - 11 - 20 - 38 - 43",
	"Tue 22 May 2018": "01 - 11 - 37 - 41 - 48",
	"Tue 15 May 2018": "04 - 16 - 20 - 31 - 39",
	"Tue 08 May 2018": "17 - 25 - 35 - 39 - 44",
	"Tue 01 May 2018": "06 - 15 - 17 - 42 - 48",
	"Fri 27 Apr 2018": "12 - 24 - 40 - 41 - 46",
	"Fri 20 Apr 2018": "03 - 16 - 25 - 39 - 44",
	"Fri 13 Apr 2018": "05 - 25 - 34 - 48 - 50",
	"Fri 06 Apr 2018": "01 - 29 - 33 - 45 - 47",
	"Fri 30 Mar 2018": "12 - 17 - 28 - 35 - 47",
	"Fri 23 Mar 2018": "05 - 07 - 11 - 46 - 50",
	"Fri 16 Mar 2018": "04 - 17 - 24 - 27 - 31",
	"Fri 09 Mar 2018": "09 - 14 - 21 - 32 - 44",
	"Fri 02 Mar 2018": "02 - 07 - 34 - 45 - 48",
	"Tue 27 Feb 2018": "03 - 31 - 41 - 48 - 50",
	"Tue 20 Feb 2018": "06 - 14 - 19 - 25 - 29",
}

type Ticket struct {
	Main  []int `json:"main"`
	Extra []int `json:"extra"`
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

func saveTicket(date string, ticket []int, extra []int) {
	result, err := client.Collection("tickets").Doc("inspiration").Get(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	var t Ticket

	mapstructure.Decode(result.Data(), &t)

	log.Println(t.Main)

	defer client.Close()

}

func main() {
	// get all dates first

	// crawl := func(e *colly.HTMLElement) (string, Ticket) {
	// 	date := queryElementText(e, ".table_cell.table_cell_1 .table_cell_block")
	// 	main := queryElementText(e, ".table_cell.table_cell_3 .table_cell_block")
	// 	extra := queryElementText(e, ".table_cell.table_cell_4 .table_cell_block")

	// 	t := Ticket{
	// 		Main:  ticketToNumbers(main),
	// 		Extra: ticketToNumbers(extra),
	// 	}

	// 	return date, t
	// }

	// c := colly.NewCollector()
	// c.OnHTML(".list_table.list_table_presentation.table_row_odd", crawl)
	// c.OnHTML(".list_table.list_table_presentation.table_row_even", crawl)
	// c.Visit(historyURL)

	// result, err := client.Collection("tickets").Doc("inspiration").Set(ctx, ticket)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	result, err := client.Collection("tickets").Doc("inspiration").Get(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println(result)

}
