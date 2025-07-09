package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	timeout            = 5 * time.Minute
	dateFormat         = "2006-01-02T15:04:05Z07:00"
	lastUpdateFileName = "owner-last-update"
	lastUpdate         = "2018-10-29T15:00:00Z" + "\n"
)

// TimestampUpdater updates the timestamp since the last update.
type TimestampUpdater struct {
	path      string
	timestamp string
}

// NewTimestampUpdater writes last update file to path if it not exists and
// returns a struct that uses this file as a data source.
func NewTimestampUpdater(path string) (*TimestampUpdater, error) {
	timestampUpdater := &TimestampUpdater{path: path}
	pathToFile := filepath.Join(timestampUpdater.path, lastUpdateFileName)
	if err := initFile(pathToFile, []byte(lastUpdate)); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}
	timestampUpdater.timestamp = string(b)
	return timestampUpdater, nil
}

// Update writes the recent time to last update file if timeout is exceeded.
func (tu *TimestampUpdater) Update() error {
	dateString := strings.TrimSuffix(tu.timestamp, "\n")
	loaded, err := time.Parse(dateFormat, dateString)
	if err != nil {
		return err
	}
	now := time.Now()
	afterTimeout := now.After(loaded.Add(timeout))
	if afterTimeout {
		err := os.WriteFile(filepath.Join(tu.path, lastUpdateFileName), []byte(now.Format(dateFormat)+"\n"), 0o644)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("timeout is set to %v to relieve server load, try in %v again",
		timeout, -(now.Sub(loaded) - timeout).Round(time.Second))
}
