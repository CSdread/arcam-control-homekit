package arcam

const (
	TransmissionStart byte = 0x21
	TransmissionEnd   byte = 0x0D
)

var ReceiverModels = map[string]interface{}{
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

type ZoneNumber byte

const (
	ZoneOne ZoneNumber = 1
	ZoneTwo            = 2
)

type PowerStatus int

const (
	PowerStatusActive    PowerStatus = 0x01
	PowerStatusNotActive PowerStatus = 0x00
)

type MuteState int

const (
	MuteStateMuted    MuteState = 0x00
	MuteStateNotMuted MuteState = 0x01
)

type Command byte

const (
	PowerCommand                    Command = 0x00
	DisplayBirghtness                       = 0x01
	Headphones                              = 0x02
	FMGenre                                 = 0x03
	SoftwareVersion                         = 0x04
	RestoreFactoryDefaultSettings           = 0x05
	SaveRestoreSecureCopyOfSettings         = 0x06
	SimulateRC5IRCommand                    = 0x08
	DisplayInformationType                  = 0x09
	RequestCurrentSource                    = 0x1D
	HeadphoneOverride                       = 0x1F
	SelectAnalogueDigital                   = 0x0B
	SetRequestVolume                        = 0x0D
	RequestMuteStatus                       = 0x0E
	RequestDirectModeStatus                 = 0x0F
	RequestDecodeModeStatus2ch              = 0x10
	RequestDecodeModeStatusMCH              = 0x11
	RequestRDSInfo                          = 0x12
	RequestVideoOutputResolution            = 0x13
	RequestMenuStatus                       = 0x14
	RequestTunerPreset                      = 0x15
	Tune                                    = 0x16
	RequestDABStation                       = 0x18
	ProgramTypeCategory                     = 0x19
	DLSPDTInfo                              = 0x1A
	RequestPresetDetails                    = 0x1B
	NetworkPlaybackStatus                   = 0x1C
	IMAXEnhanced                            = 0x0C
	TrebleEqualisation                      = 0x35
	BassEqualisation                        = 0x36
	RoomEqualisation                        = 0x37
	DolbyAudio                              = 0x38
	Balance                                 = 0x3B
	SubwooferTrim                           = 0x3F
	LipsyncDelay                            = 0x40
	Compression                             = 0x41
	RequestIncomingVideoParameters          = 0x42
	RequestIncomingAudioFormat              = 0x43
	RequestIncomingAudioSampleRate          = 0x44
	SetRequestSubStereoTrim                 = 0x45
	SetRequestZone1OSDOnOff                 = 0x4E
	SetRequestVideoOutputSwitching          = 0x4F
	SetRequestInputName                     = 0x20
	FMScanUpDown                            = 0x23
	DABScan                                 = 0x24
	Heartbeat                               = 0x25
	Reboot                                  = 0x26
	BluetoothStatus                         = 0x50
	Setup                                   = 0x27
	RoomEQName                              = 0x34
	NowPlayingInfo                          = 0x64
	InputConfig                             = 0x28
	GeneralSetup                            = 0x29
	SpeakerTypes                            = 0x2A
	SpeakerDistances                        = 0x2B
	SpeakerLevels                           = 0x2C
	VideoInputs                             = 0x2D
	HDMISettings                            = 0x2E
	ZoneSettings                            = 0x2F // not  AVR5, AVR10
	Network                                 = 0x30
	Bluetooth                               = 0x32
	EngineeringMenu                         = 0x33
)

type AVRC5CommandCode struct {
	Data1 byte
	Data2 byte
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
	BD                               = AVRC5CommandCode{0x10, 0x62}
	CD                               = AVRC5CommandCode{0x10, 0x76}
	STB                              = AVRC5CommandCode{0x10, 0x64}
	UHD                              = AVRC5CommandCode{0x10, 0x7D}
	BT                               = AVRC5CommandCode{0x10, 0x7A}
	Display                          = AVRC5CommandCode{0x10, 0x3A}
	PowerOn                          = AVRC5CommandCode{0x10, 0x7B}
	PowerOff                         = AVRC5CommandCode{0x10, 0x7C}
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

type InputSource byte

const (
	InputCD       InputSource = 0x01
	InputBD                   = 0x02
	InputAV                   = 0x03
	InputSAT                  = 0x04
	InputPVR                  = 0x05
	InputUHD                  = 0x06
	InputAUX                  = 0x08
	InputDISPLAY              = 0x09
	InputTUNERFM              = 0x0B
	InputTUNERDAB             = 0x0C
	InputNET                  = 0x0E
	InputSTB                  = 0x10
	InputGAME                 = 0x11
	InputBT                   = 0x12
)

type Answer byte

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
	Data        []byte
}

type Request struct {
	Zone    ZoneNumber
	Command Command
	Data    []byte
}

var InputDisplayNameMap = map[InputSource]string{
	InputCD:       "CD",
	InputBD:       "BD",
	InputAV:       "AV",
	InputSAT:      "Sat",
	InputPVR:      "PVR",
	InputUHD:      "UHD",
	InputAUX:      "Aux",
	InputDISPLAY:  "Display",
	InputTUNERFM:  "FM",
	InputTUNERDAB: "DAB",
	InputNET:      "Net",
	InputSTB:      "STB",
	InputGAME:     "Game",
	InputBT:       "BT",
}

var InputSourceCommandMap = map[InputSource]AVRC5CommandCode{
	InputCD:       CD,
	InputBD:       BD,
	InputAV:       AV,
	InputSAT:      Sat,
	InputPVR:      PVR,
	InputUHD:      UHD,
	InputAUX:      Aux,
	InputDISPLAY:  Display,
	InputTUNERFM:  FM,
	InputTUNERDAB: DAB,
	InputNET:      Net,
	InputSTB:      STB,
	InputGAME:     Game,
	InputBT:       BT,
}
