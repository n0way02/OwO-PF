package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/Veraticus/findingway/internal/ffxiv"

	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	Token string

	Session  *discordgo.Session
	Channels []*Channel `yaml:"channels"`
}

type Channel struct {
	Name        string   `yaml:"name"`
	ID          string   `yaml:"id"`
	Duty        string   `yaml:"duty"`
	DataCentres []string `yaml:"dataCentres"`
}

func (d *Discord) Start() error {
	s, err := discordgo.New("Bot " + d.Token)
	if err != nil {
		return fmt.Errorf("Could not start Discord: %f", err)
	}
	s.ShouldRetryOnRateLimit = false

	err = s.Open()
	if err != nil {
		return fmt.Errorf("Could not open Discord session: %f", err)
	}

	// Configura o status do bot
	err = s.UpdateGameStatus(0, "os dados do PF no Discord")
	if err != nil {
		return fmt.Errorf("Could not set bot status: %f", err)
	}

	d.Session = s
	return nil
}

func (d *Discord) CleanChannel(channelId string) error {
	messages, err := d.Session.ChannelMessages(channelId, 100, "", "", "")
	if err != nil {
		return fmt.Errorf("Could not list messages: %f", err)
	}
	messageIds := []string{}
	for _, message := range messages {
		messageIds = append(messageIds, message.ID)
	}
	err = d.Session.ChannelMessagesBulkDelete(channelId, messageIds)
	if err != nil {
		return fmt.Errorf("Could not bulk delete messages: %f", err)
	}

	return nil
}

func (d *Discord) PostListings(channelId string, listings *ffxiv.Listings, duty string, dataCentre string) error {
	scopedListings := listings.ForDutyAndDataCentre(duty, dataCentre)

	mostRecent, err := scopedListings.MostRecentUpdated()
	if err != nil {
		return fmt.Errorf("Could not find most recently updated duty: %w", err)
	}
	if mostRecent != nil {
		mostRecentUpdated, err := mostRecent.UpdatedAt()
		if err != nil {
			return fmt.Errorf("Could not find most recently updatedAt: %w", err)
		}
		if mostRecentUpdated.After(time.Now().Add(-4 * time.Minute)) {
			scopedListings, err = scopedListings.UpdatedWithinLast(4 * time.Minute)
			if err != nil {
				return fmt.Errorf("Could not find most recently listings: %w", err)
			}
		}
	}

	headerEmbed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s PFs (%v)", duty, dataCentre),
		Type:  discordgo.EmbedTypeRich,
		Color: 0x6600ff,
		Description: fmt.Sprintf(
			"Encontrados %v PF's %v",
			len(scopedListings.Listings),
			fmt.Sprintf("<t:%v:R>", time.Now().Unix()),
		),
		Footer: &discordgo.MessageEmbedFooter{
			Text: strings.Repeat("\u3000", 20),
		},
	}
	headerMessageSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{headerEmbed},
	}
	_, err = d.Session.ChannelMessageSendComplex(channelId, headerMessageSend)
	if err != nil {
		return fmt.Errorf("Could not send header: %w", err)
	}

	fields := []*discordgo.MessageEmbedField{}
	for i, listing := range scopedListings.Listings {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.Creator,
			Value:  listing.PartyDisplay(),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetTags(),
			Value:  listing.GetDescription(),
			Inline: true,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   listing.GetExpires(),
			Value:  listing.GetUpdated(),
			Inline: true,
		})

		// Send a message every 5 listings
		if (i+1)%5 == 0 {
			err = d.sendMessage(channelId, fields)
			if err != nil {
				return fmt.Errorf("Could not send message: %w", err)
			}
			fields = []*discordgo.MessageEmbedField{}
		}
	}

	// Ensure we send any remaining messages
	if len(fields) != 0 {
		err = d.sendMessage(channelId, fields)
		if err != nil {
			return fmt.Errorf("Could not send message: %w", err)
		}
	}

	return nil
}

func (d *Discord) sendMessage(channelId string, fields []*discordgo.MessageEmbedField) error {
	embed := &discordgo.MessageEmbed{
		Type:   discordgo.EmbedTypeRich,
		Color:  0xec2d92,
		Fields: fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: strings.Repeat("\u3000", 20),
		},
	}
	messageSend := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
	_, err := d.Session.ChannelMessageSendComplex(channelId, messageSend)
	if err != nil {
		return err
	}

	return nil
}
