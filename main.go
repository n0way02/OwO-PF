package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Veraticus/findingway/internal/discord"
	"github.com/Veraticus/findingway/internal/scraper"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Channels []*discord.Channel `yaml:"channels"`
}

func main() {
	// Read sensitive values from environment variables
	discordToken, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok || strings.TrimSpace(discordToken) == "" {
		panic("DISCORD_TOKEN environment variable is required but not set")
	}

	once, ok := os.LookupEnv("ONCE")
	if !ok {
		once = "false"
	}
	discordToken = strings.TrimSpace(discordToken)

	d := &discord.Discord{
		Token: discordToken,
	}

	// Carrega a configuração
	config, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("Could not read config.yaml: %w", err))
	}

	var cfg Config
	err = yaml.Unmarshal(config, &cfg)
	if err != nil {
		panic(fmt.Errorf("Could not parse config.yaml: %w", err))
	}

	d.Channels = cfg.Channels

	err = d.Start()
	defer d.Session.Close()
	if err != nil {
		panic(fmt.Errorf("Could not instantiate Discord: %f", err))
	}

	// Registra os comandos do bot
	err = d.RegisterCommands()
	if err != nil {
		panic(fmt.Errorf("Could not register commands: %w", err))
	}

	scraper := &scraper.Scraper{Url: "https://xivpf.com/"}

	fmt.Printf("Starting findingway...\n")
	for {
		totalWait := 40 * time.Second
		fmt.Printf("Scraping source...\n")

		// Faz o scraping uma única vez para todas as raids
		listings, err := scraper.Scrape()
		if err != nil {
			fmt.Printf("Scraper error: %f\n", err)
			continue
		}
		fmt.Printf("Got %v listings.\n", len(listings.Listings))
		fmt.Printf("Sending to %v channels...\n", len(d.Channels))

		// Para cada canal, limpa e envia os dados específicos daquela raid
		for _, c := range d.Channels {
			fmt.Printf("Processing Discord for %v (%v)...\n", c.Name, c.Duty)

			// Limpa o canal apenas antes de enviar os dados específicos
			fmt.Printf("Cleaning Discord for %v (%v)...\n", c.Name, c.Duty)
			err = d.CleanChannel(c.ID)
			if err != nil {
				fmt.Printf("Discord error cleaning channel: %f\n", err)
				continue // Pula para o próximo canal se houver erro na limpeza
			}

			// Envia os dados para todos os data centers deste canal
			for _, dataCentre := range c.DataCentres {
				fmt.Printf("Updating Discord for %v (%v) - %v...\n", c.Name, c.Duty, dataCentre)
				err = d.PostListings(c.ID, listings, c.Duty, dataCentre)
				if err != nil {
					fmt.Printf("Discord error updating messages: %f\n", err)
				}
			}
		}

		if once != "false" {
			os.Exit(0)
		}
		fmt.Printf("Sleeping for %v...\n", totalWait)
		time.Sleep(totalWait)
	}
}
