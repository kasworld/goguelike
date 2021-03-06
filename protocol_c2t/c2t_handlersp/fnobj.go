// Code generated by "genprotocol.exe -ver=fc37e02b6858cffd9591410bf9ff4f28fcf1782014d44a7d0e102918f2b1f57d -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_handlersp

import (
	"fmt"

	"github.com/kasworld/goguelike/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_obj"
	"github.com/kasworld/goguelike/protocol_c2t/c2t_packet"
)

// obj base demux fn map

var DemuxRsp2ObjFnMap = [...]func(me interface{}, hd c2t_packet.Header, body interface{}) error{
	c2t_idcmd.Invalid:           objRecvRspFn_Invalid,           // Invalid make empty packet error
	c2t_idcmd.Login:             objRecvRspFn_Login,             // Login
	c2t_idcmd.Heartbeat:         objRecvRspFn_Heartbeat,         // Heartbeat
	c2t_idcmd.Chat:              objRecvRspFn_Chat,              // Chat
	c2t_idcmd.AchieveInfo:       objRecvRspFn_AchieveInfo,       // AchieveInfo
	c2t_idcmd.Rebirth:           objRecvRspFn_Rebirth,           // Rebirth
	c2t_idcmd.MoveFloor:         objRecvRspFn_MoveFloor,         // MoveFloor tower cmd
	c2t_idcmd.AIPlay:            objRecvRspFn_AIPlay,            // AIPlay
	c2t_idcmd.VisitFloorList:    objRecvRspFn_VisitFloorList,    // VisitFloorList floor info of visited
	c2t_idcmd.Meditate:          objRecvRspFn_Meditate,          // Meditate rest and recover HP,SP
	c2t_idcmd.KillSelf:          objRecvRspFn_KillSelf,          // KillSelf
	c2t_idcmd.Move:              objRecvRspFn_Move,              // Move move 8way near tile
	c2t_idcmd.Attack:            objRecvRspFn_Attack,            // Attack attack near 1 tile
	c2t_idcmd.AttackWide:        objRecvRspFn_AttackWide,        // AttackWide attack near 3 tile
	c2t_idcmd.AttackLong:        objRecvRspFn_AttackLong,        // AttackLong attack 3 tile to direction
	c2t_idcmd.Pickup:            objRecvRspFn_Pickup,            // Pickup pickup carryobj
	c2t_idcmd.Drop:              objRecvRspFn_Drop,              // Drop drop carryobj
	c2t_idcmd.Equip:             objRecvRspFn_Equip,             // Equip equip equipable carryobj
	c2t_idcmd.UnEquip:           objRecvRspFn_UnEquip,           // UnEquip unequip equipable carryobj
	c2t_idcmd.DrinkPotion:       objRecvRspFn_DrinkPotion,       // DrinkPotion
	c2t_idcmd.ReadScroll:        objRecvRspFn_ReadScroll,        // ReadScroll
	c2t_idcmd.Recycle:           objRecvRspFn_Recycle,           // Recycle sell carryobj
	c2t_idcmd.EnterPortal:       objRecvRspFn_EnterPortal,       // EnterPortal
	c2t_idcmd.ActTeleport:       objRecvRspFn_ActTeleport,       // ActTeleport
	c2t_idcmd.AdminTowerCmd:     objRecvRspFn_AdminTowerCmd,     // AdminTowerCmd generic cmd
	c2t_idcmd.AdminFloorCmd:     objRecvRspFn_AdminFloorCmd,     // AdminFloorCmd generic cmd
	c2t_idcmd.AdminActiveObjCmd: objRecvRspFn_AdminActiveObjCmd, // AdminActiveObjCmd generic cmd
	c2t_idcmd.AdminFloorMove:    objRecvRspFn_AdminFloorMove,    // AdminFloorMove Next Before floorUUID
	c2t_idcmd.AdminTeleport:     objRecvRspFn_AdminTeleport,     // AdminTeleport random pos in floor
	c2t_idcmd.AdminAddExp:       objRecvRspFn_AdminAddExp,       // AdminAddExp  add arg to battle exp
	c2t_idcmd.AdminPotionEffect: objRecvRspFn_AdminPotionEffect, // AdminPotionEffect buff by arg potion type
	c2t_idcmd.AdminScrollEffect: objRecvRspFn_AdminScrollEffect, // AdminScrollEffect buff by arg Scroll type
	c2t_idcmd.AdminCondition:    objRecvRspFn_AdminCondition,    // AdminCondition add arg condition for 100 turn
	c2t_idcmd.AdminAddPotion:    objRecvRspFn_AdminAddPotion,    // AdminAddPotion add arg potion to inven
	c2t_idcmd.AdminAddScroll:    objRecvRspFn_AdminAddScroll,    // AdminAddScroll add arg scroll to inven
	c2t_idcmd.AdminAddMoney:     objRecvRspFn_AdminAddMoney,     // AdminAddMoney add arg money to inven
	c2t_idcmd.AdminAddEquip:     objRecvRspFn_AdminAddEquip,     // AdminAddEquip add random equip to inven
	c2t_idcmd.AdminForgetFloor:  objRecvRspFn_AdminForgetFloor,  // AdminForgetFloor forget current floor map
	c2t_idcmd.AdminFloorMap:     objRecvRspFn_AdminFloorMap,     // AdminFloorMap complete current floor map

}

// Invalid make empty packet error
func objRecvRspFn_Invalid(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspInvalid_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Login
func objRecvRspFn_Login(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspLogin_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Heartbeat
func objRecvRspFn_Heartbeat(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspHeartbeat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Chat
func objRecvRspFn_Chat(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspChat_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AchieveInfo
func objRecvRspFn_AchieveInfo(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAchieveInfo_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Rebirth
func objRecvRspFn_Rebirth(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspRebirth_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// MoveFloor tower cmd
func objRecvRspFn_MoveFloor(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspMoveFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AIPlay
func objRecvRspFn_AIPlay(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAIPlay_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// VisitFloorList floor info of visited
func objRecvRspFn_VisitFloorList(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspVisitFloorList_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Meditate rest and recover HP,SP
func objRecvRspFn_Meditate(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspMeditate_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// KillSelf
func objRecvRspFn_KillSelf(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspKillSelf_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Move move 8way near tile
func objRecvRspFn_Move(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspMove_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Attack attack near 1 tile
func objRecvRspFn_Attack(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAttack_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AttackWide attack near 3 tile
func objRecvRspFn_AttackWide(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAttackWide_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AttackLong attack 3 tile to direction
func objRecvRspFn_AttackLong(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAttackLong_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Pickup pickup carryobj
func objRecvRspFn_Pickup(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspPickup_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Drop drop carryobj
func objRecvRspFn_Drop(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspDrop_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Equip equip equipable carryobj
func objRecvRspFn_Equip(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// UnEquip unequip equipable carryobj
func objRecvRspFn_UnEquip(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspUnEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// DrinkPotion
func objRecvRspFn_DrinkPotion(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspDrinkPotion_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// ReadScroll
func objRecvRspFn_ReadScroll(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspReadScroll_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// Recycle sell carryobj
func objRecvRspFn_Recycle(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspRecycle_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// EnterPortal
func objRecvRspFn_EnterPortal(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspEnterPortal_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// ActTeleport
func objRecvRspFn_ActTeleport(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspActTeleport_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminTowerCmd generic cmd
func objRecvRspFn_AdminTowerCmd(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminTowerCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminFloorCmd generic cmd
func objRecvRspFn_AdminFloorCmd(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminFloorCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminActiveObjCmd generic cmd
func objRecvRspFn_AdminActiveObjCmd(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminActiveObjCmd_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminFloorMove Next Before floorUUID
func objRecvRspFn_AdminFloorMove(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminFloorMove_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminTeleport random pos in floor
func objRecvRspFn_AdminTeleport(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminTeleport_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminAddExp  add arg to battle exp
func objRecvRspFn_AdminAddExp(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminAddExp_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminPotionEffect buff by arg potion type
func objRecvRspFn_AdminPotionEffect(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminPotionEffect_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminScrollEffect buff by arg Scroll type
func objRecvRspFn_AdminScrollEffect(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminScrollEffect_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminCondition add arg condition for 100 turn
func objRecvRspFn_AdminCondition(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminCondition_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminAddPotion add arg potion to inven
func objRecvRspFn_AdminAddPotion(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminAddPotion_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminAddScroll add arg scroll to inven
func objRecvRspFn_AdminAddScroll(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminAddScroll_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminAddMoney add arg money to inven
func objRecvRspFn_AdminAddMoney(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminAddMoney_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminAddEquip add random equip to inven
func objRecvRspFn_AdminAddEquip(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminAddEquip_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminForgetFloor forget current floor map
func objRecvRspFn_AdminForgetFloor(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminForgetFloor_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}

// AdminFloorMap complete current floor map
func objRecvRspFn_AdminFloorMap(me interface{}, hd c2t_packet.Header, body interface{}) error {
	robj, ok := body.(*c2t_obj.RspAdminFloorMap_data)
	if !ok {
		return fmt.Errorf("packet mismatch %v", body)
	}
	return fmt.Errorf("Not implemented %v", robj)
}
