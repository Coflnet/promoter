package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func Promote() error {
	path := config.RepositoryFolder
	var promoteablePaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {
		if !IsPathPromoteable(path) {
			return nil
		}
		promoteablePaths = append(promoteablePaths, path)
		return nil
	})

	if err != nil {
		return err
	}

	if len(promoteablePaths) == 0 || len(promoteablePaths) > 1 {
		if len(promoteablePaths) > 1 {
			for i, p := range promoteablePaths {
				log.Warn().Msgf("found %d file %s", i, p)
			}
		}
		return fmt.Errorf("the amount of promoteable paths is %d, which does not work has to be exact 1", len(promoteablePaths))
	}

	return promoteFile(promoteablePaths[0])
}

func promoteFile(path string) error {
	log.Info().Msgf("%s is promoteable", path)

	tmpFile, err := CreateTempFile(path)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	oldFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer oldFile.Close()

	scanner := bufio.NewScanner(oldFile)
	for scanner.Scan() {
		err = WriteLineToTempFile(scanner.Text(), tmpFile)
		if err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	err = SwitchTempFileWithRealFile(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("%s was promoted", path)
	return nil
}

func WriteLineToTempFile(line string, newFile *os.File) error {

	line = ModifyImageTagIfPossible(line)

	_, err := newFile.WriteString(line + "\n")
	return err
}

func CreateTempFile(path string) (*os.File, error) {
	return os.Create(path + "_tmp")
}

func SwitchTempFileWithRealFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return os.Rename(path+"_tmp", path)
}

func ModifyImageTagIfPossible(line string) string {
	if !strings.Contains(line, "image:") {
		return line
	}

	if !strings.Contains(line, config.ImageName) {
		log.Warn().Msgf("line does contain \"image:\" but not the given image %s, line: %s", config.ImageName, line)
		return line
	}

	parts := strings.Split(line, ":")
	newLine := parts[0] + ":" + parts[1] + ":" + config.NewTag

	return newLine
}

func IsPathPromoteable(path string) bool {
	return strings.Contains(
		strings.ToLower(path),
		strings.ToLower(config.Filename)+".yaml",
	)
}
