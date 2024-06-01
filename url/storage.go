package url

type Storage interface {
	Shorten(url string, expSecond int64) (string, error)
	ShortLinkInfo(sid string) (*UrlDetailInfo, error)
	Unshorten(sid string) (string, error)
}
