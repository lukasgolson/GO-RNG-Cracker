package application

import (
	"awesomeProject/internal/tree"
	"fmt"
	"github.com/lukasgolson/FileArray/serialization"
	"log"
	"math/rand"
	"os"
	"sync"
)

func processPartition(seedSpaceLow, seedSpaceHigh, fileCount int64, sequenceHigh, sequenceOffset int, graphPath string, randSource *rand.Rand) error {
	numberOfSeeds := seedSpaceHigh - seedSpaceLow

	if fileCount > numberOfSeeds {
		fileCount = numberOfSeeds
	}

	seedsPerFile := numberOfSeeds / fileCount

	for fileIndex := int64(0); fileIndex < fileCount; fileIndex++ {
		var bkTree, err = tree.NewOrLoad(graphPath+fmt.Sprintf("-%d", fileIndex), false)
		if err != nil {
			return err
		}
		startSeed := seedSpaceLow + (fileIndex * seedsPerFile)
		endSeed := startSeed + seedsPerFile

		loadedSeedPosition := serialization.Length(seedSpaceLow) + bkTree.Length()

		if bkTree.Length() > 0 {
			fmt.Println("Loading existing tree... Start seed", startSeed, "end seed:", endSeed, "previous end seed:", loadedSeedPosition)

			if loadedSeedPosition > serialization.Length(startSeed) && loadedSeedPosition < serialization.Length(endSeed) {

				fmt.Println("Previous tree is in range. Everything is fine.")

				const seedOverlap = 5 // the number of seeds that overlap between trees to ensure that we don't miss any seeds when loading existing trees.

				newStartSeed := int64(loadedSeedPosition) - seedOverlap

				if newStartSeed < startSeed {
					newStartSeed = startSeed
				}

				startSeed = newStartSeed

			} else if loadedSeedPosition < serialization.Length(startSeed) {
				return fmt.Errorf("previous tree ends before our start seed. Ensure that you are using the the same settings as you used to generate the previous tree")
			} else if loadedSeedPosition > serialization.Length(endSeed) {
				return fmt.Errorf("previous tree ends after our end seed")
			}
		} else {
			fmt.Println("No existing tree found. Creating new tree with name", graphPath, "... Start seed", startSeed, "end seed:", endSeed)
			err := bkTree.PreExpand(serialization.Length(seedsPerFile))
			if err != nil {
				return err
			}
		}

		for seed := startSeed; seed <= endSeed; seed++ {
			sequence := GenerateRandomSequence(seed, 32, sequenceHigh, sequenceOffset, randSource)

			err := bkTree.Add([32]byte(sequence), int32(seed))
			if err != nil {
				return err
			}
		}

		err = bkTree.Close()
	}

	return nil
}

func Initialize(coreCount, fileCount int, seedCount int64, sequenceHigh, sequenceOffset int, dataDirectories []string) error {

	if fileCount < 1 {
		return fmt.Errorf("file count must be at least 1")
	}

	if seedCount < 1 {
		return fmt.Errorf("seed count must be at least 1")
	}

	if coreCount < 1 {
		return fmt.Errorf("core count must be at least 1")
	}

	if fileCount%coreCount != 0 {
		return fmt.Errorf("file count must be divisible by the core count")
	}

	if int64(fileCount) > seedCount {
		return fmt.Errorf("file count must be less than or equal to the seed count")
	}

	for _, dir := range dataDirectories {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var wg sync.WaitGroup

	partitionSize := seedCount / int64(coreCount)
	filesPerPartition := int64(fileCount) / int64(coreCount)

	overlapPerCore := seedCount % int64(coreCount)
	for p := int64(0); p < int64(coreCount); p++ {
		lo := partitionSize * p
		hi := partitionSize * (p + 1)

		if p < overlapPerCore {
			hi++
		}

		dirIndex := int(p) % len(dataDirectories)
		dataDirectory := dataDirectories[dirIndex]

		wg.Add(1)

		go func(seedSpaceLow, seedSpaceHigh int64, partitionID int64) {
			defer wg.Done()
			randSource := rand.New(rand.NewSource(0))
			dir := fmt.Sprintf("%s/graph-%d", dataDirectory, partitionID)

			if err := processPartition(seedSpaceLow, seedSpaceHigh, filesPerPartition, sequenceHigh, sequenceOffset, dir, randSource); err != nil {
				log.Printf("Error processing partition: %v\n", err)
			}
		}(lo, hi, p)
	}

	wg.Wait()

	return nil
}
