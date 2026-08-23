[RU](#beepbot-russian-version)

Inspired by `funnebot` by `@Chazoshtare`



# beepbot

beepbot is a lightweight, interactive Twitch sound bot that lets your chat trigger custom sound memes, generate text-to-speech (TTS) voices in multiple languages, and apply audio effects.

> ℹ️ **2 New Effects (v1.6.0):** Added Ring Modulation (`rm`) (robotic resonance) and Tape Stop (`ts`) (smooth vinyl slow-down). Completely refactored Stutter (`st`) and Delay (`dl`) for high-speed performance.

> ℹ️ **Audio Engine Update (v1.5.0):** Added a new audio player that supports custom output device routing (e.g., virtual audio cables for OBS isolation). [Learn more](#audio-device-setup-en)

> ℹ️ **Translation (v1.4.0):** Now you can translate the text by appending the `-tr` modifier to the language code:
> * `!m en hello chat` — read the text in English.
> * `!m ru-tr hello my friend` — translate the text to Russian ("привет, мой друг") and speak it with the Russian voice.
> Note that translating (`-tr`) may expand your text, potentially cutting it off earlier. [More about limit](#tts-limit-en)

---

## Setup & Launch

1. Open `config.env` with a text editor, enter your Twitch channel name (`CHANNEL=your_channel_name`), and optionally set your starting volume (`VOLUME=100`, range 0-200).
2. Place your sound files in **`.wav`** or **`.mp3`** format (44100 Hz recommended) into the `sounds` folder. The filename (excluding the extension) automatically becomes the chat command.
3. Run the executable file.
4. When updating to a new version, you only need to replace the old `beepbot.exe` file with the new one. Do not overwrite your configured `config.env` file or the `sounds` folder to avoid losing your data.

> * **File Duration:** Use short sounds (1–10s). The bot caches all audio into RAM for instant, lag-free playback. Long music tracks will quickly overload your computer's RAM.
> * The release package already includes a `sounds` folder with a sample `beep.wav` file. You can run the bot immediately and test it in your chat using the `!m beep` command.

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
* **Sequential Chain (using spaces):** `!m sound1-sp150 en hello sound2` (plays sped-up sound1, then reads "hello" in English, and finally plays sound2).

---

## Audio Effects

Viewers can modify any sound or TTS by adding parameters separated by a hyphen `-` (order does not matter):

| Parameter | Effect | Range | Description |
| --- | --- | --- | --- |
| `sp[value]` | Speed | 10 - 200 | Playback speed and pitch (Default: `100`.`sp150` is faster, `sp50` is slower). |
| `cs[value]` | Cut start | 0 - 100 | Cuts the specified percentage of the sound from the start (e.g., `cs20`). |
| `ce[value]` | Cut end | 0 - 100 | Cuts the specified percentage of the sound from the end (e.g., `ce20`). |
| `rs` | Reverse | — | Plays the sound backward. |
| `lq` | Low Quality | — | Applies an 8-bit retro sound effect (bitcrushing). |
| `st[value]` | Stutter | 1 - 8 / 60-300 | Applies a rapid stutter effect to the beginning of the sound (Default: 3 repetitions and 140ms interval for plain `st`; customize with number of repetitions`st5` or interval `st_200` or both `st5_200`) |
| `er` | Ear Rape | — | Applies an extreme volume overdrive. |
| `dl` | Delay | — | Applies a decaying echo effect. |
| `vb` | Vibrato | — | Applies a pitch-vibrating effect. |
| `rm[value]` | Ring Modulation | 1 - 100 | Applies a metallic, robotic ring effect (Default: `50` for plain `rm`, or append a number: `rm20` - deeper, `rm80` - brighter). |
| `ts` | Tape Stop | — | Smoothly decelerates and pitches down the end of the sound (like stopping a vinyl player). |
| `ga` | Gacha | — | Randomly adds unused effects. The number of added effects depends on how many you already specified (if you have already specified 3 or more, no effects are added unless you trigger a rare 5% jackpot, which adds 1 more) |
| `tr` | Translation | — | **TTS only.** Translates the text into the target language (e.g., `ru-tr hello`). |

*(Examples: `!m ru-sp150 hello`, `!m omg-ga`)*

> ℹ️ *Note:* Trimming (cs/ce) is always applied to the original sound first, before any other effects are processed.

> ℹ️ *Note:* Tape Stop (`ts`) might work unexpectedly on long TTS messages due to Google's automatic trailing silence.

<a name="tts-limit-en"></a>

>  **TTS Length Limit (200 chars):** Due to using free web API it has a strict 200-character limit per request. To bypass this limit, chain multiple TTS commands sequentially in one message:
> * ❌ `!m ru long_text_300_chars` (Bad - will be truncated to 200 chars).
> *  `!m ru text_150_chars ru text_150_chars` (Excellent - plays fully without truncation).

---

## Admin Commands (Broadcaster & Moderators Only)

| Command | Description |
| --- | --- |
| `!m mute` / `unmute` | Mutes / unmutes the bot (instantly stops audio, clears the queue). |
| `!m qon` / `qoff` | Enables / disables sequential queue (if `qoff`, sounds will overlap concurrently). Saved automatically (config: QUEUE).|
| `!m eron` / `eroff` | Enables / disables global ear safety (strictly blocks the `er` effect). Saved automatically (config: ER). |
| `!m stop` | Instantly stops currently playing sound and clears the entire queue. |
| `!m skip` | Instantly interrupts current sound and plays the next queued item. |
| `!m vol [value]` | Sets the master volume of the bot (range: 0-200, default: 100). Saved automatically (config: VOLUME). |

***

<a name="beepbot-russian-version"></a>

# beepbot

beepbot — это легкий интерактивный Twitch-бот, который позволяет зрителям запускать звуковые мемы, озвучивать текст (TTS) на разных языках и накладывать аудиоэффекты.

> ℹ️ **2 новых эффекта (v1.6.0):** Добавлены Ring Modulation (`rm`) (роботизированный резонанс) и Tape Stop (`ts`) (плавная остановка винила). Полностью переработаны Stutter (`st`) и Delay (`dl`) для максимальной производительности.

> ℹ️ **Обновление аудио-движка (v1.5.0):** Добавлен новый плеер с поддержкой выбора устройства вывода звука (например, для изоляции дорожки бота в OBS через виртуальный кабель). [Подробнее](#audio-device-setup-ru)

> ℹ️ **Переводчик (v1.4.0):** Теперь вы можете переводить текст, добавив модификатор `-tr` к коду языка:
> * `!m en hello chat` — озвучить текст на английском.
> * `!m ru-tr hello my friend` — автоматически перевести английский текст на русский («привет, мой друг») и озвучить его русским голосом.
> Учтите, что перевод (`-tr`) может удлинить ваш текст, из-за чего он обрежется раньше. [Подробнее про лимит](#tts-limit-ru)

---

## Настройка и запуск

1. Откройте файл `config.env` текстовым редактором, впишите имя вашего Twitch-канала (`CHANNEL=имя_вашего_канала`) и, по желанию, стартовую громкость (`VOLUME=100`, диапазон 0-200).
2. Положите свои аудиофайлы в формате **`.wav`** или **`.mp3`** (рекомендуется частота 44100 Гц) в папку `sounds`. Название файла (без расширения) становится командой вызова.
3. Запустите исполняемый файл бота.
4. При выходе новой версии достаточно заменить старый файл `beepbot.exe` на новый. Не перезаписывайте уже настроенный файл `config.env` и папку `sounds`, чтобы не потерять свои данные.

> * **Длительность звуков:** Используйте короткие звуки (1–10 сек). Бот хранит аудио в ОЗУ для мгновенного воспроизведения. Длинные треки быстро перегрузят оперативную память вашего компьютера.
> * Релизный архив уже содержит папку `sounds` с тестовым файлом `beep.wav`. Вы можете сразу запустить бота и проверить его работу в чате командой `!m beep`.

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
* **Цепочка (последовательно через пробел):** `!m sound1-sp150 ru привет sound2` (сначала проиграется ускоренный sound1, затем по-русски озвучится слово «привет», а в конце запустится sound2).

---

## Доступные аудиоэффекты

Эффекты добавляются через дефис `-` после имени звука или кода языка (порядок не имеет значения):

| Параметр | Эффект | Диапазон | Описание |
| --- | --- | --- | --- |
| `sp[число]` | Скорость | 10 - 200 | Скорость и высота воспроизведения (норма: `100`. `sp150` — быстрее и выше, `sp50` — медленнее и ниже). |
| `cs[число]` | Срез начала | 0 - 100 | Отрезать указанный процент звука с начала (например, `cs20`). |
| `ce[число]` | Срез конца | 0 - 100 | Отрезать указанный процент звука с конца (например, `ce20`). |
| `rs` | Реверс | — | Воспроизвести звук задом наперед. |
| `lq` | Лоу-фай | — | Эффект 8-битного ретро-звука (биткрашинг). |
| `er` | Перегруз | — | Экстремальный перегруз громкости (Ear Rape). |
| `st[число]` | Заикание | 1 - 8 / 60 - 300 | Эффект быстрого заикания в самом начале звука (по умолчанию: 3 повторения, 140 мс для обычного `st`; настраивается через указание количества повторений `st5` или величины интервала `st_200`, или вместе `st5_200`) |
| `dl` | Эхо (Delay) | — | Эффект плавного затухающего эхо. |
| `vb` | Вибрация | — | Эффект плавного дрожания частоты (Vibrato). |
| `rm[число]` | Ring Modulation | 1 - 100 | Добавляет металлический, роботизированный резонанс (по умолчанию: `50` при вводе `rm`, либо добавьте число: `rm20`— глубже, `rm80` — ярче). |
| `ts` | Tape Stop | — | Плавно замедляет скорость и высоту тона в самом конце звука (эффект остановки виниловой пластинки). |
| `ga` | Гача (Gacha) | — | Случайно добавляет неиспользованные эффекты. Количество зависит от того, сколько эффектов уже применено к звуку (если применено 3 или более, то не добавится ничего, кроме редкого 5% шанса сорвать джекпот и получить +1 эффект). |
| `tr` | Перевод | — | **Только для TTS.** Переводит текст на указанный язык (например, `ru-tr hello`). |

*(Примеры: `!m ru-sp150 привет`, `!m omg-ga`)*

> ℹ️ *Примечание:* Обрезка (cs/ce) всегда применяется к исходному звуку первой, до наложения любых других эффектов.

> ℹ️ *Примечание:* Эффект Tape Stop (`ts`) может работать непредсказуемо на длинных ТТС-сообщениях из-за автоматического добавления тишины в конце фраз со стороны Google.

<a name="tts-limit-ru"></a>

>  **Лимит длины TTS (200 симв.):** Из-за использования бесплатного API лимит озвучки для одного куска текста — строго 200 символов. Чтобы обойти это ограничение, склеивайте команды цепочкой:
> * ❌ `!m ru-sp150 длинный_текст_300_символов` (Плохо — обрежется до 200 симв.).
> *  `!m ru-sp150 текст_150_символов ru-sp150 текст_150_символов` (Отлично — проиграется полностью).

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