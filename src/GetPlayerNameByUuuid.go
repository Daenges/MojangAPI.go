package mojangapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetPlayerNameByUuid(playerUuid string) (string, error) {
	resp, err := http.Get("https://api.mojang.com/user/profiles/" + playerUuid + "/names")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Create an interface, since the first list entry is different.
	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Just return the current name and discard the rest.
	if len(result) > 0 {
		return fmt.Sprintf("%v", result[len(result)-1]["name"]), nil
	} else {
		panic("Request failed.")
	}
}

// getPlayerNameByUuidAsync() is the async version of getPlayerNameByUuid().
// Here you need to pass a channel and optionally a bool, that determines wether the channel should be closed (default) after the function finished.
// Call this function as a Go-Routine and receive the result over your provided channel.
func GetPlayerNameByUuidAsync(playerUuid string, channel chan *string, closeChan ...bool) {

	if result, err := GetPlayerNameByUuid(playerUuid); err != nil {
		log.Printf("UUID: %v raised exception: %v", playerUuid, err)

	} else {
		channel <- &result
	}

	if (len(closeChan) > 0 && closeChan[0]) || len(closeChan) == 0 {
		close(channel)
	}
}
