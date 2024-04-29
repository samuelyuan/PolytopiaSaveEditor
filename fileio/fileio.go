package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	fileOffsetMap = make(map[string]int)
	tribeUnitMap  = make(map[int][]UnitLocationData)
)

type MapHeaderInput struct {
	Version1           uint32
	Version2           uint32
	TotalActions       uint16
	CurrentTurn        uint32
	CurrentPlayerIndex uint8
	MaxUnitId          uint32
	UnknownByte1       uint8
	Seed               uint32
	TurnLimit          uint32
	Unknown1           [11]byte
	GameMode1          uint8
	GameMode2          uint8
}

type MapHeaderOutput struct {
	MapHeaderInput    MapHeaderInput
	MapName           string
	MapSquareSize     int
	DisabledTribesArr []int
	UnlockedTribesArr []int
	GameDifficulty    int
	NumOpponents      int
	UnknownArr        []byte
	SelectedTribes    map[int]int
	MapWidth          int
	MapHeight         int
}

type TileDataHeader struct {
	WorldCoordinates   [2]uint32
	Terrain            uint16
	Climate            uint16
	Altitude           int16
	Owner              uint8
	Capital            uint8
	CapitalCoordinates [2]int32
}

type TileData struct {
	Header            TileDataHeader
	Terrain           int
	Climate           int
	Owner             int
	Capital           int
	ResourceExists    bool
	ResourceType      int
	ImprovementExists bool
	ImprovementType   int
	HasCity           bool
	CityName          string
	CityData          *CityData
	Unit              *UnitData
	BufferUnitData    []byte
	PlayerVisibility  []uint8
	HasRoad           bool
	HasWaterRoute     bool
	Unknown           []byte
}

type CityData struct {
	CityLevel              int
	CurrentPopulation      int
	TotalPopulation        int
	UnknownShort1          int
	ParkBonus              int
	UnknownShort2          int
	UnknownShort3          int
	ConnectedPlayerCapital int
	UnknownByte1           int
	CityName               string
	FlagBeforeRewards      int
	CityRewards            []int
	RebellionFlag          int
	RebellionBuffer        []byte
}

type PlayerData struct {
	Id                   int
	Name                 string
	AccountId            string
	AutoPlay             bool
	StartTileCoordinates [2]int
	Tribe                int
	UnknownByte1         int
	UnknownInt1          int
	UnknownArr1          []int
	Currency             int
	Score                int
	UnknownInt2          int
	NumCities            int
	AvailableTech        []int
	EncounteredPlayers   []int
	Tasks                []PlayerTaskData
	TotalUnitsKilled     int
	TotalUnitsLost       int
	TotalTribesDestroyed int
	OverrideColor        []byte
	UnknownByte2         byte
	UniqueImprovements   []int
	DiplomacyArr         []DiplomacyData
	DiplomacyMessages    []DiplomacyMessage
	DestroyedByTribe     int
	DestroyedTurn        int
	UnknownBuffer2       []byte
}

type UnitData struct {
	Id                 uint32
	Owner              uint8
	UnitType           uint16
	Unknown            [8]byte // seems to be all zeros
	CurrentCoordinates [2]int32
	HomeCoordinates    [2]int32
	Health             uint16 // should be divided by 10 to get value ingame
	PromotionLevel     uint16
	Experience         uint16
	Moved              bool
	Attacked           bool
	Flipped            bool
	CreatedTurn        uint16
}

type ImprovementData struct {
	Level       uint16
	Founded     uint16
	Unknown1    [6]byte
	BaseScore   uint16
	Unknown2    [6]byte
	UnknownInt1 uint16
	Unknown3    [3]byte
}

type PlayerTaskData struct {
	Type   int
	Buffer []byte
}

type DiplomacyMessage struct {
	MessageType int
	Sender      int
}

type DiplomacyData struct {
	PlayerId               uint8
	DiplomacyRelationState uint8
	LastAttackTurn         int32
	EmbassyLevel           uint8
	LastPeaceBrokenTurn    int32
	FirstMeet              int32
	EmbassyBuildTurn       int32
	PreviousAttackTurn     int32
}

type PolytopiaSaveOutput struct {
	MapHeight       int
	MapWidth        int
	OwnerTribeMap   map[int]int
	InitialTileData [][]TileData
	TileData        [][]TileData
	MaxTurn         int
	PlayerData      []PlayerData
	FileOffsetMap   map[string]int
	TribeCityMap    map[int][]CityLocationData
	TribeUnitMap    map[int][]UnitLocationData
}

type CityLocationData struct {
	X        int
	Y        int
	CityName string
}

type UnitLocationData struct {
	X        int
	Y        int
	UnitType int
}

func readVarString(reader *io.SectionReader, varName string) string {
	variableLength := uint8(0)
	if err := binary.Read(reader, binary.LittleEndian, &variableLength); err != nil {
		log.Fatal("Failed to load variable length: ", err)
	}

	stringValue := make([]byte, variableLength)
	if err := binary.Read(reader, binary.LittleEndian, &stringValue); err != nil {
		log.Fatal(fmt.Sprintf("Failed to load string value. Variable length: %v, name: %s. Error:", variableLength, varName), err)
	}

	return string(stringValue[:])
}

func unsafeReadUint32(reader *io.SectionReader) uint32 {
	unsignedIntValue := uint32(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint32: ", err)
	}
	return unsignedIntValue
}

func unsafeReadInt32(reader *io.SectionReader) int32 {
	signedIntValue := int32(0)
	if err := binary.Read(reader, binary.LittleEndian, &signedIntValue); err != nil {
		log.Fatal("Failed to load int32: ", err)
	}
	return signedIntValue
}

func unsafeReadUint16(reader *io.SectionReader) uint16 {
	unsignedIntValue := uint16(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint16: ", err)
	}
	return unsignedIntValue
}

func unsafeReadInt16(reader *io.SectionReader) int16 {
	signedIntValue := int16(0)
	if err := binary.Read(reader, binary.LittleEndian, &signedIntValue); err != nil {
		log.Fatal("Failed to load int16: ", err)
	}
	return signedIntValue
}

func unsafeReadUint8(reader *io.SectionReader) uint8 {
	unsignedIntValue := uint8(0)
	if err := binary.Read(reader, binary.LittleEndian, &unsignedIntValue); err != nil {
		log.Fatal("Failed to load uint8: ", err)
	}
	return unsignedIntValue
}

func readFixedList(streamReader *io.SectionReader, listSize int) []byte {
	buffer := make([]byte, listSize)
	if err := binary.Read(streamReader, binary.LittleEndian, &buffer); err != nil {
		log.Fatal("Failed to load buffer: ", err)
	}
	return buffer
}

func readCityData(streamReader *io.SectionReader) CityData {
	cityLevel := unsafeReadUint32(streamReader)
	currentPopulation := unsafeReadInt16(streamReader)
	totalPopulation := unsafeReadUint16(streamReader)
	unknownShort1 := unsafeReadInt16(streamReader)
	parkBonus := unsafeReadInt16(streamReader)
	unknownShort2 := unsafeReadInt16(streamReader)
	unknownShort3 := unsafeReadInt16(streamReader)
	connectedPlayerCapital := unsafeReadUint8(streamReader)
	unknownByte1 := unsafeReadUint8(streamReader)
	cityName := readVarString(streamReader, "CityName")

	flagBeforeRewards := unsafeReadUint8(streamReader)
	if flagBeforeRewards != 0 {
		log.Fatal("flagBeforeRewards isn't 0")
	}

	cityRewardsSize := unsafeReadUint16(streamReader)
	cityRewards := make([]int, cityRewardsSize)
	for i := 0; i < int(cityRewardsSize); i++ {
		cityReward := unsafeReadUint16(streamReader)
		cityRewards[i] = int(cityReward)
	}

	rebellionFlag := unsafeReadUint16(streamReader)
	var rebellionBuffer []byte
	if rebellionFlag != 0 {
		rebellionBuffer = readFixedList(streamReader, 2)
	}

	return CityData{
		CityLevel:              int(cityLevel),
		CurrentPopulation:      int(currentPopulation),
		TotalPopulation:        int(totalPopulation),
		UnknownShort1:          int(unknownShort1),
		ParkBonus:              int(parkBonus),
		UnknownShort2:          int(unknownShort2),
		UnknownShort3:          int(unknownShort3),
		ConnectedPlayerCapital: int(connectedPlayerCapital),
		UnknownByte1:           int(unknownByte1),
		CityName:               cityName,
		FlagBeforeRewards:      int(flagBeforeRewards),
		CityRewards:            cityRewards,
		RebellionFlag:          int(rebellionFlag),
		RebellionBuffer:        rebellionBuffer,
	}
}

func readTileData(streamReader *io.SectionReader, tileData [][]TileData, mapWidth int, mapHeight int) {
	allUnitData := make([]UnitData, 0)

	for i := 0; i < int(mapHeight); i++ {
		for j := 0; j < int(mapWidth); j++ {
			tileStartKey := buildTileStartKey(j, i)
			updateFileOffsetMap(fileOffsetMap, streamReader, tileStartKey)

			tileDataHeader := TileDataHeader{}
			if err := binary.Read(streamReader, binary.LittleEndian, &tileDataHeader); err != nil {
				log.Fatal("Failed to load tileDataHeader: ", err)
			}

			// Sanity check
			if int(tileDataHeader.WorldCoordinates[0]) != j || int(tileDataHeader.WorldCoordinates[1]) != i {
				log.Fatal(fmt.Sprintf("File reached unexpected location. Iteration (%v, %v) isn't equal to world coordinates (%v, %v)",
					i, j, tileDataHeader.WorldCoordinates[0], tileDataHeader.WorldCoordinates[1]))
			}

			resourceExistsFlag := unsafeReadUint8(streamReader)
			resourceType := -1
			if resourceExistsFlag == 1 {
				resourceType = int(unsafeReadUint16(streamReader))
			}

			tileImprovementStartKey := buildTileImprovementStartKey(j, i)
			updateFileOffsetMap(fileOffsetMap, streamReader, tileImprovementStartKey)
			improvementExistsFlag := unsafeReadUint8(streamReader)
			improvementType := -1
			if improvementExistsFlag == 1 {
				improvementType = int(unsafeReadUint16(streamReader))
			}

			// If tile is city, read differently
			var cityData CityData
			if tileDataHeader.Owner > 0 && resourceType == -1 && improvementType == 1 {
				// No resource, but has improvement that is city
				cityData = readCityData(streamReader)
			} else if improvementType != -1 {
				// Has improvement and read improvement data
				improvementData := ImprovementData{}
				if err := binary.Read(streamReader, binary.LittleEndian, &improvementData); err != nil {
					log.Fatal("Failed to load buffer: ", err)
				}
			}

			tileImprovementEndKey := buildTileImprovementEndKey(j, i)
			updateFileOffsetMap(fileOffsetMap, streamReader, tileImprovementEndKey)

			// Read unit data
			hasUnitFlag := unsafeReadUint8(streamReader)
			var unitDataPtr *UnitData
			bufferUnitData := make([]byte, 0)
			if hasUnitFlag == 1 {
				unitLocationKey := buildUnitLocationKey(int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1]))
				updateFileOffsetMap(fileOffsetMap, streamReader, unitLocationKey)

				unitData := UnitData{}
				if err := binary.Read(streamReader, binary.LittleEndian, &unitData); err != nil {
					log.Fatal("Failed to load buffer: ", err)
				}
				unitDataPtr = &unitData

				updateTribeUnitMap(tileDataHeader, unitData)

				hasOtherUnitFlag := unsafeReadUint8(streamReader)
				if hasOtherUnitFlag == 1 {
					// If unit embarks or disembarks, a new unit is created in the backend, but it's still the same unit in the game
					previousUnitLocationKey := buildPreviousUnitLocationKey(int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1]))
					updateFileOffsetMap(fileOffsetMap, streamReader, previousUnitLocationKey)

					previousUnitData := UnitData{}
					if err := binary.Read(streamReader, binary.LittleEndian, &previousUnitData); err != nil {
						log.Fatal("Failed to load buffer: ", err)
					}

					_ = unsafeReadUint8(streamReader)
					bufferUnitData2 := readFixedList(streamReader, 7)
					bufferUnitData3 := readFixedList(streamReader, 7)
					if bufferUnitData2[0] == 1 {
						bufferUnitData3Remainder := readFixedList(streamReader, 4)
						bufferUnitData3 = append(bufferUnitData3, bufferUnitData3Remainder...)
					}
					bufferUnitData = append(bufferUnitData, bufferUnitData2...)
					bufferUnitData = append(bufferUnitData, bufferUnitData3...)
				} else {
					bufferUnitFlag := unsafeReadUint8(streamReader)
					bufferSize := 6
					if bufferUnitFlag == 1 {
						bufferSize = 8
					}

					bufferUnitData = readFixedList(streamReader, bufferSize)
				}
			}

			tileVisibilityLocationKey := buildTileVisibilityLocationKey(int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1]))
			updateFileOffsetMap(fileOffsetMap, streamReader, tileVisibilityLocationKey)
			playerVisibilityListSize := unsafeReadUint8(streamReader)
			playerVisibilityList := readFixedList(streamReader, int(playerVisibilityListSize))

			tileRoadKey := buildTileRoadKey(int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1]))
			updateFileOffsetMap(fileOffsetMap, streamReader, tileRoadKey)
			hasRoad := unsafeReadUint8(streamReader)

			hasWaterRoute := unsafeReadUint8(streamReader)
			unknown := readFixedList(streamReader, 4)

			hasCity := false
			cityName := ""
			var cityDataPtr *CityData
			if improvementType == 1 {
				hasCity = true

				if tileDataHeader.Owner > 0 {
					cityName = cityData.CityName
					cityDataPtr = &cityData
				}
			}

			tileData[i][j] = TileData{
				Header:            tileDataHeader,
				Terrain:           int(tileDataHeader.Terrain),
				Climate:           int(tileDataHeader.Climate),
				Owner:             int(tileDataHeader.Owner),
				Capital:           int(tileDataHeader.Capital),
				HasCity:           hasCity,
				CityName:          cityName,
				CityData:          cityDataPtr,
				Unit:              unitDataPtr,
				BufferUnitData:    bufferUnitData,
				PlayerVisibility:  playerVisibilityList,
				HasRoad:           hasRoad != 0,
				HasWaterRoute:     hasWaterRoute != 0,
				Unknown:           unknown,
				ResourceExists:    resourceExistsFlag != 0,
				ResourceType:      resourceType,
				ImprovementExists: improvementExistsFlag != 0,
				ImprovementType:   improvementType,
			}

			tileEndKey := buildTileEndKey(int(tileDataHeader.WorldCoordinates[0]), int(tileDataHeader.WorldCoordinates[1]))
			updateFileOffsetMap(fileOffsetMap, streamReader, tileEndKey)

			if tileData[i][j].Unit != nil {
				allUnitData = append(allUnitData, *tileData[i][j].Unit)
			}
		}
	}
}

func readMapHeader(streamReader *io.SectionReader) MapHeaderOutput {
	mapHeaderInput := MapHeaderInput{}
	if err := binary.Read(streamReader, binary.LittleEndian, &mapHeaderInput); err != nil {
		log.Fatal("Failed to load MapHeaderInput: ", err)
	}

	mapName := readVarString(streamReader, "MapName")

	// map dimenions is a square: squareSize x squareSize
	updateFileOffsetMap(fileOffsetMap, streamReader, "SquareSizeKey")
	squareSize := int(unsafeReadUint32(streamReader))

	disabledTribesSize := unsafeReadUint16(streamReader)
	disabledTribesArr := make([]int, disabledTribesSize)
	for i := 0; i < int(disabledTribesSize); i++ {
		disabledTribesArr[i] = int(unsafeReadUint16(streamReader))
	}

	unlockedTribesSize := unsafeReadUint16(streamReader)
	unlockedTribesArr := make([]int, unlockedTribesSize)
	for i := 0; i < int(unlockedTribesSize); i++ {
		unlockedTribesArr[i] = int(unsafeReadUint16(streamReader))
	}

	gameDifficulty := unsafeReadUint16(streamReader)
	numOpponents := unsafeReadUint32(streamReader)
	unknownArr := readFixedList(streamReader, 5+int(unlockedTribesSize))

	selectedTribeSkinSize := unsafeReadUint32(streamReader)
	selectedTribeSkins := make(map[int]int)
	for i := 0; i < int(selectedTribeSkinSize); i++ {
		tribe := unsafeReadUint16(streamReader)
		skin := unsafeReadUint16(streamReader)
		selectedTribeSkins[int(tribe)] = int(skin)
	}

	updateFileOffsetMap(fileOffsetMap, streamReader, "MapWidth")
	mapWidth := unsafeReadUint16(streamReader)
	updateFileOffsetMap(fileOffsetMap, streamReader, "MapHeight")
	mapHeight := unsafeReadUint16(streamReader)
	if mapWidth == 0 && mapHeight == 0 {
		mapWidth = unsafeReadUint16(streamReader)
		mapHeight = unsafeReadUint16(streamReader)
	}

	return MapHeaderOutput{
		MapHeaderInput:    mapHeaderInput,
		MapName:           mapName,
		MapSquareSize:     squareSize,
		DisabledTribesArr: disabledTribesArr,
		UnlockedTribesArr: unlockedTribesArr,
		GameDifficulty:    int(gameDifficulty),
		NumOpponents:      int(numOpponents),
		UnknownArr:        unknownArr,
		SelectedTribes:    selectedTribeSkins,
		MapWidth:          int(mapWidth),
		MapHeight:         int(mapHeight),
	}
}

func readPlayerData(streamReader *io.SectionReader) PlayerData {
	playerId := unsafeReadUint8(streamReader)
	playerName := readVarString(streamReader, "playerName")
	playerAccountId := readVarString(streamReader, "playerAccountId")
	autoPlay := unsafeReadUint8(streamReader)
	startTileCoordinates1 := unsafeReadInt32(streamReader)
	startTileCoordinates2 := unsafeReadInt32(streamReader)
	tribe := unsafeReadUint16(streamReader)
	unknownByte1 := unsafeReadUint8(streamReader)
	unknownInt1 := unsafeReadUint32(streamReader)

	playerArr1Key := buildPlayerArr1Key(int(playerId))
	updateFileOffsetMap(fileOffsetMap, streamReader, playerArr1Key)
	unknownArrLen1 := unsafeReadUint16(streamReader)
	unknownArr1 := make([]int, 0)
	for i := 0; i < int(unknownArrLen1); i++ {
		value1 := unsafeReadUint8(streamReader)
		value2 := readFixedList(streamReader, 4)
		unknownArr1 = append(unknownArr1, int(value1), int(value2[0]), int(value2[1]), int(value2[2]), int(value2[3]))
	}

	playerCurrencyKey := buildPlayerCurrencyKey(int(playerId))
	updateFileOffsetMap(fileOffsetMap, streamReader, playerCurrencyKey)
	currency := unsafeReadUint32(streamReader)

	score := unsafeReadUint32(streamReader)
	unknownInt2 := unsafeReadUint32(streamReader)
	numCities := unsafeReadUint16(streamReader)

	techArrayLen := unsafeReadUint16(streamReader)
	techArray := make([]int, techArrayLen)
	for i := 0; i < int(techArrayLen); i++ {
		techType := unsafeReadUint16(streamReader)
		techArray[i] = int(techType)
	}

	encounteredPlayersLen := unsafeReadUint16(streamReader)
	encounteredPlayers := make([]int, 0)
	for i := 0; i < int(encounteredPlayersLen); i++ {
		playerId := unsafeReadUint8(streamReader)
		encounteredPlayers = append(encounteredPlayers, int(playerId))
	}

	numTasks := unsafeReadInt16(streamReader)
	taskArr := make([]PlayerTaskData, int(numTasks))
	for i := 0; i < int(numTasks); i++ {
		taskType := unsafeReadInt16(streamReader)

		var buffer []byte
		if taskType == 1 || taskType == 5 { // Task type 1 is Pacifist, type 5 is Killer
			buffer = readFixedList(streamReader, 6) // Extra buffer contains a uint32
		} else if taskType >= 1 && taskType <= 8 {
			buffer = readFixedList(streamReader, 2)
		} else {
			log.Fatal("Invalid task type:", taskType)
		}
		taskArr[i] = PlayerTaskData{
			Type:   int(taskType),
			Buffer: buffer,
		}
	}

	totalKills := unsafeReadInt32(streamReader)
	totalLosses := unsafeReadInt32(streamReader)
	totalTribesDestroyed := unsafeReadInt32(streamReader)
	overrideColor := readFixedList(streamReader, 4)
	unknownByte2 := unsafeReadUint8(streamReader)

	playerUniqueImprovementsSize := unsafeReadUint16(streamReader)
	playerUniqueImprovements := make([]int, int(playerUniqueImprovementsSize))
	for i := 0; i < int(playerUniqueImprovementsSize); i++ {
		improvement := unsafeReadUint16(streamReader)
		playerUniqueImprovements[i] = int(improvement)
	}

	diplomacyArrLen := unsafeReadUint16(streamReader)
	diplomacyArr := make([]DiplomacyData, int(diplomacyArrLen))
	for i := 0; i < len(diplomacyArr); i++ {
		diplomacyData := DiplomacyData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &diplomacyData); err != nil {
			log.Fatal("Failed to load diplomacyData: ", err)
		}
		diplomacyArr[i] = diplomacyData
	}

	diplomacyMessagesSize := unsafeReadUint16(streamReader)
	diplomacyMessagesArr := make([]DiplomacyMessage, int(diplomacyMessagesSize))
	for i := 0; i < int(diplomacyMessagesSize); i++ {
		messageType := unsafeReadUint8(streamReader)
		sender := unsafeReadUint8(streamReader)

		diplomacyMessagesArr[i] = DiplomacyMessage{
			MessageType: int(messageType),
			Sender:      int(sender),
		}
	}

	destroyedByTribe := unsafeReadUint8(streamReader)
	destroyedTurn := unsafeReadUint32(streamReader)
	unknownBuffer2 := readFixedList(streamReader, 14)

	return PlayerData{
		Id:                   int(playerId),
		Name:                 playerName,
		AccountId:            playerAccountId,
		AutoPlay:             int(autoPlay) != 0,
		StartTileCoordinates: [2]int{int(startTileCoordinates1), int(startTileCoordinates2)},
		Tribe:                int(tribe),
		UnknownByte1:         int(unknownByte1),
		UnknownInt1:          int(unknownInt1),
		UnknownArr1:          unknownArr1,
		Currency:             int(currency),
		Score:                int(score),
		UnknownInt2:          int(unknownInt2),
		NumCities:            int(numCities),
		AvailableTech:        techArray,
		EncounteredPlayers:   encounteredPlayers,
		Tasks:                taskArr,
		TotalUnitsKilled:     int(totalKills),
		TotalUnitsLost:       int(totalLosses),
		TotalTribesDestroyed: int(totalTribesDestroyed),
		OverrideColor:        overrideColor,
		UnknownByte2:         unknownByte2,
		UniqueImprovements:   playerUniqueImprovements,
		DiplomacyArr:         diplomacyArr,
		DiplomacyMessages:    diplomacyMessagesArr,
		DestroyedByTribe:     int(destroyedByTribe),
		DestroyedTurn:        int(destroyedTurn),
		UnknownBuffer2:       unknownBuffer2,
	}
}

func readAllPlayerData(streamReader *io.SectionReader) []PlayerData {
	numPlayersKey := buildNumPlayersKey()
	updateFileOffsetMap(fileOffsetMap, streamReader, numPlayersKey)
	numPlayers := unsafeReadUint16(streamReader)
	fmt.Println("Num players:", numPlayers)
	allPlayerData := make([]PlayerData, int(numPlayers))

	for i := 0; i < int(numPlayers); i++ {
		playerStartKey := buildPlayerStartKey(i)
		updateFileOffsetMap(fileOffsetMap, streamReader, playerStartKey)
		playerData := readPlayerData(streamReader)
		allPlayerData[i] = playerData
	}
	return allPlayerData
}

func updateTribeUnitMap(tileDataHeader TileDataHeader, unitData UnitData) {
	tribeOwner := int(unitData.Owner)
	_, ok := tribeUnitMap[tribeOwner]
	if !ok {
		tribeUnitMap[tribeOwner] = make([]UnitLocationData, 0)
	}
	unitLocationData := UnitLocationData{
		X:        int(tileDataHeader.WorldCoordinates[0]),
		Y:        int(tileDataHeader.WorldCoordinates[1]),
		UnitType: int(unitData.UnitType),
	}
	tribeUnitMap[tribeOwner] = append(tribeUnitMap[tribeOwner], unitLocationData)
}

func buildOwnerTribeMap(allPlayerData []PlayerData) map[int]int {
	ownerTribeMap := make(map[int]int)

	for i := 0; i < len(allPlayerData); i++ {
		playerData := allPlayerData[i]
		mappedTribe, ok := ownerTribeMap[playerData.Id]
		if ok {
			log.Fatal(fmt.Sprintf("Owner to tribe map has duplicate player id %v already mapped to %v", playerData.Id, mappedTribe))
		}
		ownerTribeMap[playerData.Id] = playerData.Tribe
	}

	return ownerTribeMap
}

func buildTribeCityMap(currentMapHeaderOutput MapHeaderOutput, tileData [][]TileData) map[int][]CityLocationData {
	tribeCityMap := make(map[int][]CityLocationData)
	for i := 0; i < int(currentMapHeaderOutput.MapHeight); i++ {
		for j := 0; j < int(currentMapHeaderOutput.MapWidth); j++ {
			if tileData[i][j].HasCity {
				tribeOwner := tileData[i][j].Owner
				_, ok := tribeCityMap[tribeOwner]
				if !ok {
					tribeCityMap[tribeOwner] = make([]CityLocationData, 0)
				}
				cityLocationData := CityLocationData{
					X:        int(tileData[i][j].Header.WorldCoordinates[0]),
					Y:        int(tileData[i][j].Header.WorldCoordinates[1]),
					CityName: tileData[i][j].CityName,
				}
				tribeCityMap[tribeOwner] = append(tribeCityMap[tribeOwner], cityLocationData)
			}
		}
	}
	return tribeCityMap
}

func ReadPolytopiaDecompressedFile(inputFilename string) (*PolytopiaSaveOutput, error) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
		return nil, err
	}
	fi, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fileLength := fi.Size()
	streamReader := io.NewSectionReader(inputFile, int64(0), fileLength)

	fileOffsetMap = make(map[string]int)
	tribeUnitMap = make(map[int][]UnitLocationData)

	// Read initial map state
	initialMapHeaderOutput := readMapHeader(streamReader)

	initialTileData := make([][]TileData, initialMapHeaderOutput.MapHeight)
	for i := 0; i < initialMapHeaderOutput.MapHeight; i++ {
		initialTileData[i] = make([]TileData, initialMapHeaderOutput.MapWidth)
	}
	readTileData(streamReader, initialTileData, initialMapHeaderOutput.MapWidth, initialMapHeaderOutput.MapHeight)
	ownerTribeMap := buildOwnerTribeMap(readAllPlayerData(streamReader))

	_ = readFixedList(streamReader, 3)

	// Read current map state
	currentMapHeaderOutput := readMapHeader(streamReader)

	tileData := make([][]TileData, currentMapHeaderOutput.MapHeight)
	for i := 0; i < currentMapHeaderOutput.MapHeight; i++ {
		tileData[i] = make([]TileData, currentMapHeaderOutput.MapWidth)
	}
	readTileData(streamReader, tileData, currentMapHeaderOutput.MapWidth, currentMapHeaderOutput.MapHeight)
	playerData := readAllPlayerData(streamReader)
	ownerTribeMap = buildOwnerTribeMap(playerData)

	tribeCityMap := buildTribeCityMap(currentMapHeaderOutput, tileData)

	output := &PolytopiaSaveOutput{
		MapHeight:       currentMapHeaderOutput.MapHeight,
		MapWidth:        currentMapHeaderOutput.MapWidth,
		OwnerTribeMap:   ownerTribeMap,
		InitialTileData: initialTileData,
		TileData:        tileData,
		MaxTurn:         int(currentMapHeaderOutput.MapHeaderInput.CurrentTurn),
		PlayerData:      playerData,
		FileOffsetMap:   fileOffsetMap,
		TribeCityMap:    tribeCityMap,
		TribeUnitMap:    tribeUnitMap,
	}
	return output, nil
}
