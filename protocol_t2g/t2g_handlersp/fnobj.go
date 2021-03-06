// Code generated by "genprotocol.exe -ver=4761917555d922e5951aa6c3aee428ba5acda663bf281621426dc91134d9ffa7 -basedir=protocol_t2g -prefix=t2g -statstype=int"

package t2g_handlersp

import (
	"fmt"

	"github.com/kasworld/goguelike/protocol_t2g/t2g_idcmd"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
)

// obj base demux fn map

var DemuxRsp2ObjFnMap = [...]func(me interface{}, hd t2g_packet.Header, body interface{}) error{
	t2g_idcmd.Invalid:   objRecvRspFn_Invalid,   // Invalid make empty packet error
	t2g_idcmd.Register:  objRecvRspFn_Register,  // Register
	t2g_idcmd.Heartbeat: objRecvRspFn_Heartbeat, // Heartbeat
	t2g_idcmd.HighScore: objRecvRspFn_HighScore, // HighScore

}

// Invalid make empty packet error
func objRecvRspFn_Invalid(me interface{}, hd t2g_packet.Header, body interface{}) error {
	robj, ok := body.(*t2g_obj.RspInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Register
func objRecvRspFn_Register(me interface{}, hd t2g_packet.Header, body interface{}) error {
	robj, ok := body.(*t2g_obj.RspRegister_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Heartbeat
func objRecvRspFn_Heartbeat(me interface{}, hd t2g_packet.Header, body interface{}) error {
	robj, ok := body.(*t2g_obj.RspHeartbeat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// HighScore
func objRecvRspFn_HighScore(me interface{}, hd t2g_packet.Header, body interface{}) error {
	robj, ok := body.(*t2g_obj.RspHighScore_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}
