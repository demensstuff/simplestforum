package delivery

import (
	"encoding/json"
	"simplestforum/internal/domain/entity"
)

// RESTResponse wraps any REST response.
type RESTResponse struct {
	Data *CustomResponse `json:"data"`
}

// CustomResponse represents the general format for any response.
type CustomResponse struct {
	IsOK    bool        `json:"is_ok"`
	Payload interface{} `json:"payload"`
}

// buildResponse creates a response.
func buildResponse(isOK bool, payload interface{}, shouldWrap bool) []byte {
	var preparedData interface{}

	if shouldWrap {
		preparedData = RESTResponse{
			Data: &CustomResponse{
				IsOK:    isOK,
				Payload: payload,
			},
		}
	} else {
		preparedData = CustomResponse{
			IsOK:    isOK,
			Payload: payload,
		}
	}

	data, err := json.Marshal(preparedData)

	if err != nil {
		panic("Cannot encode the response!")
	}

	return data
}

// BuildSuccessResponse creates a response which indicates success.
func BuildSuccessResponse(payload interface{}, shouldWrap bool) []byte {
	return buildResponse(true, payload, shouldWrap)
}

// BuildErrorResponse creates an erroneous  response.
func BuildErrorResponse(sess entity.Session, err error, shouldWrap bool) []byte {
	return buildResponse(false, entity.SessionInfoToError(sess, err), shouldWrap)
}
