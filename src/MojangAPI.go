package main

import "fmt"

func main() {
	// ch := make(chan *MinecaftPlayerID)

	// go getPlayerByNameAsync("Daenges", ch, false)

	// // Mache Irgendwas

	// for output := range ch {
	// 	fmt.Println(output.Uuid)
	// }

	test, err := getUuidByName("Cyb3rko")
	if err != nil {
		fmt.Print("FEHLER")
	}

	ergebnis, _ := getPlayerSkinInfoByUuid(test)

	fmt.Print(ergebnis)

	fmt.Println("Ende")
}
