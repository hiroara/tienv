package target

import (
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	configFile string
	directory  string
	path       string
}

func GetConfigWithRestore(directory, backupExt string) (c *Config) {
	c = GetConfig(directory)
	Restore(c, backupExt)
	return
}
func GetConfig(directory string) *Config {
	return &Config{configFile: "app/config.json", directory: directory}
}

func (c *Config) GetFilePath() (path string, err error) {
	if c.path != "" {
		path = c.path
		return
	}
	path, err = filepath.Abs(filepath.Join(c.directory, c.configFile))
	c.path = path
	return
}

func (c *Config) ConvertEnv(target, as string) (encoded []byte, err error) {
	path, err := c.GetFilePath()
	if err != nil {
		return
	}

	file, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(0400))
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	json, err := simplejson.NewJson(content)
	if err != nil {
		return
	}

	env, err := json.Get("env:" + as).Map()
	json.Set("env:"+target, env)

	return json.EncodePretty()
}
