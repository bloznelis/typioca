<h1 align=center> typioca </h1>

<p align=center> Minimal, terminal based typing speed tester. </p>
<p align=center> <img alt="GitHub release (latest SemVer)" src="https://img.shields.io/github/v/release/bloznelis/typioca"> <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/bloznelis/typioca/ci.yml"> </p>

> **Tapioca** (/ˌtæpiˈoʊkə/) is a starch extracted from the storage roots of the cassava plant. Pearl tapioca is a common ingredient in Asian desserts...and sweet drinks such as **bubble tea**.

![](https://github.com/bloznelis/typioca/blob/master/img/typioca.gif)

## Features
  * Time or word/sentence count based typing speed tests
  * Proper WPM results based on https://www.speedtypingonline.com/typing-equations
  * Multiple word/sentence lists made out of classical books to spice your test up
  * Cursor aware word lines
  * Interactive menu
  * ctrl+w support
  * SSH server `typioca serve`
  * Dynamic word lists
  * Custom word lists
  * Linux/Mac/Win support

## Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/typioca.svg)](https://repology.org/project/typioca/versions)

### AUR

```
yay -S typioca-git
```

### Go

```
go install github.com/bloznelis/typioca@latest
```

**Note:** This will install typioca in `$GOBIN`, which defaults to `$GOPATH/bin` or `$HOME/go/bin` if the GOPATH environment variable is not set.

### Homebrew

```
brew install typioca
```

### Nix

```
nix-shell -p typioca
```

### openSUSE Tumbleweed

```
zypper install typioca
```

### Scoop

```
scoop bucket add extras
scoop install typioca
```

### Void Linux

```
xbps-install typioca
```

### Winget

```
winget install bloznelis.typioca
```

### X-CMD

```
x install typioca
```

### Building from source
  1. Checkout the code
  2. `make build`
  3. `./execs/typioca`

#### Prerequisites
  * `make`
  * `go`

## Custom wordlists
1. Create your word list in a new line separated manner:
```
custom
words
are
the
best
```
or in the same JSON format as the official ones:
```json
{
  "metadata" : {
    "name" : "My words",
    "size" : 5,
    "packagedAt" : "1970-01-01T00:00:00Z",
    "version" : 1
  },
  "words": [ "custom", "words", "are", "the", "best" ]
}
```
2. Place your configuration to platform specific location:

| Platform | **User configuration**                                                                     |
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
![ship it](https://user-images.githubusercontent.com/33397865/176735281-5c2b34cb-5b19-43c1-9954-92c0583c4cc5.png)

**Note:** Notice that custom wordlist controls are greyed-out, personal configuration must be handled via the file only.

---
![1](https://user-images.githubusercontent.com/33397865/176732388-11b66a1e-1d20-420f-a583-5d95241444d6.png)
![3](https://user-images.githubusercontent.com/33397865/176732403-9c64e277-f533-4bf3-96a5-a26303b37b60.png)
![2](https://user-images.githubusercontent.com/33397865/176732395-73c6c922-6a0d-4576-90bb-1f77e2c9b065.png)
![4](https://user-images.githubusercontent.com/33397865/176732415-aac89b54-15d3-4b10-8408-fac997b97085.png)

### Acknowledgments
Built with [bubbletea](https://github.com/charmbracelet/bubbletea)

🧋
