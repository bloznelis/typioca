# typioca
Minimal, terminal based typing speed tester.


> **Tapioca** (/ˌtæpiˈoʊkə/) is a starch extracted from the storage roots of the cassava plant. Pearl tapioca is a common ingredient in Asian desserts...and sweet drinks such as **bubble tea**.

![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/bloznelis/typioca)
---

![](https://github.com/bloznelis/typioca/blob/master/img/typioca.gif)

## Features
  * Time or word/sentence count based typing speed tests
  * Proper WPM results based on https://www.speedtypingonline.com/typing-equations
  * Multiple word/sentence lists made out of classical books to spice your test up
  * Cursor aware word lines
  * Interactive menu
  * ctrl+w support ;)
  * SSH server `typioca serve`
  * Dynamic word lists
  * Custom word lists
  * Linux/Mac/Win support
  
## Installation
### AUR
`yay -S typioca-git`

### Homebrew
1. `brew tap bloznelis/tap`
2. `brew install typioca`

### Go
`go install github.com/bloznelis/typioca@latest`

**Note:** This will install typioca in `$GOBIN`, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the GOPATH environment variable is not set.

### Building from source
  1. Checkout the code
  2. `make build`
  3. `./execs/typioca`

#### Prerequisites
  * `make`
  * `go`

## Custom wordlists
1. Create your word list in the same JSON format as the official ones [example](https://raw.githubusercontent.com/bloznelis/typioca/master/words/storage/words/common-english.json).
   - **Note:** for new-line separated word lists (like [this one](https://raw.githubusercontent.com/powerlanguage/word-lists/master/1000-most-common-words.txt)), for your convenience, you can use [this Clojure script](https://github.com/bloznelis/typioca/blob/master/words/common-word-list.clj). Explanation how to use it can be found [here](https://github.com/bloznelis/typioca/tree/master/words).
3. Place your configuration to platform specific location:

| Platfrom | **User configuration**                                                                     |
|----------|--------------------------------------------------------------------------------------------|
| Windows  | `%APPDATA%\typioca\typioca.conf` or `C:\Users\%USER%\AppData\Roaming\typioca\typioca.conf` |
| Linux    | `$XDG_CONFIG_HOME/typioca/typioca.conf` or `$HOME/.config/typioca/typioca.conf`            |
| macOS    | `$HOME/Library/Application Support/typioca/typioca.conf`                                   |

Config example (it is [TOML](https://github.com/toml-lang/toml)):
```toml
[[words]]
  name      = "Best hits '22"
  enabled   = false
  sentences = false
  path      = "/home/words/best-hits-22.json"
[[words]]
  name      = "Even better hits '23"
  enabled   = true
  sentences = false
  path      = "/home/words/better-hits-23.json"
```
3. Use your words!

![image](https://user-images.githubusercontent.com/33397865/176517001-967c10f4-0489-451c-a140-ea2aae839aa5.png)

**Note:** Notice that custom wordlist controls are greyed-out, personal configuration must be handled via the file only.

---
![full-menu-cropped](https://user-images.githubusercontent.com/33397865/172426966-d1295987-4df3-4681-a651-b01f3f80be42.png)
![full-test-cropped](https://user-images.githubusercontent.com/33397865/172427152-b71979e4-8c67-4427-98e0-116c6518071f.png)
![full-results-cropped](https://user-images.githubusercontent.com/33397865/172427164-b19f1bb5-43a7-47d6-a833-c343e519f447.png)

### Acknowledgments
Built with [bubbletea](https://github.com/charmbracelet/bubbletea)

🧋
