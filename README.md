# typioca

Минималистичный терминальный тренажёр скорости печати.
Форк проекта [bloznelis/typioca](https://github.com/bloznelis/typioca), база —
коммит `baf4170`, лицензия MIT.

Приложение работает локально, не требует браузера и регистрации. Встроенные
наборы слов доступны без сети; сеть нужна только для загрузки дополнительных
наборов и раскладок.

## Запуск

Требования: Go 1.22.1 или новее, GNU Make, терминал с поддержкой ANSI.

```bash
git clone https://github.com/gitmaxlla/typioca-code.git
cd typioca-code/src
make build
./execs/typioca
./execs/typioca --version
```

Проверка работоспособности:

1. Запустить `./execs/typioca` — открывается главное меню.
2. Выбрать `Time run`, `15`, `Common words`, нажать Enter.
3. Ввести несколько слов и убедиться, что появился экран результатов с WPM,
   CPM и точностью.

Протокол проверки: 09.11.2026, Linux amd64, Go 1.25.6, `@gitmaxlla`:
`make build` и `--version` выполнены. Интерактивный запуск требует настоящего TTY;
проверка на другой машине пока не проводилась.

## Команды

```bash
./execs/typioca                 # локальный TUI
./execs/typioca --version       # версия без запуска TUI
./execs/typioca serve           # SSH-сервер на :2229
./execs/typioca completion bash # скрипт автодополнения в stdout
./execs/typioca help serve      # справка по SSH-серверу
```

`serve` принимает `--bind`, `--port` и `--key`; при отсутствии файла ключа
создаётся Ed25519-ключ. Клиенту нужен PTY (`ssh -t`). Встроенной
аутентификации нет, поэтому сервер нельзя публиковать в интернете без
дополнительной защиты. Контракт команды — в
[`docs/interface-contract.md`](docs/interface-contract.md).

`help` и `completion` — стандартные команды Cobra; `completion` печатает скрипт
автодополнения в stdout и ничего не записывает на диск.

## Статус требований

Полная версия — в [`docs/01-requirements.md`](docs/01-requirements.md).

| Номер | Краткая формулировка | Приоритет | Состояние |
| --- | --- | --- | --- |
| FR-1 | Тест на время | обязательно | уже есть |
| FR-2 | Тест по количеству слов | обязательно | уже есть |
| FR-3 | Тест по количеству предложений | обязательно | уже есть |
| FR-4 | Выбор набора и раскладки | обязательно | уже есть |
| FR-5 | Ввод и исправление текста | обязательно | уже есть |
| FR-6 | Просмотр результатов | обязательно | уже есть |
| FR-7 | Управление наборами и раскладками | обязательно | уже есть |
| FR-8 | Тест через SSH | желательно | уже есть |
| FR-9 | Практика на коде: наборы языков и подсветка синтаксиса | обязательно | нужно сделать |
| FR-10 | Добавление своего набора слов | обязательно | уже есть |

## Синхронизация с upstream

```bash
git remote add upstream https://github.com/bloznelis/typioca.git
git fetch upstream
git switch master
git merge upstream/master
git push origin master
```

Конфликты разрешаются вручную; после слияния проверяется `make build-all`.

## Документация

- [`docs/00-project.md`](docs/00-project.md) — паспорт проекта.
- [`docs/01-requirements.md`](docs/01-requirements.md) — требования и модели.
- [`docs/interface-contract.md`](docs/interface-contract.md) — контракт `typioca serve`.
- [`docs/diagrams/`](docs/diagrams/) — исходники и PNG диаграмм.
- [`docs/diagrams/code-map.png`](docs/diagrams/code-map.png) — карта кода для разработчика.
- [`docs/demo-1min.md`](docs/demo-1min.md) — сценарий демо на одну минуту.

## Лицензия

MIT, Copyright (c) 2022 Lukas Bloznelis. Файл [`LICENSE`](LICENSE) сохранён
без изменений.
