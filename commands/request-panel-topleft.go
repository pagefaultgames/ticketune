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
	"log"

	utils "github.com/pagefaultgames/ticketune/utils"

	"github.com/amatsagu/tempest"
)

const requestPanelCommandDescription = "Ping and ask the user to provide a screenshot of the login page with the usernames panel open"

const requestPanelCommandMsg = "Could you please provide a screenshot of the login page __with the usernames panel open or the error code it might display__?\n" +
	"To try opening the usernames panel, click __on the gear in the top left corner__ - see this image for clarification!"

var RequestPanelCommand = tempest.Command{
	Name:                "request-panel-topleft",
	Description:         requestPanelCommandDescription,
	RequiredPermissions: tempest.ADMINISTRATOR_PERMISSION_FLAG,
	Options:             []tempest.CommandOption{NO_PING_OPTION},
	SlashCommandHandler: requestPanelCommandImpl,
	Contexts:            []tempest.InteractionContextType{tempest.GUILD_CONTEXT_TYPE},
}

func requestPanelCommandImpl(itx *tempest.CommandInteraction) {
	// Get the user associated with this thread (this handles responding to the interaction on error)
	userID, err := utils.GetUserFromThread(itx)
	if err != sql.ErrNoRows && err != nil {
		return
	}
	responseMsg := "The user has been reminded to provide a screenshot with the usernames panel open."
	noPing, _ := utils.GetOption[bool](itx, "no-ping", false)
	msgContent := requestPanelCommandMsg

	switch {
	case !noPing && err == nil:
		log.Println("Error fetching user for thread:", err)
		responseMsg = "I couldn't find a user associated with this thread in my database, so I can't ping them." +
			"However, I've sent the login message to the thread."
	case !noPing:
		msgContent = "Hi <@" + userID.String() + ">!\n" + requestPanelCommandMsg
	}

	msg := tempest.Message{
		Flags: tempest.IS_COMPONENTS_V2_MESSAGE_FLAG,
		Components: []tempest.LayoutComponent{
			tempest.ContainerComponent{
				Type: tempest.CONTAINER_COMPONENT_TYPE,
				Components: []tempest.AnyComponent{
					tempest.TextDisplayComponent{
						Type:    tempest.TEXT_DISPLAY_COMPONENT_TYPE,
						Content: msgContent,
					},
					tempest.MediaGalleryComponent{
						Type: tempest.MEDIA_GALLERY_COMPONENT_TYPE,
						Items: []tempest.MediaGalleryItem{{
							Media: tempest.UnfurledMediaItem{
								URL: "https://cdn.discordapp.com/attachments/1284892701023666176/1433485614598062170/new_logscreen_obvious.png",
							},
							Description: "Image showing the location of the usernames panel",
						}},
					},
				},
			},
		},
	}

	_, err = itx.Client.SendMessage(
		itx.ChannelID,
		msg,
		nil,
	)
	if err != nil {
		itx.SendLinearReply("Something went wrong trying to send the message: "+err.Error(), true)
		return
	}

	itx.SendLinearReply(responseMsg, true)
}
