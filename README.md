# Turnkey: A Kindle actions orchestrator

I want to perform different actions on the Kindle, like turning pages, or
changing the rotation, when certain events happen, like a connected BLE ring
making a specific gesture.

This project relies heavily on [KindleBT][] and [GoKindleBT][] for the
Bluetooth connectivity in the events input part of the program.

This project is still a WIP, but so far the following actions have been
already implemented:

* Page turns: through [X11][] as well as [evdev][]
* Orientation: through a combination of [X11][] and [lipc][]

The events input side is still heavily under development, with so far only
Bluetooth being supported. Right now, the only Bluetooth stack supported is
[KindleBT][]. To do this, we use [a gRPC API][daemon API].

This has two reasons: first, the KindleBT stack, which relies on the existing
[Bluedroid][] stack used by Kindles, can only run under the `bluetooth` user.
This means that as soon as we drop privileges to that user, we won't be able
to interact with X11 and other system resources.

Additionally, there's been some further development in this area and
[another Bluetooth stack by zampierilucas][zampierilucas' BT stack], based on
Google's [Bumble][], was developed since I started working on KindleBT. Using
this approach will allow us to support multiple Bluetooth stacks through a
common gRCP API.

You can find the current KindleBT gRPC server under [daemons/kindlebt/][].

[KindleBT]: https://github.com/Sighery/kindlebt
[GoKindleBT]: https://github.com/Sighery/gokindlebt
[X11]: https://www.x.org/releases/current/doc/
[evdev]: https://docs.kernel.org/input/input.html
[lipc]: https://github.com/clintharrison/liblipc-go
[daemon API]: ./daemons/proto/btdaemon.proto
[Bluedroid]: https://android.googlesource.com/platform/system/bt/+/f8881fe
[zampierilucas' BT stack]: https://github.com/zampierilucas/kindle-hid-passthrough
[Bumble]: https://github.com/google/bumble
[daemons/kindlebt/]: ./daemons/kindlebt/
