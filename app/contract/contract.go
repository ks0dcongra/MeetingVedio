package contract

const (
	Success = iota
	ERROR_Upload_Video
	ERROR_Download_Video
)

var Message = map[int]string{
	Success:              "Upload video successfully",
	ERROR_Upload_Video:   "Upload video fail",
	ERROR_Download_Video: "Download video fail",
}
