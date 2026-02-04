package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/Sighery/turnkey/actions"
	btapi "github.com/Sighery/turnkey/services/btdaemon"
	lipcapi "github.com/Sighery/turnkey/services/lipc"
	wtitle "github.com/Sighery/turnkey/services/windowtitles"
	x11api "github.com/Sighery/turnkey/services/x11"
)

var (
	btdaemonAddr = flag.String("btdaemonAddr", "0.0.0.0:50010", "ip:port of the BT daemon")
)

var (
	OnVal  = []byte{0x4F, 0x46, 0x46}
	OffVal = []byte{0x4F, 0x4E}
)

func main() {
	log.Printf("Hello World\n")

	flag.Parse()

	btdaemon, err := btapi.NewBtdaemonClient(*btdaemonAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer btdaemon.Close()

	ctx := context.TODO()

	if err := btdaemon.IsReady(ctx); err != nil {
		log.Fatal(err)
	}

	bleConn, err := btdaemon.Connect(ctx, "2C:CF:67:B8:DC:3F")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to BLE? Id:", bleConn)

	charVal, err := btdaemon.ReadChar(ctx, bleConn, "ff120000000000000000000000000000")
	if err != nil {
		log.Fatal(err)
	}

	if btapi.IsASCIIPrintable(charVal) {
		fmt.Printf("Char response: %s\n", string(charVal))
	} else {
		fmt.Printf("Char response: %v\n", charVal)
	}

	var toWrite []byte
	if bytes.Equal(charVal, OffVal) {
		toWrite = OnVal
	} else {
		toWrite = OffVal
	}
	res, err := btdaemon.WriteChar(ctx, bleConn, "ff120000000000000000000000000000", true, toWrite)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Write response: %v\n", res)

	i := 0
	sub, err := btdaemon.NotifyChar(ctx, bleConn, "ff120000000000000000000000000000")
	if err != nil {
		log.Fatal(err)
	}

	for v := range sub.C {
		i += 1
		fmt.Println("Notification, got: %v\n", v)

		if i == 5 {
			sub.Stop()
		}
	}

	dbusConn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	lipc := lipcapi.NewLipcClient(dbusConn)

	orientationProvider := actions.NewLipcOrientationProvider(lipc)
	orientation, err := orientationProvider.CurrentOrientation(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Lipc Orientation? %d\n", orientation)

	x11, err := x11api.OpenX11()
	if err != nil {
		log.Fatal(err)
	}
	defer x11.Close()

	inputDevice, err := actions.FindTouchscreenInputDevice()
	if err != nil {
		log.Fatal(err)
	}
	ptProvider, err := actions.NewInputDevicePageTurnerProvider(inputDevice, orientationProvider)
	if err != nil {
		log.Fatal(err)
	}
	defer ptProvider.Close()
	ptProvider2 := actions.NewX11PageTurnerProvider(x11)
	defer ptProvider2.Close()

	activeWindow, err := x11.ActiveWindow()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Active window is %v\n", activeWindow)

	windowName, err := x11.WindowName(activeWindow)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Window name %s\n", windowName)

	title := wtitle.KindleTitle(windowName)
	orientationVal, found := title.Get(wtitle.ComponentOrientation)
	if found {
		fmt.Printf("Found component O! Val: %s\n", orientationVal)
		if title.IsApplication(wtitle.ApplicationReader) {
			fmt.Println("We are in the reader app!")

			fmt.Println("Replacing...")
			title := title.Set(wtitle.ComponentOrientation, "ULDR")
			fmt.Printf("Title is now... %s\n", title)

			x11.SetWindowName(activeWindow, string(title))
			fmt.Println("Changed window name")
		}
	} else {
		log.Printf("Not found component O\n")
	}

	for _, i := range []int{1, 2, 3, 0} {
		orientation = actions.ScreenOrientation(i * 90)
		fmt.Printf("Setting orientation to %d\n", orientation)
		err = orientationProvider.SetOrientation(ctx, orientation)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(5500 * time.Millisecond)

		fmt.Println("Turning page forward WITH TOUCH")
		ptProvider.TurnPage(ctx, actions.TurnDirectionForward)
		time.Sleep(3500 * time.Millisecond)

		fmt.Println("Turning page backward WITH X11")
		ptProvider2.TurnPage(ctx, actions.TurnDirectionBackward)
		time.Sleep(3500 * time.Millisecond)

		fmt.Println("Turning page forward WITH X11")
		ptProvider2.TurnPage(ctx, actions.TurnDirectionForward)
		time.Sleep(3500 * time.Millisecond)

		fmt.Println("Turning page backward WITH TOUCH")
		ptProvider.TurnPage(ctx, actions.TurnDirectionBackward)
		time.Sleep(3500 * time.Millisecond)
	}
}
