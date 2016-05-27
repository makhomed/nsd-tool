package zonelist

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/makhomed/nsd-tool/config"
)

func Generate(conf *config.Config, pattern string, filename string) error {
	files, err := ioutil.ReadDir(conf.ZoneDir)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	prefix := filepath.Base(filename)
	newZoneList, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return err
	}
	defer func(newZoneListName string) {
		_, err := os.Stat(newZoneListName)
		if os.IsNotExist(err) {
			return
		}
		os.Remove(newZoneListName)
	}(newZoneList.Name())

	for _, file := range files {
		_, err = io.WriteString(newZoneList, "\n")
		if err != nil {
			return err
		}
		_, err = io.WriteString(newZoneList, "zone:\n")
		if err != nil {
			return err
		}
		_, err = io.WriteString(newZoneList, "    name: \"" + file.Name() + "\"\n")
		if err != nil {
			return err
		}
		_, err = io.WriteString(newZoneList, "    include-pattern: \"" + pattern + "\"\n")
		if err != nil {
			return err
		}
	}
	newZoneList.Close()
	var mode os.FileMode
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		mode = 0644
	} else {
		mode = stat.Mode()
	}
	os.Chmod(newZoneList.Name(), mode)
	err = os.Rename(newZoneList.Name(), filename)
	if err != nil {
		return err
	}
	return nil
}
