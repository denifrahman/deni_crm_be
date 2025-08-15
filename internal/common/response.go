package common

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(data interface{}) BaseResponse {
	return BaseResponse{
		Status: "ok",
		Data:   data,
	}
}

func SuccessWithMeta(data, meta interface{}) BaseResponse {
	return BaseResponse{
		Status: "ok",
		Data:   data,
		Meta:   meta,
	}
}

func Error(message string) BaseResponse {
	return BaseResponse{
		Status:  "error",
		Message: message,
	}
}
