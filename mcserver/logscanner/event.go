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
    DefaultGameModeEventType
    PreparingStartRegionEventType
    PreparingSpawnAreaProgressEventType
    InitializationDoneEventType
    TickSyncMessageEventType
    ServerStopIssuedEventType
    ServerStoppingEventType
    SavingChunksEventType
    PlayerConnectEventType
    GameModeChangedEventType
    PlayerOppedEventType
    PlayerDeOppedEventType
    PlayerIssuedCommandEventType
    PlayerChatEventType
    PlayerDisconnectEventType
    PlayerTeleportEventType
    CraftBukkitVersionInfoEventType
    PluginMessageEventType

    MiscEventType
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

// Type DefaultGameModeEvent represents an event informing you of the default game mode.
type DefaultGameModeEvent struct {
    baseEvent
    gameMode int
}

// Function Type returns the type of the event.
func (event *DefaultGameModeEvent) Type() EventType {
    return DefaultGameModeEventType
}

// Function GameMode returns the default game mode.
func (event *DefaultGameModeEvent) GameMode() int {
    return event.gameMode
}

// -------------------------------------------------------------------------------------------------

// Type PreparingStartRegionEvent represents an event informing you that the start region is being
// prepared.
type PreparingStartRegionEvent struct {
    baseEvent
    levelNumber int
    seed        int64
}

// Function Type returns the type of the event.
func (event *PreparingStartRegionEvent) Type() EventType {
    return PreparingStartRegionEventType
}

// Function LevelNumber returns the number of the level being prepared.
func (event *PreparingStartRegionEvent) LevelNumber() int {
    return event.levelNumber
}

// Function Seed returns the seed of the level being prepared (if present).
func (event *PreparingStartRegionEvent) Seed() int64 {
    return event.seed
}

// -------------------------------------------------------------------------------------------------

// Type PreparingSpawnAreaProgressEvent represents an event informing you of the progress of the
// spawn area preparation.
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

// Type InitializationDoneEvent represents an event informing you that initialization was completed.
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

// Type TickSyncMessageEvent represents an event informing you that the server may be running
// slowly.
type TickSyncMessageEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *TickSyncMessageEvent) Type() EventType {
    return TickSyncMessageEventType
}

// -------------------------------------------------------------------------------------------------

// Type ServerStopIssuedEvent represents an event informing you that a user issued the /stop
// command.
type ServerStopIssuedEvent struct {
    baseEvent
    player string
}

// Function Type returns the type of the event.
func (event *ServerStopIssuedEvent) Type() EventType {
    return ServerStopIssuedEventType
}

// Function Player returns the user who issued the server stop.
func (event *ServerStopIssuedEvent) Player() string {
    return event.player
}

// -------------------------------------------------------------------------------------------------

// Type ServerStoppingEvent represents an event informing you that the server is stopping.
type ServerStoppingEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *ServerStoppingEvent) Type() EventType {
    return ServerStoppingEventType
}

// -------------------------------------------------------------------------------------------------

// Type SavingChunksEvent represents an event informing you that the chunks are being saved.
type SavingChunksEvent struct {
    baseEvent
}

// Function Type returns the type of the event.
func (event *SavingChunksEvent) Type() EventType {
    return SavingChunksEventType
}

// -------------------------------------------------------------------------------------------------

// Type PlayerConnectEvent represents a player logging in.
type PlayerConnectEvent struct {
    baseEvent
    player  string
    address string
    eid     int
    world   string
    x       float32
    y       float32
    z       float32
}

// Function Type returns the type of the event.
func (event *PlayerConnectEvent) Type() EventType {
    return PlayerConnectEventType
}

// Function Player returns the username of the player who logged in.
func (event *PlayerConnectEvent) Player() string {
    return event.player
}

// Function Address returns the IP address that the player logged in from.
func (event *PlayerConnectEvent) Address() string {
    return event.address
}

// Function EntityID returns the entity ID assigned to the player when they logged in.
func (event *PlayerConnectEvent) EntityID() int {
    return event.eid
}

// Function World returns the world the player spawned in.
func (event *PlayerConnectEvent) World() string {
    return event.world
}

// Function X returns the x coordinate of the location where the player spawned.
func (event *PlayerConnectEvent) X() float32 {
    return event.x
}

// Function Y returns the x coordinate of the location where the player spawned.
func (event *PlayerConnectEvent) Y() float32 {
    return event.y
}

// Function Z returns the x coordinate of the location where the player spawned.
func (event *PlayerConnectEvent) Z() float32 {
    return event.z
}

// -------------------------------------------------------------------------------------------------

// Type GameModeChangedEvent represents a player's game mode being changed.
type GameModeChangedEvent struct {
    baseEvent
    sender   string
    player   string
    gameMode int
}

// Function Type returns the type of the event.
func (event *GameModeChangedEvent) Type() EventType {
    return GameModeChangedEventType
}

// Function Sender returns the name of the player who issued the command.
func (event *GameModeChangedEvent) Sender() string {
    return event.sender
}

// Function Player returns the name of the player whose game mode was changed.
func (event *GameModeChangedEvent) Player() string {
    return event.player
}

// Function GameMode returns the name of the level that is being prepared.
func (event *GameModeChangedEvent) GameMode() int {
    return event.gameMode
}

// -------------------------------------------------------------------------------------------------

// Type PlayerOppedEvent represents a player being opped.
type PlayerOppedEvent struct {
    baseEvent
    sender string
    player string
}

// Function Type returns the type of the event.
func (event *PlayerOppedEvent) Type() EventType {
    return PlayerOppedEventType
}

// Function Sender returns the name of the player who issued the command.
func (event *PlayerOppedEvent) Sender() string {
    return event.sender
}

// Function Player returns the name of the player who was opped.
func (event *PlayerOppedEvent) Player() string {
    return event.player
}

// -------------------------------------------------------------------------------------------------

// Type PlayerDeOppedEvent represents a player being de-opped.
type PlayerDeOppedEvent struct {
    baseEvent
    sender string
    player string
}

// Function Type returns the type of the event.
func (event *PlayerDeOppedEvent) Type() EventType {
    return PlayerDeOppedEventType
}

// Function Sender returns the name of the player who issued the command.
func (event *PlayerDeOppedEvent) Sender() string {
    return event.sender
}

// Function Player returns the name of the player who was de-opped.
func (event *PlayerDeOppedEvent) Player() string {
    return event.player
}

// -------------------------------------------------------------------------------------------------

// Type PlayerIssuedCommandEvent represents a player issuing a non-standard server command.
type PlayerIssuedCommandEvent struct {
    baseEvent
    player        string
    commandString string
}

// Function Type returns the type of the event.
func (event *PlayerIssuedCommandEvent) Type() EventType {
    return PlayerIssuedCommandEventType
}

// Function Player returns the name of the player who issued the command.
func (event *PlayerIssuedCommandEvent) Player() string {
    return event.player
}

// Function CommandString returns the command that was issued, without the leading slash.
func (event *PlayerIssuedCommandEvent) CommandString() string {
    return event.commandString
}

// -------------------------------------------------------------------------------------------------

// Type PlayerChatEvent represents a player sending a chat message.
type PlayerChatEvent struct {
    baseEvent
    player  string
    message string
}

// Function Type returns the type of the event.
func (event *PlayerChatEvent) Type() EventType {
    return PlayerChatEventType
}

// Function Player returns the name of the player who sent the message.
func (event *PlayerChatEvent) Player() string {
    return event.player
}

// Function CommandString returns the chat message.
func (event *PlayerChatEvent) Message() string {
    return event.message
}

// -------------------------------------------------------------------------------------------------

// Type DisconnectReason represents a reason for disconnecting.
type DisconnectReason uint

// Disconnect reason constants.
const (
    _ DisconnectReason = iota

    GenericReason
    EndOfStream
    Quitting
)

// Type PlayerDisconnectEvent represents a player disconnecting.
type PlayerDisconnectEvent struct {
    baseEvent
    player string
    reason DisconnectReason
}

// Function Type returns the type of the event.
func (event *PlayerDisconnectEvent) Type() EventType {
    return PlayerDisconnectEventType
}

// Function Player returns the name of the player who disconnected.
func (event *PlayerDisconnectEvent) Player() string {
    return event.player
}

// Function Reason returns the disconnect reason.
func (event *PlayerDisconnectEvent) Reason() DisconnectReason {
    return event.reason
}

// -------------------------------------------------------------------------------------------------

// Type PlayerTeleportEvent represents a player teleporting.
type PlayerTeleportEvent struct {
    baseEvent
    sender string
    from   string
    to     string
}

// Function Type returns the type of the event.
func (event *PlayerTeleportEvent) Type() EventType {
    return PlayerTeleportEventType
}

// Function Sender returns the name of the player who issued the command.
func (event *PlayerTeleportEvent) Sender() string {
    return event.sender
}

// Function From returns the name of the player who was teleported.
func (event *PlayerTeleportEvent) From() string {
    return event.from
}

// Function To returns the name of the player who was teleported to.
func (event *PlayerTeleportEvent) To() string {
    return event.to
}

// -------------------------------------------------------------------------------------------------

// Type CraftBukkitVersionInfoEvent represents CraftBukkit's version information message.
type CraftBukkitVersionInfoEvent struct {
    baseEvent
    version    string
    mcVersion  string
    apiVersion string
}

// Function Type returns the type of the event.
func (event *CraftBukkitVersionInfoEvent) Type() EventType {
    return CraftBukkitVersionInfoEventType
}

// Function Version returns the CraftBukkit version.
func (event *CraftBukkitVersionInfoEvent) Version() string {
    return event.version
}

// Function MCVersion returns the MineCraft version.
func (event *CraftBukkitVersionInfoEvent) MCVersion() string {
    return event.mcVersion
}

// Function APIVersion returns the CraftBukkit API version.
func (event *CraftBukkitVersionInfoEvent) APIVersion() string {
    return event.apiVersion
}

// -------------------------------------------------------------------------------------------------

// Type PluginMessageEvent represents a log message from a plugin
type PluginMessageEvent struct {
    baseEvent
    pluginName string
    message    string
}

// Function Type returns the type of the event.
func (event *PluginMessageEvent) Type() EventType {
    return PluginMessageEventType
}

// Function PluginName returns the name of the plugin.
func (event *PluginMessageEvent) PluginName() string {
    return event.pluginName
}

// Function Message returns the actual message.
func (event *PluginMessageEvent) Message() string {
    return event.message
}

// -------------------------------------------------------------------------------------------------

// Type MiscEvent represents an event line that could not be interpreted.
type MiscEvent struct {
    baseEvent
    message string
}

// Function Type returns the type of the event.
func (event *MiscEvent) Type() EventType {
    return MiscEventType
}

// Function Message returns the actual message.
func (event *MiscEvent) Message() string {
    return event.message
}
