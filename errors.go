package pf

import "net/http"

type httpError int

func (e httpError) Error() string {
	return http.StatusText(int(e))
}

var (
	ErrBadRequest                   error = httpError(400)
	ErrUnauthorized                 error = httpError(401)
	ErrPaymentRequired              error = httpError(402)
	ErrForbidden                    error = httpError(403)
	ErrNotFound                     error = httpError(404)
	ErrMethodNotAllowed             error = httpError(405)
	ErrNotAcceptable                error = httpError(406)
	ErrProxyAuthRequired            error = httpError(407)
	ErrRequestTimeout               error = httpError(408)
	ErrConflict                     error = httpError(409)
	ErrGone                         error = httpError(410)
	ErrLengthRequired               error = httpError(411)
	ErrPreconditionFailed           error = httpError(412)
	ErrRequestEntityTooLarge        error = httpError(413)
	ErrRequestURITooLong            error = httpError(414)
	ErrUnsupportedMediaType         error = httpError(415)
	ErrRequestedRangeNotSatisfiable error = httpError(416)
	ErrExpectationFailed            error = httpError(417)
	ErrImATeapot                    error = httpError(418)
	ErrMisdirectedRequest           error = httpError(421)
	ErrUnprocessableEntity          error = httpError(422)
	ErrLocked                       error = httpError(423)
	ErrFailedDependency             error = httpError(424)
	ErrTooEarly                     error = httpError(425)
	ErrUpgradeRequired              error = httpError(426)
	ErrPreconditionRequired         error = httpError(428)
	ErrTooManyRequests              error = httpError(429)
	ErrRequestHeaderFieldsTooLarge  error = httpError(431)
	ErrUnavailableForLegalReasons   error = httpError(451)

	ErrInternalServerError           error = httpError(500)
	ErrNotImplemented                error = httpError(501)
	ErrBadGateway                    error = httpError(502)
	ErrServiceUnavailable            error = httpError(503)
	ErrGatewayTimeout                error = httpError(504)
	ErrHTTPVersionNotSupported       error = httpError(505)
	ErrVariantAlsoNegotiates         error = httpError(506)
	ErrInsufficientStorage           error = httpError(507)
	ErrLoopDetected                  error = httpError(508)
	ErrNotExtended                   error = httpError(510)
	ErrNetworkAuthenticationRequired error = httpError(511)
)
