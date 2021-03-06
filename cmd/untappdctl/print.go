package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mdlayher/untappd"
)

// printBadges turns a slice of *untappd.Badge structs into a human-friendly
// output format, and prints it to stdout.
func printBadges(badges []*untappd.Badge) {
	tw := tabWriter()

	// Print field header
	fmt.Fprintln(tw, "ID\tName\tEarned\tCheckinID")

	// Function to be invoked for each badge and badge level
	printFn := func(b *untappd.Badge) {
		y, m, d := b.Earned.Date()

		fmt.Fprintf(tw, "%d\t%s\t%s\t%d\n",
			b.ID,
			b.Name,
			fmt.Sprintf("%04d-%02d-%02d", y, m, d),
			b.CheckinID,
		)
	}

	// Print out each badge
	for _, b := range badges {
		printFn(b)

		// Print out each badge level
		for _, bb := range b.Levels {
			printFn(bb)
		}
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// printBeers turns a slice of *untappd.Beer structs into a human-friendly
// output format, and prints it to stdout.
func printBeers(beers []*untappd.Beer) {
	tw := tabWriter()

	// Print field header
	fmt.Fprintln(tw, "ID\tName\tBrewery\tStyle\tABV\tIBU")

	// Print out each beer
	for _, b := range beers {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%0.1f\t%03d\n",
			b.ID,
			b.Name,
			b.Brewery.Name,
			b.Style,
			b.ABV,
			b.IBU,
		)
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// printBreweries turns a slice of *untappd.Brewery structs into a human-friendly
// output format, and prints it to stdout.
func printBreweries(breweries []*untappd.Brewery) {
	tw := tabWriter()

	// Print field header
	fmt.Fprintln(tw, "ID\tName\tLocation")

	// Print out each brewery
	for _, b := range breweries {
		// Use country as default location
		l := b.Country

		// Add extended information if available
		if b.Location.City != "" && b.Location.State != "" {
			l = fmt.Sprintf("%s, %s, %s",
				b.Location.City,
				b.Location.State,
				b.Country,
			)
		}

		fmt.Fprintf(tw, "%d\t%s\t%s\n",
			b.ID,
			b.Name,
			l,
		)
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// printCheckins turns a slice of *untappd.Checkin structs into a human-friendly
// output format, and prints it to stdout.
func printCheckins(checkins []*untappd.Checkin) {
	tw := tabWriter()

	// Print field header
	fmt.Fprintln(tw, "ID\tName\tBrewery\tRating\tBadges\tToasts\tComments\tComment")

	// Print out each checkin
	for _, c := range checkins {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%0.2f\t%d\t%d\t%d\t%s\n",
			c.ID,
			c.Beer.Name,
			c.Brewery.Name,
			c.UserRating,
			len(c.Badges),
			len(c.Toasts),
			len(c.Comments),
			c.Comment,
		)
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// printUsers turns a slice of *untappd.User structs into a human-friendly
// output format, and prints it to stdout.  The info parameter allows
// extended information to be printed for user info.
func printUsers(users []*untappd.User, info bool) {
	tw := tabWriter()

	header := "ID\tUserName\tName"
	if info {
		header += "\tCheckins\tBadges\tBeers"
	}

	// Print field header
	fmt.Fprintln(tw, header)

	// Print out each user
	for _, u := range users {
		fmt.Fprintf(tw, "%d\t%s\t%s %s",
			u.UID,
			u.UserName,
			u.FirstName,
			u.LastName,
		)

		if info {
			fmt.Fprintf(tw, "\t%d\t%d\t%d",
				u.Stats.TotalCheckins,
				u.Stats.TotalBadges,
				u.Stats.TotalBeers,
			)
		}
		fmt.Fprintf(tw, "\n")
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// printVenues turns a slice of *untappd.Venue structs into a human-friendly
// output format, and prints it to stdout.
func printVenues(venues []*untappd.Venue) {
	tw := tabWriter()

	// Print field header
	fmt.Fprintln(tw, "ID\tName\tCategory\tPublic\tLocation")

	// Print out each checkin
	for _, v := range venues {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%t\t%s\n",
			v.ID,
			v.Name,
			v.Category,
			v.Public,
			fmt.Sprintf("%s, %s, %s",
				v.Location.City,
				v.Location.State,
				v.Location.Country,
			),
		)
	}

	// Flush buffered output
	if err := tw.Flush(); err != nil {
		log.Fatal(err)
	}
}

// tabWriter returns a *tabwriter.Writer appropriately configured
// for tabular output.
func tabWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
}
