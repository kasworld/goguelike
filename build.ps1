$PROTOCOL_T2G_VERSION=makesha256sum protocol_t2g/*.enum protocol_t2g/t2g_obj/protocol_noti.go protocol_t2g/t2g_obj/protocol_cmd.go

$PROTOCOL_C2T_VERSION=makesha256sum protocol_c2t/*.enum protocol_c2t/c2t_obj/protocol_objects.go protocol_c2t/c2t_obj/protocol_noti.go protocol_c2t/c2t_obj/protocol_admin.go protocol_c2t/c2t_obj/protocol_aoact.go protocol_c2t/c2t_obj/protocol_cmd.go

$Data_VERSION=makesha256sum config/gameconst/gameconst.go config/gameconst/serviceconst.go config/gamedata/*.go enum/*.enum

echo "Protocol T2G Version:" ${PROTOCOL_T2G_VERSION}
echo "Protocol C2T Version:" ${PROTOCOL_C2T_VERSION}
echo "Data Version:" ${Data_VERSION}
