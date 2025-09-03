package commands

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/amatsagu/tempest"
	"github.com/google/go-github/v74/github"
	githubClient "github.com/pagefaultgames/ticketune/github-client"
)

// Register the command (add this to your command registration logic)
var NewIssueCommand = tempest.Command{
	Name:                "new-issue",
	Description:         "Create a new, blank GitHub issue.",
	SlashCommandHandler: newIssueCommand,
}

// TODO: Make the slash command generate a form to fill out the issue title/body/labels
// and then use that to create the issue
// Also allow uploading images / screenshots.

// var issueChan = make(chan github.IssueRequest, 100)

// newIssueCommand handles the /new-issue command
func newIssueCommand(itx *tempest.CommandInteraction) {
	// titleVal, present := itx.GetOptionValue("title")
	// title, ok := titleVal.(string)
	// if !present || !ok || title == "" {
	// 	itx.SendLinearReply("Error: title is required", true)
	// 	return
	// }

	_ = itx.SendLinearReply("Sending a request to create the issue...", false)

	go func() {
		// We have 30 seconds to respond
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// TODO: replace with actual owner/org
		owner := "owner"
		repo := "repo-name"
		issueRequest := &github.IssueRequest{
			Title:  github.Ptr("my new title"),
			Body:   github.Ptr("created by ticketune"),
			Labels: &[]string{"Bug"},
		}
		issue, resp, err := githubClient.Client.Issues.Create(ctx, owner, repo, issueRequest)
		if err != nil {
			if resp != nil && resp.Rate.Remaining == 0 {
				itx.SendLinearFollowUp("GitHub rate limit exceeded. Please try again later.", true)
				return
			} else if resp != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Failed to read GitHub API response body: %v", err)
				}
				log.Printf("GitHub API response body: %s", string(body))
			}
			itx.SendLinearFollowUp("Failed to create issue: "+err.Error(), true)
			return
		}
		itx.SendLinearFollowUp("Issue created: "+issue.GetHTMLURL(), true)
	}()
}
