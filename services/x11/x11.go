// Heavily inspired by Clint's implementation:
// https://github.com/clintharrison/bueno/blob/aebf3574fc7e10937635195da5d38ca427b254b3/xkb/xkb.go

package x11

/*
#include <stdlib.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/keysym.h>
*/
import "C"

import (
	"fmt"
	"log/slog"
	"runtime"
	"unsafe"
)

type XKeysym uint32

const (
	XKPageUp   XKeysym = 0xFF55
	XKPageDown XKeysym = 0xFF56
)

func (x XKeysym) String() string {
	return map[XKeysym]string{
		XKPageUp:   "Page Up",
		XKPageDown: "Page Down",
	}[x]
}

//export onX11Error
func onX11Error(display *C.Display, error *C.XErrorEvent) {
	slog.Error("X11 error", "display", display, "error", error)
}

type X11 struct {
	display    *C.Display
	rootWindow C.Window
}

func OpenX11() (*X11, error) {
	runtime.LockOSThread()

	x := &X11{}

	x.display = C.XOpenDisplay(nil)
	if x.display == nil {
		return nil, fmt.Errorf("Failed to open X display")
	}

	x.rootWindow = C.XDefaultRootWindow(x.display)
	return x, nil
}

func (x *X11) Close() error {
	C.XCloseDisplay(x.display)
	runtime.UnlockOSThread()
	return nil
}

func (x *X11) ActiveWindow() (C.Window, error) {
	atomName := C.CString("_NET_ACTIVE_WINDOW")
	defer C.free(unsafe.Pointer(atomName))

	atom := C.XInternAtom(x.display, atomName, C.True)
	if atom == 0 {
		return 0, fmt.Errorf("Couldn't find _NET_ACTIVE_WINDOW atom")
	}

	var actualType C.Atom
	var actualFormat C.int
	var nItems C.ulong
	var bytesAfter C.ulong
	var prop *C.uchar
	status := C.XGetWindowProperty(
		x.display,
		x.rootWindow,
		atom,
		0,
		1,
		C.False,
		C.XA_WINDOW,
		&actualType,
		&actualFormat,
		&nItems,
		&bytesAfter,
		&prop,
	)

	if status != C.Success || nItems == 0 || actualFormat == 0 {
		err := fmt.Errorf(
			"Failed to get active window. Status %d, nItems %d, actualFormat %d",
			status, nItems, actualFormat,
		)
		slog.Warn("XGetWindowProperty failed", "error", err)
		return 0, err
	}

	propData := C.GoBytes(unsafe.Pointer(prop), C.int(nItems)*(actualFormat/8))
	defer C.XFree(unsafe.Pointer(prop))

	if len(propData) != 4 {
		err := fmt.Errorf(
			"Expected 32-bit value in _NET_ACTIVE_WINDOW prop (i.e., actualFormat=32). "+
				"Length %d, actualFormat %d", len(propData), actualFormat,
		)
		slog.Warn("Unexpected _NET_ACTIVE_WINDOW prop", "error", err)
		return 0, err
	}

	activeWindowId := uint32(propData[0]) | (uint32(propData[1]) << 8) |
		(uint32(propData[2]) << 16) | (uint32(propData[3]) << 24)
	slog.Debug(
		"Got active window", "activeWindowId", activeWindowId, "actualType", actualType,
		"actualFormat", actualFormat, "nItems", nItems, "bytesAfter", bytesAfter,
	)
	return C.Window(activeWindowId), nil
}

func (x *X11) WindowName(window C.Window) (string, error) {
	var name *C.char
	status := C.XFetchName(x.display, window, &name)
	if status == 0 || name == nil {
		return "", fmt.Errorf("XFetchName call failed %d", status)
	}
	defer C.XFree(unsafe.Pointer(name))
	return C.GoString(name), nil
}

func (x *X11) SetWindowName(window C.Window, name string) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))

	C.XStoreName(x.display, window, cstr)
}

func (x *X11) SendKeyPress(window C.Window, key XKeysym) error {
	keycode := C.XKeysymToKeycode(x.display, C.KeySym(key))

	evt := C.XKeyEvent{
		display:   x.display,
		window:    window,
		subwindow: C.None,
		keycode:   C.uint(keycode),
		// TODO: handle state (e.g. modifier keys)? is that ever used for anything on kindle?
		state:       0,
		root:        x.rootWindow,
		same_screen: C.True,
		_type:       C.KeyPress,
		// Unsure if these need to be set.
		// https://github.com/jordansissel/xdotool/blob/33092d8a74d60c9ad3ab39c4f05b90e047ea51d8/xdo.c#L1517-L1518
		x:      C.int(1),
		y:      C.int(1),
		x_root: C.int(1),
		y_root: C.int(1),
	}

	C.XSendEvent(x.display, window, C.True, C.KeyPressMask, (*C.XEvent)(unsafe.Pointer(&evt)))

	evt._type = C.KeyRelease
	C.XSendEvent(x.display, window, C.True, C.KeyPressMask, (*C.XEvent)(unsafe.Pointer(&evt)))

	// xdotool doesn't know if this is needed.
	// https://github.com/jordansissel/xdotool/blob/33092d8a74d60c9ad3ab39c4f05b90e047ea51d8/xdo.c#L1103-L1104
	//
	// It is. Otherwise, we'll end up with an event lingering in the X event queue,
	// and this ends up with an off-by-one error if we send multiple key events in quick succession,
	// and a spurious key event at program exit?
	C.XFlush(x.display)

	return nil
}
