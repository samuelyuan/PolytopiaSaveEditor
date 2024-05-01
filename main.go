package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/samuelyuan/PolytopiaSaveEditor/fileio"
)

func main() {
	inputPtr := flag.String("input", "", "Input filename")
	outputPtr := flag.String("output", "", "Output filename")
	modePtr := flag.String("mode", "decompress", "Output mode")
	xPtr := flag.Int("x", -1, "x")
	yPtr := flag.Int("y", -1, "y")
	// color
	rPtr := flag.Int("r", -1, "Color r")
	gPtr := flag.Int("g", -1, "Color g")
	bPtr := flag.Int("b", -1, "Color b")
	cityNamePtr := flag.String("cityname", "", "City name")
	oldValuePtr := flag.String("oldvalue", "", "Old value")
	newValuePtr := flag.String("value", "", "New value")

	flag.Parse()

	inputFilename := *inputPtr
	outputFilename := *outputPtr
	mode := *modePtr
	fmt.Println("Input filename: ", inputFilename)
	fmt.Println("Output filename: ", outputFilename)
	fmt.Println("Mode:", mode)

	if mode == "decompress" {
		fileio.DecompressFile(inputFilename)
		return
	} else if mode == "compress" {
		fileio.CompressFile(inputFilename)
		return
	}

	saveOutput, err := fileio.ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	if mode == "modify-unit-tribe" || mode == "modify-unit-type" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		if mode == "modify-unit-tribe" {
			fileio.ModifyUnitTribe(inputFilename, targetX, targetY, updatedValue)
		} else if mode == "modify-unit-type" {
			offset := fileio.GetUnitLocationFileOffset(targetX, targetY)
			fileio.WriteUint16AtFileOffset(inputFilename, offset+5, updatedValue)
		}
		fmt.Println(fmt.Sprintf("Target is at (%v, %v), command: %v, updated value: %v", targetX, targetY, mode, updatedValue))
	} else if mode == "reveal-tile" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.RevealTileForTribe(inputFilename, targetX, targetY, updatedValue)
		fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, updatedValue))
	} else if mode == "modify-tile-capital" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		if updatedValue >= 255 {
			log.Fatal("Value must be less than 255")
		}
		updatedTile := saveOutput.TileData[targetY][targetX]
		updatedTile.Capital = updatedValue
		fileio.WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
		fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have capital %v", targetX, targetY, updatedValue))
	} else if mode == "modify-tile-terrain" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.ModifyTileTerrain(inputFilename, targetX, targetY, updatedValue)
		fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have terrain %v", targetX, targetY, updatedValue))
	} else if mode == "modify-tile-owner" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		updatedTile := saveOutput.TileData[targetY][targetX]
		updatedTile.Owner = updatedValue
		fileio.WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
		fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have owner %v", targetX, targetY, updatedValue))
	} else if mode == "modify-tile-road" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		if updatedValue != 0 && updatedValue != 1 {
			log.Fatal("New value must be 0 or 1")
		}

		updatedTile := saveOutput.TileData[targetY][targetX]
		if updatedValue == 1 {
			updatedTile.HasRoad = true
		} else {
			updatedTile.HasRoad = false
		}
		fileio.WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
		fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have road %v", targetX, targetY, updatedValue))
	} else if mode == "add-city" {
		targetX := *xPtr
		targetY := *yPtr
		cityName := *cityNamePtr
		tribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.AddCityToTile(inputFilename, targetX, targetY, cityName, tribe)
		fmt.Println(fmt.Sprintf("Created city %v at (%v, %v) for player %v", cityName, targetX, targetY, tribe))
	} else if mode == "add-player" {
		fileio.AddPlayer(inputFilename)
		fmt.Println("Added new player to game")
	} else if mode == "reset-tile" {
		targetX := *xPtr
		targetY := *yPtr

		fileio.ResetTile(inputFilename, targetX, targetY)
		fmt.Println(fmt.Sprintf("Reset tile (%v, %v)", targetX, targetY))
	} else if mode == "list-cities" {
		for tribe, cities := range saveOutput.TribeCityMap {
			fmt.Printf("Tribe %v has %v cities:\n", tribe, len(cities))
			for i := 0; i < len(cities); i++ {
				fmt.Printf("City %v: %+v\n", i, cities[i])
			}
		}
	} else if mode == "list-units" {
		for tribe, units := range saveOutput.TribeUnitMap {
			fmt.Printf("Tribe %v has %v units:\n", tribe, len(units))
			for i := 0; i < len(units); i++ {
				fmt.Printf("Unit %v: %+v\n", i, units[i])
			}
		}
	} else if mode == "list-players" {
		for i := 0; i < len(saveOutput.PlayerData); i++ {
			playerData := saveOutput.PlayerData[i]
			fmt.Printf("Player id: %v, name: %v, tribe: %v, override color: %v\n",
				playerData.Id, playerData.Name, playerData.Tribe, playerData.OverrideColor)
		}
	} else if mode == "modify-player-color" {
		playerId, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		colorR := *rPtr
		colorG := *gPtr
		colorB := *bPtr

		fileio.ModifyPlayerColor(inputFilename, playerId, colorR, colorG, colorB)
		fmt.Println(fmt.Sprintf("Set player %v color to RGB(%v, %v, %v)", playerId, colorR, colorG, colorB))
	} else if mode == "convert-tribe" {
		oldValue, err := strconv.Atoi(*oldValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		tribeUnits, ok := saveOutput.TribeUnitMap[oldValue]
		if !ok {
			log.Fatal(fmt.Sprintf("Tribe %v doesn't exist", oldValue))
		}

		for i := 0; i < len(tribeUnits); i++ {
			fileio.ModifyUnitTribe(inputFilename, tribeUnits[i].X, tribeUnits[i].Y, updatedValue)
		}
		fmt.Println(fmt.Sprintf("Changed all units under tribe %v to tribe %v. Total of %v units converted.", oldValue, updatedValue, len(tribeUnits)))
	} else if mode == "convert-all-units" {
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		for tribe, tribeUnits := range saveOutput.TribeUnitMap {
			if tribe == updatedValue {
				continue
			}

			fmt.Println(fmt.Sprintf("Converting all units from tribe %v to tribe %v. Total of %v units converted.", tribe, updatedValue, len(tribeUnits)))
			for i := 0; i < len(tribeUnits); i++ {
				fileio.ModifyUnitTribe(inputFilename, tribeUnits[i].X, tribeUnits[i].Y, updatedValue)
			}
		}
		fmt.Println(fmt.Sprintf("Changed all units to be under tribe %v", updatedValue))
	} else if mode == "reveal-all-tiles" {
		newTribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.RevealAllTiles(inputFilename, newTribe)
	} else if mode == "expand-map-rows" {
		newRowDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.ExpandRows(inputFilename, newRowDimensions)
	} else if mode == "expand-map-cols" {
		newRowDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.ExpandColumns(inputFilename, newRowDimensions)
	} else if mode == "expand-map" {
		newSquareSizeDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fileio.ExpandTiles(inputFilename, newSquareSizeDimensions)
	} else if mode == "reset-game" {
		currentMapHeight := len(saveOutput.TileData)
		currentMapWidth := len(saveOutput.TileData[0])

		initialMapHeight := len(saveOutput.InitialTileData)
		initialMapWidth := len(saveOutput.InitialTileData[0])
		fileio.WriteMapToFile(inputFilename, saveOutput.InitialTileData, currentMapWidth-1, currentMapHeight-1)
		fileio.ModifyMapDimensions(inputFilename, initialMapWidth, initialMapHeight)
	} else if mode == "copy-map" {
		// tests byte conversion of map data
		// make sure new data is equal to old data
		mapHeight := len(saveOutput.TileData)
		mapWidth := len(saveOutput.TileData[0])
		fileio.WriteMapToFile(inputFilename, saveOutput.TileData, mapWidth-1, mapHeight-1)
	} else {
		log.Fatal("Invalid mode:", mode)
	}
}
