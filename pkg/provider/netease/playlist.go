package netease

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/winterssy/ghttp"
	"mxget/pkg/api"
	"mxget/pkg/utils"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.Collection, error) {
	_playlistId, err := strconv.Atoi(playlistId)
	if err != nil {
		return nil, err
	}

	resp, err := a.GetPlaylistRaw(ctx, _playlistId)
	if err != nil {
		return nil, err
	}

	n := resp.Playlist.Total
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}

	tracks := resp.Playlist.Tracks
	// 修复网易云音乐数量限制
	trackSize := len(tracks)
	if n > trackSize {
		extra := n - trackSize
		trackIds := make([]int, extra)
		for i, j := trackSize, 0; i < n; i, j = i+1, j+1 {
			trackIds[j] = resp.Playlist.TrackIds[i].Id
		}

		queue := make(chan []*Song)
		wg := new(sync.WaitGroup)
		for i := 0; i < extra; i += songRequestLimit {
			if ctx.Err() != nil {
				break
			}

			songIds := trackIds[i:utils.Min(i+songRequestLimit, extra)]
			wg.Add(1)
			go func() {
				resp, err := a.GetSongsRaw(ctx, songIds...)
				if err != nil {
					wg.Done()
					return
				}
				queue <- resp.Songs
			}()
		}

		go func() {
			for s := range queue {
				tracks = append(tracks, s...)
				wg.Done()
			}
		}()
		wg.Wait()
	}

	a.patchSongsURL(ctx, songDefaultBR, tracks...)
	a.patchSongsLyric(ctx, tracks...)
	songs := translate(tracks...)
	return &api.Collection{
		Id:     strconv.Itoa(resp.Playlist.Id),
		Name:   strings.TrimSpace(resp.Playlist.Name),
		PicURL: resp.Playlist.PicURL,
		Songs:  songs,
	}, nil
}

// 获取歌单
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId int) (*PlaylistResponse, error) {
	data := map[string]int{
		"id": playlistId,
		"n":  100000,
	}

	resp := new(PlaylistResponse)
	req, _ := ghttp.NewRequest(ghttp.MethodPost, apiGetPlaylist)
	req.SetForm(weapi(data))
	req.SetContext(ctx)
	r, err := a.SendRequest(req)
	if err == nil {
		err = r.JSON(resp)
	}
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
	}

	return resp, nil
}
