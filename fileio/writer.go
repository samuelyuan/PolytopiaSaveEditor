package fileio

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func WriteUint8AtFileOffset(inputFilename string, offset int, value int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if value >= 256 {
		log.Fatal("Value is too large for uint8")
	}
	if _, err := inputFile.WriteAt([]byte{uint8(value)}, int64(offset)); err != nil {
		log.Fatal("Failed to write uint8 to file:", err)
	}
}

func WriteUint16AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 65536 {
		log.Fatal("Value is too large for uint16")
	}
	byteArrUnitType := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArrUnitType, uint16(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

func WriteUint32AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 4294967295 {
		log.Fatal("Value is too large for uint32")
	}
	byteArrUnitType := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArrUnitType, uint32(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

func GetFileRemainingData(inputFile *os.File, offset int) []byte {
	if _, err := inputFile.Seek(int64(offset), 0); err != nil {
		log.Fatal(err)
	}
	remainder, err := ioutil.ReadAll(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	return remainder
}

func ConvertUint32Bytes(value int) []byte {
	byteArr := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArr, uint32(value))
	return byteArr
}

func ConvertUint16Bytes(value int) []byte {
	byteArr := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArr, uint16(value))
	return byteArr
}

func ConvertVarString(value string) []byte {
	byteArr := make([]byte, 0)
	byteArr = append(byteArr, byte(len(value)))
	byteArr = append(byteArr, []byte(value)...)
	return byteArr
}

func ModifyTileTerrain(inputFilename string, targetX int, targetY int, updatedValue int) {
	offset, ok := fileOffsetMap[buildTileStartKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile start key on x: %v, y: %v. Command not run.", targetX, targetY))
	}

	// write terrain
	WriteUint16AtFileOffset(inputFilename, offset+8, updatedValue)

	// write altitude
	altitude := 0
	if updatedValue == 1 { // water altitude is -1
		altitude = -1
	} else if updatedValue == 2 { // ocean altitude is -2
		altitude = -2
	} else if updatedValue == 3 || updatedValue == 5 { // flat tile altitude is 1
		altitude = 1
	} else if updatedValue == 4 { // mountain altitude is 2
		altitude = 2
	}
	WriteUint16AtFileOffset(inputFilename, offset+12, altitude)
}

func ModifyTileOwner(inputFilename string, targetX int, targetY int, updatedValue int) {
	offset, ok := fileOffsetMap[buildTileStartKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile start key on x: %v, y: %v. Command not run.", targetX, targetY))
	}

	// write owner
	WriteUint8AtFileOffset(inputFilename, offset+14, updatedValue)
}

func ModifyTileRoad(inputFilename string, targetX int, targetY int, updatedValue int) {
	offset, ok := fileOffsetMap[buildTileRoadKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile road key on x: %v, y: %v. Command not run.", targetX, targetY))
	}

	if updatedValue != 0 && updatedValue != 1 {
		log.Fatal("New value must be 0 or 1")
	}

	// write road
	WriteUint8AtFileOffset(inputFilename, offset, updatedValue)
}

func ModifyUnitTribe(inputFilename string, targetX int, targetY int, updatedValue int) {
	offset := GetUnitLocationFileOffset(targetX, targetY)
	WriteUint8AtFileOffset(inputFilename, offset+4, updatedValue)

	offsetPreviousUnit, ok := fileOffsetMap[buildPreviousUnitLocationKey(targetX, targetY)]
	if ok {
		WriteUint8AtFileOffset(inputFilename, offsetPreviousUnit+4, updatedValue)
	}
}

func BuildEmptyTile(x int, y int) []byte {
	terrain := 3 // flat
	climate := 1
	altitude := 1
	owner := 0
	capital := 0
	capitalX := -1
	capitalY := -1

	headerBytes := make([]byte, 0)
	headerBytes = append(headerBytes, ConvertUint32Bytes(x)...)
	headerBytes = append(headerBytes, ConvertUint32Bytes(y)...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(terrain)...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(climate)...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(altitude)...) // should be int16
	headerBytes = append(headerBytes, byte(owner))
	headerBytes = append(headerBytes, byte(capital))
	headerBytes = append(headerBytes, ConvertUint32Bytes(capitalX)...) // should be int32
	headerBytes = append(headerBytes, ConvertUint32Bytes(capitalY)...) // should be int32

	remainingTileData := make([]byte, 10)
	for i := 0; i < len(remainingTileData); i++ {
		remainingTileData[i] = 0
	}

	allTileData := append(headerBytes, remainingTileData...)
	return allTileData
}

func ModifyMapDimensions(inputFilename string, width int, height int) {
	minSquareSize := width
	if minSquareSize > height {
		minSquareSize = height
	}
	squareSizeOffset, ok := fileOffsetMap["SquareSizeKey"]
	if !ok {
		log.Fatal("Error: No square size key. Command not run.")
	}
	WriteUint32AtFileOffset(inputFilename, squareSizeOffset, minSquareSize)

	widthOffset, ok := fileOffsetMap["MapWidth"]
	if !ok {
		log.Fatal("Error: No map width key. Command not run.")
	}
	WriteUint16AtFileOffset(inputFilename, widthOffset, width)

	heightOffset, ok := fileOffsetMap["MapHeight"]
	if !ok {
		log.Fatal("Error: No map height key. Command not run.")
	}
	WriteUint16AtFileOffset(inputFilename, heightOffset, height)
}

func BuildTileHeaderTribeCity(targetX int, targetY int, tribe int) []byte {
	// world coordinates, terrain, climate, altitude are the same for all players
	// the difference is the tribe that owns this city

	tileHeaderTribeBytes := make([]byte, 0)
	// new owner
	tileHeaderTribeBytes = append(tileHeaderTribeBytes, byte(tribe))

	// set capital to 0 unless this city is designated as capital city
	tileHeaderTribeBytes = append(tileHeaderTribeBytes, byte(0))

	// write city coordinates
	tileHeaderTribeBytes = append(tileHeaderTribeBytes, ConvertUint32Bytes(targetX)...)
	tileHeaderTribeBytes = append(tileHeaderTribeBytes, ConvertUint32Bytes(targetY)...)

	return tileHeaderTribeBytes
}

func BuildCity(cityName string) []byte {
	cityBytes := make([]byte, 0)
	cityBytes = append(cityBytes, byte(1))

	// improvement is 1 for cities
	improvement := 1
	cityBytes = append(cityBytes, ConvertUint16Bytes(improvement)...)

	cityLevel := 1
	cityBytes = append(cityBytes, ConvertUint16Bytes(cityLevel)...)

	foundedTurn := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(foundedTurn)...)

	currentPopulation := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(currentPopulation)...)

	totalPopulation := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(totalPopulation)...)

	// unknown
	cityBytes = append(cityBytes, ConvertUint16Bytes(1)...)

	baseScore := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(baseScore)...)

	// unknown
	cityBytes = append(cityBytes, ConvertUint16Bytes(1)...)
	// unknown
	cityBytes = append(cityBytes, ConvertUint16Bytes(0)...)

	// new city is not connected to capital
	connectedPlayerCapital := 0
	cityBytes = append(cityBytes, byte(connectedPlayerCapital))

	hasCity := 1
	cityBytes = append(cityBytes, byte(hasCity))

	cityBytes = append(cityBytes, ConvertVarString(cityName)...)

	foundedTribe := 0
	cityBytes = append(cityBytes, byte(foundedTribe))

	cityRewardsSize := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(cityRewardsSize)...)

	rebellionFlag := 0
	cityBytes = append(cityBytes, ConvertUint16Bytes(rebellionFlag)...)

	return cityBytes
}

func AddCityToTile(inputFilename string, targetX int, targetY int, cityName string, tribe int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	// Overwrite tile header tribe
	tileStartOffset, ok := fileOffsetMap[buildTileStartKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile start key on x: %v, y: %v. Command not run.", targetX, targetY))
	}
	tileHeaderTribe := BuildTileHeaderTribeCity(targetX, targetY, tribe)
	if _, err := inputFile.WriteAt(tileHeaderTribe, int64(tileStartOffset+14)); err != nil {
		log.Fatal(err)
	}

	// Overwrite improvement data and set city
	offsetTileImprovementEnd, ok := fileOffsetMap[buildTileImprovementEndKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile end key on x: %v, y: %v. Command not run.", targetX, targetY))
	}
	remainder := GetFileRemainingData(inputFile, offsetTileImprovementEnd)

	offsetTileImprovementStart, ok := fileOffsetMap[buildTileImprovementStartKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile start key on x: %v, y: %v. Command not run.", targetX, targetY))
	}
	currentOffset := offsetTileImprovementStart
	cityBytes := BuildCity(cityName)
	if _, err := inputFile.WriteAt(cityBytes, int64(currentOffset)); err != nil {
		log.Fatal(err)
	}

	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(currentOffset+len(cityBytes))); err != nil {
		log.Fatal(err)
	}
}

func ResetTile(inputFilename string, targetX int, targetY int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	offsetTileEnd, ok := fileOffsetMap[buildTileEndKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile end key on x: %v, y: %v. Command not run.", targetX, targetY))
	}
	remainder := GetFileRemainingData(inputFile, offsetTileEnd)

	offsetTileStart, ok := fileOffsetMap[buildTileStartKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No tile start key on x: %v, y: %v. Command not run.", targetX, targetY))
	}

	currentOffset := offsetTileStart
	tileBytes := BuildEmptyTile(targetX, targetY)
	if _, err := inputFile.WriteAt(tileBytes, int64(offsetTileStart)); err != nil {
		log.Fatal(err)
	}
	currentOffset += len(tileBytes)

	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(currentOffset)); err != nil {
		log.Fatal(err)
	}
}

func ExpandRows(inputFilename string, newRowDimensions int) {
	if newRowDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	fmt.Println(fmt.Sprintf("Old dimensions width: %v, height: %v", saveOutput.MapWidth, saveOutput.MapHeight))

	if newRowDimensions <= saveOutput.MapHeight {
		log.Fatal(fmt.Sprintf("New row dimensions are less than existing dimensions, new value: %v, existing height: %v",
			newRowDimensions, saveOutput.MapHeight))
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	offset, ok := fileOffsetMap[buildTileEndKey(saveOutput.MapWidth-1, saveOutput.MapHeight-1)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: can't find endpoint of tile x: %v, y: %v. Command not run.",
			saveOutput.MapWidth-1, saveOutput.MapHeight-1))
	}

	remainder := GetFileRemainingData(inputFile, offset)
	currentOffset := offset
	allNewTileBytes := make([]byte, 0)
	for y := saveOutput.MapHeight; y < newRowDimensions; y++ {
		for x := 0; x < saveOutput.MapWidth; x++ {
			tileBytes := BuildEmptyTile(x, y)
			allNewTileBytes = append(allNewTileBytes, tileBytes...)
			currentOffset += len(tileBytes)
		}
	}
	if _, err := inputFile.WriteAt(allNewTileBytes, int64(offset)); err != nil {
		log.Fatal(err)
	}
	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(currentOffset)); err != nil {
		log.Fatal(err)
	}

	ModifyMapDimensions(inputFilename, saveOutput.MapWidth, newRowDimensions)

	finalSaveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	fmt.Println(fmt.Sprintf("New dimensions, width: %v, height: %v", finalSaveOutput.MapWidth, finalSaveOutput.MapHeight))
}

func ExpandColumns(inputFilename string, newColDimensions int) {
	if newColDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	fmt.Println(fmt.Sprintf("Old dimensions width: %v, height: %v", saveOutput.MapWidth, saveOutput.MapHeight))

	if newColDimensions <= saveOutput.MapWidth {
		log.Fatal(fmt.Sprintf("New column dimensions are less than existing dimensions, new value: %v, existing width: %v",
			newColDimensions, saveOutput.MapWidth))
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}
	for y := saveOutput.MapHeight - 1; y >= 0; y-- {
		offset, ok := fileOffsetMap[buildTileEndKey(saveOutput.MapWidth-1, y)]
		if !ok {
			log.Fatal(fmt.Sprintf("Error: can't find endpoint of tile x: %v, y: %v. Command not run.",
				saveOutput.MapWidth-1, saveOutput.MapHeight-1))
		}

		remainder := GetFileRemainingData(inputFile, offset)
		currentOffset := offset
		allNewTileBytes := make([]byte, 0)
		for x := saveOutput.MapWidth; x < newColDimensions; x++ {
			tileBytes := BuildEmptyTile(x, y)
			allNewTileBytes = append(allNewTileBytes, tileBytes...)
			currentOffset += len(tileBytes)
		}
		if _, err := inputFile.WriteAt(allNewTileBytes, int64(offset)); err != nil {
			log.Fatal(err)
		}
		// shift remaining data to the right
		if _, err := inputFile.WriteAt(remainder, int64(currentOffset)); err != nil {
			log.Fatal(err)
		}
	}

	ModifyMapDimensions(inputFilename, newColDimensions, saveOutput.MapHeight)

	finalSaveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	fmt.Println(fmt.Sprintf("New dimensions, width: %v, height: %v", finalSaveOutput.MapWidth, finalSaveOutput.MapHeight))
}

func ExpandTiles(inputFilename string, newSquareSizeDimensions int) {
	if newSquareSizeDimensions >= 256 {
		log.Fatal("Updated value is over 256")
	}

	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	if newSquareSizeDimensions <= saveOutput.MapWidth || newSquareSizeDimensions <= saveOutput.MapHeight {
		log.Fatal(fmt.Sprintf("New dimensions are less than existing dimensions, new value: %v, existing width: %v, height: %v",
			newSquareSizeDimensions, saveOutput.MapWidth, saveOutput.MapHeight))
	}

	ExpandColumns(inputFilename, newSquareSizeDimensions)
	ExpandRows(inputFilename, newSquareSizeDimensions)
}

func WriteEmptyRow(inputFilename string, maxX int, maxY int) {
	offset, ok := fileOffsetMap[buildTileEndKey(maxX, maxY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: can't find endpoint of tile x: %v, y: %v. Command not run.", maxX, maxY))
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	remainder := GetFileRemainingData(inputFile, offset)

	currentOffset := offset
	for x := 0; x <= maxX; x++ {
		tileBytes := BuildEmptyTile(x, maxY+1)

		if _, err := inputFile.WriteAt(tileBytes, int64(currentOffset)); err != nil {
			log.Fatal(err)
		}

		currentOffset += len(tileBytes)
	}

	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(currentOffset)); err != nil {
		log.Fatal(err)
	}
}

func WriteEmptyColumn(inputFilename string, maxX int, maxY int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	for y := maxY; y >= 0; y-- {
		offset, ok := fileOffsetMap[buildTileEndKey(maxX, y)]
		if !ok {
			log.Fatal(fmt.Sprintf("Error: can't find endpoint of tile x: %v, y: %v. Command not run.", maxX, maxY))
		}

		remainder := GetFileRemainingData(inputFile, offset)

		tileBytes := BuildEmptyTile(maxX+1, y)

		if _, err := inputFile.WriteAt(tileBytes, int64(offset)); err != nil {
			log.Fatal(err)
		}
		// shift remaining data to the right
		if _, err := inputFile.WriteAt(remainder, int64(offset+len(tileBytes))); err != nil {
			log.Fatal(err)
		}
	}
}

func RevealAllTiles(inputFilename string, newTribe int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	for i := saveOutput.MapHeight - 1; i >= 0; i-- {
		for j := saveOutput.MapWidth - 1; j >= 0; j-- {
			targetX := j
			targetY := i
			RevealTileForTribe(inputFilename, targetX, targetY, newTribe)
			fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, newTribe))
		}
	}
}

func RevealTileForTribe(inputFilename string, targetX int, targetY int, newTribe int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	fmt.Printf("Reading tile visibility data for (%v, %v)\n", targetX, targetY)
	offset, ok := fileOffsetMap[buildTileVisibilityLocationKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No visibility data on tile x: %v, y: %v. Command not run.", targetX, targetY))
	}

	visibilitySize := make([]byte, 1)
	_, err = inputFile.ReadAt(visibilitySize, int64(offset))
	if err != nil {
		log.Fatal("Failed to load visibility size:", err)
	}

	visibilityData := make([]byte, visibilitySize[0])
	bytesRead, err := inputFile.ReadAt(visibilityData, int64(offset)+1)
	if err != nil {
		log.Fatal("Failed to load visibility data:", err)
	}
	if bytesRead != int(visibilitySize[0]) {
		log.Fatal(fmt.Sprintf("Not enough visibility data loaded. Expected %v but only read %v bytes.", visibilitySize, bytesRead))
	}

	fmt.Println("Existing visibility data:", visibilityData)
	for i := 0; i < len(visibilityData); i++ {
		if int(visibilityData[i]) == newTribe {
			fmt.Printf("Tile is already visible to tribe %v. No change will be made to visibility data.\n", newTribe)
			return
		}
	}

	remainder := GetFileRemainingData(inputFile, offset+1+len(visibilityData))

	newVisibilityData := append(visibilityData, byte(newTribe))
	writeVisibilityData := append([]byte{uint8(len(newVisibilityData))}, newVisibilityData...)
	if _, err := inputFile.WriteAt(writeVisibilityData, int64(offset)); err != nil {
		log.Fatal(err)
	}
	if _, err := inputFile.WriteAt(remainder, int64(offset+len(writeVisibilityData))); err != nil {
		log.Fatal(err)
	}
}

func BuildEmptyPlayer(index int, playerName string, overrideColor color.RGBA) []byte {
	allPlayerData := make([]byte, 0)

	if index >= 254 {
		log.Fatal("Over 255 players")
	}

	allPlayerData = append(allPlayerData, byte(index))
	allPlayerData = append(allPlayerData, ConvertVarString(playerName)...)
	accountId := "00000000-0000-0000-0000-000000000000"
	allPlayerData = append(allPlayerData, ConvertVarString(accountId)...)
	autoPlay := 1 // true for bot
	allPlayerData = append(allPlayerData, byte(autoPlay))

	// player start coordinates
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(0)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(0)...)

	tribe := 2
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(tribe)...)

	// unknown
	allPlayerData = append(allPlayerData, byte(1))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(2)...)

	// unknown array
	newArraySize := index + 1
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(newArraySize)...)
	for i := 1; i <= int(newArraySize); i++ {
		playerId := i
		if i == newArraySize {
			playerId = 255
		}
		allPlayerData = append(allPlayerData, byte(playerId))
		// unknown
		allPlayerData = append(allPlayerData, ConvertUint32Bytes(0)...)
	}

	startingCurrency := 5
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(startingCurrency)...)
	score := 0
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(score)...)
	// unknown int
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(0)...)
	numCities := 1
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(numCities)...)

	techArrSize := 0
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(techArrSize)...)
	encounteredPlayersSize := 0
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(encounteredPlayersSize)...)
	numTasks := 0
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(numTasks)...)

	totalKills := 0
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(totalKills)...)
	totalLosses := 0
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(totalLosses)...)
	totalTribesDestroyed := 0
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(totalTribesDestroyed)...)
	allPlayerData = append(allPlayerData, []byte{overrideColor.B, overrideColor.G, overrideColor.A, 0}...)

	remainingPlayerData := make([]byte, 12)
	for i := 0; i < len(remainingPlayerData); i++ {
		remainingPlayerData[i] = 0
	}
	allPlayerData = append(allPlayerData, remainingPlayerData...)

	unknownBuffer2 := []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255}
	allPlayerData = append(allPlayerData, unknownBuffer2...)

	return allPlayerData
}

func generateRandomColor() color.RGBA {
	rand.Seed(time.Now().UnixNano())
	return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
}

func convertIntArrToByteArr(intArr []int) []byte {
	byteArr := make([]byte, len(intArr))
	for i := 0; i < len(byteArr); i++ {
		if intArr[i] > 255 {
			log.Fatal(fmt.Sprintf("Integer array index: %v, value: %v is over 255", i, intArr[i]))
		}
		byteArr[i] = byte(intArr[i])
	}
	return byteArr
}

func BuildNewPlayerUnknownArr(oldUnknownArr1 []int, newPlayerId int) []byte {
	existingLen := len(oldUnknownArr1)
	if existingLen % 5 != 0 {
		log.Fatal("Invalid array length. There was an error in processing.")
	}

	oldPlayerCount := existingLen / 5
	oldMaximumPlayerId := oldPlayerCount - 1 // excludes player 255 nature
	if oldMaximumPlayerId >= newPlayerId {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. No need to add player id %v.",
			oldPlayerCount, oldPlayerCount - 1, newPlayerId))
		return convertIntArrToByteArr(oldUnknownArr1)
	} else {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. New player id %v needs to be included.",
			oldPlayerCount, oldPlayerCount - 1, newPlayerId))
	}

	dataInsert := []int{newPlayerId, 0, 0, 0, 0}
	// assumes player 255 is always last
	existingPlayers := oldUnknownArr1[0 : existingLen - 5]
	naturePlayer := make([]int, 5)
	copy(naturePlayer, oldUnknownArr1[existingLen - 5 : existingLen])

	newUnknownArr1 := append(existingPlayers, dataInsert...)
	newUnknownArr1 = append(newUnknownArr1, naturePlayer...)
	return convertIntArrToByteArr(newUnknownArr1)
}

func convertPlayerIndexToId(playerIndex int, totalPlayers int) int {
	if playerIndex == totalPlayers - 1 {
		return 255
	} else {
		return playerIndex + 1
	}
}

func ModifyAllExistingPlayerUnknownArr(inputFilename string) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	newPlayerCount := len(saveOutput.PlayerData)
	fmt.Println(fmt.Sprintf("New player count: %v", newPlayerCount))

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	for i := len(saveOutput.PlayerData) - 1; i >= 0; i-- {
		playerId := convertPlayerIndexToId(i, len(saveOutput.PlayerData))
		fmt.Println("Player index:", i, ", corresponding player id", playerId)

		playerCurrencyKey := buildPlayerCurrencyKey(playerId)
		playerCurrencyOffset, ok := fileOffsetMap[playerCurrencyKey]
		if !ok {
			log.Fatal(fmt.Sprintf("Error: No player currency key: %v. Command not run.", playerCurrencyKey))
		}
		remainder := GetFileRemainingData(inputFile, playerCurrencyOffset)

		playerArr1Key := buildPlayerArr1Key(playerId)
		playerArr1Offset, ok := fileOffsetMap[playerArr1Key]
		if !ok {
			log.Fatal(fmt.Sprintf("Error: No player unknown arr1 key: %v. Command not run.", playerArr1Key))
		}

		newPlayerId := newPlayerCount - 1
		newUnknownArr1 := BuildNewPlayerUnknownArr(saveOutput.PlayerData[i].UnknownArr1, newPlayerId)
		if _, err := inputFile.WriteAt(ConvertUint16Bytes(newPlayerCount), int64(playerArr1Offset)); err != nil {
			log.Fatal(err)
		}
		if _, err := inputFile.WriteAt(newUnknownArr1, int64(playerArr1Offset+2)); err != nil {
			log.Fatal(err)
		}
		if _, err := inputFile.WriteAt(remainder, int64(playerArr1Offset+2+len(newUnknownArr1))); err != nil {
			log.Fatal(err)
		}
	}
}

func AddPlayer(inputFilename string) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	oldPlayerCount := len(saveOutput.PlayerData)
	fmt.Println(fmt.Sprintf("Old num players: %v", oldPlayerCount))

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	numPlayersKey := buildNumPlayersKey()
	numPlayersOffset, ok := fileOffsetMap[numPlayersKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No player start key: %v. Command not run.", numPlayersKey))
	}
	newPlayerCount := oldPlayerCount + 1
	WriteUint16AtFileOffset(inputFilename, numPlayersOffset, newPlayerCount)

	lastPlayerStartKey := buildPlayerStartKey(oldPlayerCount - 1)
	lastPlayerOffset, ok := fileOffsetMap[lastPlayerStartKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No last player start key: %v. Command not run.", lastPlayerStartKey))
	}
	remainder := GetFileRemainingData(inputFile, lastPlayerOffset)

	// existing index will be 1, 2, 3, ..., oldPlayerCount-1, 255 (size is oldPlayerCount)
	// new index list will be 1, 2, 3, ..., oldPlayerCount-1, oldPlayerCount, 255 (size is oldPlayerCount + 1)
	playerName := fmt.Sprintf("Player%v", oldPlayerCount)
	overrideColor := generateRandomColor()
	newPlayerBytes := BuildEmptyPlayer(oldPlayerCount, playerName, overrideColor)
	if _, err := inputFile.WriteAt(newPlayerBytes, int64(lastPlayerOffset)); err != nil {
		log.Fatal(err)
	}
	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(lastPlayerOffset+len(newPlayerBytes))); err != nil {
		log.Fatal(err)
	}

	ModifyAllExistingPlayerUnknownArr(inputFilename)
}
