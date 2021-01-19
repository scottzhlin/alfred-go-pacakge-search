package main

import (
	"fmt"

	aw "github.com/deanishe/awgo"
	alfred_go_pacakge_search "github.com/scottzhlin/alfred-go-pacakge-search"
)

var alfredWorkFlow *aw.Workflow

func init() {
	alfredWorkFlow = aw.New()
}

func run() {
	if len(alfredWorkFlow.Args()) == 0 {
		return
	}

	// Disable UIDs so Alfred respects our sort order. Without this,
	// it may bump read/unpublished books to the top of results, but
	// we want to force them to always be below unread books.
	alfredWorkFlow.Configure(aw.SuppressUIDs(true))

	query := alfredWorkFlow.Args()[0]
	results, err := alfred_go_pacakge_search.SearchGoPackages(query, 10)
	if err != nil {
		alfredWorkFlow.FatalError(err)
		alfredWorkFlow.SendFeedback()
		return
	}

	for _, pkg := range results {
		alfredWorkFlow.NewItem(pkg.Name).
			Subtitle(pkg.Path).
			Arg(pkg.Path).
			UID(fmt.Sprintf("%s:%s", pkg.Name, pkg.Path)).
			Valid(true)
	}

	alfredWorkFlow.WarnEmpty("No matching items", "Try a different query ?")
	alfredWorkFlow.SendFeedback()
}
