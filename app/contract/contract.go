package contract

const (
	Success = iota
	ERROR_Upload_Video
	ERROR_Download_Video
	ERROR_Download_Video_NotExist
	ERROR_Video_To_Hex
	ERROR_Upload_Video_NotExist
	ERROR_Video_Format
)

var Message = map[int]string{
	Success:                       "Upload video successfully",
	ERROR_Upload_Video:            "Upload video fail",
	ERROR_Download_Video:          "Download video fail",
	ERROR_Download_Video_NotExist: "videoID not found",
	ERROR_Video_To_Hex:            "video bytes to hex fail",
	ERROR_Upload_Video_NotExist:   "video not found",
	ERROR_Video_Format:            "video format has some problem",
}
