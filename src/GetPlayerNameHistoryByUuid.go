package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type NameHistory struct {
	FirstName      string
	FollowingNames []NameChangeEntry
}

type NameChangeEntry struct {
	Name        string
	ChangedToAt time.Time
}

// getPlayerNameHistoryByUuid() returns a players name history.
func getPlayerNameHistoryByUuid(playerUuid string) (*NameHistory, error) {
	resp, err := http.Get("https://api.mojang.com/user/profiles/" + playerUuid + "/names")
	if err != nil {
		return &NameHistory{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &NameHistory{}, err
	}

	// Create an interface, since the first list entry is different.
	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return &NameHistory{}, err
	}

	// Create NameHistory object, that is returned as a pointer.
	var nh NameHistory
	// The first entry (Name) does not contain a date.
	nh.FirstName = fmt.Sprintf("%v", result[0]["name"])
	// Add the following names with their corresponding names into the list.
	// Time is in Milliseconds, but time.unix works with Nanoseconds -> Multiplication
	for i := 1; i < len(result); i++ {
		nh.FollowingNames = append(nh.FollowingNames,
			NameChangeEntry{fmt.Sprintf("%v", result[i]["name"]),
				time.Unix(0, int64(result[i]["changedToAt"].(float64))*int64(time.Millisecond))})
	}

	return &nh, nil
}

// getPlayerNameHistoryByUuidAsync() is the async version of getPlayerNameHistoryByUuid().
// Here you need to pass a channel and optionally a bool, that determines wether the channel should be closed (default) after the function finished.
// Call this function as a Go-Routine and receive the result over your provided channel.
func getPlayerNameHistoryByUuidAsync(playerUuid string, channel chan *NameHistory, closeChan ...bool) {

	if result, err := getPlayerNameHistoryByUuid(playerUuid); err != nil {
		log.Printf("UUID: %v raised exception: %v", playerUuid, err)

	} else {
		channel <- result
	}

	if (len(closeChan) > 0 && closeChan[0]) || len(closeChan) == 0 {
		close(channel)
	}
}
