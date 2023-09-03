package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Promote() error {
	paths := config.RepositoryFolders

	for folder := range paths {
		promoted, err := PromotePath(folder)

		if err != nil {
			return err
		}

		if promoted {
			break
		}
	}

	return nil
}

func PromotePath(path string) (bool, error) {
	helmPromoted := false
	var promoteablePaths []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, _ error) error {

		// already updated a helm chart
		if helmPromoted {
			return nil
		}

		// check if the current directory is a helm chart
		if isHelmChart(path) {
			correctChart, err := isCorrectHelmChart(path, config.Filename)
			if err != nil {
				return err
			}

			if correctChart {
				slog.Debug("found the promotable helm chart %s", path)
				helmPromoted = true
				return promoteHelmChart(path)
			}
		}

		// check for old school update
		if !IsPathPromoteable(path) {
			return nil
		}
		promoteablePaths = append(promoteablePaths, path)
		return nil
	})

	if err != nil {
		return false, err
	}

	if helmPromoted {
		return true, nil
	}

	if len(promoteablePaths) == 0 || len(promoteablePaths) > 1 {
		if len(promoteablePaths) > 1 {
			for i, p := range promoteablePaths {
				slog.Debug("found %d file %s", i, p)
			}
		}
		return false, nil
	}

	return true, promoteFile(promoteablePaths[0])
}

func promoteHelmChart(path string) error {

	// yaml file path
	yamlFile := filepath.Dir(path) + "/values.yaml"

	// search the tag line
	file, err := os.Open(yamlFile)
	if err != nil {
		return err
	}

	tmpFile, err := CreateTempFile(yamlFile)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "harbor.flou.dev") {
			parts := strings.Split(line, ":")
			repoParts := strings.Split(parts[1], "/")
			res := fmt.Sprintf("%s: %s/%s", parts[0], "muehlhansfl", repoParts[2])

			slog.Debug("found the repository line %s; replace with: %s", line, res)
			_, err = tmpFile.WriteString(res + "\n")
			continue
		}

		if !strings.Contains(line, "tag:") {
			_, err = tmpFile.WriteString(line + "\n")
			continue
		}

		// add a random number to the tag
		slog.Debug(fmt.Sprintf("found the tag line %s", line))
		slog.Debug(fmt.Sprintf("use the new tag: %s", config.NewTag))

		// replace the tag

		// the amount of spaces before the tag key
		// this is already a string and can be concated with the new tag
		spaces := amountOfSpacesBeforeTag(line)
		slog.Debug(fmt.Sprintf("found %d spaces before the tag key", len(spaces)))

		newLine := spaces + "tag: " + config.NewTag
		_, err = tmpFile.WriteString(newLine + "\n")
	}

	// delete the old yaml fiel and replica it with the new one
	SwitchTempFileWithRealFile(yamlFile)

	return nil
}

func amountOfSpacesBeforeTag(line string) string {
	for i, c := range line {
		if c != ' ' {
			return line[:i]
		}
	}
	return ""
}

func isHelmChart(path string) bool {
	return filepath.Base(path) == "Chart.yaml"
}

func isCorrectHelmChart(path, project string) (bool, error) {

	// read chart.yaml file
	file, err := os.Open(path)
	if err != nil {
		return true, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, "name:") {
			continue
		}

		val := strings.Split(line, ":")[1]
		val = strings.ReplaceAll(val, " ", "")
		val = strings.ReplaceAll(val, "-", "")

		sanitiedProject := strings.ReplaceAll(project, "-", "")
		sanitiedProject = strings.ReplaceAll(sanitiedProject, " ", "")

		return val == sanitiedProject, nil
	}

	return false, nil
}

func promoteFile(path string) error {
	slog.Debug(fmt.Sprintf("%s is promoteable", path))

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

	slog.Debug(fmt.Sprintf("%s was promoted", path))
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
	slog.Debug("deleting file", "file", path)
	err := os.Remove(path)
	if err != nil {
		return err
	}

	slog.Debug(fmt.Sprintf("renaming file %s to %s", path+"_tmp", path))
	return os.Rename(path+"_tmp", path)
}

func ModifyImageTagIfPossible(line string) string {

	if strings.Contains(line, "repository:") {
		parts := strings.Split(line, ":")
		repoParts := strings.Split(parts[1], "/")
		res := fmt.Sprintf("%s: %s/%s", parts[0], "muehlhansfl", repoParts[2])

		slog.Debug(fmt.Sprintf("found the repository line %s; replace with: %s", line, res))
		return res
	}

	if !strings.Contains(line, "image:") {
		return line
	}

	if !strings.Contains(line, config.ImageName) {
		slog.Warn(fmt.Sprintf("line does contain \"image:\" but not the given image %s, line: %s", config.ImageName, line))
		return line
	}

	parts := strings.Split(line, ":")
	newLine := parts[0] + ":" + parts[1] + ":" + config.NewTag

	return newLine
}

func IsPathPromoteable(path string) bool {
	path = strings.ToLower(path)
	path = strings.Replace(path, "-", "", -1)
	comp := strings.ToLower(config.Filename) + ".yaml"
	comp = strings.Replace(comp, "-", "", -1)

	return strings.Contains(
		path,
		comp,
	)
}
