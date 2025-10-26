package gamelogic

type Command string

const (
	CmdServerPause  = "pause"
	CmdServerResume = "resume"
	CmdServerQuit   = "quit"
	CmdServerHelp   = "help"
)

func getAllServerCommands() map[Command]struct{} {
	return map[Command]struct{}{
		CmdServerPause:  {},
		CmdServerResume: {},
		CmdServerQuit:   {},
		CmdServerHelp:   {},
	}
}

const (
	CmdClientMove    = "move"
	CmdClientSpawn   = "spawn"
	CmdClientStatus  = "status"
	CmdClientSpam    = "spam"
	CmdClientQuit    = "quit"
	CmdClientHelp    = "help"
	CmdClientInvalid = "invalid"
)

func getAllClientCommands() map[Command]struct{} {
	return map[Command]struct{}{
		CmdClientMove:    {},
		CmdClientSpawn:   {},
		CmdClientStatus:  {},
		CmdClientSpam:    {},
		CmdClientQuit:    {},
		CmdClientHelp:    {},
		CmdClientInvalid: {},
	}
}
