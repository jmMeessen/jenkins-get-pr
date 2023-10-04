/*
Copyright © 2023 Jean-Marc Meessen jean-marc@meessen-web.org

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v55/github"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	//See https://github.com/schollz/progressbar
	// "github.com/schollz/progressbar/v3"
)

// quotaCmd represents the quota command
var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Gets the current GitHub API quota status",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		get_quota()
	},
}

func init() {
	rootCmd.AddCommand(quotaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quotaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quotaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ---
// Retrieves the GitHub API Quota
func get_quota() {
	limit, remaining := get_quota_data()
	fmt.Printf("V3 Limit: %d \nV3 Remaining %d \n\n", limit, remaining)

	limit_v4, remaining_v4, resetTimeString, secondsToGo := get_quota_data_v4()

	fmt.Printf("V4 Limit: %d \nV4 Remaining: %d \nV4 Reset time: %s (in %d secs)\n", limit_v4, remaining_v4, resetTimeString, secondsToGo)
}

// Retrieves the GitHub Quota.
func get_quota_data() (limit int, remaining int) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)

	client := github.NewClient(nil).WithAuthToken(ghToken)

	limitsData, _, err := client.RateLimits(context.Background())
	if err != nil {
		log.Printf("Error getting limit: %v", err)
		return 0, 0
	}
	return limitsData.Core.Limit, limitsData.Core.Remaining
}

/*
query {
  viewer {
    login
  }
  rateLimit {
    limit
    cost
    remaining
    resetAt
  }
}
*/

var quotaQuery struct {
	Viewer struct {
		Login string
	}
	RateLimit struct {
		Limit     int
		Cost      int
		Remaining int
		ResetAt   time.Time
	}
}

func get_quota_data_v4() (limit int, remaining int, resetAt string, secondsToReset int) {
	// retrieve the token value from the specified environment variable
	// ghTokenVar is global and set by the CLI parser
	ghToken := loadGitHubToken(ghTokenVar)
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	err := client.Query(context.Background(), &quotaQuery, nil)
	if err != nil {
		//FIXME: Better error handling
		log.Panic(err)
	}

	// pretty print the reset time (UTC)
	reset_time := quotaQuery.RateLimit.ResetAt
	resetTimeString := reset_time.Format(time.RFC1123)

	// compute how many seconds are before reset
	now := time.Now()
	diff := reset_time.Sub(now)
	secondsToGo := int(diff.Seconds())

	return quotaQuery.RateLimit.Limit, quotaQuery.RateLimit.Remaining, resetTimeString, secondsToGo
}

// // Get's the V4 quota, checks whether there is enough quota. If not will wait for the reset
// func checkIfSufficientQuota(expectedLoad int) {
// 	// initialize we  are called outside the normal flow
// 	initLoggers()

// 	limit, remaining, resetAt, secondsToReset := get_quota_data_v4()
// 	if isRootDebug || isDebugGet {
// 		loggers.debug.Printf("Quota: %d/%d (%d secs -> %s\n", remaining, limit, secondsToReset, resetAt)
// 		loggers.debug.Printf("Requesting to process %d\n", expectedLoad)
// 	}

// 	globalIsBigFile = false

// 	if expectedLoad >= limit {
// 		if isRootDebug || isDebugGet {
// 			loggers.debug.Printf("Expected load (%d) is higher then limit (%d)\n", expectedLoad, limit)
// 		}
// 		fmt.Printf("Expected load (%d) is higher then limit (%d)\n Crossing fingers and continuing...\n", expectedLoad, limit)
// 		globalIsBigFile = true
// 		return
// 	}

// 	if (expectedLoad + 20) > remaining {
// 		//Not enough resources, we need to wait
// 		waitForReset(secondsToReset)
// 	}
// 	// Else we do nothing as we are good to go.
// }

// // Wait for a certain number of seconds
// func waitForReset(secondsToReset int) {
// 	//TODO: check input value

// 	bar := progressbar.NewOptions(secondsToReset,
// 		progressbar.OptionShowBytes(false),
// 		progressbar.OptionSetDescription("Waiting for quota reset   "),
// 		progressbar.OptionSetPredictTime(false),
// 		progressbar.OptionShowBytes(false),
// 		progressbar.OptionFullWidth(),
// 		progressbar.OptionShowCount(),
// 		progressbar.OptionClearOnFinish(),
// 	)

// 	for i := 0; i < secondsToReset; i++ {
// 		err := bar.Add(1)
// 		if err != nil {
// 			log.Printf("Unexpected error updating progress bar (%v)\n", err)
// 		}
// 		time.Sleep(1 * time.Second)
// 	}

// 	// Clear the progress bar
// 	bar.Reset()
// 	err := bar.Finish()
// 	if err != nil {
// 		log.Printf("Unexpected error clearing progress bar (%v)\n", err)
// 	}
// }
