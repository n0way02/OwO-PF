package discord

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

// ID do usuário que pode usar o comando addraid
const ALLOWED_USER_ID = "153106158775697410"

// Estrutura para armazenar as raids disponíveis
type RaidConfig struct {
	Raids []*Raid `yaml:"raids"`
}

type Raid struct {
	Name  string `yaml:"name"`  // Nome completo da raid
	Value string `yaml:"value"` // Valor curto para referência
}

func (d *Discord) RegisterCommands() error {
	// Carrega as raids disponíveis
	raids, err := d.loadRaids()
	if err != nil {
		return fmt.Errorf("erro ao carregar raids: %w", err)
	}

	// Cria as choices para o comando setpfchannel
	var raidChoices []*discordgo.ApplicationCommandOptionChoice
	for _, raid := range raids.Raids {
		raidChoices = append(raidChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  raid.Name,
			Value: raid.Value,
		})
	}

	// Comando para configurar canal
	setChannelCommand := &discordgo.ApplicationCommand{
		Name:                     "setpfchannel",
		Description:              "Define um canal para receber PFs de uma raid específica",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0], // Requer permissão de administrador
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "canal",
				Description: "O canal onde os PFs serão enviados",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "raid",
				Description: "Selecione a raid",
				Required:    true,
				Choices:     raidChoices,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "datacenter1",
				Description: "Primeiro Data Center",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Aether", Value: "Aether"},
					{Name: "Primal", Value: "Primal"},
					{Name: "Crystal", Value: "Crystal"},
					{Name: "Dynamis", Value: "Dynamis"},
					{Name: "Chaos", Value: "Chaos"},
					{Name: "Light", Value: "Light"},
					{Name: "Materia", Value: "Materia"},
					{Name: "Meteor", Value: "Meteor"},
					{Name: "Mana", Value: "Mana"},
					{Name: "Gaia", Value: "Gaia"},
					{Name: "Elemental", Value: "Elemental"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "datacenter2",
				Description: "Segundo Data Center (opcional)",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Aether", Value: "Aether"},
					{Name: "Primal", Value: "Primal"},
					{Name: "Crystal", Value: "Crystal"},
					{Name: "Dynamis", Value: "Dynamis"},
					{Name: "Chaos", Value: "Chaos"},
					{Name: "Light", Value: "Light"},
					{Name: "Materia", Value: "Materia"},
					{Name: "Meteor", Value: "Meteor"},
					{Name: "Mana", Value: "Mana"},
					{Name: "Gaia", Value: "Gaia"},
					{Name: "Elemental", Value: "Elemental"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "datacenter3",
				Description: "Terceiro Data Center (opcional)",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Aether", Value: "Aether"},
					{Name: "Primal", Value: "Primal"},
					{Name: "Crystal", Value: "Crystal"},
					{Name: "Dynamis", Value: "Dynamis"},
					{Name: "Chaos", Value: "Chaos"},
					{Name: "Light", Value: "Light"},
					{Name: "Materia", Value: "Materia"},
					{Name: "Meteor", Value: "Meteor"},
					{Name: "Mana", Value: "Mana"},
					{Name: "Gaia", Value: "Gaia"},
					{Name: "Elemental", Value: "Elemental"},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "datacenter4",
				Description: "Quarto Data Center (opcional)",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Aether", Value: "Aether"},
					{Name: "Primal", Value: "Primal"},
					{Name: "Crystal", Value: "Crystal"},
					{Name: "Dynamis", Value: "Dynamis"},
					{Name: "Chaos", Value: "Chaos"},
					{Name: "Light", Value: "Light"},
					{Name: "Materia", Value: "Materia"},
					{Name: "Meteor", Value: "Meteor"},
					{Name: "Mana", Value: "Mana"},
					{Name: "Gaia", Value: "Gaia"},
					{Name: "Elemental", Value: "Elemental"},
				},
			},
		},
	}

	// Comando para adicionar nova raid
	addRaidCommand := &discordgo.ApplicationCommand{
		Name:                     "addraid",
		Description:              "Adiciona uma nova raid à lista de raids disponíveis",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0], // Requer permissão de administrador
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "nome",
				Description: "Nome completo da raid (ex: The Omega Protocol (Ultimate))",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "valor",
				Description: "Valor curto para referência (ex: top)",
				Required:    true,
			},
		},
	}

	// Comando para remover canal
	removeChannelCommand := &discordgo.ApplicationCommand{
		Name:                     "removepfchannel",
		Description:              "Remove um canal da lista de canais que recebem PFs",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0], // Requer permissão de administrador
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "canal",
				Description: "O canal que você deseja remover da lista",
				Required:    true,
			},
		},
	}

	// Comando para remover raid
	removeRaidCommand := &discordgo.ApplicationCommand{
		Name:                     "removeraid",
		Description:              "Remove uma raid da lista de raids disponíveis",
		DefaultMemberPermissions: &[]int64{discordgo.PermissionAdministrator}[0], // Requer permissão de administrador
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "raid",
				Description: "Selecione a raid que deseja remover",
				Required:    true,
				Choices:     raidChoices,
			},
		},
	}

	// Registra os comandos
	commands := []*discordgo.ApplicationCommand{setChannelCommand, addRaidCommand, removeChannelCommand, removeRaidCommand}
	for _, cmd := range commands {
		_, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("erro ao registrar comando %s: %w", cmd.Name, err)
		}
	}

	// Adiciona o handler para os comandos
	d.Session.AddHandler(d.handleCommands)

	return nil
}

func (d *Discord) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "setpfchannel":
		d.handleSetPFChannel(s, i)
	case "addraid":
		d.handleAddRaid(s, i)
	case "removepfchannel":
		d.handleRemovePFChannel(s, i)
	case "removeraid":
		d.handleRemoveRaid(s, i)
	}
}

func (d *Discord) handleSetPFChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	channel := options[0].ChannelValue(s)
	raidValue := options[1].StringValue()

	// Coleta todos os data centers selecionados
	dataCenters := []string{}
	for i := 2; i < len(options); i++ {
		dc := options[i].StringValue()
		if dc != "" {
			dataCenters = append(dataCenters, dc)
		}
	}

	// Remove duplicatas dos data centers
	uniqueDCs := make(map[string]bool)
	uniqueDataCenters := []string{}
	for _, dc := range dataCenters {
		if !uniqueDCs[dc] {
			uniqueDCs[dc] = true
			uniqueDataCenters = append(uniqueDataCenters, dc)
		}
	}

	// Mapa de valores para nomes completos das raids
	raidNames := map[string]string{
		"m1s":     "AAC Cruiserweight M1 (Savage)",
		"m2s":     "AAC Cruiserweight M2 (Savage)",
		"m3s":     "AAC Cruiserweight M3 (Savage)",
		"m4s":     "AAC Cruiserweight M4 (Savage)",
		"chaotic": "The Cloud of Darkness (Chaotic)",
		"unreal":  "Hells' Kier (Unreal)",
		"ucob":    "The Unending Coil of Bahamut (Ultimate)",
		"uwu":     "The Weapon's Refrain (Ultimate)",
		"tea":     "The Epic of Alexander (Ultimate)",
		"dsr":     "Dragonsong's Reprise (Ultimate)",
		"top":     "The Omega Protocol (Ultimate)",
		"fru":     "Futures Rewritten (Ultimate)",
	}

	// Cria um novo canal
	newChannel := &Channel{
		ID:          channel.ID,
		Name:        channel.Name,         // Usa o nome do canal do Discord
		Duty:        raidNames[raidValue], // Usa o nome completo da raid
		DataCentres: uniqueDataCenters,    // Usa os data centers selecionados
	}

	// Adiciona o novo canal à lista
	d.Channels = append(d.Channels, newChannel)

	// Salva a configuração atualizada
	err := d.saveConfig()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao salvar configuração: %v", err),
			},
		})
		return
	}

	// Responde ao usuário
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Canal %s configurado com sucesso para receber PFs de %s nos Data Centers: %s!",
				channel.Mention(),
				raidNames[raidValue],
				strings.Join(uniqueDataCenters, ", ")),
		},
	})
}

func (d *Discord) handleAddRaid(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Verifica se o usuário é o autorizado
	if i.Member.User.ID != ALLOWED_USER_ID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Você não tem permissão para usar este comando!",
				Flags:   discordgo.MessageFlagsEphemeral, // A mensagem só será visível para quem executou o comando
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options
	raidName := options[0].StringValue()
	raidValue := options[1].StringValue()

	// Carrega as raids existentes
	raids, err := d.loadRaids()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao carregar raids: %v", err),
			},
		})
		return
	}

	// Verifica se o valor já existe
	for _, raid := range raids.Raids {
		if raid.Value == raidValue {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Já existe uma raid com o valor '%s'", raidValue),
				},
			})
			return
		}
	}

	// Adiciona a nova raid
	raids.Raids = append(raids.Raids, &Raid{
		Name:  raidName,
		Value: raidValue,
	})

	// Salva as raids atualizadas
	err = d.saveRaids(raids)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao salvar raids: %v", err),
			},
		})
		return
	}

	// Responde ao usuário imediatamente
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Raid '%s' adicionada com sucesso! Atualizando comandos...", raidName),
		},
	})

	// Atualiza os comandos em background
	go func() {
		err := d.RegisterCommands()
		if err != nil {
			// Envia uma mensagem de erro no canal
			s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Erro ao atualizar os comandos: %v", err))
		}
	}()
}

func (d *Discord) handleRemovePFChannel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	channel := options[0].ChannelValue(s)

	// Procura o canal na lista
	found := false
	var newChannels []*Channel
	var removedChannel *Channel

	for _, c := range d.Channels {
		if c.ID == channel.ID {
			found = true
			removedChannel = c
		} else {
			newChannels = append(newChannels, c)
		}
	}

	if !found {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("O canal %s não está na lista de canais configurados!", channel.Mention()),
			},
		})
		return
	}

	// Atualiza a lista de canais
	d.Channels = newChannels

	// Salva a configuração atualizada
	err := d.saveConfig()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao salvar configuração: %v", err),
			},
		})
		return
	}

	// Responde ao usuário
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Canal %s removido com sucesso! Ele não receberá mais PFs de %s.",
				channel.Mention(),
				removedChannel.Duty),
		},
	})
}

func (d *Discord) handleRemoveRaid(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Verifica se o usuário é o autorizado
	if i.Member.User.ID != ALLOWED_USER_ID {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Você não tem permissão para usar este comando!",
				Flags:   discordgo.MessageFlagsEphemeral, // A mensagem só será visível para quem executou o comando
			},
		})
		return
	}

	options := i.ApplicationCommandData().Options
	raidValue := options[0].StringValue()

	// Carrega as raids existentes
	raids, err := d.loadRaids()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao carregar raids: %v", err),
			},
		})
		return
	}

	// Procura a raid para remover
	found := false
	var newRaids []*Raid
	var removedRaid *Raid

	for _, raid := range raids.Raids {
		if raid.Value == raidValue {
			found = true
			removedRaid = raid
		} else {
			newRaids = append(newRaids, raid)
		}
	}

	if !found {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Não foi encontrada uma raid com o valor '%s'", raidValue),
			},
		})
		return
	}

	// Atualiza a lista de raids
	raids.Raids = newRaids

	// Salva as raids atualizadas
	err = d.saveRaids(raids)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Erro ao salvar raids: %v", err),
			},
		})
		return
	}

	// Responde ao usuário imediatamente
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Raid '%s' removida com sucesso! Atualizando comandos...", removedRaid.Name),
		},
	})

	// Atualiza os comandos em background
	go func() {
		err := d.RegisterCommands()
		if err != nil {
			// Envia uma mensagem de erro no canal
			s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Erro ao atualizar os comandos: %v", err))
		}
	}()
}

func (d *Discord) loadRaids() (*RaidConfig, error) {
	data, err := os.ReadFile("raids.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			// Se o arquivo não existe, retorna uma configuração vazia
			return &RaidConfig{Raids: []*Raid{}}, nil
		}
		return nil, fmt.Errorf("erro ao ler arquivo de raids: %w", err)
	}

	var config RaidConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo de raids: %w", err)
	}

	return &config, nil
}

func (d *Discord) saveRaids(config *RaidConfig) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("erro ao converter raids para YAML: %w", err)
	}

	err = os.WriteFile("raids.yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo de raids: %w", err)
	}

	return nil
}

func (d *Discord) saveConfig() error {
	// Cria uma estrutura temporária apenas com os dados que queremos salvar
	type Config struct {
		Channels []*Channel `yaml:"channels"`
	}

	config := Config{
		Channels: d.Channels,
	}

	// Converte a estrutura para YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("erro ao converter para YAML: %w", err)
	}

	// Salva no arquivo
	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %w", err)
	}

	return nil
}
