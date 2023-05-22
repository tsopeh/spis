package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/gocolly/colly/v2"
)

type MenuItem struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	OrderBy  int        `json:"orderBy"`
	Level    int        `json:"level"`
	Count    int        `json:"count"`
	Children []MenuItem `json:"children"`
}

func menuResponseToUrls(items *[]MenuItem) []string {
	acc := []string{}
	for _, item := range *items {
		if item.Children == nil || len(item.Children) == 0 {
			acc = append(acc, "https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/RegistarServlet?subareaid="+strconv.Itoa(item.ID))
		} else {
			acc = append(acc, menuResponseToUrls(&item.Children)...)
		}
	}
	return acc
}

func createMenuCollector(documentCollector *colly.Collector) *colly.Collector {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		menuItems := []MenuItem{}
		err := json.Unmarshal([]byte(r.Body), &menuItems)
		if err != nil {
			log.Fatalln(err)
		}
		urls := menuResponseToUrls(&menuItems)
		log.Println("Successfully retrieved the list of " + strconv.Itoa(len(urls)) + " menu items.")
		for _, url := range urls {
			documentCollector.Visit(url)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `createMenuCollector` for URL", r.Request.URL.String(), err)
	})

	return c
}
