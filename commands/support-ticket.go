/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package commands

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pagefaultgames/ticketune/constants"
	"github.com/pagefaultgames/ticketune/db"
	"github.com/pagefaultgames/ticketune/types"
	"github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

// This command sends a message with an "Open Ticket" button to the channel where the command was invoked
func supportTicketCmdImpl(itx *tempest.CommandInteraction) error {
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

	channel, presence := itx.GetOptionValue("channel")

	channelID := itx.ChannelID
	if presence {
		channelID, _ = tempest.StringToSnowflake(channel.(string))
	}

	_, err := itx.Client.SendMessage(channelID, msg, nil)
	if err != nil {
		itx.SendReply(tempest.ResponseMessageData{
			Content: "Failed to send ticket message" + err.Error(),
		}, true, nil)
		return nil
	} else {
		itx.SendReply(tempest.ResponseMessageData{
			Content: "Ticket message sent!",
		}, true, nil)
		return nil
	}
}

// Create the command
var CreateSupportTicketCommand = tempest.Command{
	Name:                "send-ticket-message",
	Description:         "Send a message with an Open Ticket button to the specified channel",
	SlashCommandHandler: supportTicketCmdImpl,
	// By default, only let admins use the command.
	// This is not a hard restriction imposed by the bot, but is simply the permission level
	// that shows up in the Discord UI. Guild admins can tweak this command to allow other roles to use it.
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
	GuildID:             constants.DISCORD_GUILD_ID,
	Options: []tempest.CommandOption{{
		Type:         tempest.CHANNEL_OPTION_TYPE,
		Name:         "channel",
		Description:  "Channel to send the ticket message to",
		Required:     true,
		ChannelTypes: []tempest.ChannelType{tempest.GUILD_TEXT_CHANNEL_TYPE},
	}},
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
	return err == nil
}

// Return whether there is an open, non-locked ticket for the user
func checkIfOpenTicketExists(client *tempest.Client, userID tempest.Snowflake) (bool, tempest.Snowflake, error) {
	// Get the thread ID from the database
	tid, err := db.Get().GetUserThread(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}

		return false, 0, err
	}

	// Check if the thread still exists
	channel, err := utils.GetChannelFromID(client, tid)
	// If error is not nil, either channel does not exist, or we couldn't unmarshal the response
	// in either case, proceed as though there is no existing ticket
	if err != nil {
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

// Respond to the interaction with an ephemeral message containing the link to the created thread
func sendAlreadyCreatedTicketMessage(itx *tempest.ComponentInteraction, threadID tempest.Snowflake) error {
	err := itx.AcknowledgeWithMessage(tempest.ResponseMessageData{
		Content: fmt.Sprintf("I found an existing ticket, try using this: <#%d>", threadID),
	}, true)
	if err != nil {
		return err
	}

	return nil
}

// "I was unable to create your support ticket. Please try again..."
var couldNotCreateThread = fmt.Sprintf(
	"I was unable to create your support ticket. Please try again.\n"+
		"If this issue persists, please reach out to someone in <#%d>.",
	constants.BOT_TROUBLESHOOTING_CHANNEL_ID,
)

// "I was unable to add you to the support ticket..."
var couldNotAddToThread = fmt.Sprintf(
	"I created support thread, but something went wrong while trying to give you access to it. "+
		"Please reach out to someone in <#%d> for help, and mention that I was unable to give you access to your password reset ticket.",
	constants.BOT_TROUBLESHOOTING_CHANNEL_ID,
)

// "something went wrong while trying to send the instructions..."
var couldNotSendInstruction = fmt.Sprintf(
	"I created your support ticket and added you to it, "+
		"but something went wrong while trying to send the instructions. Please reach out to someone in <#%d>, and mention that I could not send the instructions.",
	constants.BOT_TROUBLESHOOTING_CHANNEL_ID)

// "Something went wrong, I couldn't get your user ID..."
var couldNotGetUserID = fmt.Sprintf(
	"Something went wrong, I couldn't get your user ID. Please try again, and if the issue persists, reach out to someone in <#%d>.",
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

	user := itx.Member.User
	userID := user.ID

	// Discard any errors from checkIfOpenTicketExists, proceeding as though no ticket
	// exists.
	exists, tid, _ := checkIfOpenTicketExists(itx.Client, userID)
	if exists {
		sendAlreadyCreatedTicketMessage(&itx, tid)
		return
	}

	threadID, err := createThread(itx.Client, constants.TICKET_CHANNEL_ID, fmt.Sprintf("Password Help - %s", user.Username))
	if err != nil {
		log.Println("failed to create thread", err)
		// Notify the user that we failed to create the thread
		acknowledgeErrorMessage(&itx, couldNotCreateThread)
		return
	}

	// Set the user thread if we were able to create it, regardless if we successfully added them.
	// This ensures users cannot spam the button to create multiple threads, even if the bot ran into some issue..
	err = db.Get().SetUserThread(userID, threadID)
	// TODO: When this happens, send a message to some channel saying something went wrong with DB
	if err != nil {
		log.Println("failed to save thread to database", err)
	}

	// Give the user permission to view and send messages in threads in the ticket channel
	err = giveUserTicketChannelPerms(itx.Client, userID)
	if err != nil {
		log.Println("failed to give user ticket channel perms", err)
		acknowledgeErrorMessage(&itx, couldNotAddToThread)
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
		// This code path means that the bot was not able to reply with a simple message.
		// There's nothing we can do to communicate with the user, but they would have still had a ticket opened.
		// Proceed to try to send the instructions message, but log the error
		log.Println("failed to send post ticket created message", err)
	}

	// TODO: Change this to a modal?
	err = sendSupportTicketMessage(itx.Client, threadID, user)
	if err != nil {
		log.Println("failed to send instruction message", err)
		acknowledgeErrorMessage(&itx, couldNotSendInstruction)
	}
}

// Respond to the interaction with an ephemeral message containing the link to the created thread
func sendPostTicketCreatedMessage(itx *tempest.ComponentInteraction, threadID tempest.Snowflake) error {
	err := itx.AcknowledgeWithMessage(tempest.ResponseMessageData{
		Content: fmt.Sprintf("A new ticket has been created: <#%d>", threadID),
	}, true)
	if err != nil {
		return err
	}

	return nil
}

// Create a thread in the given channel, with the provided name
// https://discord.com/developers/docs/resources/channel#start-thread-without-message
func createThread(client *tempest.Client, channelID tempest.Snowflake, threadName string) (tempest.Snowflake, error) {
	// TODO: Maybe add the X-Audit-Log-Reason header here
	// if we can figure out how, that contains the interaction
	// that caused this thread to be created
	response, err := client.Rest.Request(
		http.MethodPost,
		fmt.Sprintf("/channels/%d/threads", channelID),
		types.CreateThreadWithoutMessageParams{ // Create a new, private thread
			Name:      threadName,
			Invitable: false,
		},
	)
	if err != nil {
		return tempest.Snowflake(0), err
	}

	var thread types.Channel
	err = json.Unmarshal(response, &thread)
	if err != nil {
		return tempest.Snowflake(0), err
	}

	return thread.ID, nil
}

func addMemberToThread(client *tempest.Client, threadID, userID tempest.Snowflake) error {
	_, err := client.Rest.Request(
		http.MethodPut,
		fmt.Sprintf("/channels/%d/thread-members/%d", threadID, userID),
		nil,
	)
	if err != nil {
		return err
	}

	return nil
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
							"Hello %s!\n"+
								"Please provide a screenshot of the login page __with the usernames panel open__.\n"+
								"You need to __click on the gear in the top left corner__ (see attached image for where to find that)!\n"+
								"**Please keep in mind that we are real people volunteering our time, so please don't ping us over and over. "+
								"When someone is free, they'll reach out to help you, but until then, please be patient and wait until "+
								"we get back to you.**\n\n"+
								"This process will link your PokéRogue account with the Discord account you used to open this ticket, allowing you to log in without "+
								"needing your password.\n"+
								"We have __no way__ to access, check, change, or reset your password.\n"+
								"Also, **NEVER** give out personal details such as passwords anywhere and to anyone, including in these threads.",
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
	if err != nil {
		return err
	}

	return nil
}

// Give the user ID permissions to view, send messages in threads, and read message history in the ticket channel
func giveUserTicketChannelPerms(client *tempest.Client, userID tempest.Snowflake) error {
	_, err := client.Rest.Request(
		http.MethodPut,
		fmt.Sprintf("/channels/%d/permissions/%d", constants.TICKET_CHANNEL_ID, userID),
		types.EditChannelPermissionsParams{
			Allow: tempest.SEND_MESSAGES_IN_THREADS_PERMISSION_FLAG | tempest.VIEW_CHANNEL_PERMISSION_FLAG | tempest.READ_MESSAGE_HISTORY_PERMISSION_FLAG,
			Type:  types.MEMBER_TYPE,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
