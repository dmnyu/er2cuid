package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "/Volumes/ACMBornDigital/Archivematica-Staging/in-process/ingest/UA_RG_9-8/UA_RG_9-8-batch-1"
	//check that the root exists and is a dir
	fi, err := os.Stat(root)
	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		panic(fmt.Errorf("%s is not a directory", root))
	}

	//locate the tsv
	mdDir := filepath.Join(root, "metadata")
	wo, err := getWO(mdDir)
	if err != nil {
		panic(err)
	}

	ers, err := getERMap(wo)
	if err != nil {
		panic(err)
	}

	outputLog, err := os.Create(filepath.Join(mdDir, "er2cuid.log"))
	if err != nil {
		panic(err)
	}

	log.SetOutput(outputLog)

	for k, v := range ers {
		orig := filepath.Join(root, k)
		updated := filepath.Join(root, v)
		err := os.Rename(orig, updated)
		if err != nil {
			log.Printf("[ERROR] %s when renameing %s %s", err.Error(), orig, updated)
			panic(err)
		}
		log.Printf("[INFO] %s renamed %s", k, v)
	}
}

func getWO(mdDir string) (*string, error) {
	mdFiles, err := os.ReadDir(mdDir)
	if err != nil {
		panic(err)
	}

	for _, file := range mdFiles {
		f := file.Name()
		if strings.Contains(f, "aspace_wo.tsv") {
			wo := filepath.Join(mdDir, f)
			return &wo, nil
		}
	}
	return nil, fmt.Errorf("No Workorder Found")
}

func getERMap(wo *string) (map[string]string, error) {
	ers := map[string]string{}
	woFile, err := os.Open(*wo)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(woFile)
	scanner.Scan()
	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), "\t")
		er := strings.ReplaceAll(cols[4], "\"", "")
		cuid := strings.ReplaceAll(cols[7], "\"", "")
		ers[er] = cuid
	}
	return ers, nil
}
