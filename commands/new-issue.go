package commands

import (
	"context"
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"github.com/amatsagu/tempest"
	"github.com/google/go-github/v74/github"
	githubClient "github.com/pagefaultgames/ticketune/github-client"
)

// Register the command (add this to your command registration logic)
var NewIssueCommand = tempest.Command{
	Name:                "new-issue",
	Type:                tempest.MESSAGE_COMMAND_TYPE,
	SlashCommandHandler: newIssueCommand, // Despite the field name, this is the intended way to handle message commands
}

var ErrNoResolvedData = errors.New("no resolved data in interaction")
var ErrNoMessage = errors.New("no message found in resolved data")

const CreateIssueModalId = "create-gh-issue-modal"

// Duration after which we time out the issue creation request
const issueTimeout = time.Minute

// TODO: Make the slash command generate a form to fill out the issue title/body/labels
// and then use that to create the issue
// Also allow uploading images / screenshots.

func newIssueMessageVariant(itx *tempest.CommandInteraction) (messageContents string, err error) {
	resolved := itx.Data.Resolved
	// If this is nil, discord did something wrong.
	if resolved == nil {
		itx.SendLinearReply("Error: Message missing", true)
		return "", ErrNoResolvedData
	}

	msg, ok := resolved.Messages[itx.Data.TargetID]
	if !ok {
		itx.SendLinearReply("Error: Message not found", true)
		return "", ErrNoMessage
	}

	// Link to message that generated the issue
	messageLink := "https://discord.com/channels/" + itx.GuildID.String() + "/" + msg.ChannelID.String() + "/" + msg.ID.String()

	// get length of message link, truncate content of message to fit within 4000 characters minus link
	// need to use runes to properly handle unicode

	// Discord supports unicode
	// Truncate the message to fit within the limit
	messageContents = msg.Content
	const maxIssueBodyLength = 4000
	msgRunes := []rune(msg.Content)
	msgLength := len(msgRunes)
	linkText := "[Related Discord message](" + messageLink + ")\n\n"
	availableLength := maxIssueBodyLength - len([]rune(linkText)) // leave space for the link text

	// Available length *should* always be >= 0, as link text is unlikely to exceed 4000 characters, unless discord gave us huge snowflakes
	if msgLength > availableLength && availableLength >= 0 {
		messageContents = linkText + string(msgRunes[:availableLength])
	} else if availableLength >= 0 {
		messageContents = linkText + msg.Content
	}

	return messageContents, nil
}

func sendIssueModal(itx *tempest.CommandInteraction, prefillBody string) {
	// get the contents of the message they right clicked on..

	err := itx.SendModal(tempest.ResponseModalData{
		CustomID: CreateIssueModalId,
		Title:    "Create GitHub Issue",
		Components: []tempest.LayoutComponent{
			// Can have at most 5 components.
			tempest.LabelComponent{
				Type:        tempest.LABEL_COMPONENT_TYPE,
				Label:       "Issue Title",
				Description: "Summarize the issue in a few words, (omit the [Bug] prefix)",
				Component: tempest.TextInputComponent{
					Type:      tempest.TEXT_INPUT_COMPONENT_TYPE,
					CustomID:  "issue-title",
					Style:     tempest.SHORT_TEXT_INPUT_STYLE,
					MaxLength: 200, // Github title limit is 256. Leave extra space for tags
					Required:  true,
				},
			},
			tempest.LabelComponent{
				Type:  tempest.LABEL_COMPONENT_TYPE,
				Label: "Category Labels (select up to 4)",
				Component: tempest.StringSelectComponent{
					Type:        tempest.STRING_SELECT_COMPONENT_TYPE,
					CustomID:    "issue-labels",
					MinValues:   1,
					MaxValues:   4,
					Placeholder: "Select 1 or more categories the bug falls under",
					Options: []tempest.SelectMenuOption{
						{Label: "Move", Value: "Move", Description: "Issues with a Pokémon move"},
						{Label: "Ability", Value: "Ability", Description: "Issues with abilities"},
						{Label: "Item", Value: "Item", Description: "Issues with items"},
						{Label: "Sprite/Animation", Value: "Sprite/Animation", Description: "Issues with sprites or animations"},
						{Label: "UI/UX", Value: "UI/UX", Description: "User interface / user experience issues"},
						{Label: "Mystery Encounter", Value: "Mystery Encounter", Description: "Issues with a mystery encounter"},
						{Label: "Audio", Value: "Audio", Description: "Issues with sound effects or music"},
						{Label: "Challenges", Value: "Challenges", Description: "Challenge mode(s) related"},
						{Label: "Miscellaneous", Value: "Miscellaneous", Description: "None of the other categories fit"},
						{Label: "Beta", Value: "Beta", Description: "Only present on Beta (do not select unless it is known the issue does not happen on main)"},
					},
				},
			},
			tempest.LabelComponent{
				Type:        tempest.LABEL_COMPONENT_TYPE,
				Label:       "Issue Description",
				Description: "Describe the issue (GitHub flavored markdown supported)",
				Component: tempest.TextInputComponent{
					Type:      tempest.TEXT_INPUT_COMPONENT_TYPE,
					CustomID:  "issue-description",
					Style:     tempest.PARAGRAPH_TEXT_INPUT_STYLE,
					Value:     prefillBody,
					Required:  true,
					MaxLength: 3800, // Need space for the URL that will be inserted.
				},
			},
			tempest.LabelComponent{
				Type:        tempest.LABEL_COMPONENT_TYPE,
				Label:       "Steps to reproduce",
				Description: "Describe the steps to reproduce this bug",
				Component: tempest.TextInputComponent{
					Type:     tempest.TEXT_INPUT_COMPONENT_TYPE,
					CustomID: "issue-steps",
					Style:    tempest.PARAGRAPH_TEXT_INPUT_STYLE,
					Required: false,
				},
			},
			tempest.LabelComponent{
				Type:        tempest.LABEL_COMPONENT_TYPE,
				Label:       "Additional context",
				Description: "Add any other context about the problem here",
				Component: tempest.TextInputComponent{
					Type:     tempest.TEXT_INPUT_COMPONENT_TYPE,
					CustomID: "issue-additional-context",
					Style:    tempest.PARAGRAPH_TEXT_INPUT_STYLE,
					Required: false,
				},
			},
		},
	})
	if err != nil {
		itx.SendLinearReply("Error: Failed to send issue modal: "+err.Error(), true)
		return
	}
}

// newIssueCommand handles the /new-issue Message command
func newIssueCommand(itx *tempest.CommandInteraction) error {
	var isMessage bool = itx.Interaction != nil && itx.Data.TargetID != 0

	prefillBody := ""

	// switch for now, later we can add special logic for slash command variant
	if isMessage {
		prefillBody, _ = newIssueMessageVariant(itx)
	} /* else {
	 TODO: uncomment and implement special logic for slash command variant
	}  */

	sendIssueModal(itx, prefillBody)

	return nil
}

// Helper function to extract the component from a modal response's label
// Returns the component if found, (otherwise the zero value)
func getLabelComponent[T tempest.StringSelectComponent | tempest.TextInputComponent](itx tempest.ModalInteraction, expectedIndex int) T {
	var zero T
	if len(itx.Data.Components) <= expectedIndex {
		return zero
	}
	label, ok := itx.Data.Components[expectedIndex].(tempest.LabelComponent)
	if !ok {
		return zero
	}
	// Get the child component
	component, ok := label.Component.(T)
	if !ok {
		return zero
	}
	return component
}

// Error indicating that the issue creation request timed out
var ErrIssueTimeout = errors.New("issue creation timed out")

func HandleNewIssueModal(mitx tempest.ModalInteraction) {
	// This means that the interaction was not in a guild, which should not be possible unless discord is broken
	if mitx.Member == nil {
		_ = mitx.AcknowledgeWithLinearMessage("Error: Unable to identify user", true)
		return
	}
	title := getLabelComponent[tempest.TextInputComponent](mitx, 0).Value
	if title == "" {
		mitx.AcknowledgeWithLinearMessage("Error: Unable to find issue title", true)
		return
	} else {
		title = "[Bug] " + strings.TrimSpace(title)
	}

	issueLabels := getLabelComponent[tempest.StringSelectComponent](mitx, 1).Values
	issueLabels = append([]string{"Triage"}, issueLabels...)

	description := getLabelComponent[tempest.TextInputComponent](mitx, 2).Value

	// description is required, so this should never happen unless discord is broken
	if description == "" {
		mitx.AcknowledgeWithLinearMessage("Error: Unable to find issue description", true)
		return
	}

	stepsComponent := getLabelComponent[tempest.TextInputComponent](mitx, 3)
	steps := stepsComponent.Value
	if steps == "" {
		steps = "_No response_"
	}

	additionalContext := getLabelComponent[tempest.TextInputComponent](mitx, 4).Value
	if additionalContext == "" {
		steps = "_No response_"
	}

	log.Printf("Creating issue with title: %s, labels: %v", title, issueLabels)

	err := mitx.Defer(true)
	if err != nil {
		log.Print("Failed to defer issue modal: ", err)
	}

	if err != nil {
		log.Print("Failed to send follow-up message: ", err)
	}

	var issueBody string
	// If member is nil, discord broke or someone made this command available in DMs (which they shouldn't have)
	// If user is nil, discord broke, since it can only be nil in MESSAGE_CREATE and MESSAGE_UPDATE events
	// 	source: https://discord.com/developers/docs/resources/guild#guild-member-object-guild-member-structure
	if mitx.Member != nil && mitx.Member.User != nil {
		issueBody = "Bug report initiated by Discord user **" + mitx.Member.User.Username + "**\n"
	}
	issueBody +=
		"### Describe the bug\n\n" + description +
			"\n\n### Reproduction\n\n" + steps +
			"\n\n### Additional context\n\n" + additionalContext

	go func() {
		// We have 30 seconds to respond
		ctx, cancel := context.WithTimeoutCause(context.Background(), issueTimeout, ErrIssueTimeout)
		defer cancel()
		context.AfterFunc(ctx, func() {
			if context.Cause(ctx) == ErrIssueTimeout {
				mitx.SendLinearFollowUp("Issue creation timed out. Please try again later.", true)
			}
		})
		// context.
		// TODO: replace with actual owner/org
		owner := "pagefaultgames"
		repo := "pokerogue"
		issueRequest := &github.IssueRequest{
			// Assuming discord respected our 200 character limit,
			// this will be less than the 256 character github limit for titles
			Title:  github.Ptr(title),
			Body:   github.Ptr(issueBody),
			Type:   github.Ptr("bug"),
			Labels: &issueLabels,
		}
		issue, resp, err := githubClient.Client.Issues.Create(ctx, owner, repo, issueRequest)
		if err != nil {
			if resp != nil && resp.Rate.Remaining == 0 {
				mitx.SendLinearFollowUp("GitHub rate limit exceeded. Please try again later.", true)
				return
			} else if resp != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Failed to read GitHub API response body: %v", err)
				}
				log.Printf("GitHub API response body: %s", string(body))
			}
			mitx.SendLinearFollowUp("Failed to create issue: "+err.Error(), true)
			return
		}
		mitx.SendLinearFollowUp("Issue created: "+issue.GetHTMLURL(), true)
	}()
}
