package fileio

import (
	"image/color"
	"reflect"
	"testing"
)

func TestBuildTile(t *testing.T) {
	resultBytes := BuildEmptyTile(12, 34)

	expectedBytes := []byte{12, 0, 0, 0, 34, 0, 0, 0,
		3, 0, 1, 0, 1, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertCityDataToBytes(t *testing.T) {
	cityData := CityData{
		CityLevel:              3,
		CurrentPopulation:      1,
		TotalPopulation:        6,
		UnknownShort1:          1,
		ParkBonus:              0,
		UnknownShort2:          1,
		UnknownShort3:          -2,
		ConnectedPlayerCapital: 1,
		HasCityName:            1,
		CityName:               "Test",
		FoundedTribe:           0,
		CityRewards:            []int{4, 7},
		RebellionFlag:          0,
		RebellionBuffer:        []byte{},
	}
	resultBytes := ConvertCityDataToBytes(cityData)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 6, 0, 1, 0, 0, 0, 1, 0, 254, 255, 1, 1, 4, 84, 101, 115, 116, 0, 2, 0, 4, 0, 7, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertImprovementDataToBytes(t *testing.T) {
	improvementData := ImprovementData{
		Level:                  1,
		FoundedTurn:            0,
		CurrentPopulation:      0,
		TotalPopulation:        0,
		UnknownShort1:          1,
		BaseScore:              0,
		Unknown2:               [2]uint16{0, 0},
		ConnectedPlayerCapital: 0,
		HasCityName:            0,
		FoundedTribe:           0,
		RewardsSize:            0,
		RebellionFlag:          0,
	}
	resultBytes := ConvertImprovementDataToBytes(improvementData)
	expectedBytes := []byte{1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertUnitDataToBytes(t *testing.T) {
	unitData := UnitData{
		Id:                 4,
		Owner:              4,
		UnitType:           2,
		Unknown:            [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		CurrentCoordinates: [2]int32{8, 2},
		HomeCoordinates:    [2]int32{8, 2},
		Health:             100,
		PromotionLevel:     0,
		Experience:         0,
		Moved:              false,
		Attacked:           false,
		Flipped:            false,
		CreatedTurn:        0,
	}
	resultBytes := ConvertUnitDataToBytes(unitData)
	expectedBytes := []byte{4, 0, 0, 0, 4, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 0, 0, 0, 2, 0, 0, 0, 8, 0, 0, 0, 2, 0, 0, 0, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertEmptyTileDataToBytes(t *testing.T) {
	tileData := TileData{
		WorldCoordinates:   [2]int{3, 1},
		Terrain:            3,
		Climate:            8,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  false,
		ImprovementType:    -1,
		HasCity:            false,
		CityName:           "",
		CityData:           nil,
		ImprovementData:    nil,
		Unit:               nil,
		BufferUnitData:     []uint8{},
		PlayerVisibility:   []uint8{},
		HasRoad:            false,
		HasWaterRoute:      false,
		Unknown:            []uint8{0, 0, 0, 0},
	}
	resultBytes := ConvertTileToBytes(tileData)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertTileDataToBytes(t *testing.T) {
	tileData := TileData{
		WorldCoordinates:   [2]int{3, 1},
		Terrain:            3,
		Climate:            8,
		Altitude:           1,
		Owner:              0,
		Capital:            0,
		CapitalCoordinates: [2]int{-1, -1},
		ResourceExists:     false,
		ResourceType:       -1,
		ImprovementExists:  true,
		ImprovementType:    1,
		HasCity:            true,
		CityName:           "",
		CityData:           nil,
		ImprovementData: &ImprovementData{
			Level:                  1,
			FoundedTurn:            0,
			CurrentPopulation:      0,
			TotalPopulation:        0,
			UnknownShort1:          1,
			BaseScore:              0,
			Unknown2:               [2]uint16{1, 0},
			ConnectedPlayerCapital: 0,
			HasCityName:            0,
			FoundedTribe:           0,
			RewardsSize:            0,
			RebellionFlag:          0,
		},
		Unit:             nil,
		BufferUnitData:   []uint8{},
		PlayerVisibility: []uint8{},
		HasRoad:          false,
		HasWaterRoute:    false,
		Unknown:          []uint8{0, 0, 0, 0},
	}

	resultBytes := ConvertTileToBytes(tileData)
	expectedBytes := []byte{3, 0, 0, 0, 1, 0, 0, 0, 3, 0, 8, 0, 1, 0, 0, 0,
		// coordinates
		255, 255, 255, 255, 255, 255, 255, 255,
		// resource
		0,
		// improvement
		1, 1, 0,
		1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestConvertPlayerDataToBytes(t *testing.T) {

	playerData := PlayerData{
		Id:                   1,
		Name:                 "TestPlayer",
		AccountId:            "00000000-0000-0000-0000-000000000000",
		AutoPlay:             true,
		StartTileCoordinates: [2]int{6, 22},
		Tribe:                15,
		UnknownByte1:         1,
		UnknownInt1:          1,
		UnknownArr1: []int{1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
			6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
			11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
			16, 41, 125, 0, 0, 255, 0, 0, 0, 0},
		Currency:           900,
		Score:              10000,
		UnknownInt2:        0,
		NumCities:          11,
		AvailableTech:      []int{0, 8, 15, 10, 39, 18, 13, 1, 4, 14, 20},
		EncounteredPlayers: []int{7, 11, 3, 5, 10},
		Tasks: []PlayerTaskData{
			{Type: 6, Buffer: []byte{1, 1}},
			{Type: 5, Buffer: []byte{1, 1, 10, 0, 0, 0}},
			{Type: 8, Buffer: []byte{1, 0}},
			{Type: 3, Buffer: []byte{1, 1}},
		},
		TotalUnitsKilled:     28,
		TotalUnitsLost:       32,
		TotalTribesDestroyed: 1,
		OverrideColor:        []byte{153, 0, 255, 255},
		UnknownByte2:         0,
		UniqueImprovements:   []int{27},
		DiplomacyArr: []DiplomacyData{
			{PlayerId: 1, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 2, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 3, DiplomacyRelationState: 0, LastAttackTurn: 21, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 7, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 4, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 5, DiplomacyRelationState: 0, LastAttackTurn: 19, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 8, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 6, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 7, DiplomacyRelationState: 0, LastAttackTurn: 20, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 0, EmbassyBuildTurn: -100, PreviousAttackTurn: 21},
			{PlayerId: 8, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 9, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 10, DiplomacyRelationState: 0, LastAttackTurn: 15, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 13, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 11, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: 0, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 12, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 13, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 14, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 15, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 16, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 255, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
			{PlayerId: 0, DiplomacyRelationState: 0, LastAttackTurn: -100, EmbassyLevel: 0, LastPeaceBrokenTurn: -100, FirstMeet: -100, EmbassyBuildTurn: -100, PreviousAttackTurn: -100},
		},
		DiplomacyMessages: []DiplomacyMessage{},
		DestroyedByTribe:  0,
		DestroyedTurn:     0,
		UnknownBuffer2:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255},
	}
	resultBytes := ConvertPlayerDataToBytes(playerData)
	expectedBytes := []byte{1,
		// Player name
		10, 84, 101, 115, 116, 80, 108, 97, 121, 101, 114,
		// Account Id
		36, 48, 48, 48, 48, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		1, 6, 0, 0, 0, 22, 0, 0, 0, 15, 0, 1, 1, 0, 0, 0,
		// Unknown Array 1
		17, 0,
		1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
		6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
		11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
		16, 41, 125, 0, 0, 255, 0, 0, 0, 0,
		// currency
		132, 3, 0, 0,
		// score
		16, 39, 0, 0,
		0, 0, 0, 0,
		// num cities
		11, 0,
		// tech
		11, 0, 0, 0, 8, 0, 15, 0, 10, 0, 39, 0, 18, 0, 13, 0, 1, 0, 4, 0, 14, 0, 20, 0,
		// encountered players
		5, 0, 7, 11, 3, 5, 10,
		// tasks
		4, 0, 6, 0, 1, 1, 5, 0, 1, 1, 10, 0, 0, 0, 8, 0, 1, 0, 3, 0, 1, 1,
		28, 0, 0, 0,
		32, 0, 0, 0,
		1, 0, 0, 0,
		// override color
		153, 0, 255, 255,
		0,
		// improvements
		1, 0, 27, 0,
		18, 0, 1, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 2, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 3, 0, 21, 0, 0, 0, 0, 156, 255, 255, 255, 7, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 4, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 5, 0, 19, 0, 0, 0, 0, 156, 255, 255, 255, 8, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 6, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 7, 0, 20, 0, 0, 0, 0, 156, 255, 255, 255, 0, 0, 0, 0, 156, 255, 255, 255, 21, 0, 0, 0, 8, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 9, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 10, 0, 15, 0, 0, 0, 0, 156, 255, 255, 255, 13, 0, 0, 0, 156, 255, 255, 255, 156, 255, 255, 255, 11, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 0, 0, 0, 0, 156, 255, 255, 255, 156, 255, 255, 255, 12, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 13, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 14, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 15, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 16, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 255, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 0, 0, 156, 255, 255, 255, 0, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 156, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255,
	}

	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestBuildTileHeaderTribeCity(t *testing.T) {
	targetX := 12
	targetY := 34
	tribe := 1
	resultBytes := BuildTileHeaderTribeCity(targetX, targetY, tribe)

	expectedBytes := []byte{1, 0, 12, 0, 0, 0, 34, 0, 0, 0}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}

func TestBuildCity(t *testing.T) {
	cityName := "Test"
	cityBytes := BuildCity(cityName)

	expectedBytes := []byte{1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 4, 84, 101, 115, 116, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(cityBytes, expectedBytes) {
		t.Fatalf(`result = %v, expected = %v`, cityBytes, expectedBytes)
	}
}

func TestBuildNewPlayerUnknownArr(t *testing.T) {
	oldArr := []int{1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
		6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
		11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
		16, 41, 125, 0, 0, 255, 0, 0, 0, 0}
	resultBytesNoChange := BuildNewPlayerUnknownArr(oldArr, 16)
	expectedBytesNoChange := []byte{1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
		6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
		11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
		16, 41, 125, 0, 0, 255, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytesNoChange, expectedBytesNoChange) {
		t.Fatalf(`No change failed. Result = %v, expected = %v`, resultBytesNoChange, expectedBytesNoChange)
	}

	resultBytesWithChange := BuildNewPlayerUnknownArr(oldArr, 17)
	expectedBytesWithChange := []byte{1, 0, 0, 0, 0, 2, 80, 69, 0, 0, 3, 88, 29, 1, 0, 4, 39, 95, 0, 0, 5, 222, 34, 1, 0,
		6, 218, 77, 1, 0, 7, 134, 250, 0, 0, 8, 243, 153, 0, 0, 9, 131, 143, 0, 0, 10, 180, 147, 0, 0,
		11, 74, 89, 0, 0, 12, 7, 125, 0, 0, 13, 74, 69, 0, 0, 14, 66, 163, 0, 0, 15, 165, 216, 0, 0,
		16, 41, 125, 0, 0, 17, 0, 0, 0, 0, 255, 0, 0, 0, 0}
	if !reflect.DeepEqual(resultBytesWithChange, expectedBytesWithChange) {
		t.Fatalf(`Change to include player 17,  failed. Result = %v, expected = %v`, resultBytesWithChange, expectedBytesWithChange)
	}
}

func TestBuildPlayer(t *testing.T) {
	resultBytes := BuildEmptyPlayer(17, "Player17", color.RGBA{100, 150, 200, 255})
	expectedBytes := []byte{17,
		// Player name
		8, 80, 108, 97, 121, 101, 114, 49, 55,
		// 00000000-0000-0000-0000-000000000000
		36, 48, 48, 48, 48, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 45, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48,
		1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 1, 2, 0, 0, 0,
		18, 0, 1, 0, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 0, 4, 0, 0, 0, 0, 5, 0, 0, 0, 0, 6, 0, 0, 0, 0, 7, 0, 0, 0, 0,
		8, 0, 0, 0, 0, 9, 0, 0, 0, 0, 10, 0, 0, 0, 0, 11, 0, 0, 0, 0, 12, 0, 0, 0, 0, 13, 0, 0, 0, 0, 14, 0, 0, 0, 0,
		15, 0, 0, 0, 0, 16, 0, 0, 0, 0, 17, 0, 0, 0, 0, 255, 0, 0, 0, 0,
		5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		// override color
		200, 150, 255, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255, 255, 255, 255, 255, 0, 0, 255, 255, 255, 255,
	}

	if !reflect.DeepEqual(len(resultBytes), len(expectedBytes)) {
		t.Fatalf(`Size not equal. Result = %v (size = %v), expected = %v (size = %v)`,
			resultBytes, len(resultBytes), expectedBytes, len(expectedBytes))
	}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}
