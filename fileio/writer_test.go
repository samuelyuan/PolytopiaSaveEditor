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
    t.Fatalf(`Change to include player 17 failed. Result = %v, expected = %v`, resultBytesWithChange, expectedBytesWithChange)
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
		t.Fatalf(`Size not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
	if !reflect.DeepEqual(resultBytes, expectedBytes) {
		t.Fatalf(`Contents not equal. Result = %v, expected = %v`, resultBytes, expectedBytes)
	}
}
