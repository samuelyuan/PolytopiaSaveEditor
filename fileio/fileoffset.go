package fileio

import (
	"fmt"
	"io"
	"log"
)

func buildMapStartKey() string {
	return "MapStart"
}

func buildMapEndKey() string {
	return "MapEnd"
}

func buildUnitLocationKey(x int, y int) string {
	return fmt.Sprintf("UnitLocation%v,%v", x, y)
}

func buildPreviousUnitLocationKey(x int, y int) string {
	return fmt.Sprintf("PreviousUnitLocation%v,%v", x, y)
}

func buildTileVisibilityLocationKey(x int, y int) string {
	return fmt.Sprintf("TileVisibility%v,%v", x, y)
}

func buildTileStartKey(x int, y int) string {
	return fmt.Sprintf("TileStart%v,%v", x, y)
}

func buildTileEndKey(x int, y int) string {
	return fmt.Sprintf("TileEnd%v,%v", x, y)
}

func buildTileImprovementStartKey(x int, y int) string {
	return fmt.Sprintf("TileImprovementStart%v,%v", x, y)
}

func buildTileImprovementEndKey(x int, y int) string {
	return fmt.Sprintf("TileImprovementEnd%v,%v", x, y)
}

func buildAllPlayersStartKey() string {
	return "AllPlayersStart"
}

func buildAllPlayersEndKey() string {
	return "AllPlayersEnd"
}

func buildPlayerStartKey(index int) string {
	return fmt.Sprintf("PlayerStart%v", index)
}

func buildPlayerArr1Key(playerId int) string {
	return fmt.Sprintf("PlayerArr1-Id%v", playerId)
}

func buildPlayerCurrencyKey(playerId int) string {
	return fmt.Sprintf("PlayerCurrency-Id%v", playerId)
}

func GetUnitLocationFileOffset(targetX int, targetY int) int {
	offset, ok := fileOffsetMap[buildUnitLocationKey(targetX, targetY)]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: No unit on tile x: %v, y: %v. Command not run.", targetX, targetY))
	}
	return offset
}

func updateFileOffsetMap(fileOffsetMap map[string]int, streamReader *io.SectionReader, unitLocationKey string) {
	fileOffset, err := streamReader.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Fatal(err)
	}
	fileOffsetMap[unitLocationKey] = int(fileOffset)
}
