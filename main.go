package main

import (
	"log"
	"net/http"
	"os"

	command "ticketune-bot/commands"
	ticketune_db "ticketune-bot/ticketune-db"

	"ticketune-bot/constants"

	tempest "github.com/amatsagu/tempest"
)

func main() {
	constants.InitConstants()
	// open (or create) the database
	log.Println("Initializing the database...")
	if err := ticketune_db.InitDB(); err != nil {
		log.Fatalln("failed to initialize the database", err)
	}

	log.Println("Creating new Tempest client...")
	client := tempest.NewClient(tempest.ClientOptions{
		Token:     os.Getenv("DISCORD_BOT_TOKEN"),
		PublicKey: os.Getenv("DISCORD_PUBLIC_KEY"),
	})

	addr := os.Getenv("LISTENING_ADDRESS")
	testServerID, err := tempest.StringToSnowflake(os.Getenv("DISCORD_GUILD_ID"))
	if err != nil {
		log.Fatalln("failed to parse env variable to snowflake", err)
	}

	// Register a simple ping command
	client.RegisterCommand(command.PingCommand)
	client.RegisterCommand(command.CreateSupportTicketCommand)
	client.RegisterComponent([]string{"open-ticket-button"}, command.OpenTicketButtonCallback)

	err = client.SyncCommandsWithDiscord([]tempest.Snowflake{testServerID}, nil, false)
	if err != nil {
		log.Fatalln("failed to sync local commands storage with Discord API", err)
	}

	http.HandleFunc("POST /discord/callback", client.DiscordRequestHandler)

	log.Printf("Serving application at: %s/discord/callback\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("something went terribly wrong", err)
	}
}
