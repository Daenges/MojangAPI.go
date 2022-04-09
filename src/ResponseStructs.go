package main

import "time"

// Information and direct links to a players skin.
type SkinInfo struct {
	Skin      string
	SkinModel string
	Cape      string
}

// All available infomation about a Minecraft Player
type MinecraftPlayer struct {
	NameHistory map[string]time.Time
	Uuid        string
	Skin        SkinInfo
}
