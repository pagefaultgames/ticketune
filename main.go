/*
 * SPDX-FileCopyrightText: 2025 Pagefault Games
 * SPDX-FileContributor: SirzBenjie
 * SPDX-FileContributor: patapancakes
 * SPDX-FileContributor: Lugiadrien
 *
 * SPDX-License-Identifier: AGPL-3.0-or-later
 */

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pagefaultgames/ticketune/commands"
	"github.com/pagefaultgames/ticketune/db"

	"github.com/amatsagu/tempest"
)

func main() {
	// open (or create) the database
	log.Println("Initializing the database...")
	err := db.Init()
	if err != nil {
		log.Fatal("failed to initialize the database", err)
	}

	log.Println("Creating new Tempest client...")
	client := tempest.NewClient(tempest.ClientOptions{
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
	})

	addr := os.Getenv("LISTENING_ADDRESS")
	guildID, err := tempest.StringToSnowflake(os.Getenv("DISCORD_GUILD_ID"))
	if err != nil {
		log.Fatal("failed to parse env variable to snowflake", err)
	}

	// Register a simple ping command
	client.RegisterCommand(commands.PingCommand)
	client.RegisterCommand(commands.CreateSupportTicketCommand)
	client.RegisterComponent([]string{"open-ticket-button"}, commands.OpenTicketButtonCallback)
	client.RegisterCommand(commands.GetUserTicketCommand)
	client.RegisterCommand(commands.CloseCommand)
	client.RegisterCommand(commands.TryDiscordCommand)
	client.RegisterCommand(commands.FailDiscordCommand)
	client.RegisterCommand(commands.NoSaveCommmand)
	client.RegisterCommand(commands.RequestPanelCommand)
	client.RegisterCommand(commands.OldAccountCommandGroup)
	client.RegisterSubCommand(commands.OldAccountDefault, commands.OldAccountCommandGroup.Name)
	client.RegisterSubCommand(commands.OldAccountSpecific, commands.OldAccountCommandGroup.Name)
	client.RegisterCommand(commands.SayCommand)
	client.RegisterCommand(commands.WhichAccountCommand)
	client.RegisterCommand(commands.TechIssuesCommand)
	client.RegisterCommand(commands.PingSpamCommand)

	err = client.SyncCommandsWithDiscord([]tempest.Snowflake{guildID}, nil, false)
	if err != nil {
		log.Fatal("failed to sync local commands storage with Discord API", err)
	}

	http.HandleFunc("POST /discord/callback", client.DiscordRequestHandler)

	log.Printf("Serving application at: %s/discord/callback\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("something went terribly wrong", err)
	}
}
