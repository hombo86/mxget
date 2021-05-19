package provider

import (
	"errors"

	"mxget/pkg/api"
	"mxget/pkg/provider/baidu"
	"mxget/pkg/provider/kugou"
	"mxget/pkg/provider/kuwo"
	"mxget/pkg/provider/migu"
	"mxget/pkg/provider/netease"
	"mxget/pkg/provider/tencent"
)

func GetClient(platform string) (api.Provider, error) {
	switch platform {
	case "netease", "nc":
		return netease.Client(), nil
	case "tencent", "qq":
		return tencent.Client(), nil
	case "migu", "mg":
		return migu.Client(), nil
	case "kugou", "kg":
		return kugou.Client(), nil
	case "kuwo", "kw":
		return kuwo.Client(), nil
	case "qianqian", "baidu", "bd":
		return baidu.Client(), nil
	default:
		return nil, errors.New("unexpected music platform")
	}
}

func GetDesc(platform string) string {
	switch platform {
	case "netease", "nc":
		return "netease cloud music"
	case "tencent", "qq":
		return "qq music"
	case "migu", "mg":
		return "migu music"
	case "kugou", "kg":
		return "kugou music"
	case "kuwo", "kw":
		return "kuwo music"
	case "qianqian", "baidu", "bd":
		return "qianqian music"
	default:
		return "unknown"
	}
}
