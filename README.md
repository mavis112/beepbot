[RU](#beepbot-russian-version)

Inspired by `funnebot` by `@Chazoshtare`



# beepbot

beepbot is an interactive Twitch sound bot that plays custom sounds, generates TTS voices in multiple languages, and applies audio effects.

> [!NOTE]
> **Support for funnebot command syntax (v1.6.3):** Added full support for funnebot command syntax. Also, a protective hyphen rule has been introduced for the original syntax to prevent accidental sound playback inside TTS text. [Read more](#hyphen-rule-en)

> [!NOTE]
> **2 New Effects (v1.6.0):** Added Ring Modulation (`rm`) and Tape Stop (`ts`). Completely refactored Stutter (`st`) and Delay (`dl`).

---

## Setup & Launch

1. Open the `config.env` file with a text editor, enter your Twitch channel name (`CHANNEL=your_channel_name`).
2. Place your sound files in **`.wav`** or **`.mp3`** format (44100 Hz recommended) into the `sounds` folder. The filename (excluding the extension) automatically becomes the chat command.
3. Run the executable file of the bot.
4. When updating to a new version, you only need to replace the old `beepbot.exe` file with the new one. Do not overwrite your configured `config.env` file or the `sounds` folder to avoid losing your data.

> [!CAUTION]
> **File Duration:** Use short sounds (1–10s). Long sounds will quickly overload your computer's RAM.

> [!TIP]
> The release package already includes a `sounds` folder with a sample `beep.wav` file. You can run the bot immediately and test it in your chat using the `!m beep` command.

---

<a name="audio-device-setup-en"></a>
## Audio Device Selection (v1.5.0)

Starting from version 1.5.0, a new optional parameter `AUDIO_DEVICE` is added to `config.env`. You can use it to specify the audio output device where the bot's sounds will be routed.

* **Default:** If the `AUDIO_DEVICE` parameter is empty or entirely missing from `config.env`, the system default playback device will be used.
* **List Devices:** Set `AUDIO_DEVICE=list` to print the list of all available system audio outputs to the console on startup.
* **Partial Matching:** You do not need to write the exact name. The search is case-insensitive and matches substrings (e.g., `AUDIO_DEVICE="Cable"` will successfully connect `"Cable Input (VB-Audio Virtual Cable)"`).
* **Auto-Fallback:** If the configured device is not found, the bot prints a warning and automatically falls back to the system default device.
* **Legacy Engine:** If you experience any audio issues on Windows or are running the bot on Linux, set `AUDIO_DEVICE="oto"` to force the stable legacy audio engine.

---

## Chat Commands

> [!NOTE]
> Starting from version 1.6.3, you can write commands using the `funnebot` syntax (including sequential chains via `+!` and effects like `f`, `s`, `c`, `sk`, `r`, etc.). Effects not present in `beepbot` (currently `rv` and `vl`) will be ignored.

The main command for viewers is:
`!m [sound_name_or_language_code]-[effects]`

* `!m rand` — play a random sound from the `sounds` folder.

### 1. Text-To-Speech (TTS)
Specify the language code before the text you want to read:
* `!m en hello chat` — read the text in English.
* `!m jp ohayo` — read the text in Japanese.

[Full list of supported languages](tts/languages.go)

### 2. Combining Sounds & Speech
* **Simultaneous Mix (using `+`):** `!m sound1+sound2-rs` (both sounds will play at the exact same time, reversed).
* **Sequential Chain (using spaces):** `!m sound1-sp150 en hello sound2-` (plays sped-up sound1, then reads "hello" in English, and finally plays sound2).

<a name="hyphen-rule-en"></a>
## Hyphen Rule for Commands inside TTS

To prevent common words inside TTS text (like `a`, `yes`, `no`) from accidentally triggering sounds of the same name, a simple rule now applies:

After a language code, any command (sound or language switch) will only trigger if it contains a hyphen (`-`).

> [!IMPORTANT]
> * If you need a plain sound without effects after TTS text, append a hyphen at the end: use `alert-` instead of `alert`.
> * If a sound has effects (like `alert-sp150`), no changes are needed.

Examples:

Suppose we have sounds `car` and `alert`:

- `!m en I have a nice car` — The bot speaks the sentence normally (the word "car" is spoken, the sound "car" is not triggered).
- `!m en have a nice car-` — Speaks the text up to the word car, then plays the `car` sound.
- `!m en have a nice car- alert` — Speaks the text up to the word car, then plays the `car` and `alert` sounds sequentially (even though the `alert` sound has no effects, adding a hyphen is not required since a sound was played before it, not TTS).

---

## Audio Effects

Viewers can modify any sound or TTS by adding parameters separated by a hyphen `-` (order does not matter):

| Parameter | Effect | Range | Description |
| --- | --- | --- | --- |
| `sp[value]` | Speed | 10 - 200 | Playback speed and pitch (Default: `100`. `sp150` is faster, `sp50` is slower). |
| `cs[value]` | Cut Start | 0 - 100 | Cuts the specified percentage of the sound from the start (e.g., `cs20`). |
| `ce[value]` | Cut End | 0 - 100 | Cuts the specified percentage of the sound from the end (e.g., `ce20`). |
| `rs` | Reverse | — | Plays the sound backward. |
| `lq` | Low Quality | — | Degrades the sound quality through bitcrushing. |
| `er` | Ear Rape | — | Extreme volume overdrive. |
| `st[value]` | Stutter | 1 - 8 / 60 - 300 | Rapid stutter effect at the beginning of the sound (Default: 3 repetitions, 140ms interval for plain `st`; customize with repetitions `st5`, interval `st_200`, or both `st5_200`). |
| `dl` | Delay | — | Applies a decaying echo effect. |
| `vb` | Vibrato | — | Applies a pitch-vibrating effect. |
| `rm[value]` | Ring Modulation | 1 - 100 | Applies a metallic, robotic ring effect (Default: `50` for plain `rm`, or append a number: `rm20` - deeper, `rm80` - brighter). |
| `ts` | Tape Stop | — | Smoothly decelerates and pitches down the end of the sound (like stopping a cassette player). |
| `ga` | Gacha | — | Randomly adds unused effects. The number of added effects depends on how many you already specified (if you have already specified 3 or more, no effects are added unless you trigger a rare 5% jackpot, which adds 1 more). |
| `tr` | Translate | — | **TTS only.** Translates the text into the target language (e.g., `ru-tr hello`). |

*(Examples: `!m ru-sp150 hello`, `!m omg-ga`)*

> [!NOTE]
> Trimming (`cs`/`ce`) is always applied to the original sound first, before any other effects are processed.

> [!WARNING]
> Tape Stop (`ts`) might work unexpectedly on long TTS messages due to Google's automatic trailing silence.

<a name="tts-limit-en"></a>

> [!IMPORTANT]
> **TTS Length Limit (200 chars):** Due to using free web API it has a strict 200-character limit per request. To bypass this limit, chain multiple TTS commands sequentially in one message:
> * ❌ `!m ru long_text_300_chars` (Bad - will be truncated to 200 chars).
> * `!m ru text_150_chars ru- text_150_chars` (Excellent - plays fully without truncation).

---

## Admin Commands (Broadcaster & Moderators Only)

| Command | Description |
| --- | --- |
| `!m mute` / `unmute` | Mutes / unmutes the bot (instantly stops audio, clears the queue). |
| `!m qon` / `qoff` | Enables / disables sequential queue (if `qoff`, sounds will overlap concurrently). Saved automatically (config: QUEUE).|
| `!m eron` / `erfoot` | Enables / disables global ear safety (strictly blocks the `er` effect). Saved automatically (config: ER). |
| `!m stop` | Instantly stops currently playing sound and clears the entire queue. |
| `!m skip` | Instantly interrupts current sound and plays the next queued item. |
| `!m vol [value]` | Sets the master volume of the bot (range: 0-200, default: 100). Saved automatically (config: VOLUME). |


***

<a name="beepbot-russian-version"></a>

# beepbot

beepbot — это интерактивный Twitch-бот, который позволяет проигрывать звуковые файлы, озвучивать текст на разных языках и накладывать аудиоэффекты.

> [!NOTE]
> **Поддержка синтаксиса команд funnebot (v1.6.3):** Добавлена полная поддержка синтаксиса команд funnebot. Также для оригинального синтаксиса введено защитное правило дефиса для предотвращения ложного запуска звуков внутри текста TTS. [Подробнее](#hyphen-rule-ru)

> [!NOTE]
> **2 новых эффекта (v1.6.0):** Добавлены Ring Modulation (`rm`) и Tape Stop (`ts`). Полностью переработаны Stutter (`st`) и Delay (`dl`).

---

## Настройка и запуск

1. Откройте файл `config.env` текстовым редактором, впишите имя вашего Twitch-канала (`CHANNEL=имя_вашего_канала`).
2. Положите свои аудиофайлы в формате **`.wav`** или **`.mp3`** (рекомендуется частота 44100 Гц) в папку `sounds`. Название файла (без расширения) становится командой вызова.
3. Запустите исполняемый файл бота.
4. При выходе новой версии достаточно заменить старый файл `beepbot.exe` на новый. Не перезаписывайте уже настроенный файл `config.env` и папку `sounds`, чтобы не потерять свои данные.

> [!CAUTION]
> **Длительность звуков:** Используйте короткие звуки (1–10 сек). Длинные звуки быстро перегрузят оперативную память вашего компьютера.

> [!TIP]
> Релизный архив уже содержит папку `sounds` с тестовым файлом `beep.wav`. Вы можете сразу запустить бота и проверить его работу в чате командой `!m beep`.

---

<a name="audio-device-setup-ru"></a>
## Выбор аудиоустройства (v1.5.0)

Начиная с версии 1.5.0 в `config.env` добавляется новый параметр `AUDIO_DEVICE`, в котором опционально можно указать аудиоустройство, в которое будут направляться звуки бота.

* **По умолчанию:** Если параметр `AUDIO_DEVICE` пустой или вообще отсутствует в файле `config.env`, то будет использовано системное устройство по умолчанию.
* **Список устройств:** Установите `AUDIO_DEVICE=list`, чтобы вывести список всех доступных аудиовыходов в консоль при запуске бота.
* **Частичное совпадение:** Не обязательно писать имя целиком. Поиск нечувствителен к регистру и ищет подстроку (например, `AUDIO_DEVICE="Cable"` успешно подключит `"Cable Input (VB-Audio Virtual Cable)"`).
* **Авто-откат:** Если указанное устройство не найдено, бот выведет предупреждение и автоматически переключится на системный динамик по умолчанию.
* **Старый движок:** Если у вас возникли проблемы со звуком на Windows или вы запускаете бота на Linux, установите `AUDIO_DEVICE="oto"`, чтобы принудительно запустить стабильный старый аудио-движок.

---

## Синтаксис команд в чате

> [!NOTE]
> Начиная с версии 1.6.3, вы можете писать команды в синтаксисе `funnebot` (включая последовательные цепочки через `+!` и эффекты `f`, `s`, `c`, `sk`, `r` и т.д.). Эффекты, которых нет в `beepbot` (на данный момент это `rv` и `vl`), будут проигнорированы.

Основная команда для зрителей:
`!m [имя_звука_или_код_языка]-[эффекты]`

* `!m rand` — проиграть случайный звук из папки `sounds`.

### 1. Озвучка текста (TTS)
Укажите код языка перед текстом, который нужно озвучить:
* `!m ru привет чат` — озвучить текст на русском.
* `!m jp аниме` — озвучить текст на японском.

[Полный список поддерживаемых языков](tts/languages.go)

### 2. Комбинирование (Миксы и Цепочки)
* **Микс (одновременно через `+`):** `!m sound1+sound2-rs` (звуки запустятся одновременно и оба проиграются реверсом).
* **Цепочка (последовательно через пробел):** `!m sound1-sp150 ru привет sound2-` (сначала проиграется ускоренный sound1, затем по-русски озвучится слово «привет», а в конце запустится sound2).

<a name="hyphen-rule-ru"></a>
## Правило дефиса для команд внутри TTS

Чтобы обычные слова внутри текста TTS (например, `a`, `yes`, `no`) случайно не запускали одноименные звуковые файлы, теперь действует простое правило:

После кода языка любая команда (звук или голосовая озвучка) сработает, только если она содержит дефис (`-`).

> [!IMPORTANT]
> * Если вам нужен звук без эффектов после текста TTS, добавьте дефис на конце: `alert-` вместо `alert`.
> * Если у звука есть эффекты (например, `alert-sp150`), ничего менять не нужно.

Примеры:

Допустим у нас есть звуки `car` и `alert`:

- `!m en I have a nice car` — Бот просто озвучит фразу целиком.
- `!m en have a nice car-` — Озвучит фразу до слова car, а затем сыграет звук car.
- `!m en have a nice car- alert` — Озвучит фразу до слова car, затем сыграет последовательно звуки car и alert (хоть к звуку `alert` не применено эффектов, добавлять дефис не обязательно, т.к. до него играл звук, а не ттс).

---

## Доступные аудиоэффекты

Эффекты добавляются через дефис `-` после имени звука или кода языка (порядок не имеет значения):

| Параметр | Эффект | Диапазон | Описание |
| --- | --- | --- | --- |
| `sp[число]` | Speed | 10 - 200 | Скорость и высота воспроизведения (норма: `100`. `sp150` — быстрее и выше, `sp50` — медленнее и ниже). |
| `cs[число]` | Cut Start | 0 - 100 | Отрезать указанный процент звука с начала (например, `cs20`). |
| `ce[число]` | Cut End | 0 - 100 | Отрезать указанный процент звука с конца (например, `ce20`). |
| `rs` | Reverse | — | Воспроизвести звук задом наперед. |
| `lq` | Low Quality | — | Ломает качество звука посредством биткрашинга. |
| `er` | Ear Rape | — | Экстремальный перегруз громкости. |
| `st[число]` | Stutter | 1 - 8 / 60 - 300 | Эффект быстрого заикания в самом начале звука (по умолчанию: 3 повторения, 140 мс для обычного `st`; настраивается через указание количества повторений `st5` или величины интервала `st_200`, или всё вместе `st5_200`). |
| `dl` | Delay | — | Эффект плавного затухающего эхо. |
| `vb` | Vibrato | — | Эффект плавного дрожания частоты. |
| `rm[число]` | Ring Modulation | 1 - 100 | Добавляет металлический, роботизированный резонанс (по умолчанию: `50` при вводе `rm`, либо добавьте число: `rm20` — глубже, `rm80` — ярче). |
| `ts` | Tape Stop | — | Плавно замедляет скорость и высоту тона в самом конце звука (эффект остановки кассетного плеера). |
| `ga` | Gacha | — | Случайно добавляет неиспользованные эффекты. Количество зависит от того, сколько эффектов уже применено к звуку (если применено 3 или более, то не добавится ничего, кроме редкого 5% шанса сорвать джекпот и получить +1 эффект). |
| `tr` | Translate | — | **Только для TTS.** Переводит текст на указанный язык (например, `ru-tr hello`). |

*(Примеры: `!m ru-sp150 привет`, `!m omg-ga`)*

> [!NOTE]
> Обрезка (cs/ce) всегда применяется к исходному звуку первой, до наложения любых других эффектов.

> [!WARNING]
> Эффект Tape Stop (`ts`) может работать непредсказуемо на длинных ТТС-сообщениях из-за автоматического добавления тишины в конце фраз со стороны Google.

<a name="tts-limit-ru"></a>

> [!IMPORTANT]
> **Лимит длины TTS (200 симв.):** Из-за использования бесплатного API лимит озвучки для одного куска текста — строго 200 символов. Чтобы обойти это ограничение, склеивайте команды цепочкой:
> * ❌ `!m ru-sp150 длинный_текст_300_символов` (Плохо — обрежется до 200 симв.).
> * `!m ru-sp150 текст_150_символов ru-sp150 текст_150_символов` (Отлично — проиграется полностью).

---

## Команды модерирования (Для Стримера и Модераторов)

| Команда | Описание |
| --- | --- |
| `!m mute` / `unmute` | Заглушить / включить бота (при mute текущие звуки обрываются, очередь очищается). |
| `!m qon` / `qoff` | Включить / выключить очередь (при `qoff` звуки в чате накладываются параллельно). Сохраняется автоматически (config: QUEUE). |
| `!m eron` / `eroff` | Включить / выключить глобальную безопасность ушей (блокирует эффект `er` для всех). Сохраняется автоматически (config: ER). |
| `!m stop` | Прервать текущий звук и полностью очистить очередь. |
| `!m skip` | Прервать текущий звук и запустить следующий из очереди. |
| `!m vol [число]` | Устанавливает общую громкость бота (диапазон: 0-200, норма: 100). Сохраняется автоматически (config: VOLUME). |