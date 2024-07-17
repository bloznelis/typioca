package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/bloznelis/typioca/cmd/words"
	"github.com/kirsle/configdir"
)

const currentConfigVersion = 4

func ReadConfig() Config {
	var config Config
	configFile := getSystemConfigPath()

	//File does not exist?
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		config = defaultConfig()
		WriteConfig(config)
	} else {
		readConfigFile(&config, configFile)
		if config.Version != currentConfigVersion {
			config = defaultConfig()
			WriteConfig(config)
		}
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

	for idx, elem := range config.LayoutFiles {
		config.LayoutFiles[idx].synced = fileExists(elem.Path)
		config.LayoutFiles[idx].syncOk = true
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

func githubLayoutsURI(fileName string) string {
	return fmt.Sprintf("%s%s",
		"https://raw.githubusercontent.com/bloznelis/typioca/master/layouts/",
		fileName,
	)
}

func retrieveLayout(layout LayoutFile) Layout {
	if layout.Path == "" {
		return Layout{
			Name: layout.Name,
		}
	}

	f, err := os.Open(layout.Path)
	if err != nil {
		log.Println(layout.Path)
		panic(err)
	}

	var res Layout
	dec := json.NewDecoder(f)
	if err := dec.Decode(&res); err != nil {
		panic(err)
	}

	return res
}

func defaultLayoutFile(cachePath string, name string, localName string) LayoutFile {
	if localName == "" {
		return LayoutFile{
			Name: name,
		}
	}

	path := filepath.Join(cachePath, "layouts", localName)

	return LayoutFile{
		Name:      name,
		Path:      path,
		RemoteURI: githubLayoutsURI(localName),
		synced:    fileExists(path),
	}
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
	defaultLayout := defaultLayoutFile(cachePath, "Qwerty", "")

	return Config{
		TestSettingCursors: initTestSettingCursors(),
		Version:            currentConfigVersion,
		EmbededWordLists: []EmbededWordList{
			{"Common words", false, true},
			{"Frankenstein sentences", true, true},
		},
		WordLists: []WordList{
			defaultWordList(cachePath, "Frankenstein words", "frankenstein.json", true, false),

			defaultWordList(cachePath, "Dorian Gray words", "dorian-gray.json", true, false),
			defaultWordList(cachePath, "Dorian gray sentences", "dorian-gray.json", true, true),

			defaultWordList(cachePath, "Pride and Prejudice words", "pride-and-prejudice.json", true, false),
			defaultWordList(cachePath, "Pride and Prejudice sentences", "pride-and-prejudice.json", true, true),

			defaultWordList(cachePath, "Sherlock Holmes words", "sherlock-holmes.json", true, false),
			defaultWordList(cachePath, "Sherlock Holmes sentences", "sherlock-holmes.json", true, true),

			defaultWordList(cachePath, "Dracula words", "dracula.json", true, false),
			defaultWordList(cachePath, "Dracula sentences", "dracula.json", true, true),

			defaultWordList(cachePath, "The Yellow Wallpaper words", "the-yellow-wallpaper.json", true, false),
			defaultWordList(cachePath, "The Yellow Wallpaper sentences", "the-yellow-wallpaper.json", true, true),

			defaultWordList(cachePath, "A Tale of Two Cities words", "a-tale-of-two-cities.json", true, false),
			defaultWordList(cachePath, "A Tale of Two Cities sentences", "a-tale-of-two-cities.json", true, true),

			defaultWordList(cachePath, "The Great Gatsby words", "the-great-gatsby.json", true, false),
			defaultWordList(cachePath, "The Great Gatsby sentences", "the-great-gatsby.json", true, true),

			defaultWordList(cachePath, "The Count of Monte Cristo words", "the-count-of-monte-cristo.json", true, false),
			defaultWordList(cachePath, "The Count of Monte Cristo sentences", "the-count-of-monte-cristo.json", true, true),

			defaultWordList(cachePath, "Treasure Island words", "treasure-island.json", true, false),
			defaultWordList(cachePath, "Treasure Island sentences", "treasure-island.json", true, true),

			defaultWordList(cachePath, "Little Women words", "little-women.json", true, false),
			defaultWordList(cachePath, "Little Women sentences", "little-women.json", true, true),

			defaultWordList(cachePath, "Peter Pan words", "peter-pan.json", true, false),
			defaultWordList(cachePath, "Peter Pan sentences", "peter-pan.json", true, true),
		},
		LayoutFiles: []LayoutFile{
			defaultLayoutFile("", "Qwerty", ""),
			defaultLayoutFile(cachePath, "Colemak DH", "colemak-dh.json"),
			defaultLayoutFile(cachePath, "Gallium", "gallium.json"),
		},
		Layout: retrieveLayout(defaultLayout),
	}
}
