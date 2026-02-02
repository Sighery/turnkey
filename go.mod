module github.com/Sighery/turnkey

go 1.25.3

replace github.com/Sighery/turnkey/daemons => ./daemons

require (
	github.com/Sighery/turnkey/daemons v0.0.0-00010101000000-000000000000
	github.com/clintharrison/liblipc-go v0.0.0-20251208145554-2c265cf38c47
	github.com/godbus/dbus/v5 v5.2.0
	github.com/holoplot/go-evdev v0.0.0-20250804134636-ab1d56a1fe83
	google.golang.org/grpc v1.78.0
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
