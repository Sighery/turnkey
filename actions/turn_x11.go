package actions

import (
	"context"
	"fmt"

	wtitle "github.com/Sighery/turnkey/services/windowtitles"
	x11api "github.com/Sighery/turnkey/services/x11"
)

type X11PageTurnerProvider struct {
	client *x11api.X11
}

func NewX11PageTurnerProvider(client *x11api.X11) X11PageTurnerProvider {
	return X11PageTurnerProvider{client: client}
}

func (p X11PageTurnerProvider) Close() error {
	// Leave closing of the X11 client to the one that passed us the reference
	return nil
}

func (p X11PageTurnerProvider) TurnPage(ctx context.Context, direction PageTurnDirection) error {
	var key x11api.XKeysym
	switch direction {
	case TurnDirectionForward:
		key = x11api.XKPageDown
	case TurnDirectionBackward:
		key = x11api.XKPageUp
	default:
		return fmt.Errorf("Unknown page direction %d", direction)
	}

	window, err := p.client.ActiveWindow()
	if err != nil {
		return err
	}

	windowName, err := p.client.WindowName(window)
	if err != nil {
		return err
	}

	title := wtitle.KindleTitle(windowName)
	if !title.IsApplication(wtitle.ApplicationReader) {
		return fmt.Errorf("Cannot turn page outside the native reader application")
	}

	return p.client.SendKeyPress(window, key)
}
