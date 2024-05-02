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

func WriteAndShiftData(inputFilename string, offsetStartOriginalBlockKey string, offsetEndOriginalBlockKey string, newData []byte) {
	// Update file offsets to make sure they are up to date
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}
	fmt.Println(fmt.Sprintf("Dimensions width: %v, height: %v", saveOutput.MapWidth, saveOutput.MapHeight))

	// Open file to modify
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	offsetOriginalBlockStart, ok := fileOffsetMap[offsetStartOriginalBlockKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find start of data block with key %v. Command not run.", offsetStartOriginalBlockKey))
	}
	offsetOriginalBlockEnd, ok := fileOffsetMap[offsetEndOriginalBlockKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find end of data block with key %v. Command not run.", offsetStartOriginalBlockKey))
	}
	// Get all data after end of block
	remainder := GetFileRemainingData(inputFile, offsetOriginalBlockEnd)

	// overwrite block with new data at original block start
	if _, err := inputFile.WriteAt(newData, int64(offsetOriginalBlockStart)); err != nil {
		log.Fatal(err)
	}

	// shift remaining data and write after new data instead of original end start
	if _, err := inputFile.WriteAt(remainder, int64(offsetOriginalBlockStart+len(newData))); err != nil {
		log.Fatal(err)
	}
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

func ConvertByteList(oldArr []int) []byte {
	// the values are stored as ints but they were originally bytes
	newArr := make([]byte, len(oldArr))
	for i := 0; i < len(oldArr); i++ {
		if oldArr[i] > 255 {
			log.Fatal("Byte list has value over 255")
		}
		newArr[i] = byte(oldArr[i])
	}
	return newArr
}

func ConvertBoolToByte(value bool) byte {
	if value {
		return 1
	} else {
		return 0
	}
}

func ConvertCityDataToBytes(cityData CityData) []byte {
	data := make([]byte, 0)
	data = append(data, ConvertUint16Bytes(int(cityData.CityLevel))...)
	data = append(data, ConvertUint16Bytes(int(cityData.FoundedTurn))...)
	data = append(data, ConvertUint16Bytes(int(cityData.CurrentPopulation))...)
	data = append(data, ConvertUint16Bytes(int(cityData.TotalPopulation))...)
	data = append(data, ConvertUint16Bytes(int(cityData.UnknownShort1))...)
	data = append(data, ConvertUint16Bytes(int(cityData.ParkBonus))...)
	data = append(data, ConvertUint16Bytes(int(cityData.UnknownShort2))...)
	data = append(data, ConvertUint16Bytes(int(cityData.UnknownShort3))...)
	data = append(data, byte(cityData.ConnectedPlayerCapital))
	data = append(data, byte(cityData.HasCityName))
	data = append(data, ConvertVarString(cityData.CityName)...)
	data = append(data, byte(cityData.FoundedTribe))
	data = append(data, ConvertUint16Bytes(len(cityData.CityRewards))...)
	for i := 0; i < len(cityData.CityRewards); i++ {
		data = append(data, ConvertUint16Bytes(int(cityData.CityRewards[i]))...)
	}
	data = append(data, ConvertUint16Bytes(int(cityData.RebellionFlag))...)
	if cityData.RebellionFlag != 0 {
		data = append(data, cityData.RebellionBuffer...)
	}
	return data
}

func ConvertImprovementDataToBytes(improvementData ImprovementData) []byte {
	data := make([]byte, 0)
	data = append(data, ConvertUint16Bytes(int(improvementData.Level))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.FoundedTurn))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.CurrentPopulation))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.TotalPopulation))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.UnknownShort1))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.BaseScore))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.Unknown2[0]))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.Unknown2[1]))...)
	data = append(data, improvementData.ConnectedPlayerCapital)
	data = append(data, improvementData.HasCityName)
	data = append(data, improvementData.FoundedTribe)
	data = append(data, ConvertUint16Bytes(int(improvementData.RewardsSize))...)
	data = append(data, ConvertUint16Bytes(int(improvementData.RebellionFlag))...)
	return data
}

func ConvertUnitDataToBytes(unitData UnitData) []byte {
	data := make([]byte, 0)
	data = append(data, ConvertUint32Bytes(int(unitData.Id))...)
	data = append(data, byte(unitData.Owner))
	data = append(data, ConvertUint16Bytes(int(unitData.UnitType))...)
	data = append(data, unitData.Unknown[:]...)
	data = append(data, ConvertUint32Bytes(int(unitData.CurrentCoordinates[0]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.CurrentCoordinates[1]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.HomeCoordinates[0]))...)
	data = append(data, ConvertUint32Bytes(int(unitData.HomeCoordinates[1]))...)
	data = append(data, ConvertUint16Bytes(int(unitData.Health))...)
	data = append(data, ConvertUint16Bytes(int(unitData.PromotionLevel))...)
	data = append(data, ConvertUint16Bytes(int(unitData.Experience))...)
	data = append(data, ConvertBoolToByte(unitData.Moved))
	data = append(data, ConvertBoolToByte(unitData.Attacked))
	data = append(data, ConvertBoolToByte(unitData.Flipped))
	data = append(data, ConvertUint16Bytes(int(unitData.CreatedTurn))...)
	return data
}

func ConvertTileToBytes(tileData TileData) []byte {
	tileBytes := make([]byte, 0)

	headerBytes := make([]byte, 0)
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.WorldCoordinates[0]))...)
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.WorldCoordinates[1]))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Terrain))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Climate))...)
	headerBytes = append(headerBytes, ConvertUint16Bytes(int(tileData.Altitude))...) // should be int16
	headerBytes = append(headerBytes, byte(tileData.Owner))
	headerBytes = append(headerBytes, byte(tileData.Capital))
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.CapitalCoordinates[0]))...) // should be int32
	headerBytes = append(headerBytes, ConvertUint32Bytes(int(tileData.CapitalCoordinates[1]))...) // should be int32
	tileBytes = append(tileBytes, headerBytes...)

	if tileData.ResourceExists {
		tileBytes = append(tileBytes, byte(1))
		tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.ResourceType)...)
	} else {
		tileBytes = append(tileBytes, byte(0))
	}

	if tileData.ImprovementExists {
		tileBytes = append(tileBytes, byte(1))
		tileBytes = append(tileBytes, ConvertUint16Bytes(tileData.ImprovementType)...)
	} else {
		tileBytes = append(tileBytes, byte(0))
	}

	if tileData.Owner > 0 && tileData.ResourceType == -1 && tileData.ImprovementType == 1 {
		// No resource, but has improvement that is city
		tileBytes = append(tileBytes, ConvertCityDataToBytes(*tileData.CityData)...)
	} else if tileData.ImprovementType != -1 {
		// Has improvement and read improvement data
		tileBytes = append(tileBytes, ConvertImprovementDataToBytes(*tileData.ImprovementData)...)
	}

	// no unit
	if tileData.Unit != nil {
		tileBytes = append(tileBytes, 1)
		tileBytes = append(tileBytes, ConvertUnitDataToBytes(*tileData.Unit)...)

		if tileData.PreviousUnit != nil {
			tileBytes = append(tileBytes, 1)
			tileBytes = append(tileBytes, ConvertUnitDataToBytes(*tileData.PreviousUnit)...)

			tileBytes = append(tileBytes, 0)
			tileBytes = append(tileBytes, tileData.BufferUnitData...)
		} else {
			tileBytes = append(tileBytes, 0)
			tileBytes = append(tileBytes, byte(tileData.BufferUnitFlag))
			tileBytes = append(tileBytes, tileData.BufferUnitData...)
		}
	} else {
		tileBytes = append(tileBytes, 0)
	}

	tileBytes = append(tileBytes, byte(len(tileData.PlayerVisibility)))
	tileBytes = append(tileBytes, ConvertByteList(tileData.PlayerVisibility)...)
	tileBytes = append(tileBytes, ConvertBoolToByte(tileData.HasRoad))
	tileBytes = append(tileBytes, ConvertBoolToByte(tileData.HasWaterRoute))
	tileBytes = append(tileBytes, tileData.Unknown...)
	return tileBytes
}

func ConvertMapDataToBytes(tileData [][]TileData) []byte {
	mapHeight := len(tileData)
	mapWidth := len(tileData[0])

	allMapBytes := make([]byte, 0)
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			tileBytes := ConvertTileToBytes(tileData[i][j])
			allMapBytes = append(allMapBytes, tileBytes...)
		}
	}
	return allMapBytes
}

func ConvertDiplomacyDataToBytes(diplomacyData DiplomacyData) []byte {
	data := make([]byte, 0)
	data = append(data, byte(diplomacyData.PlayerId))
	data = append(data, byte(diplomacyData.DiplomacyRelationState))
	data = append(data, ConvertUint32Bytes(int(diplomacyData.LastAttackTurn))...)
	data = append(data, byte(diplomacyData.EmbassyLevel))
	data = append(data, ConvertUint32Bytes(int(diplomacyData.LastPeaceBrokenTurn))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.FirstMeet))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.EmbassyBuildTurn))...)
	data = append(data, ConvertUint32Bytes(int(diplomacyData.PreviousAttackTurn))...)
	return data
}

func ConvertPlayerDataToBytes(playerData PlayerData) []byte {
	allPlayerData := make([]byte, 0)

	allPlayerData = append(allPlayerData, byte(playerData.Id))
	allPlayerData = append(allPlayerData, ConvertVarString(playerData.Name)...)
	allPlayerData = append(allPlayerData, ConvertVarString(playerData.AccountId)...)
	allPlayerData = append(allPlayerData, ConvertBoolToByte(playerData.AutoPlay))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.StartTileCoordinates[0])...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.StartTileCoordinates[1])...)
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.Tribe)...)
	allPlayerData = append(allPlayerData, byte(playerData.UnknownByte1))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.UnknownInt1)...)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.UnknownArr1)/5)...)
	for i := 0; i < len(playerData.UnknownArr1); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.UnknownArr1[i]))
	}

	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.Currency)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.Score)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.UnknownInt2)...)
	allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.NumCities)...)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.AvailableTech))...)
	for i := 0; i < len(playerData.AvailableTech); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.AvailableTech[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.EncounteredPlayers))...)
	for i := 0; i < len(playerData.EncounteredPlayers); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.EncounteredPlayers[i]))
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.Tasks))...)
	for i := 0; i < len(playerData.Tasks); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.Tasks[i].Type)...)
		allPlayerData = append(allPlayerData, playerData.Tasks[i].Buffer...)
	}

	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalUnitsKilled)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalUnitsLost)...)
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.TotalTribesDestroyed)...)
	allPlayerData = append(allPlayerData, ConvertByteList(playerData.OverrideColor)...)

	allPlayerData = append(allPlayerData, playerData.UnknownByte2)

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.UniqueImprovements))...)
	for i := 0; i < len(playerData.UniqueImprovements); i++ {
		allPlayerData = append(allPlayerData, ConvertUint16Bytes(playerData.UniqueImprovements[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.DiplomacyArr))...)
	for i := 0; i < len(playerData.DiplomacyArr); i++ {
		allPlayerData = append(allPlayerData, ConvertDiplomacyDataToBytes(playerData.DiplomacyArr[i])...)
	}

	allPlayerData = append(allPlayerData, ConvertUint16Bytes(len(playerData.DiplomacyMessages))...)
	for i := 0; i < len(playerData.DiplomacyMessages); i++ {
		allPlayerData = append(allPlayerData, byte(playerData.DiplomacyMessages[i].MessageType))
		allPlayerData = append(allPlayerData, byte(playerData.DiplomacyMessages[i].Sender))
	}

	allPlayerData = append(allPlayerData, byte(playerData.DestroyedByTribe))
	allPlayerData = append(allPlayerData, ConvertUint32Bytes(playerData.DestroyedTurn)...)
	allPlayerData = append(allPlayerData, playerData.UnknownBuffer2...)

	return allPlayerData
}

func ConvertAllPlayerDataToBytes(allPlayerData []PlayerData) []byte {
	allPlayerBytes := make([]byte, 0)
	allPlayerBytes = append(allPlayerBytes, ConvertUint16Bytes(len(allPlayerData))...)
	for i := 0; i < len(allPlayerData); i++ {
		allPlayerBytes = append(allPlayerBytes, ConvertPlayerDataToBytes(allPlayerData[i])...)
	}
	return allPlayerBytes
}

func WriteTileToFile(inputFilename string, tileDataOverwrite TileData, targetX int, targetY int) {
	tileBytes := ConvertTileToBytes(tileDataOverwrite)
	WriteAndShiftData(inputFilename, buildTileStartKey(targetX, targetY), buildTileEndKey(targetX, targetY), tileBytes)
}

func WriteMapToFile(inputFilename string, tileDataOverwrite [][]TileData) {
	allTileBytes := ConvertMapDataToBytes(tileDataOverwrite)
	WriteAndShiftData(inputFilename, buildMapStartKey(), buildMapEndKey(), allTileBytes)
}

func WritePlayersToFile(inputFilename string, playersList []PlayerData) {
	allPlayerBytes := ConvertAllPlayerDataToBytes(playersList)
	WriteAndShiftData(inputFilename, buildAllPlayersStartKey(), buildAllPlayersEndKey(), allPlayerBytes)
}

func ModifyTileTerrain(inputFilename string, targetX int, targetY int, updatedValue int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	updatedTile := saveOutput.TileData[targetY][targetX]

	// write terrain
	updatedTile.Terrain = updatedValue

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
	updatedTile.Altitude = altitude

	WriteTileToFile(inputFilename, updatedTile, targetX, targetY)
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
	tileData := TileData{
		WorldCoordinates:   [2]int{x, y},
		Terrain:            3,
		Climate:            1,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  false,
		ImprovementType:    -1,
		HasCity:            false,
		CityData:           nil,
		ImprovementData:    nil,
		Unit:               nil,
		BufferUnitData:     []uint8{},
		PlayerVisibility:   []int{},
		HasRoad:            false,
		HasWaterRoute:      false,
		Unknown:            []uint8{0, 0, 0, 0},
	}
	allTileData := ConvertTileToBytes(tileData)
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

	cityData := CityData{
		CityLevel:              1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		UnknownShort1:          1,
		ParkBonus:              0,
		UnknownShort2:          1,
		UnknownShort3:          0,
		ConnectedPlayerCapital: 0,
		HasCityName:            1,
		CityName:               cityName,
		FoundedTribe:           0,
		CityRewards:            []int{},
		RebellionFlag:          0,
		RebellionBuffer:        []byte{},
	}
	cityBytes = append(cityBytes, ConvertCityDataToBytes(cityData)...)

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

	allNewTileBytes := make([]byte, 0)
	for y := saveOutput.MapHeight; y < newRowDimensions; y++ {
		for x := 0; x < saveOutput.MapWidth; x++ {
			allNewTileBytes = append(allNewTileBytes, BuildEmptyTile(x, y)...)
		}
	}

	WriteAndShiftData(inputFilename, buildMapEndKey(), buildMapEndKey(), allNewTileBytes)
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

			visibilityData := saveOutput.TileData[i][j].PlayerVisibility
			fmt.Println("Existing visibility data:", visibilityData)
			isAlreadyVisible := false
			for visibilityIndex := 0; visibilityIndex < len(visibilityData); visibilityIndex++ {
				if int(visibilityData[visibilityIndex]) == newTribe {
					fmt.Printf("Tile is already visible to tribe %v. No change will be made to visibility data.\n", newTribe)
					isAlreadyVisible = true
					break
				}
			}
			if !isAlreadyVisible {
				saveOutput.TileData[i][j].PlayerVisibility = append(saveOutput.TileData[i][j].PlayerVisibility, newTribe)
				fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, newTribe))
			}
		}
	}

	for i := saveOutput.MapHeight - 1; i >= 0; i-- {
		for j := saveOutput.MapWidth - 1; j >= 0; j-- {
			fmt.Println(fmt.Sprintf("Tile (%v, %v) visibility: %v", j, i, saveOutput.TileData[i][j].PlayerVisibility))
		}
	}

	WriteMapToFile(inputFilename, saveOutput.TileData)
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
	if index >= 254 {
		log.Fatal("Over 255 players")
	}

	// unknown array
	newArraySize := index + 1
	unknownArr1 := make([]int, 0)
	for i := 1; i <= int(newArraySize); i++ {
		playerId := i
		if i == newArraySize {
			playerId = 255
		}
		unknownArr1 = append(unknownArr1, playerId)
		// unknown
		unknownArr1 = append(unknownArr1, 0, 0, 0, 0)
	}

	playerData := PlayerData{
		Id:                   index,
		Name:                 playerName,
		AccountId:            "00000000-0000-0000-0000-000000000000",
		AutoPlay:             true,
		StartTileCoordinates: [2]int{0, 0},
		Tribe:                2, // Ai-mo
		UnknownByte1:         1,
		UnknownInt1:          2,
		UnknownArr1:          unknownArr1,
		Currency:             5,
		Score:                0,
		UnknownInt2:          0,
		NumCities:            1,
		AvailableTech:        []int{},
		EncounteredPlayers:   []int{},
		Tasks:                []PlayerTaskData{},
		TotalUnitsKilled:     0,
		TotalUnitsLost:       0,
		TotalTribesDestroyed: 0,
		OverrideColor:        []int{int(overrideColor.B), int(overrideColor.G), int(overrideColor.A), 0},
		UnknownByte2:         0,
		UniqueImprovements:   []int{},
		DiplomacyArr:         []DiplomacyData{},
		DiplomacyMessages:    []DiplomacyMessage{},
		DestroyedByTribe:     0,
		DestroyedTurn:        0,
		UnknownBuffer2:       []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255},
	}

	return ConvertPlayerDataToBytes(playerData)
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
	if existingLen%5 != 0 {
		log.Fatal(fmt.Sprintf("Invalid array length. New player unknown array length %v is not divisible by 5.", existingLen))
	}

	oldPlayerCount := existingLen / 5
	oldMaximumPlayerId := oldPlayerCount - 1 // excludes player 255 nature
	if oldMaximumPlayerId >= newPlayerId {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. No need to add player id %v.",
			oldPlayerCount, oldPlayerCount-1, newPlayerId))
		return convertIntArrToByteArr(oldUnknownArr1)
	} else {
		fmt.Println(fmt.Sprintf("Existing player count is %v, which includes players 1 to %v. New player id %v needs to be included.",
			oldPlayerCount, oldPlayerCount-1, newPlayerId))
	}

	dataInsert := []int{newPlayerId, 0, 0, 0, 0}
	// assumes player 255 is always last
	existingPlayers := oldUnknownArr1[0 : existingLen-5]
	naturePlayer := make([]int, 5)
	copy(naturePlayer, oldUnknownArr1[existingLen-5:existingLen])

	newUnknownArr1 := append(existingPlayers, dataInsert...)
	newUnknownArr1 = append(newUnknownArr1, naturePlayer...)
	return convertIntArrToByteArr(newUnknownArr1)
}

func convertPlayerIndexToId(playerIndex int, totalPlayers int) int {
	if playerIndex == totalPlayers-1 {
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

	numPlayersKey := buildAllPlayersStartKey()
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

func SwapPlayers(inputFilename string, playerId1 int, playerId2 int) {
	saveOutput, err := ReadPolytopiaDecompressedFile(inputFilename)
	if err != nil {
		log.Fatal("Failed to read save file")
	}

	 // assumes 254 is not used by any players
	unusedPlayerId := 254

	// Need to reassign so we don't merge the players
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == playerId1 {
				saveOutput.TileData[i][j].Owner = unusedPlayerId
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(playerId1) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(unusedPlayerId)
			}

			if saveOutput.TileData[i][j].PreviousUnit != nil && saveOutput.TileData[i][j].PreviousUnit.Owner == uint8(playerId1) {
				saveOutput.TileData[i][j].PreviousUnit.Owner = uint8(unusedPlayerId)
			}
		}
	}

	// Overwrite all playerId2 tiles and units with playerId1
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == playerId2 {
				saveOutput.TileData[i][j].Owner = playerId1
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(playerId2) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(playerId1)
			}

			if saveOutput.TileData[i][j].PreviousUnit != nil && saveOutput.TileData[i][j].PreviousUnit.Owner == uint8(playerId2) {
				saveOutput.TileData[i][j].PreviousUnit.Owner = uint8(playerId1)
			}
		}
	}

	// Overwrite old playerId tiles and units with playerId2
	for i := 0; i < saveOutput.MapHeight; i++ {
		for j := 0; j < saveOutput.MapWidth; j++ {
			if saveOutput.TileData[i][j].Owner == unusedPlayerId {
				saveOutput.TileData[i][j].Owner = playerId2
			}

			if saveOutput.TileData[i][j].Unit != nil && saveOutput.TileData[i][j].Unit.Owner == uint8(unusedPlayerId) {
				saveOutput.TileData[i][j].Unit.Owner = uint8(playerId2)
			}

			if saveOutput.TileData[i][j].PreviousUnit != nil && saveOutput.TileData[i][j].PreviousUnit.Owner == uint8(unusedPlayerId) {
				saveOutput.TileData[i][j].PreviousUnit.Owner = uint8(playerId2)
			}
		}
	}

	WriteMapToFile(inputFilename, saveOutput.TileData)
}
