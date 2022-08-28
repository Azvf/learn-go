package main

import (
	"get-bili/downloader"
	myfmt "get-bili/format"
)

func main() {
	request := downloader.InfoRequest{Bvids: []string{"BV1AB4y147fg", "BV1Ff4y187q9"}}
	response, err := downloader.BatchDownloadVideoInfo(request)
	if err != nil {
		panic(err)
	}

	for _, info := range response.Infos {
		myfmt.Logger.Printf("title: %s\ndesc: %s\n\n", info.Data.Title, info.Data.Desc)
	}

}
