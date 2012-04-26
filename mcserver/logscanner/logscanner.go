// Package logscanner provides parsing of Minecraft server logs.
package logscanner

import (
    "bufio"
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
    DefaultGameType              = regexp.MustCompile("^Default game type: (\\d+)$")
    PreparingStartRegion         = regexp.MustCompile("^Preparing start region for level (\\d+)$")
    PreparingSpawnAreaProgress   = regexp.MustCompile("^Preparing spawn area: (\\d+)%$")
    InitializationDone           = regexp.MustCompile("^Done \\(([0-9.]+)s\\)! For help, type \"help\" or \"\\?\"$")
    TickSyncMessage              = regexp.MustCompile("^Can't keep up! Did the system time change, or is the server overloaded\\?$")
    ServerStopIssued             = regexp.MustCompile("^(.+): Stopping the server..$")
    ServerStopping               = regexp.MustCompile("^Stopping server$")
)

// Errors
var (
    LineNotMatched               = bettererrors.New("The following line could not be interpreted:\n%s")
    UnrecognisedLogLevel         = bettererrors.New("The log level '%s' was not recognised.")
    UnrecognisedFailedToLoadFile = bettererrors.New("The term '%s' was not recognised.")
)

// Type LogScanner represents a log scanner.
type LogScanner struct {
    source   *bufio.Reader
    timezone *time.Location
}

// Function NewLogScanner creates and returns a new log scanner.
func NewLogScanner(source *bufio.Reader, timezone *time.Location) (ls *LogScanner) {
    return &LogScanner{
        source:   source,
        timezone: timezone,
    }
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
        return nil, LineNotMatched.Format(line)
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

    } else if match := DefaultGameType.FindStringSubmatch(message); match != nil {
        gameType, err := strconv.ParseUint(match[1], 10, 0)
        if err != nil {
            return nil, err
        }

        event = &DefaultGameTypeEvent{baseEvent{msg.date}, int(gameType)}

    } else if match := PreparingStartRegion.FindStringSubmatch(message); match != nil {
        levelNumber, err := strconv.ParseUint(match[1], 10, 0)
        if err != nil {
            return nil, err
        }

        event = &PreparingStartRegionEvent{baseEvent{msg.date}, int(levelNumber)}

    } else if match := PreparingSpawnAreaProgress.FindStringSubmatch(message); match != nil {
        progress, err := strconv.ParseUint(match[1], 10, 0)
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

    } else {
        return nil, LineNotMatched.Format(message)
    }

    return event, nil
}
