package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/bloznelis/typioca/cmd/words"
	"github.com/kirsle/configdir"
)

func ReadConfig() Config {
	var config Config
	configFile := getSystemConfigPath()

	//File does not exist?
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config = defaultConfig()
		WriteConfig(config)
	} else {
		readConfigFile(&config, configFile)
	}
	config = mergeConfigs(config)
	checkSync(&config)

	return config
}

func mergeConfigs(config Config) Config {
	localConfigFile := getLocalConfigPath()

	if _, err := os.Stat(localConfigFile); os.IsNotExist(err) {
	} else {
		var localConfig LocalConfig
		readLocalConfigFile(&localConfig, localConfigFile)

		config.WordLists = append(localConfig.Words, config.WordLists...)
	}

	return config
}

func checkSync(config *Config) {
	for idx, elem := range config.WordLists {
		config.WordLists[idx].synced = fileExists(elem.Path)
		config.WordLists[idx].syncOK = true
	}
}

func WriteConfig(config Config) {
	configFile := getSystemConfigPath()
	words.EnsureDir(configFile)
	fh, err := os.Create(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	var acc []WordList
	for _, elem := range config.WordLists {
		if !elem.isLocal {
			acc = append(acc, elem)
		}
	}
	config.WordLists = acc

	encoder := json.NewEncoder(fh)
	encoder.SetIndent("", "\t")
	encoder.Encode(&config)
}

func getCachePath() string {
	cachePath := configdir.LocalCache("typioca")

	err := configdir.MakePath(cachePath)
	if err != nil {
		panic(err)
	}

	return cachePath
}

func getSystemConfigPath() string {
	return getConfigPath(configdir.LocalCache("typioca"))
}

func getLocalConfigPath() string {
	return getConfigPath(configdir.LocalConfig("typioca"))
}

func getConfigPath(configDir string) string {
	err := configdir.MakePath(configDir)
	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(configDir, "typioca.conf")

	return configFile
}

func readConfigFile(config *Config, configFile string) {
	fh, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	decoder := json.NewDecoder(fh)
	decoder.Decode(&config)
}

func readLocalConfigFile(config *LocalConfig, configFile string) {
	fh, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	_, err = toml.DecodeFile(configFile, &config)

	if err != nil {
		panic(err)
	}

	for idx := range config.Words {
		config.Words[idx].isLocal = true
	}

}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func githubWordsURI(fileName string) string {
	return fmt.Sprintf("%s%s",
		"https://raw.githubusercontent.com/bloznelis/typioca/master/words/storage/words/",
		fileName,
	)
}
func githubSentencesURI(fileName string) string {
	return fmt.Sprintf("%s%s",
		"https://raw.githubusercontent.com/bloznelis/typioca/master/words/storage/sentences/",
		fileName,
	)
}

func defaultWordList(cachePath string, name string, localName string, enabled bool, sentences bool) WordList {
	var subdir string
	var uri string
	if sentences {
		subdir = "sentences"
		uri = githubSentencesURI(localName)
	} else {
		subdir = "words"
		uri = githubWordsURI(localName)
	}

	file := filepath.Join(cachePath, subdir, localName)
	return WordList{
		Sentences: sentences,
		Name:      name,
		Path:      file,
		RemoteURI: uri,
		Enabled:   enabled,
		synced:    fileExists(file),
	}
}

func defaultConfig() Config {
	cachePath := getCachePath()

	return Config{
		Version: 1,
		EmbededWordLists: []EmbededWordList{
			{"Common words", false, true},
			{"Frankenstein sentences", true, true},
		},
		WordLists: []WordList{
			defaultWordList(cachePath, "Frankenstein words", "frankenstein.json", true, false),
			defaultWordList(cachePath, "Dorian Gray words", "dorian-gray.json", true, false),
			defaultWordList(cachePath, "Pride and Prejudice words", "pride-and-prejudice.json", true, false),
			defaultWordList(cachePath, "Dorian gray sentences", "dorian-gray.json", true, true),
			defaultWordList(cachePath, "Pride and Prejudice sentences", "pride-and-prejudice.json", true, true),
		},
	}
}
