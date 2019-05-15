package commands

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gcla/gowid/widgets/vpadding"
)

// UI ...
func UI() error {

	helloworld := text.NewFromContentExt(
		text.NewContent([]text.ContentSegment{
			text.StringContent("Hello World"),
		}),
		text.Options{
			Align: gowid.HAlignMiddle{},
		},
	)

	f := gowid.RenderFlow{}

	view := vpadding.New(
		pile.New([]gowid.IContainerWidget{
			&gowid.ContainerWidget{IWidget: helloworld, D: f},
		}),
		gowid.VAlignTop{},
		f,
	)

	app, _ := gowid.NewApp(gowid.AppArgs{
		View: view,
	})

	app.SimpleMainLoop()
	return nil
}
