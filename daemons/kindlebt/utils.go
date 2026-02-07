package main

import (
	"slices"

	kbt "github.com/Sighery/gokindlebt"
)

func LittleEndianUUID(uuid kbt.CharacteristicUuid) kbt.CharacteristicUuid {
	bytes := uuid.Bytes
	slices.Reverse(bytes[:])
	return kbt.CharacteristicUuid{Bytes: bytes}
}
