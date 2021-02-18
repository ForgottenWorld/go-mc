package main

import (
	"bytes"
	"flag"
	"log"

	"github.com/google/uuid"

	"github.com/ForgottenWorld/go-mc/bot"
	"github.com/ForgottenWorld/go-mc/chat"
	_ "github.com/ForgottenWorld/go-mc/data/lang/zh-cn"
	pk "github.com/ForgottenWorld/go-mc/net/packet"
)

var address = flag.String("address", "127.0.0.1", "The server address")
var c *bot.Client

func main() {
	flag.Parse()
	c = bot.NewClient()

	//Login
	err := c.JoinServer(*address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//Register event handlers
	c.Events.GameStart = onGameStart
	c.Events.ChatMsg = onChatMsg
	c.Events.Disconnect = onDisconnect
	c.Events.PluginMessage = onPluginMessage

	//JoinGame
	err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}

func onDeath() error {
	log.Println("Died and Respawned")
	// If we exclude Respawn(...) then the player won't press the "Respawn" button upon death
	return c.Respawn()
}

func onGameStart() error {
	log.Println("Game start")
	return nil //if err isn't nil, HandleGame() will return it.
}

func onChatMsg(c chat.Message, pos byte, uuid uuid.UUID) error {
	log.Println("Chat:", c.ClearString()) // output chat message without any format code (like color or bold)
	return nil
}

func onDisconnect(c chat.Message) error {
	log.Println("Disconnect:", c)
	return nil
}

func onPluginMessage(channel string, data []byte) error {
	switch channel {
	case "minecraft:brand":
		var brand pk.String
		if err := brand.Decode(bytes.NewReader(data)); err != nil {
			return err
		}
		log.Println("Server brand is:", brand)

	default:
		log.Println("PluginMessage", channel, data)
	}
	return nil
}
