package main

import (
	"errors"
	"log"

	"github.com/dying/selfbot/config"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/handler"
	"github.com/diamondburned/arikawa/state"
)

// To run, do `BOT_TOKEN="TOKEN HERE" go run .`
func main() {
	config.Load()
	//	log.Print(config.Conf.Bot.Token)
	var token = config.Conf.Bot.Token

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
	var cmdError error
	var UnknownCommand *bot.ErrUnknownCommand
	session.ErrorLogger(cmdError)
	if errors.As(cmdError, &UnknownCommand) {
		log.Print(UnknownCommand.Error)
	}

	// Register subcommands
	session.MustRegisterSubcommand(&Debug{})

	err = session.Open()
	defer session.Close()
	if err != nil {
		log.Fatalln(err)
	}

	session.HasPrefix = bot.NewPrefix(config.Conf.Bot.Prefix)
	session.Start()

	// Subcommand demo, but this can be in another package.

	log.Println("Bot started")

	// As of this commit, Wait() will block until SIGINT or fatal. The past
	// versions close on call, but this one will block.
	// If for some reason you want the Cancel() function, manually make a new
	// context.
	select {}
}
