package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/bot"
	"github.com/diamondburned/arikawa/bot/extras/arguments"
	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
)

type Bot struct {
	// Context must not be embedded.
	Ctx *bot.Context
}

// Help prints the default help message.
func (bot *Bot) Help(m *gateway.MessageCreateEvent) (string, error) {
	me, err := bot.Ctx.Store.Me()
	if err != nil {
		return "", err
	}

	if m.Author.ID != me.ID {
		return "", nil
	}

	bot.Ctx.DeleteMessage(m.ChannelID, m.ID)
	return bot.Ctx.Help(), nil
}

// Add demonstrates the usage of typed arguments. Run it with "~add 1 2".
func (bot *Bot) Add(m *gateway.MessageCreateEvent, a, b int) error {
	content := fmt.Sprintf("%d + %d = %d", a, b, a+b)
	me, err := bot.Ctx.Store.Me()
	if err != nil {
		return err
	}

	if m.Author.ID != me.ID {
		return nil
	}

	bot.Ctx.DeleteMessage(m.ChannelID, m.ID)
	_, err = bot.Ctx.SendMessage(m.ChannelID, content, nil)
	return err
}

// Ping is a simple ping example, perhaps the most simple you could make it.
func (bot *Bot) Ping(m *gateway.MessageCreateEvent) error {
	me, err := bot.Ctx.Store.Me()
	if err != nil {
		return err
	}

	if m.Author.ID != me.ID {
		return nil
	}
	bot.Ctx.DeleteMessage(m.ChannelID, m.ID)
	_, err = bot.Ctx.SendMessage(m.ChannelID, "Pong!", nil)
	return err
}

// Dd is used to delete all the messages made by the user
func (bot *Bot) Dd(m *gateway.MessageCreateEvent) error {
	me, err := bot.Ctx.Store.Me()
	if err != nil {
		return err
	}

	if m.Author.ID != me.ID {
		return nil
	}

	msgs, err := bot.Ctx.State.Messages(m.ChannelID)
	if err != nil {
		log.Print("Error getting the messages: ", err)
	}
	bot.Ctx.DeleteMessage(m.ChannelID, m.ID)

	for _, msg := range msgs {

		if msg.Author.ID == me.ID {
			err := bot.Ctx.Session.DeleteMessage(m.ChannelID, msg.ID)
			if err != nil {
				log.Print("Error deleting a message: ", err)
			}
		}
	}

	return nil
}

// Embed is a simple embed creator. Its purpose is to demonstrate the usage of
// the ParseContent interface, as well as using the stdlib flag package.
func (bot *Bot) Embed(
	m *gateway.MessageCreateEvent, f *arguments.Flag) (*discord.Embed, error) {

	fs := arguments.NewFlagSet()

	var (
		title  = fs.String("title", "", "Title")
		author = fs.String("author", "", "Author")
		footer = fs.String("footer", "", "Footer")
		color  = fs.String("color", "#FFFFFF", "Color in hex format #hhhhhh")
	)

	if err := f.With(fs.FlagSet); err != nil {
		return nil, err
	}

	if len(fs.Args()) < 1 {
		return nil, fmt.Errorf("Usage: embed [flags] content...\n" + fs.Usage())
	}

	// Check if the color string is valid.
	if !strings.HasPrefix(*color, "#") || len(*color) != 7 {
		return nil, errors.New("Invalid color, format must be #hhhhhh")
	}

	// Parse the color into decimal numbers.
	colorHex, err := strconv.ParseInt((*color)[1:], 16, 64)
	if err != nil {
		return nil, err
	}

	// Make a new embed
	embed := discord.Embed{
		Title:       *title,
		Description: strings.Join(fs.Args(), " "),
		Color:       discord.Color(colorHex),
	}

	if *author != "" {
		embed.Author = &discord.EmbedAuthor{
			Name: *author,
		}
	}
	if *footer != "" {
		embed.Footer = &discord.EmbedFooter{
			Text: *footer,
		}
	}

	return &embed, err
}

// Btc is a simple btc follower
func (bot *Bot) Btc(m *gateway.MessageCreateEvent, f *arguments.Flag) (*discord.Embed, error) {
	me, err := bot.Ctx.Store.Me()
	if err != nil {
		return nil, err
	}

	if m.Author.ID != me.ID {
		return nil, nil
	}
	bot.Ctx.DeleteMessage(m.ChannelID, m.ID)
	var currency = f.String()
	if currency == "" {
		// Empty message, ignore
		currency = "usd"
	}

	rate, err := GetPrice(currency)
	if err != nil {
		return nil, err
	}
	// Make a new embed
	embed := discord.Embed{
		Title:       "Bitcoin Price",
		Description: "Current Bitcoin per " + strings.ToUpper(currency) + " price:",
		Fields: []discord.EmbedField{
			discord.EmbedField{
				Name:   strings.ToUpper(currency),
				Value:  rate,
				Inline: true,
			},
		},
		Color: discord.Color(0xf4a435),
	}

	embed.Author = &discord.EmbedAuthor{
		Name: m.Author.Username,
	}

	embed.Footer = &discord.EmbedFooter{
		Text: "Provided by Coinbase",
	}

	return &embed, err
}
