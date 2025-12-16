package main

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/godbus/dbus/v5"

	"github.com/Sighery/turnkey/actions"
	lipcapi "github.com/Sighery/turnkey/services/lipc"
	wtitle "github.com/Sighery/turnkey/services/windowtitles"
	x11api "github.com/Sighery/turnkey/services/x11"
)

func UseBluetoothPrivileges() error {
	user, err := user.Lookup("bluetooth")
	if err != nil {
		return err
	}
	userId, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}
	err = syscall.Setuid(userId)
	if err != nil {
		return err
	}
	return nil
}


func main() {
	log.Printf("Hello World\n")

	dbusConn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal(err)
	}
	lipc := lipcapi.NewLipcClient(dbusConn)

	ctx := context.TODO()

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
