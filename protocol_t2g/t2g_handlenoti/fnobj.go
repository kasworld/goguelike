// Code generated by "genprotocol.exe -ver=4761917555d922e5951aa6c3aee428ba5acda663bf281621426dc91134d9ffa7 -basedir=protocol_t2g -prefix=t2g -statstype=int"

package t2g_handlenoti

import (
	"fmt"

	"github.com/kasworld/goguelike/protocol_t2g/t2g_idnoti"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
)

// obj base demux fn map

var DemuxNoti2ObjFnMap = [...]func(me interface{}, hd t2g_packet.Header, body interface{}) error{
	t2g_idnoti.Invalid: objRecvNotiFn_Invalid, // Invalid make empty packet error

}

// Invalid make empty packet error
func objRecvNotiFn_Invalid(me interface{}, hd t2g_packet.Header, body interface{}) error {
	robj, ok := body.(*t2g_obj.NotiInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}
