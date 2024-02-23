package enmu

type MessageType uint8

const (
	CONNECT     MessageType = 1
	CONNACK     MessageType = 2
	PUBLISH     MessageType = 3
	PUBACK      MessageType = 4
	PUBREC      MessageType = 5
	PUBREL      MessageType = 6
	PUBCOMP     MessageType = 7
	SUBSCRIBE   MessageType = 8
	SUBACK      MessageType = 9
	UNSUBSCRIBE MessageType = 10
	UNSUBACK    MessageType = 11
	PINGREQ     MessageType = 12
	PINGRESP    MessageType = 13
	DISCONNECT  MessageType = 14
	AUTH        MessageType = 15 // AUTH v5.0新增
)
