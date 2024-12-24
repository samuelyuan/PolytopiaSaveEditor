package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	polytopiamapmodel "github.com/samuelyuan/polytopiamapmodelgo"
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
	playerIdPtr := flag.String("playerid", "", "Player Id")

	flag.Parse()

	inputFilename := *inputPtr
	outputFilename := *outputPtr
	mode := *modePtr
	fmt.Println("Input filename: ", inputFilename)
	fmt.Println("Output filename: ", outputFilename)
	fmt.Println("Mode:", mode)

	if mode == "decompress" {
		polytopiamapmodel.DecompressFile(inputFilename)
		return
	} else if mode == "compress" {
		polytopiamapmodel.CompressFile(inputFilename, outputFilename)
		return
	}

	saveOutput, err := polytopiamapmodel.ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	fileInfo := polytopiamapmodel.FileInfo{
		InputFilename: inputFilename,
		GameVersion:   int(saveOutput.GameVersion),
	}

	if mode == "modify-unit-tribe" || mode == "modify-unit-type" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		if mode == "modify-unit-tribe" {
			polytopiamapmodel.ModifyUnitTribe(fileInfo, targetX, targetY, updatedValue)
		} else if mode == "modify-unit-type" {
			polytopiamapmodel.ModifyUnitType(fileInfo, targetX, targetY, updatedValue)
		}
		fmt.Println(fmt.Sprintf("Target is at (%v, %v), command: %v, updated value: %v", targetX, targetY, mode, updatedValue))
	} else if mode == "reveal-tile" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.RevealTileForTribe(fileInfo, targetX, targetY, updatedValue)
		fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, updatedValue))
	} else if mode == "set-new-tile-capital" {
		targetX := *xPtr
		targetY := *yPtr
		updatedTribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		newCityName := *cityNamePtr

		polytopiamapmodel.SetTileCapital(fileInfo, targetX, targetY, newCityName, updatedTribe)
	} else if mode == "modify-tile-terrain" {
		targetX := *xPtr
		targetY := *yPtr
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.ModifyTileTerrain(fileInfo, targetX, targetY, updatedValue)
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
		polytopiamapmodel.WriteTileToFile(fileInfo, updatedTile, targetX, targetY)
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
		polytopiamapmodel.WriteTileToFile(fileInfo, updatedTile, targetX, targetY)
		fmt.Println(fmt.Sprintf("Modified tile (%v, %v) to have road %v", targetX, targetY, updatedValue))
	} else if mode == "add-city" {
		targetX := *xPtr
		targetY := *yPtr
		cityName := *cityNamePtr
		tribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.AddCityToTile(fileInfo, targetX, targetY, cityName, tribe)
		fmt.Println(fmt.Sprintf("Created city %v at (%v, %v) for player %v", cityName, targetX, targetY, tribe))
	} else if mode == "add-player" {
		polytopiamapmodel.AddPlayer(inputFilename)
		fmt.Println("Added new player to game")
	} else if mode == "reset-tile" {
		targetX := *xPtr
		targetY := *yPtr

		polytopiamapmodel.ResetTile(fileInfo, targetX, targetY)
		fmt.Println(fmt.Sprintf("Reset tile (%v, %v)", targetX, targetY))
	} else if mode == "list-cities" {
		for tribe, cities := range saveOutput.TribeCityMap {
			fmt.Printf("Tribe %v has %v cities:\n", tribe, len(cities))
			for i := 0; i < len(cities); i++ {
				fmt.Printf("City %v: %+v\n", i, cities[i])
			}
		}
	} else if mode == "list-units" {
		tribeUnitMap := polytopiamapmodel.BuildTribeUnitMap(saveOutput)
		for tribe, units := range tribeUnitMap {
			fmt.Printf("Tribe %v has %v units:\n", tribe, len(units))
			for i := 0; i < len(units); i++ {
				fmt.Printf("Unit %v: %+v\n", i, units[i])
			}
		}
	} else if mode == "list-players" {
		for i := 0; i < len(saveOutput.PlayerData); i++ {
			playerData := saveOutput.PlayerData[i]
			fmt.Printf("Player id: %v, name: %v, tribe: %v, override color: %v\n",
				playerData.PlayerId, playerData.Name, playerData.Tribe, playerData.OverrideColor)
		}
	} else if mode == "modify-player-color" {
		playerId, err := strconv.Atoi(*playerIdPtr)
		if err != nil {
			log.Fatal(err)
		}
		colorR := *rPtr
		colorG := *gPtr
		colorB := *bPtr

		for i := 0; i < len(saveOutput.PlayerData); i++ {
			if saveOutput.PlayerData[i].PlayerId == playerId {
				saveOutput.PlayerData[i].OverrideColor = []int{int(colorB), int(colorG), int(colorR), 0}
				break
			}
		}

		polytopiamapmodel.WritePlayersToFile(inputFilename, saveOutput.PlayerData)
		fmt.Println(fmt.Sprintf("Set player %v color to RGB(%v, %v, %v)", playerId, colorR, colorG, colorB))
	} else if mode == "modify-player-tribe" {
		playerId, err := strconv.Atoi(*playerIdPtr)
		if err != nil {
			log.Fatal(err)
		}
		newTribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(saveOutput.PlayerData); i++ {
			if saveOutput.PlayerData[i].PlayerId == playerId {
				saveOutput.PlayerData[i].Tribe = newTribe
				break
			}
		}

		polytopiamapmodel.WritePlayersToFile(inputFilename, saveOutput.PlayerData)
		fmt.Println(fmt.Sprintf("Set player %v tribe to %v", playerId, newTribe))
	} else if mode == "modify-player-name" {
		playerId, err := strconv.Atoi(*playerIdPtr)
		if err != nil {
			log.Fatal(err)
		}
		newName := *newValuePtr

		for i := 0; i < len(saveOutput.PlayerData); i++ {
			if saveOutput.PlayerData[i].PlayerId == playerId {
				saveOutput.PlayerData[i].Name = newName
				break
			}
		}

		polytopiamapmodel.WritePlayersToFile(inputFilename, saveOutput.PlayerData)
		fmt.Println(fmt.Sprintf("Set player %v newName to %v", playerId, newName))
	} else if mode == "convert-tribe" {
		oldTribe, err := strconv.Atoi(*oldValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		newTribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.ConvertTribe(fileInfo, oldTribe, newTribe)
	} else if mode == "convert-all-units" {
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		newTribe := updatedValue
		totalConverted := 0
		tribeUnitMap := polytopiamapmodel.BuildTribeUnitMap(saveOutput)
		for tribe, tribeUnits := range tribeUnitMap {
			if tribe == updatedValue {
				continue
			}

			fmt.Println(fmt.Sprintf("Converting all units from tribe %v to tribe %v. Total of %v units converted.", tribe, updatedValue, len(tribeUnits)))
			totalConverted += len(tribeUnits)
			for i := 0; i < len(tribeUnits); i++ {
				targetX := tribeUnits[i].X
				targetY := tribeUnits[i].Y

				updatedTile := saveOutput.TileData[targetY][targetX]
				if updatedTile.Unit != nil {
					updatedTile.Unit.Owner = uint8(newTribe)
				}
				if updatedTile.PassengerUnit != nil {
					updatedTile.PassengerUnit.Owner = uint8(newTribe)
				}
				fmt.Println(fmt.Sprintf("Converted unit on (%v, %v) from tribe %v to %v", targetX, targetY, tribe, newTribe))

				saveOutput.TileData[targetY][targetX] = updatedTile
			}
		}

		polytopiamapmodel.WriteMapToFile(fileInfo, saveOutput.TileData)
		fmt.Println(fmt.Sprintf("Changed all units to be under tribe %v. Converted total of %v units.", updatedValue, totalConverted))
	} else if mode == "reveal-all-tiles" {
		newTribe, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.RevealAllTiles(fileInfo, newTribe)
	} else if mode == "expand-map-rows" {
		newRowDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.ExpandRows(fileInfo, newRowDimensions)
	} else if mode == "expand-map-cols" {
		newRowDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.ExpandColumns(fileInfo, newRowDimensions)
	} else if mode == "expand-map" {
		newSquareSizeDimensions, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		polytopiamapmodel.ExpandTiles(fileInfo, newSquareSizeDimensions)
	} else if mode == "swap-players" {
		oldValue, err := strconv.Atoi(*oldValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fmt.Sprintf("Swap player id %v with id %v", oldValue, updatedValue))
		polytopiamapmodel.SwapPlayers(fileInfo, oldValue, updatedValue)
		fmt.Println(fmt.Sprintf("Swapped players %v and %v", oldValue, updatedValue))
	} else if mode == "reset-game" {
		initialMapHeight := len(saveOutput.InitialTileData)
		initialMapWidth := len(saveOutput.InitialTileData[0])
		polytopiamapmodel.WriteMapToFile(fileInfo, saveOutput.InitialTileData)
		for i := 0; i < len(saveOutput.InitialPlayerData); i++ {
			saveOutput.InitialPlayerData[i].AvailableTech = []int{0}
		}
		polytopiamapmodel.WritePlayersToFile(inputFilename, saveOutput.InitialPlayerData)
		polytopiamapmodel.ModifyMapDimensions(inputFilename, initialMapWidth, initialMapHeight)
	} else if mode == "copy-data" {
		// tests byte conversion of map data and player data
		// make sure new data is equal to old data
		polytopiamapmodel.WriteMapToFile(fileInfo, saveOutput.TileData)
		polytopiamapmodel.WritePlayersToFile(inputFilename, saveOutput.PlayerData)
	} else if mode == "import-json" {
		// parameters: -value=importJsonFilename
		importJsonFilename := *newValuePtr
		fmt.Println(fmt.Sprintf("Importing data from %v", importJsonFilename))
		polytopiaJson := polytopiamapmodel.ImportPolytopiaDataFromJson(importJsonFilename)
		polytopiamapmodel.WriteMapToFile(fileInfo, polytopiaJson.TileData)
		polytopiamapmodel.WritePlayersToFile(inputFilename, polytopiaJson.PlayerData)
		polytopiamapmodel.WriteMapHeaderToFile(inputFilename, polytopiaJson.MapHeaderOutput)
		fmt.Println(fmt.Sprintf("Updated file %v with imported json", inputFilename))
	} else if mode == "export-json" {
		// parameters: -value=exportFilename
		exportJsonFilename := *newValuePtr
		polytopiamapmodel.ExportPolytopiaJsonFile(saveOutput, exportJsonFilename)
		fmt.Println(fmt.Sprintf("Exported json from save state %v to %v", inputFilename, exportJsonFilename))
	} else {
		log.Fatal("Invalid mode:", mode)
	}
}
