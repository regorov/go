package logscanner

import (
    "time"
)

// Type EventType represents the type of an event.
type EventType uint

// Event type constants
const (
    _ EventType = iota
    StartingServerEventType
    LoadingPropertiesEventType
    ServerPropertiesDoesNotExistEventType
    GeneratingNewPropertiesFileEventType
    StartingServerAddressEventType
    FailedToLoadFileEventType
    PreparingLevelEventType
    DefaultGameTypeEventType
    PreparingStartRegionEventType
    PreparingSpawnAreaProgressEventType
    InitializationDoneEventType
    TickSyncMessageEventType
    ServerStopIssuedEventType
)

// Interface Event represents a logged event.
type Event interface {
    Date() time.Time
    Type() EventType
}

// Type baseEvent includes common fields for all events.
type baseEvent struct {
    date time.Time
}

// Function Date returns the event's date.
func (event *baseEvent) Date() time.Time {
    return event.date
}

// -------------------------------------------------------------------------------------------------

// Type StartingServerEvent represents a server-start event.
type StartingServerEvent struct {
    baseEvent
    version string
}

// Function Type returns the type of the event.
func (event *StartingServerEvent) Type() EventType {
    return StartingServerEventType
}

// Function Version returns the Minecraft version being run.
func (event *StartingServerEvent) Version() string {
    return event.version
}

// -------------------------------------------------------------------------------------------------

// Type LoadingPropertiesEvent represents a properties-load event.
type LoadingPropertiesEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *LoadingPropertiesEvent) Type() EventType {
    return LoadingPropertiesEventType
}

// -------------------------------------------------------------------------------------------------

// Type ServerPropertiesDoesNotExistEvent represents the event when server.properties does not exist
// upon server start.
type ServerPropertiesDoesNotExistEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *ServerPropertiesDoesNotExistEvent) Type() EventType {
    return ServerPropertiesDoesNotExistEventType
}

// -------------------------------------------------------------------------------------------------

// Type GeneratingNewPropertiesFileEvent represents the event where a new server.properties must be
// generated.
type GeneratingNewPropertiesFileEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *GeneratingNewPropertiesFileEvent) Type() EventType {
    return GeneratingNewPropertiesFileEventType
}

// -------------------------------------------------------------------------------------------------

// Type StartingServerAddressEvent represents the event where the server is being started.
type StartingServerAddressEvent struct {
    baseEvent
    address string
}

// Function Type returns the type of the event.
func (event *StartingServerAddressEvent) Type() EventType {
    return StartingServerAddressEventType
}

// Function Address returns the address the server has started on.
func (event *StartingServerAddressEvent) Address() string {
    return event.address
}

// -------------------------------------------------------------------------------------------------

// Failed to load file constant type.
type FailedToLoadFileType uint

// Failed to load file constants.
const (
    _ FailedToLoadFileType = iota
    FailedToLoadBanList
    FailedToLoadIpBanList
    FailedToLoadOperatorsList
    FailedToLoadWhiteList
)

// Type FailedToLoadFileEvent represents the event where a file failed to load.
type FailedToLoadFileEvent struct {
    baseEvent
    what        FailedToLoadFileType
    javaError   string
    javaMessage string
}

// Function Type returns the type of the event.
func (event *FailedToLoadFileEvent) Type() EventType {
    return FailedToLoadFileEventType
}

// Function What returns what failed to load.
func (event *FailedToLoadFileEvent) What() FailedToLoadFileType {
    return event.what
}

// Function JavaError returns the name of the Java error that occurred.
func (event *FailedToLoadFileEvent) JavaError() string {
    return event.javaError
}

// Function JavaMessage returns the associated Java error message.
func (event *FailedToLoadFileEvent) JavaMessage() string {
    return event.javaMessage
}

// -------------------------------------------------------------------------------------------------

// Type PreparingLevelEvent represents the event where a level is being prepared.
type PreparingLevelEvent struct {
    baseEvent
    levelName string
}

// Function Type returns the type of the event.
func (event *PreparingLevelEvent) Type() EventType {
    return PreparingLevelEventType
}

// Function LevelName returns the name of the level that is being prepared.
func (event *PreparingLevelEvent) LevelName() string {
    return event.levelName
}

// -------------------------------------------------------------------------------------------------

// Type DefaultGameTypeEvent represents an event informing you of the default game type.
type DefaultGameTypeEvent struct {
    baseEvent
    gameType int
}

// Function Type returns the type of the event.
func (event *DefaultGameTypeEvent) Type() EventType {
    return DefaultGameTypeEventType
}

// Function GameType returns the name of the level that is being prepared.
func (event *DefaultGameTypeEvent) GameType() int {
    return event.gameType
}

// -------------------------------------------------------------------------------------------------

// Type PreparingStartRegionEvent represents an event informing you of the default game type.
type PreparingStartRegionEvent struct {
    baseEvent
    levelNumber int
}

// Function Type returns the type of the event.
func (event *PreparingStartRegionEvent) Type() EventType {
    return PreparingStartRegionEventType
}

// Function LevelNumber returns the number of the level being prepared.
func (event *PreparingStartRegionEvent) LevelNumber() int {
    return event.levelNumber
}

// -------------------------------------------------------------------------------------------------

// Type PreparingSpawnAreaProgressEvent represents an event informing you of the default game type.
type PreparingSpawnAreaProgressEvent struct {
    baseEvent
    progress int
}

// Function Type returns the type of the event.
func (event *PreparingSpawnAreaProgressEvent) Type() EventType {
    return PreparingSpawnAreaProgressEventType
}

// Function Progess returns the percentage progress.
func (event *PreparingSpawnAreaProgressEvent) Progress() int {
    return event.progress
}

// -------------------------------------------------------------------------------------------------

// Type InitializationDoneEvent represents an event informing you of the default game type.
type InitializationDoneEvent struct {
    baseEvent
    secs float32
}

// Function Type returns the type of the event.
func (event *InitializationDoneEvent) Type() EventType {
    return InitializationDoneEventType
}

// Function Secs returns the time it took to initialize.
func (event *InitializationDoneEvent) Secs() float32 {
    return event.secs
}

// -------------------------------------------------------------------------------------------------

// Type TickSyncMessageEvent represents an event informing you of the default game type.
type TickSyncMessageEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *TickSyncMessageEvent) Type() EventType {
    return TickSyncMessageEventType
}

// -------------------------------------------------------------------------------------------------

// Type ServerStopIssuedEvent represents an event informing you of the default game type.
type ServerStopIssuedEvent struct {
    baseEvent
    user string
}

// Function Type returns the type of the event.
func (event *ServerStopIssuedEvent) Type() EventType {
    return ServerStopIssuedEventType
}

// Function User returns the user who issued the server stop.
func (event *ServerStopIssuedEvent) User() string {
    return event.user
}
