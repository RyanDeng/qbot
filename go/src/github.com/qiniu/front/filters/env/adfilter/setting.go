package adfilter

var (
	Env *Setting
)

type Setting struct {
	StartTime    string `conf:"start_time"`
	EndTime      string `conf:"end_time"`
	SyncInterval string `conf:"sync_interval"`
}
