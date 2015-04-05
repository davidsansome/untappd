package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mdlayher/untappd"
)

const (
	// appName is the name of this binary.
	appName = "untappdctl"
)

func main() {
	// Initialize new CLI app
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "query and display information from Untappd APIv4"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Matt Layher",
			Email: "mdlayher@gmail.com",
		},
	}

	// Add global flags for Untappd API client ID and client secret
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "client_id",
			Usage:  "client ID parameter for Untappd APIv4",
			EnvVar: "UNTAPPD_ID",
		},
		cli.StringFlag{
			Name:   "client_secret",
			Usage:  "client secret parameter for Untappd APIv4",
			EnvVar: "UNTAPPD_SECRET",
		},
	}

	// Add commands mirroring available untappd.Client services
	app.Commands = []cli.Command{
		userCommand(),
	}

	// Print all log output to stderr, so stdout only contains Untappd data
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
	log.SetPrefix(appName + "> ")

	app.Run(os.Args)
}

// untappdClient creates an initialized *untappd.Client using the client ID
// and secret from global CLI context.
func untappdClient(ctx *cli.Context) *untappd.Client {
	c, err := untappd.NewClient(
		ctx.GlobalString("client_id"),
		ctx.GlobalString("client_secret"),
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

// printRateLimit is a helper method which displays the remaining rate limit
// header for each HTTP request.
func printRateLimit(res *http.Response) {
	const header = "X-Ratelimit-Remaining"
	if v := res.Header.Get(header); v != "" {
		log.Printf("%s: %s", header, v)
	}
}

// mustStringArg is a helper method which checks for a string argument in the
// CLI context, and prints a help message if it is not found.
func mustStringArg(ctx *cli.Context, name string) string {
	a := ctx.Args().First()
	if a == "" {
		log.Fatalf("missing argument: %s", name)
	}

	return a
}

// offsetLimitSort retrieves a triple of offset, limit, and sort parameters
// from CLI context, as accepted by the Untappd API.
func offsetLimitSort(ctx *cli.Context) (int, int, untappd.Sort) {
	offset, limit := ctx.Int("offset"), ctx.Int("limit")

	// If no sort found, ignore sanity checks
	sort := ctx.String("sort")
	if sort == "" {
		return offset, limit, untappd.Sort("")
	}

	// Ensure sort type is valid
	for _, s := range untappd.Sorts() {
		// Return on valid sort
		if sort == string(s) {
			return offset, limit, s
		}
	}

	// Die on invalid sort, and show options
	log.Fatalf("invalid sort type %q (options: %s)", sort, untappd.Sorts())
	return offset, limit, untappd.Sort("")
}
