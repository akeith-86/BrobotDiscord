package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var (
	Token     string
	BotPrefix string
	GptToken  string
	config    *configStruct
)

type configStruct struct {
	Token     string `json : "Token"`
	BotPrefix string `json : "BotPrefix"`
	GptToken  string `json : "GptToken"`
}

func ReadConfig() error {
	fmt.Println("Reading config file")
	file, err := os.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix
	GptToken = config.GptToken

	return nil

}

var BotId string
var broBot *discordgo.Session

func Start() {
	broBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := broBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	broBot.AddHandler(messageHandler)

	err = broBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running !")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	var resp string
	cmdCheck := strings.Split(m.Content, " ")

	if m.Author.ID == BotId {
		return
	}

	fmt.Println(m.Author.Username + ": " + m.Content)

	if m.Content == BotPrefix+"help" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "This is Brobot, a helpful bot who's also a bro. Here are some commands that I have! \n !help: Shows this output \n !ping: Tests if I'm alive! \n !bros: Assembles the bro team. \n <:brobot:1065746958481895474>: <:brobot:1065746958481895474> \n !catjam: Summon a cool cat to jam with \n !gpt <question>: Ask me a question, I might have an answer!")
	}

	if m.Content == BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if m.Content == BotPrefix+"bros" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "BROS ASSEMBLE! @everyone")
	}

	if m.Content == "<:brobot:1065746958481895474>" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "<:brobot:1065746958481895474>")
	}

	if m.Content == BotPrefix+"catjam" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "<a:catjam:1065746950839869551>")
	}

	if cmdCheck[0] == BotPrefix+"gpt" {

		userPrompt := strings.Join(strings.Split(m.Content, " ")[1:], " ")

		resp = gpt(userPrompt)

		_, _ = s.ChannelMessageSend(m.ChannelID, resp)
	}
}

func gpt(userPrompt string) string {
	gptClient := gogpt.NewClient(GptToken)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 1000,
		Prompt:    fmt.Sprint(userPrompt),
	}

	resp, err := gptClient.CreateCompletion(ctx, req)
	if err != nil {
		fmt.Println(err.Error())
		return "An error has occured!"
	}

	return resp.Choices[0].Text
}

func main() {
	err := ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	<-make(chan struct{})
	return
}
