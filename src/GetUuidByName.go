package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// MinecraftPlayer contains username and uuid of a Minecraft player.
type MinecaftPlayerID struct {
	Name string `json:"name"`
	Uuid string `json:"id"`
}

// getPlayerByName() requests the username and uuid of a Minecraft player by sending its name to the MojangAPI.
// It returns a MinecraftPlayer struct.
func getUuidByName(playerName string) (string, error) {
	resp, err := http.Get("https://api.mojang.com/users/profiles/minecraft/" + playerName)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var playerID map[string]string
	if err := json.Unmarshal(body, &playerID); err != nil {
		return "", err
	}

	return playerID["id"], nil
}

// getPlayerByNameAsync() is the async version of getPlayerByName().
// Here you need to pass a channel and optionally a bool, that determines wether the channel should be closed (default) after the function finished.
// Call this function as a Go-Routine and receive the result over your provided channel.
func getUuidByNameAsync(playerName string, channel chan string, closeChan ...bool) {

	if result, err := getUuidByName(playerName); err != nil {
		log.Printf("UUID: %v raised exception: %v", playerName, err)

	} else {
		channel <- result
	}

	if (len(closeChan) > 0 && closeChan[0]) || len(closeChan) == 0 {
		close(channel)
	}
}
