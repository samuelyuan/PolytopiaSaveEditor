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
		updatedValue, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}

		for i := saveOutput.MapHeight - 1; i >= 0; i-- {
			for j := saveOutput.MapWidth - 1; j >= 0; j-- {
				targetX := j
				targetY := i
				fileio.RevealTileForTribe(inputFilename, targetX, targetY, updatedValue)
				fmt.Println(fmt.Sprintf("Revealed (%v, %v) for tribe %v", targetX, targetY, updatedValue))
			}
		}
	} else {
		log.Fatal("Invalid mode:", mode)
	}
}
