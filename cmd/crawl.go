package cmd

import (
	"net/url"
	"os"

	"github.com/iangregson/spider-go/crawler"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// crawlCmd represents the crawl command
var crawlCmd = &cobra.Command{
	Use:     "crawl <URL>",
	Args:    cobra.ExactArgs(1),
	Short:   "Crawl a given URL",
	Long:    `Given a starting URL, visit and print each URL - and a list of its links - found on the same domain.`,
	Example: "spider-go crawl http://crawler-test.com/",
	Run: func(cmd *cobra.Command, args []string) {
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		c := crawler.New(concurrency)

		u, err := url.Parse(args[0])
		if err != nil {
			log.Error().Err(err).Msg("Please supply a valid fully qualified URL.")
			os.Exit(1)
		}

		c.Crawl(u)
		log.Info().Msgf("Completed crawl successfully. Visited %d URLs.", c.VisitedCount())
	},
}

func init() {
	rootCmd.AddCommand(crawlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crawlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// crawlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	crawlCmd.PersistentFlags().Int("concurrency", 1, "Set the maximum concurrency e.g. 16")
}
