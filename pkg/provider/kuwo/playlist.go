package kuwo

import (
	"context"
	"errors"
	"fmt"
	"github.com/winterssy/ghttp"
	"github.com/winterssy/glog"
	"strconv"
	"strings"
	"time"

	"mxget/pkg/api"
)

func (a *API) GetPlaylist(ctx context.Context, playlistId string) (*api.Collection, error) {
	resp, err := a.GetPlaylistRaw(ctx, playlistId, 1, 9999)
	if err != nil {
		return nil, err
	}

	n := len(resp.Data.MusicList)
	if n == 0 {
		return nil, errors.New("get playlist: no data")
	}
	a.patchSongsURL(ctx, songDefaultBR, resp.Data.MusicList...)
	a.patchSongsLyric(ctx, resp.Data.MusicList...)
	songs := translate(resp.Data.MusicList...)

	glog.Info("songs>>>",len(songs))

	return &api.Collection{
		Id:     strconv.Itoa(resp.Data.Id),
		Name:   strings.TrimSpace(resp.Data.Name),
		PicURL: resp.Data.Img700,
		Songs:  songs,
	}, nil
}

// 获取歌单，page: 页码； pageSize: 每页数量，如果要获取全部请设置为较大的值
func (a *API) GetPlaylistRaw(ctx context.Context, playlistId string, page int, pageSize int) (*PlaylistResponse, error) {
	resp := new(PlaylistResponse)
	respTemp := new(PlaylistResponse)

	var params ghttp.Params
	var req *ghttp.Request
	var respSource *ghttp.Response
	var err error

	for i := 1; i <= 100; i++ {
		params = ghttp.Params{
			"pid": playlistId,
			"pn":  i,
			"rn":  100,
		}
		respTemp = new(PlaylistResponse)
		req, _ = ghttp.NewRequest(ghttp.MethodGet, apiGetPlaylist)
		req.SetQuery(params)
		req.SetContext(ctx)
		respSource, err = a.SendRequest(req)
		if err == nil {
			err = respSource.JSON(respTemp)
		}
		if err != nil {
			// 重试
			i--
			glog.Info(fmt.Errorf("send request: %s", err))
			continue
			//return nil, err
		}
		if len(respTemp.Data.MusicList) > 0 {
			resp.Data.MusicList = append(resp.Data.MusicList, respTemp.Data.MusicList...)
		} else {
			break
		}

		if respTemp.Code != 200 {
			glog.Info(fmt.Errorf("get playlist: %s", resp.errorMessage()))
			continue
			//return nil, fmt.Errorf("get playlist: %s", resp.errorMessage())
		}
		time.Sleep(1)
	}

	glog.Info(fmt.Sprintf("playlist size: %d ",len(resp.Data.MusicList)))
	return resp, nil
}
