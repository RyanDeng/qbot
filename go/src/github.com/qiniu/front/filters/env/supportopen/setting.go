package supportopen

var (
	Env *Setting
)

type Setting struct {
	StartTime string `conf:"start_time"`
	EndTime   string `conf:"end_time"`
	AdminUids string `conf:"admin_uids"`
}
