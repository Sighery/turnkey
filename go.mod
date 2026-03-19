module github.com/Sighery/turnkey

go 1.25.3

replace github.com/Sighery/turnkey/daemons => ./daemons

require (
	github.com/Sighery/turnkey/daemons v0.0.0-00010101000000-000000000000
	github.com/clintharrison/liblipc-go v0.0.0-20251208145554-2c265cf38c47
	github.com/godbus/dbus/v5 v5.2.0
	github.com/holoplot/go-evdev v0.0.0-20250804134636-ab1d56a1fe83
	google.golang.org/grpc v1.79.3
)

require (
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
