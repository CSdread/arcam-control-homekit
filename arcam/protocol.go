package arcam

const (
	TransmissionStart = 0x21
	TransmissionEnd   = 0x0D
)

type ZoneNumber int

var RecieverModels = map[string]interface{}{
	"AVR5":  struct{}{},
	"AVR10": struct{}{},
	"AVR11": struct{}{},
	"AVR20": struct{}{},
	"AVR21": struct{}{},
	"AVR30": struct{}{},
	"AVR31": struct{}{},
	"AV40":  struct{}{},
	"AVR41": struct{}{},
}

const (
	ZoneOne ZoneNumber = 1
	ZoneTwo            = 2
)

type Command int64

const (
	// System Commands
	PowerCommand = iota
	DisplayBirghtness
	Headphones
	FMGenre
	SoftwareVersion
	RestoreFactoryDefaultSettings
	SaveRestoreSecureCopyOfSettings
	SimulateRC5IRCommand
	DisplayInformationType
	RequestCurrentSource
	HeadphoneOverride
	// Input Command
	SelectAnalogueDigital
	// Output Commands
	SetRequestVolume
	RequestMuteStatus
	RequestDirectModeStatus
	RequestDecodeModeStatus2ch
	RequestDecodeModeStatusMCH
	RequestRDSInfo
	RequestVideoOutputResolution
	// Menu Commands
	RequestMenuStatus
	RequestTunerPreset
	Tune
	RequestDABStation
	ProgramTypeCategory
	DLSPDTInfo
	RequestPresetDetails
	NetworkPlaybackStatus
	IMAXEnhanced
	// Setup Adjustment Commands
	TrebleEqualisation
	BassEqualisation
	RoomEqualisation
	DolbyAudio
	Balance
	SubwooferTrim
	LipsyncDelay
	Compression
	RequestIncomingVideoParameters
	RequestIncomingAudioFormat
	RequestIncomingAudioSampleRate
	SetRequestSubStereoTrim
	SetRequestZone1OSDOnOff
	SetRequestVideoOutputSwitching
	SetRequestInputName
	FMScanUpDown
	DABScan
	Heartbeat
	Reboot
	BluetoothStatus
	Setup
	RoomEQName
	NowPlayingInfo
	InputConfig
	GeneralSetup
	SpeakerTypes
	SpeakerDistances
	SpeakerLevels
	VideoInputs
	HDMISettings
	ZoneSettings // not AVR5, AVR10
	Network
	Bluetooth
	EngineeringMenu
)

type AVRC5CommandCode struct {
	Data1 int64
	Data2 int64
}

var (
	Standby                          = AVRC5CommandCode{0x10, 0x0C}
	Eject                            = AVRC5CommandCode{0x10, 0x2D}
	One                              = AVRC5CommandCode{0x10, 0x01}
	Two                              = AVRC5CommandCode{0x10, 0x01}
	Three                            = AVRC5CommandCode{0x10, 0x01}
	Four                             = AVRC5CommandCode{0x10, 0x01}
	Five                             = AVRC5CommandCode{0x10, 0x01}
	Six                              = AVRC5CommandCode{0x10, 0x01}
	Seven                            = AVRC5CommandCode{0x10, 0x01}
	Eight                            = AVRC5CommandCode{0x10, 0x01}
	Nine                             = AVRC5CommandCode{0x10, 0x01}
	AccessLipsyncDelayControl        = AVRC5CommandCode{0x10, 0x32}
	Zero                             = AVRC5CommandCode{0x10, 0x00}
	CycleBetweenVFDInformationPanels = AVRC5CommandCode{0x10, 0x37}
	Rewind                           = AVRC5CommandCode{0x10, 0x79}
	FastForward                      = AVRC5CommandCode{0x10, 0x34}
	SkipBack                         = AVRC5CommandCode{0x10, 0x21}
	SkipForward                      = AVRC5CommandCode{0x10, 0x0B}
	Stop                             = AVRC5CommandCode{0x10, 0x36}
	Play                             = AVRC5CommandCode{0x10, 0x35}
	Pause                            = AVRC5CommandCode{0x10, 0x30}
	DiscRecordDTSDialogControl       = AVRC5CommandCode{0x10, 0x5A}
	MENUEnterSystemStatusMenu        = AVRC5CommandCode{0x10, 0x52}
	NavigateUp                       = AVRC5CommandCode{0x10, 0x56}
	PopUpDolbyVolumeOnOff            = AVRC5CommandCode{0x10, 0x46}
	NavigateLeft                     = AVRC5CommandCode{0x10, 0x51}
	OK                               = AVRC5CommandCode{0x10, 0x57}
	NavigateRight                    = AVRC5CommandCode{0x10, 0x50}
	AudioRoomEQOnOff                 = AVRC5CommandCode{0x10, 0x1E}
	NavigateDown                     = AVRC5CommandCode{0x10, 0x55}
	RTNAccessSubwooferTrimControl    = AVRC5CommandCode{0x10, 0x33}
	HOME                             = AVRC5CommandCode{0x10, 0x2B}
	Mute                             = AVRC5CommandCode{0x10, 0x0D}
	IncreaseVolume                   = AVRC5CommandCode{0x10, 0x10}
	MODECycleBetweenDecodingModes    = AVRC5CommandCode{0x10, 0x20}
	DISPChangeVFDBrightness          = AVRC5CommandCode{0x10, 0x3B}
	ActivateDirectMode               = AVRC5CommandCode{0x10, 0x0A}
	DecreaseVolume                   = AVRC5CommandCode{0x10, 0x11}
	Red                              = AVRC5CommandCode{0x10, 0x29}
	Green                            = AVRC5CommandCode{0x10, 0x2A}
	Yellow                           = AVRC5CommandCode{0x10, 0x2B}
	Blue                             = AVRC5CommandCode{0x10, 0x37}
	Radio                            = AVRC5CommandCode{0x10, 0x5B}
	Aux                              = AVRC5CommandCode{0x10, 0x63}
	Net                              = AVRC5CommandCode{0x10, 0x5C}
	AV                               = AVRC5CommandCode{0x10, 0x5E}
	Sat                              = AVRC5CommandCode{0x10, 0x1B}
	PVR                              = AVRC5CommandCode{0x10, 0x60}
	Game                             = AVRC5CommandCode{0x10, 0x61}
	ChangeControlToNextZone          = AVRC5CommandCode{0x10, 0x5F}
	AccessBassControl                = AVRC5CommandCode{0x10, 0x27}
	AccessSpeakerTrimControl         = AVRC5CommandCode{0x10, 0x25}
	AccessTrebleControl              = AVRC5CommandCode{0x10, 0x0E}
	Random                           = AVRC5CommandCode{0x10, 0x4C}
	Repeat                           = AVRC5CommandCode{0x10, 0x31}
	DirectModeOn                     = AVRC5CommandCode{0x10, 0x4E}
	DirectModeOff                    = AVRC5CommandCode{0x10, 0x4F}
	MultiChannel                     = AVRC5CommandCode{0x10, 0x6A}
	Stereo                           = AVRC5CommandCode{0x10, 0x6B}
	DolbySurround                    = AVRC5CommandCode{0x10, 0x6E}
	DTSNeo6Cinema                    = AVRC5CommandCode{0x10, 0x6F}
	DTSNeo6Music                     = AVRC5CommandCode{0x10, 0x70}
	DTSNeuralX                       = AVRC5CommandCode{0x10, 0x71}
	Reserved                         = AVRC5CommandCode{0x10, 0x72}
	VirtualHeight                    = AVRC5CommandCode{0x10, 0x73}
	FiveSevenChStereo                = AVRC5CommandCode{0x10, 0x45}
	DolbyDEX                         = AVRC5CommandCode{0x10, 0x17}
	AuroMatic3D                      = AVRC5CommandCode{0x10, 0x47}
	AuroNative                       = AVRC5CommandCode{0x10, 0x67}
	Auro2D                           = AVRC5CommandCode{0x10, 0x68}
	MuteOn                           = AVRC5CommandCode{0x10, 0x1A}
	MuteOff                          = AVRC5CommandCode{0x10, 0x78}
	FM                               = AVRC5CommandCode{0x10, 0x1c}
	DAB                              = AVRC5CommandCode{0x10, 0x48}
	LipSyncPlus5ms                   = AVRC5CommandCode{0x10, 0x0F}
	LipSyncMinus5ms                  = AVRC5CommandCode{0x10, 0x65}
	SubTrimPlusHalfDb                = AVRC5CommandCode{0x10, 0x69}
	SubTrimMinusHalfDb               = AVRC5CommandCode{0x10, 0x6C}
	DisplayOff                       = AVRC5CommandCode{0x10, 0x1F}
	DisplayL1                        = AVRC5CommandCode{0x10, 0x22}
	DisplayL2                        = AVRC5CommandCode{0x10, 0x23}
	BalanceLeft                      = AVRC5CommandCode{0x10, 0x26}
	BalanceRight                     = AVRC5CommandCode{0x10, 0x28}
	BassPlus1                        = AVRC5CommandCode{0x10, 0x2C}
	BassMinus1                       = AVRC5CommandCode{0x10, 0x38}
	TreblePlus1                      = AVRC5CommandCode{0x10, 0x2E}
	TrebleMinus1                     = AVRC5CommandCode{0x10, 0x66}
	SetZone2ToFollowZone1            = AVRC5CommandCode{0x10, 0x14}
	Zone2PowerOn                     = AVRC5CommandCode{0x17, 0x7B} // Not AVR5, AVR10
	Zone2PowerOff                    = AVRC5CommandCode{0x17, 0x7C} // Not AVR5, AVR10
	Zone2VolumePlus                  = AVRC5CommandCode{0x17, 0x01} // Not AVR5, AVR10
	Zone2VolumeMinus                 = AVRC5CommandCode{0x17, 0x02} // Not AVR5, AVR10
	Zone2Mute                        = AVRC5CommandCode{0x17, 0x03} // Not AVR5, AVR10
	ZoneMuteOn                       = AVRC5CommandCode{0x17, 0x04} // Not AVR5, AVR10
	ZoneMuteOff                      = AVRC5CommandCode{0x17, 0x05} // Not AVR5, AVR10
	Zone2CD                          = AVRC5CommandCode{0x17, 0x06}
	Zone2BD                          = AVRC5CommandCode{0x17, 0x07}
	Zone2STB                         = AVRC5CommandCode{0x17, 0x08}
	Zone2AV                          = AVRC5CommandCode{0x17, 0x09}
	Zone2Game                        = AVRC5CommandCode{0x17, 0x0B}
	Zone2Aux                         = AVRC5CommandCode{0x17, 0x0D}
	Zone2PVR                         = AVRC5CommandCode{0x17, 0x0F}
	Zone2FM                          = AVRC5CommandCode{0x17, 0x0E}
	Zone2DAB                         = AVRC5CommandCode{0x17, 0x10}
	Zone2USB                         = AVRC5CommandCode{0x17, 0x12}
	Zone2NET                         = AVRC5CommandCode{0x17, 0x13}
	Zone2SAT                         = AVRC5CommandCode{0x17, 0x14}
	Zone2UHD                         = AVRC5CommandCode{0x17, 0x17}
	Zone2BT                          = AVRC5CommandCode{0x17, 0x16}
	SelectHDMIOut1                   = AVRC5CommandCode{0x17, 0x49}
	SelectHDMIOut2                   = AVRC5CommandCode{0x17, 0x4A}
	SelectHDMIOut1And2               = AVRC5CommandCode{0x17, 0x4B}
)

type Answer int64

const (
	StatusUpdate             Answer = 0x00
	ZoneInvalid              Answer = 0x82
	CommandNotRecognized     Answer = 0x83
	ParameterNotRecognized   Answer = 0x84
	CommandInvalidAtThisTime Answer = 0x85
	InvalidDataLength        Answer = 0x86
)

type Response struct {
	Zone        ZoneNumber
	CommandCode Command
	AnswerCode  Answer
	Data        []string
}

type Transmission struct {
	Zone    ZoneNumber
	Command string
	Data    []string
}
