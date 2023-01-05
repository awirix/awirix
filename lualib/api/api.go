package api

import (
	"github.com/vivi-app/vivi/luadoc"
)

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        "api",
		Description: "Vivi api. Used to interact with the system. For example, to open an app or watch a video.",
		Funcs: []*luadoc.Func{
			{
				Name:        "watch",
				Description: "Opens a video in the default video player. It can be a local file or a url.",
				Value:       watch,
				Params: []*luadoc.Param{
					{
						Name:        "target",
						Description: "The video to watch. It can be a local file or a url.",
						Type:        luadoc.String,
					},
				},
			},
			{
				Name:        "open",
				Description: "Opens a file in the default app. It can be a local file or a url.",
				Value:       openDefault,
				Params: []*luadoc.Param{
					{
						Name:        "target",
						Description: "The file to open. It can be a local file or a url.",
						Type:        luadoc.String,
					},
					{
						Name:        "app",
						Description: "The app to open the file with. If not specified, the default app will be used.",
						Type:        luadoc.String,
						Opt:         true,
					},
				},
			},
			{
				Name:        "open_data",
				Description: "Saves data to a temp file and opens it in the default app.",
				Value:       openData,
				Params: []*luadoc.Param{
					{
						Name:        "data",
						Description: "The data to save.",
						Type:        luadoc.String,
					},
					{
						Name:        "ext",
						Description: "The extension of the file to save.",
						Type:        luadoc.String,
					},
					{
						Name:        "app",
						Description: "The app to open the file with. If not specified, the default app will be used.",
						Type:        luadoc.String,
						Opt:         true,
					},
				},
			},
			{
				Name:        "save",
				Description: "Saves data to a file under extension downloads directory.",
				Value:       save,
				Params: []*luadoc.Param{
					{
						Name:        "data",
						Description: "The data to save.",
						Type:        luadoc.String,
					},
					{
						Name:        "path",
						Description: "The path to save the file to. It will be joined with extension downloads directory.",
						Type:        luadoc.String,
					},
				},
			},
		},
	}
}
