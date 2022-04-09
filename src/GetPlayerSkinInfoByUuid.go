package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type PlayerSkinInfo struct {
	LinkCapeImage  string
	LinkSkinImage  string
	ModelTypeSteve bool
}

type PlayerProfileJson struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	TextureInfo []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"properties"`
}

type PlayerSkinJSON struct {
	Timestamp uint64 `json:"timestamp"`
	ProfileId string `json:"profileId"`
	Textures  struct {
		SKIN struct {
			Url      string `json:"url"`
			metadata interface{}
		} `json:"SKIN"`
		CAPE struct {
			Url string `json:"url"`
		} `json:"CAPE"`
	} `json:"textures"`
}

// getPlayerSkinInfoByUuid() returns a struct, containing the Player model type and links to the 2D skin Image/Cape of a player.
// The actual skin information is Base64 encrypted and part of the profile request. Hence we need two json conversions.
func getPlayerSkinInfoByUuid(playerUuid string) (*PlayerSkinInfo, error) {
	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/profile/" + playerUuid)
	if err != nil {
		return &PlayerSkinInfo{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &PlayerSkinInfo{}, err
	}

	var profile PlayerProfileJson
	err = json.Unmarshal(body, &profile)
	if err != nil {
		return &PlayerSkinInfo{}, err
	}

	// Extract and translate the Base64 encrypted texture information from the profile.
	textureJsonString, err := base64.StdEncoding.DecodeString(profile.TextureInfo[0].Value)
	if err != nil {
		return &PlayerSkinInfo{}, err
	}

	var textureJson PlayerSkinJSON
	if err := json.Unmarshal(textureJsonString, &textureJson); err != nil {
		return &PlayerSkinInfo{}, err
	}

	return &PlayerSkinInfo{
		LinkSkinImage:  textureJson.Textures.SKIN.Url,
		LinkCapeImage:  textureJson.Textures.CAPE.Url,
		ModelTypeSteve: textureJson.Textures.SKIN.metadata == nil,
	}, nil
}

// getPlayerSkinInfoByUuidAsync() is the async version of getPlayerSkinInfoByUuid().
// Here you need to pass a channel and optionally a bool, that determines wether the channel should be closed (default) after the function finished.
// Call this function as a Go-Routine and receive the result over your provided channel.
func getPlayerSkinInfoByUuidAsync(playerUuid string, channel chan *PlayerSkinInfo, closeChan ...bool) {

	if result, err := getPlayerSkinInfoByUuid(playerUuid); err != nil {
		log.Printf("UUID: %v raised exception: %v", playerUuid, err)

	} else {
		channel <- result
	}

	if (len(closeChan) > 0 && closeChan[0]) || len(closeChan) == 0 {
		close(channel)
	}
}
