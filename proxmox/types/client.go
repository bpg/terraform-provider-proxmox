package types

import "context"

type Client interface {
	DoRequest(
		ctx context.Context,
		method, path string,
		requestBody, responseBody interface{},
	) error

	AddPrefix(path string) string
}
