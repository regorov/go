// Package logscanner provides parsing of Minecraft server logs.
package logscanner

import (
    "github.com/kierdavis/go/bettererrors"
    "regexp"
    "strconv"
    "time"
)

// Regular expressions
var (
    LogLine = regexp.MustCompile("^(\\d{4})-(\\d{2})-(\\d{2}) (\\d{2}):(\\d{2}):(\\d{2}) \\[(.+?)\\] (.+)$")

    StartingServer               = regexp.MustCompile("^Starting minecraft server version (.+)$")
    LoadingProperties            = regexp.MustCompile("^Loading properties$")
    ServerPropertiesDoesNotExist = regexp.MustCompile("^server.properties does not exist$")
    GeneratingNewPropertiesFile  = regexp.MustCompile("^Generating new properties file$")
    StartingServerAddress        = regexp.MustCompile("^Starting Minecraft server on (.+)$")
    FailedToLoadFile             = regexp.MustCompile("^Failed to load (.+): (.+): (.+)$")
    PreparingLevel               = regexp.MustCompile("^Preparing level \"(.+)\"$")
    DefaultGameMode              = regexp.MustCompile("^Default game type: (\\d+)$")
    PreparingStartRegion         = regexp.MustCompile("^Preparing start region for level (\\d+)(?: \\(Seed: ([-0-9]+)\\))?$")
    PreparingSpawnAreaProgress   = regexp.MustCompile("^Preparing spawn area: (\\d+)%$")
    InitializationDone           = regexp.MustCompile("^Done \\(([0-9.]+)s\\)! For help, type \"help\" or \"\\?\"$")
    TickSyncMessage              = regexp.MustCompile("^Can't keep up! Did the system time change, or is the server overloaded\\?$")
    ServerStopIssued             = regexp.MustCompile("^(.+): Stopping the server..$")
    ServerStopping               = regexp.MustCompile("^Stopping server$")
    SavingChunks                 = regexp.MustCompile("^Saving chunks$")
    PlayerConnect                = regexp.MustCompile("^(.+) \\[/([0-9.:]+)\\] logged in with entity id (\\d+) at \\((?:\\[(.+)\\] )?([-0-9.]+), ([-0-9.]+), ([-0-9.]+)\\)$")
    GameModeChanged              = regexp.MustCompile("^(.+): Setting (.+) to game mode (\\d+)$")
    PlayerOpped                  = regexp.MustCompile("^(.+): Opping (.+)$")
    PlayerDeOpped                = regexp.MustCompile("^(.+): De-opping (.+)$")
    PlayerIssuedCommand          = regexp.MustCompile("^(.+) issued server command: (.+)$")
    PlayerOldChat                = regexp.MustCompile("^<(.+)> (.+)$")
    PlayerDisconnect             = regexp.MustCompile("^(.+) lost connection(?:: disconnect\\.(.+))?$")
    PlayerTeleport               = regexp.MustCompile("^(.+): Teleporting (.+) to (.+).")
    CraftBukkitVersionInfo       = regexp.MustCompile("^This server is running CraftBukkit version (.+) \\(MC: (.+)\\) \\(Implementing API version (.+)\\)$")
    PluginMessage                = regexp.MustCompile("^\\[(.+)\\] (.+)$")
    FolderMigrationBegan         = regexp.MustCompile("^---- Migration of old (.+) folder required ----$")
    FolderMigrationComplete      = regexp.MustCompile("^---- Migration of old (.+) folder complete ----$")
)

// Errors
var (
    LineNotMatched               = bettererrors.New("The following line could not be interpreted:\n%s")
    UnrecognisedLogLevel         = bettererrors.New("The log level '%s' was not recognised.")
    UnrecognisedFailedToLoadFile = bettererrors.New("The term '%s' was not recognised.")
    UnrecognisedDisconnectReason = bettererrors.New("The disconnect reason '%s' was not recognised.")
)

// Interface LineReader represents a reader that has a ReadLine method, such as a bufio.Reader.
type LineReader interface {
    ReadLine() ([]byte, bool, error)
}

// Type LogScanner represents a log scanner.
type LogScanner struct {
    source    LineReader
    timezone  *time.Location
    lastLevel LogLevel
    lastDate  time.Time
}

// Function NewLogScanner creates and returns a new log scanner.
func NewLogScanner(source LineReader, timezone *time.Location) (ls *LogScanner) {
    return &LogScanner{
        source:   source,
        timezone: timezone,
    }
}

// Function Source returns the log scanner's source.
func (ls *LogScanner) Source() (source LineReader) {
    return ls.source
}

// Function SetSource sets the source of the log scanner.
func (ls *LogScanner) SetSource(source LineReader) {
    ls.source = source
}

// Function Timezone returns the log scanner's timezone.
func (ls *LogScanner) Timezone() (timezone *time.Location) {
    return ls.timezone
}

// Function SetTimezone sets the log scanner's timezone.
func (ls *LogScanner) SetTimezone(timezone *time.Location) {
    ls.timezone = timezone
}

// Function ReadLine reads and returns a complete line from the source.
func (ls *LogScanner) ReadLine() (line string, err error) {
    byteLine, isPrefix, err := ls.source.ReadLine()
    if err != nil {
        return "", err
    }

    line = string(byteLine)

    var nextLine []byte

    for isPrefix {
        nextLine, isPrefix, err = ls.source.ReadLine()
        if err != nil {
            return "", err
        }

        line += string(nextLine)
    }

    return line, nil
}

// Function ReadLogMessage reads and returns a message from the log.
func (ls *LogScanner) ReadLogMessage() (msg *LogMessage, err error) {
    line, err := ls.ReadLine()
    if err != nil {
        return nil, err
    }

    parts := LogLine.FindStringSubmatch(line)
    if parts == nil {
        //return nil, LineNotMatched.Format(line)
        // We can recover.

        return &LogMessage{
            message: line,
            level:   ls.lastLevel,
            date:    ls.lastDate,
        }, nil

    }

    year, err := strconv.ParseInt(parts[1], 10, 0)
    if err != nil {
        return nil, err
    }

    month, err := strconv.ParseInt(parts[2], 10, 0)
    if err != nil {
        return nil, err
    }

    day, err := strconv.ParseInt(parts[3], 10, 0)
    if err != nil {
        return nil, err
    }

    hour, err := strconv.ParseInt(parts[4], 10, 0)
    if err != nil {
        return nil, err
    }

    min, err := strconv.ParseInt(parts[5], 10, 0)
    if err != nil {
        return nil, err
    }

    sec, err := strconv.ParseInt(parts[6], 10, 0)
    if err != nil {
        return nil, err
    }

    date := time.Date(int(year), time.Month(month), int(day), int(hour), int(min), int(sec), 0, ls.timezone)

    var level LogLevel

    switch parts[7] {
    case "INFO":
        level = Info
    case "WARNING":
        level = Warning
    case "SEVERE":
        level = Severe
    default:
        return nil, UnrecognisedLogLevel.Format(parts[7])
    }

    message := parts[8]

    ls.lastLevel = level
    ls.lastDate = date

    return &LogMessage{
        message: message,
        level:   level,
        date:    date,
    }, nil
}

// Function ReadEvent reads and returns a logged event.
func (ls *LogScanner) ReadEvent() (event Event, err error) {
    msg, err := ls.ReadLogMessage()
    if err != nil {
        return nil, err
    }

    message := msg.message

    if match := StartingServer.FindStringSubmatch(message); match != nil {
        event = &StartingServerEvent{baseEvent{msg.date}, match[1]}

    } else if match := LoadingProperties.FindStringSubmatch(message); match != nil {
        event = &LoadingPropertiesEvent{baseEvent{msg.date}}

    } else if match := ServerPropertiesDoesNotExist.FindStringSubmatch(message); match != nil {
        event = &ServerPropertiesDoesNotExistEvent{baseEvent{msg.date}}

    } else if match := GeneratingNewPropertiesFile.FindStringSubmatch(message); match != nil {
        event = &GeneratingNewPropertiesFileEvent{baseEvent{msg.date}}

    } else if match := StartingServerAddress.FindStringSubmatch(message); match != nil {
        event = &StartingServerAddressEvent{baseEvent{msg.date}, match[1]}

    } else if match := FailedToLoadFile.FindStringSubmatch(message); match != nil {
        var what FailedToLoadFileType

        switch match[1] {
        case "ban list":
            what = FailedToLoadBanList
        case "ip ban list":
            what = FailedToLoadIpBanList
        case "operators list":
            what = FailedToLoadOperatorsList
        case "white-list":
            what = FailedToLoadWhiteList
        default:
            return nil, UnrecognisedFailedToLoadFile.Format(match[1])
        }

        event = &FailedToLoadFileEvent{baseEvent{msg.date}, what, match[2], match[3]}

    } else if match := PreparingLevel.FindStringSubmatch(message); match != nil {
        event = &PreparingLevelEvent{baseEvent{msg.date}, match[1]}

    } else if match := DefaultGameMode.FindStringSubmatch(message); match != nil {
        gameMode, err := strconv.ParseInt(match[1], 10, 0)
        if err != nil {
            return nil, err
        }

        event = &DefaultGameModeEvent{baseEvent{msg.date}, int(gameMode)}

    } else if match := PreparingStartRegion.FindStringSubmatch(message); match != nil {
        levelNumber, err := strconv.ParseInt(match[1], 10, 0)
        if err != nil {
            return nil, err
        }

        var seed int64

        if match[2] != "" {
            seed, err = strconv.ParseInt(match[2], 10, 64)
            if err != nil {
                return nil, err
            }
        }

        event = &PreparingStartRegionEvent{baseEvent{msg.date}, int(levelNumber), seed}

    } else if match := PreparingSpawnAreaProgress.FindStringSubmatch(message); match != nil {
        progress, err := strconv.ParseInt(match[1], 10, 0)
        if err != nil {
            return nil, err
        }

        event = &PreparingSpawnAreaProgressEvent{baseEvent{msg.date}, int(progress)}

    } else if match := InitializationDone.FindStringSubmatch(message); match != nil {
        secs, err := strconv.ParseFloat(match[1], 32)
        if err != nil {
            return nil, err
        }

        event = &InitializationDoneEvent{baseEvent{msg.date}, float32(secs)}

    } else if match := TickSyncMessage.FindStringSubmatch(message); match != nil {
        event = &TickSyncMessageEvent{baseEvent{msg.date}}

    } else if match := ServerStopIssued.FindStringSubmatch(message); match != nil {
        event = &ServerStopIssuedEvent{baseEvent{msg.date}, match[1]}

    } else if match := ServerStopping.FindStringSubmatch(message); match != nil {
        event = &ServerStoppingEvent{baseEvent{msg.date}}

    } else if match := SavingChunks.FindStringSubmatch(message); match != nil {
        event = &SavingChunksEvent{baseEvent{msg.date}}

    } else if match := PlayerConnect.FindStringSubmatch(message); match != nil {
        eid, err := strconv.ParseInt(match[3], 10, 0)
        if err != nil {
            return nil, err
        }

        x, err := strconv.ParseFloat(match[5], 32)
        if err != nil {
            return nil, err
        }

        y, err := strconv.ParseFloat(match[6], 32)
        if err != nil {
            return nil, err
        }

        z, err := strconv.ParseFloat(match[7], 32)
        if err != nil {
            return nil, err
        }

        event = &PlayerConnectEvent{baseEvent{msg.date}, match[1], match[2], int(eid), match[4], float32(x), float32(y), float32(z)}

    } else if match := GameModeChanged.FindStringSubmatch(message); match != nil {
        gameMode, err := strconv.ParseInt(match[3], 10, 0)
        if err != nil {
            return nil, err
        }

        event = &GameModeChangedEvent{baseEvent{msg.date}, match[1], match[2], int(gameMode)}

    } else if match := PlayerOpped.FindStringSubmatch(message); match != nil {
        event = &PlayerOppedEvent{baseEvent{msg.date}, match[1], match[2]}

    } else if match := PlayerDeOpped.FindStringSubmatch(message); match != nil {
        event = &PlayerDeOppedEvent{baseEvent{msg.date}, match[1], match[2]}

    } else if match := PlayerIssuedCommand.FindStringSubmatch(message); match != nil {
        event = &PlayerIssuedCommandEvent{baseEvent{msg.date}, match[1], match[2]}

    } else if match := PlayerOldChat.FindStringSubmatch(message); match != nil {
        event = &PlayerChatEvent{baseEvent{msg.date}, match[1], match[2]}

    } else if match := PlayerDisconnect.FindStringSubmatch(message); match != nil {
        var reason DisconnectReason

        switch match[2] {
        case "genericReason", "":
            reason = GenericReason
        case "endOfStream":
            reason = EndOfStream
        case "quitting":
            reason = Quitting
        default:
            return nil, UnrecognisedDisconnectReason.Format(match[2])
        }

        event = &PlayerDisconnectEvent{baseEvent{msg.date}, match[1], reason}

    } else if match := PlayerTeleport.FindStringSubmatch(message); match != nil {
        event = &PlayerTeleportEvent{baseEvent{msg.date}, match[1], match[2], match[3]}

    } else if match := CraftBukkitVersionInfo.FindStringSubmatch(message); match != nil {
        event = &CraftBukkitVersionInfoEvent{baseEvent{msg.date}, match[1], match[2], match[3]}

    } else if match := PluginMessage.FindStringSubmatch(message); match != nil {
        event = &PluginMessageEvent{baseEvent{msg.date}, match[1], match[2]}

    } else if match := FolderMigrationBegan.FindStringSubmatch(message); match != nil {
        for {
            msg, err := ls.ReadLogMessage()
            if err != nil {
                return nil, err
            }

            matchString := FolderMigrationComplete.FindString(msg.message)
            if matchString != "" {
                break
            }
        }

        event, err = ls.ReadEvent()
        if err != nil {
            return nil, err
        }

    } else {
        //return nil, LineNotMatched.Format(message)
        event = &MiscEvent{baseEvent{msg.date}, message}
    }

    return event, nil
}
