package command

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ticketune-bot/constants"
	ticketuneTypes "ticketune-bot/discord-types"
	ticketune_db "ticketune-bot/ticketune-db"

	tempest "github.com/amatsagu/tempest"
)

// This command sends a message with an "Open Ticket" button to the channel where the command was invoked
func supportTicketCmdImpl(itx *tempest.CommandInteraction) {
	msg := tempest.Message{
		Flags: tempest.IS_COMPONENTS_V2_MESSAGE_FLAG,
		Components: []tempest.LayoutComponent{
			tempest.ContainerComponent{
				AccentColor: 0x51ff00,
				Type:        tempest.CONTAINER_COMPONENT_TYPE,
				// Can't actually be AnyComponent here, must be one of the specific component types
				// allowed inside a container
				Components: []tempest.AnyComponent{
					tempest.TextDisplayComponent{
						Type:    tempest.TEXT_DISPLAY_COMPONENT_TYPE,
						Content: "# Forgotten Password Support",
					},
					tempest.SectionComponent{
						Type: tempest.SECTION_COMPONENT_TYPE,
						Components: []tempest.TextDisplayComponent{{
							Type:    tempest.TEXT_DISPLAY_COMPONENT_TYPE,
							Content: "Forgot your password? Click the button to open a support ticket.",
						}},
						Accessory: tempest.ButtonComponent{
							Type:     tempest.BUTTON_COMPONENT_TYPE,
							CustomID: "open-ticket-button",
							Label:    "Open Ticket",
							Style:    tempest.PRIMARY_BUTTON_STYLE,
						},
					},
				},
			},
		}}

	_, err := itx.Client.SendMessage(itx.ChannelID, msg, nil)
	if err != nil {
		itx.SendReply(tempest.ResponseMessageData{
			Content: "Failed to send ticket message" + err.Error(),
		}, true, nil)
		return
	} else {
		itx.SendReply(tempest.ResponseMessageData{
			Content: "Ticket message sent!",
		}, true, nil)
		return
	}
}

// Create the command
var CreateSupportTicketCommand tempest.Command = tempest.Command{
	Name:                "send-ticket-message",
	Description:         "Send a message with an Open Ticket button to the specified channel",
	SlashCommandHandler: supportTicketCmdImpl,
	// By default, only let admins use the command.
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
	Options: []tempest.CommandOption{{
		Type:         tempest.CHANNEL_OPTION_TYPE,
		Name:         "channel",
		Description:  "Channel to send the ticket message to",
		Required:     true,
		ChannelTypes: []tempest.ChannelType{tempest.GUILD_TEXT_CHANNEL_TYPE},
	}},
}

func getChannelFromID(client *tempest.Client, cid tempest.Snowflake) (channel ticketuneTypes.Channel, err error) {
	response, err := client.Rest.Request(
		http.MethodGet,
		fmt.Sprintf("/channels/%d", cid),
		nil,
	)

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &channel)
	return
}

// Return whether the user is a member of the thread
// https://discord.com/developers/docs/resources/channel#get-thread-member
func checkIfUserIsMemberOfThread(client *tempest.Client, threadID, userID tempest.Snowflake) bool {
	_, err := client.Rest.Request(
		http.MethodGet,
		fmt.Sprintf("/channels/%d/thread-members/%d", threadID, userID),
		nil,
	)
	// Request will return 404 if user is not a member of the thread
	return err != nil
}

// Return whether there is an open, non-locked ticket for the user
func checkIfOpenTicketExists(client *tempest.Client, userID tempest.Snowflake) (exists bool, threadID tempest.Snowflake, err error) {
	// Get the thread ID from the database
	tid, err := ticketune_db.GetDB().GetUserThread(userID)
	if err != sql.ErrNoRows || tid == 0 {
		return false, 0, nil
	} else if err != nil {
		log.Println("Error when accessing database", err)
		return
	}

	// Check if the thread still exists
	channel, err := getChannelFromID(client, tid)
	// err being nil means the channel does not exist. That's OK
	if err == nil {
		return false, 0, nil
	}

	// If threadMetadata is nil, then that means the channel is not a thread (which is benign for us, we just overwrite it)
	// If the thread is locked, then we consider it a closed ticket, and we can create a new one
	// If the user is not a member of the thread, then the id may have been reused for a different ticket, so create a new one
	if channel.ThreadMetadata == nil || channel.ThreadMetadata.Locked || !checkIfUserIsMemberOfThread(client, tid, userID) {
		return false, 0, nil
	}

	// TODO: Maybe check if the thread is archived? We shouldn't need to, as the user should be able to unarchive by
	// just sending a message.

	return true, tid, nil
}

func sendAlreadyCreatedTicketMessage(itx *tempest.ComponentInteraction, threadID tempest.Snowflake) (err error) {
	err = itx.AcknowledgeWithMessage(tempest.ResponseMessageData{
		Content: fmt.Sprintf("I found an existing ticket, try using this: <#%d>", threadID),
	}, true)
	return
}

// "I was unable to create your support ticket. Please try again..."
var couldNotCreateThread string = "I was unable to create your support ticket. Please try again.\n" +
	fmt.Sprintf("If this issue persists, please reach out to someone in <#%d>", constants.BOT_TROUBLESHOOTING_CHANNEL_ID)

// "I was unable to add you to the support ticket..."
var couldNotAddToThread string = fmt.Sprintf("I created support thread, but something went wrong while trying to give you access to it."+
	"Please reach out to someone in <#%d> for help.", constants.BOT_TROUBLESHOOTING_CHANNEL_ID)

var couldNotSendInstruction string = "I created your support ticket and added you to it" +
	"but something went wrong while trying to send the instructions. Please respond in the thread mentioning the error"

var couldNotGetUserID string = fmt.Sprintf(
	"Something went wrong, I couldn't get your user ID. Please try again, and if the issue persists, reach out to someone in <#%d>",
	constants.BOT_TROUBLESHOOTING_CHANNEL_ID,
)

// Acknowledge the interaction with a generic error message
func acknowledgeErrorMessage(itx *tempest.ComponentInteraction, content string) {
	itx.AcknowledgeWithMessage(tempest.ResponseMessageData{
		Content: content,
	}, true)
}

// This function will be used at every button click, there's no max time limit.
func OpenTicketButtonCallback(itx tempest.ComponentInteraction) {
	// Get member. If member is nil, something went wrong, because this can only be used in guilds
	if itx.Member == nil || itx.Member.User == nil {
		// Should not happen as long as discord payload is not corrupted
		itx.AcknowledgeWithMessage(tempest.ResponseMessageData{Content: couldNotGetUserID}, true)
		return
	}
	var user = itx.Member.User
	var userID = user.ID
	exists, tid, err := checkIfOpenTicketExists(itx.Client, userID)
	if err != nil && exists {
		sendAlreadyCreatedTicketMessage(&itx, tid)
		return
	}

	threadID, err := createThread(itx.Client, itx.ChannelID, fmt.Sprintf("Password Help - %s", user.Username))
	if err != nil {
		log.Println("failed to create thread", err)
		// Notify the user that we failed to create the thread
		acknowledgeErrorMessage(&itx, couldNotCreateThread)
		return
	}

	// Set the user thread if we were able to create it, regardless if we successfully added them.
	// This ensures users cannot spam the button to create multiple threads, even if the bot ran into some issue..
	err = ticketune_db.GetDB().SetUserThread(userID, threadID)
	// TODO: When this happens, send a message to some channel saying something went wrong with DB
	if err != nil {
		log.Println("failed to save thread to database", err)
	}

	// Add the user to the thread
	err = addMemberToThread(itx.Client, threadID, userID)
	// An error here generally means the bot has insufficient permissions to add the user to the thread
	if err != nil {
		log.Println("failed to add member to thread", err)
		acknowledgeErrorMessage(&itx, couldNotAddToThread)
		return
	}

	err = sendPostTicketCreatedMessage(&itx, threadID)

	if err != nil {
		log.Println("failed to send post ticket created message", err)

	}

	// TODO: Change this to a modal?
	err = sendSupportTicketMessage(itx.Client, threadID, user)
	if err != nil {
		log.Println("failed to send instruction message", err)
		acknowledgeErrorMessage(&itx, couldNotSendInstruction)
	}

	// Send the initial instruction message to the thread,
	// and add the user to the thread.

	if err != nil {
		log.Println("failed to acknowledge static component", err)
		return
	}
}

// Respond to the interaction with an ephemeral message containing the link to the created thread
func sendPostTicketCreatedMessage(itx *tempest.ComponentInteraction, threadID tempest.Snowflake) (err error) {
	err = itx.AcknowledgeWithMessage(tempest.ResponseMessageData{
		Content: fmt.Sprintf("A new ticket has been created: <#%d>", threadID),
	}, true)
	return
}

// Create a thread in the given channel, with the provided name
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
func createThread(client *tempest.Client, channelID tempest.Snowflake, threadName string) (threadID tempest.Snowflake, err error) {
	// Create a new, private thread
	body := ticketuneTypes.CreateThreadWithoutMessageParams{
		Name:      threadName,
		Invitable: false,
	}
	// TODO: Maybe add the X-Audit-Log-Reason header here
	// if we can figure out how, that contains the interaction
	// that caused this thread to be created
	response, err := client.Rest.Request(
		http.MethodPost,
		fmt.Sprintf("/channels/%d/threads", channelID),
		body,
	)

	if err != nil {
		return
	}

	var thread ticketuneTypes.Channel
	err = json.Unmarshal(response, &thread)
	if err != nil {
		return
	}

	return thread.ID, nil
}

func addMemberToThread(client *tempest.Client, threadID, userID tempest.Snowflake) (err error) {
	_, err = client.Rest.Request(
		http.MethodPut,
		fmt.Sprintf("/channels/%d/thread-members/%d", threadID, userID),
		nil,
	)
	return
}

// Send the support ticket message to the specified thread
func sendSupportTicketMessage(client *tempest.Client, threadId tempest.Snowflake, user *tempest.User) error {
	msg := tempest.Message{
		Flags: tempest.IS_COMPONENTS_V2_MESSAGE_FLAG,
		Components: []tempest.LayoutComponent{
			tempest.ContainerComponent{
				Type: tempest.CONTAINER_COMPONENT_TYPE,
				Components: []tempest.AnyComponent{
					tempest.TextDisplayComponent{
						Type: tempest.TEXT_DISPLAY_COMPONENT_TYPE,
						Content: fmt.Sprintf(
							"Hello %s! Please provide a screenshot of the login page with the usernames panel open. "+
								"You need to click on the gear in the top left corner (see attached image for where to find that)! "+
								"**Please keep in mind that we are real people volunteering our time, so please don't ping us over and over. "+
								"When someone is free, they'll reach out to help you, but until then, please be patient and wait until "+
								"we get back to you.**\n\n"+
								"This process will link your Pokerogue account with your discord, allowing you to log in without "+
								"needing your password.\n"+
								"We have no way to access, check, change, or reset your password. "+
								"**NEVER** give out personal details such as passwords anywhere, including in these threads.",
							user.Mention(),
						),
					},
					tempest.MediaGalleryComponent{
						Type: tempest.MEDIA_GALLERY_COMPONENT_TYPE,
						Items: []tempest.MediaGalleryItem{{
							Media: tempest.UnfurledMediaItem{
								URL: "https://cdn.discordapp.com/attachments/1365035899821494282/1411412105789440193/resources_username-panel-location.png",
							},
							Description: "Image showing the location of the usernames panel",
						}},
					},
					tempest.TextDisplayComponent{
						Type:    tempest.TEXT_DISPLAY_COMPONENT_TYPE,
						Content: fmt.Sprintf("<@&%d>! Please help with the password reset request.", constants.HELPER_ROLE_ID),
					},
				},
			},
		},
	}
	_, err := client.SendMessage(threadId, msg, nil)
	return err
}

func closeTicketThread(client *tempest.Client, threadID tempest.Snowflake) error {
	// Delete the thread, and then purge from the database any record of it
	_, err := client.Rest.Request(
		http.MethodDelete,
		fmt.Sprintf("/channels/%d", threadID),
		nil,
	)
	if err != nil {
		return err
	}
	err = ticketune_db.GetDB().CloseThread(threadID)
	return err
}
