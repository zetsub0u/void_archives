package loaders

import (
	"context"
	"fmt"
	"github.com/zetsub0u/void_archives/archive"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"regexp"
	"time"
)

type videoData struct {
	channelID        string
	videoID          string
	videoName        string
	videoDescription string
	videoPublishDate time.Time
}

type Youtube struct {
	service   *youtube.Service
	channelID string
}

func NewYoutube(channelID string) (*Youtube, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("GAPI_KEY")))
	if err != nil {
		return nil, fmt.Errorf("failed initializing youtube poller: %w", err)
	}
	return &Youtube{service: youtubeService, channelID: channelID}, nil
}

func (y *Youtube) GetRecentRefs(since time.Duration) (archive.Refs, error) {

	// fetch our channel subscriptions first
	s, err := y.service.Subscriptions.List([]string{"contentDetails", "snippet"}).
		ChannelId(y.channelID).MaxResults(100).Do()
	if err != nil {
		return nil, fmt.Errorf("failed fetching my subscriptions: %w", err)
	}

	allVideos := make([]videoData, 0)
	for _, sub := range s.Items {
		videos, err := y.getChannelUploads(sub.Snippet.ResourceId.ChannelId)
		if err != nil {
			return nil, fmt.Errorf("failed getting channel uploads: %w", err)
		}
		allVideos = append(allVideos, videos...)
	}

	// try to filter out things that don't look like a reference video of ma/abyss runs
	refs, err := y.filterForRefs(allVideos)
	if err != nil {
		return nil, fmt.Errorf("failed filtering refs from videos: %w", err)
	}

	return refs, nil
}

func (y *Youtube) getChannelUploads(channelID string) ([]videoData, error) {
	// get the channel info, to fetch the uploads playlist id
	ch, err := y.service.Channels.List([]string{"contentDetails"}).Id(channelID).Do()
	if err != nil {
		return nil, fmt.Errorf("failed fetching channel data: %w", err)
	}

	// get the uploads from given channel
	vidIDs, err := y.service.PlaylistItems.List([]string{"snippet", "contentDetails"}).
		PlaylistId(ch.Items[0].ContentDetails.RelatedPlaylists.Uploads).MaxResults(20).Do()
	if err != nil {
		return nil, fmt.Errorf("failed fetching videos: %w", err)
	}

	videos := make([]videoData, len(vidIDs.Items))

	for i, v := range vidIDs.Items {
		ts, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)

		if err != nil {
			return nil, fmt.Errorf("failed parsing creation timestamp of video %s: %w", v.Snippet.ResourceId.VideoId, err)
		}
		videos[i] = videoData{
			channelID:        v.Snippet.ChannelId,
			videoID:          v.Snippet.ResourceId.VideoId,
			videoName:        v.Snippet.Title,
			videoDescription: v.Snippet.Description,
			videoPublishDate: ts,
		}
	}
	return videos, nil
}

func (y *Youtube) filterForRefs(videos []videoData) (archive.Refs, error) {
	filtered := make(archive.Refs, 0)
	// dumbass "heuristics"
	var re = regexp.MustCompile(`(?i)\sma\s|\[ma\]|abyss`)
	for _, v := range videos {
		if re.Match([]byte(v.videoName)) {
			filtered = append(filtered, archive.Ref{
				URL:      fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.videoID),
				Creator:  v.channelID,
				Runs:     nil,
				Parsed:   false,
				Verified: false,
			})
		}
	}
	return filtered, nil
}
