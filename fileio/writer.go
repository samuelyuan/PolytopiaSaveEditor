package fileio

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func ModifyUnitTribe(inputFilename string, targetX int, targetY int, updatedValue int) {
	offset := GetUnitLocationFileOffset(targetX, targetY)
	WriteUint8AtFileOffset(inputFilename, offset+4, updatedValue)

	offsetPreviousUnit, ok := fileOffsetMap[buildPreviousUnitLocationKey(targetX, targetY)]
	if ok {
		WriteUint8AtFileOffset(inputFilename, offsetPreviousUnit+4, updatedValue)
	}
}

func BuildEmptyTile(x int, y int) []byte {
	flatTerrain := 3
	climate := 1
	altitude := 1
	owner := 0
	capital := 0
	capitalX := -1
	capitalY := -1
	tileHeaderBytes := BuildTileHeaderBytes(x, y, flatTerrain, climate, altitude, owner, capital, capitalX, capitalY)

	remainingTileData := make([]byte, 10)
	for i := 0; i < len(remainingTileData); i++ {
		remainingTileData[i] = 0
	}

	allTileData := append(tileHeaderBytes, remainingTileData...)
	return allTileData
}

func BuildTileHeaderBytes(x int, y int, terrain int, climate int, altitude int, owner int, capital int, capitalX int, capitalY int) []byte {
	worldCoordinates1 := make([]byte, 4)
	binary.LittleEndian.PutUint32(worldCoordinates1, uint32(x))
	worldCoordinates2 := make([]byte, 4)
	binary.LittleEndian.PutUint32(worldCoordinates2, uint32(y))
	terrainBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(terrainBytes, uint16(terrain))
	climateBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(climateBytes, uint16(climate))
	altitudeBytes := make([]byte, 2) // should be int16
	binary.LittleEndian.PutUint16(altitudeBytes, uint16(altitude))
	ownerBytes := []byte{byte(owner)}
	capitalBytes := []byte{byte(capital)}
	capitalCoordinates1 := make([]byte, 4) // should be int32
	binary.LittleEndian.PutUint32(capitalCoordinates1, uint32(capitalX))
	capitalCoordinates2 := make([]byte, 4) // should be int32
	binary.LittleEndian.PutUint32(capitalCoordinates2, uint32(capitalY))

	headerBytes := append(worldCoordinates1, worldCoordinates2...)
	headerBytes = append(headerBytes, terrainBytes...)
	headerBytes = append(headerBytes, climateBytes...)
	headerBytes = append(headerBytes, altitudeBytes...)
	headerBytes = append(headerBytes, ownerBytes...)
	headerBytes = append(headerBytes, capitalBytes...)
	headerBytes = append(headerBytes, capitalCoordinates1...)
	headerBytes = append(headerBytes, capitalCoordinates2...)

	return headerBytes
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
	// write new owner
	WriteUint8AtFileOffset(inputFilename, tileStartOffset+14, tribe)
	// set capital to 0 unless this city is designated as capital city
	WriteUint8AtFileOffset(inputFilename, tileStartOffset+15, 0)
	// write city coordinates
	WriteUint32AtFileOffset(inputFilename, tileStartOffset+16, targetX)
	WriteUint32AtFileOffset(inputFilename, tileStartOffset+20, targetY)

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

	WriteUint8AtFileOffset(inputFilename, currentOffset, 1)
	currentOffset += 1

	// city improvement is 1
	improvement := 1
	WriteUint16AtFileOffset(inputFilename, currentOffset, improvement)
	currentOffset += 2

	cityLevel := 1
	WriteUint16AtFileOffset(inputFilename, currentOffset, cityLevel)
	currentOffset += 2

	foundedTurn := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, foundedTurn)
	currentOffset += 2

	currentPopulation := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, currentPopulation)
	currentOffset += 2

	totalPopulation := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, totalPopulation)
	currentOffset += 2

	// unknown
	WriteUint16AtFileOffset(inputFilename, currentOffset, 1)
	currentOffset += 2

	baseScore := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, baseScore)
	currentOffset += 2

	// unknown
	WriteUint16AtFileOffset(inputFilename, currentOffset, 1)
	currentOffset += 2
	// unknown
	WriteUint16AtFileOffset(inputFilename, currentOffset, 0)
	currentOffset += 2

	// new city is not connected to capital
	connectedPlayerCapital := 0
	WriteUint8AtFileOffset(inputFilename, currentOffset, connectedPlayerCapital)
	currentOffset += 1

	hasCity := 1
	WriteUint8AtFileOffset(inputFilename, currentOffset, hasCity)
	currentOffset += 1

	cityNameBytes := []byte(cityName)
	WriteUint8AtFileOffset(inputFilename, currentOffset, len(cityNameBytes))
	currentOffset += 1
	if _, err := inputFile.WriteAt([]byte(cityName), int64(currentOffset)); err != nil {
		log.Fatal(err)
	}
	currentOffset += len(cityNameBytes)

	foundedTribe := 0
	WriteUint8AtFileOffset(inputFilename, currentOffset, foundedTribe)
	currentOffset += 1

	cityRewardsSize := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, cityRewardsSize)
	currentOffset += 2

	rebellionFlag := 0
	WriteUint16AtFileOffset(inputFilename, currentOffset, rebellionFlag)
	currentOffset += 2

	// shift remaining data to the right
	if _, err := inputFile.WriteAt(remainder, int64(currentOffset)); err != nil {
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

func RevealTileForTribe(inputFilename string, targetX int, targetY int, newTribe int) {
	fmt.Printf("Reading tile visibility data for (%v, %v)\n", targetX, targetY)
	offset, ok := fileOffsetMap[buildTileVisibilityLocationKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No visibility data on tile x: %v, y: %v. Command not run.", targetX, targetY))
	}

	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
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
