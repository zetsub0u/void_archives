package loaders

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
)

type videoData struct {
	Title string
	ID    string
}

func Youtube() error {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("GAPI_KEY")))
	if err != nil {
		return fmt.Errorf("failed initializing youtube poller: %w", err)
	}

	r, err := youtubeService.Channels.List([]string{"contentDetails"}).Id("UCGHXGwWtG76ltsbTuqfkE3A").Do()
	if err != nil {
		return fmt.Errorf("failed fetching videos: %w", err)
	}

	r2, err := youtubeService.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(r.Items[0].ContentDetails.RelatedPlaylists.Uploads).Do()
	if err != nil {
		return fmt.Errorf("failed fetching videos: %w", err)
	}

	var videos map[string]string
	for _, v := range r2.Items {
		r3, err := youtubeService.Videos.List([]string{"contentDetails", "snippet"}).Id(v.ContentDetails.VideoId).Do()
		if err != nil {
			return fmt.Errorf("error getting video: %w", err)
		}
		fmt.Printf("response2: %#v\n", r3)
		videos[v.ContentDetails.VideoId] = ""
	}

	fmt.Printf("response: %#v\n", r)
	fmt.Printf("response2: %#v\n", r2)

	return nil
}
