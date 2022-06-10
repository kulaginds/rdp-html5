package asn1

const (
	ClassMask            uint8 = 0xC0
	ClassUniversal       uint8 = 0x00
	ClassApplication     uint8 = 0x40
	ClassContextSpecific uint8 = 0x80
	ClassPrivate         uint8 = 0xC0
)

const (
	PCMask      uint8 = 0x20
	PCPrimitive uint8 = 0x00
	PCConstruct uint8 = 0x20
)

const (
	TagMask           uint8 = 0x1F
	TagBoolean        uint8 = 0x01
	TagInteger        uint8 = 0x02
	TagBitString      uint8 = 0x03
	TagOctetString    uint8 = 0x04
	TagObjectIdenfier uint8 = 0x06
	TagEnumerated     uint8 = 0x0A
	TagSequence       uint8 = 0x10
	TagSequenceOf     uint8 = 0x10
)
