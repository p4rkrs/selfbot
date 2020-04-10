package main

import (
	"log"
	"os"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/handler"
	"github.com/diamondburned/arikawa/state"
)

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`
func main() {
	var token = os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalln("No $BOT_TOKEN given.")
	}

	// Create a new Store
	store := state.NewDefaultStore(&state.DefaultStoreOptions{
		MaxMessages: 50,
	})

	// Create a new State with the store, necessary for selfbot
	s, err := state.NewWithStore(token, store)
	if err != nil {
		log.Fatal(err)
	}
	commands := &Bot{}

	// Create a PreHandler
	s.PreHandler = handler.New()
	s.PreHandler.Synchronous = true

	session, err := bot.New(s, commands)

	if err != nil {
		log.Fatalln(err)
	}

	// Register subcommands
	session.MustRegisterSubcommand(&Debug{})

	err = session.Open()

	if err != nil {
		log.Fatalln(err)
	}

	session.HasPrefix = bot.NewPrefix("!", "~")

	session.Start()

	// Subcommand demo, but this can be in another package.

	log.Println("Bot started")

	// As of this commit, Wait() will block until SIGINT or fatal. The past
	// versions close on call, but this one will block.
	// If for some reason you want the Cancel() function, manually make a new
	// context.
	if err := session.Wait(); err != nil {
		log.Fatalln("Gateway fatal error:", err)
	}
}
