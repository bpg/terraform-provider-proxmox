package types

import "context"

type Client interface {
	DoRequest(
		ctx context.Context,
		method, path string,
		requestBody, responseBody interface{},
	) error

	AdjustPath(path string) string
}
