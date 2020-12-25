// Code generated by "genprotocol.exe -ver=4761917555d922e5951aa6c3aee428ba5acda663bf281621426dc91134d9ffa7 -basedir=protocol_t2g -prefix=t2g -statstype=int"

package t2g_handlereq

import (
	"github.com/kasworld/goguelike/protocol_t2g/t2g_error"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_idcmd"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_obj"
	"github.com/kasworld/goguelike/protocol_t2g/t2g_packet"
)

// bytes base fn map api, unmarshal in api
var DemuxReq2BytesAPIFnMap = [...]func(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error){
	t2g_idcmd.Invalid:   bytesAPIFn_ReqInvalid,   // Invalid make empty packet error
	t2g_idcmd.Register:  bytesAPIFn_ReqRegister,  // Register
	t2g_idcmd.Heartbeat: bytesAPIFn_ReqHeartbeat, // Heartbeat
	t2g_idcmd.HighScore: bytesAPIFn_ReqHighScore, // HighScore

} // DemuxReq2BytesAPIFnMap

// Invalid make empty packet error
func bytesAPIFn_ReqInvalid(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	// robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*t2g_obj.ReqInvalid_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := t2g_packet.Header{
		ErrorCode: t2g_error.None,
	}
	sendBody := &t2g_obj.RspInvalid_data{}
	return sendHeader, sendBody, nil
}

// Register
func bytesAPIFn_ReqRegister(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	// robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*t2g_obj.ReqRegister_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := t2g_packet.Header{
		ErrorCode: t2g_error.None,
	}
	sendBody := &t2g_obj.RspRegister_data{}
	return sendHeader, sendBody, nil
}

// Heartbeat
func bytesAPIFn_ReqHeartbeat(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	// robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*t2g_obj.ReqHeartbeat_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := t2g_packet.Header{
		ErrorCode: t2g_error.None,
	}
	sendBody := &t2g_obj.RspHeartbeat_data{}
	return sendHeader, sendBody, nil
}

// HighScore
func bytesAPIFn_ReqHighScore(
	me interface{}, hd t2g_packet.Header, rbody []byte) (
	t2g_packet.Header, interface{}, error) {
	// robj, err := t2g_json.UnmarshalPacket(hd, rbody)
	// if err != nil {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", rbody)
	// }
	// recvBody, ok := robj.(*t2g_obj.ReqHighScore_data)
	// if !ok {
	// 	return hd, nil, fmt.Errorf("Packet type miss match %v", robj)
	// }
	// _ = recvBody

	sendHeader := t2g_packet.Header{
		ErrorCode: t2g_error.None,
	}
	sendBody := &t2g_obj.RspHighScore_data{}
	return sendHeader, sendBody, nil
}