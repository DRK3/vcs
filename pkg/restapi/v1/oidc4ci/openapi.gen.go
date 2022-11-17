// Package oidc4ci provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package oidc4ci

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// Model for Access Token Response.
type AccessTokenResponse struct {
	// The access token issued by the authorization server.
	AccessToken string `json:"access_token"`

	// String containing a nonce to be used to create a proof of possession of key material when requesting a credential.
	CNonce *string `json:"c_nonce,omitempty"`

	// Integer denoting the lifetime in seconds of the c_nonce.
	CNonceExpiresIn *int `json:"c_nonce_expires_in,omitempty"`

	// The lifetime in seconds of the access token.
	ExpiresIn *int `json:"expires_in,omitempty"`

	// The refresh token, which can be used to obtain new access tokens.
	RefreshToken *string `json:"refresh_token,omitempty"`

	// OPTIONAL, if identical to the scope requested by the client; otherwise, REQUIRED.
	Scope *string `json:"scope,omitempty"`

	// The type of the token issued.
	TokenType string `json:"token_type"`
}

// Model for OIDC Credential request.
type CredentialRequest struct {
	// DID to which issued credential has to be bound.
	Did string `json:"did"`

	// Format of the credential being issued.
	Format *string   `json:"format,omitempty"`
	Proof  *JWTProof `json:"proof,omitempty"`

	// Type of the credential being issued.
	Type string `json:"type"`
}

// Model for OIDC Credential response.
type CredentialResponse struct {
	// A JSON string containing a token subsequently used to obtain a Credential. MUST be present when credential is not returned.
	AcceptanceToken *string `json:"acceptance_token,omitempty"`

	// JSON string containing a nonce to be used to create a proof of possession of key material when requesting a Credential.
	CNonce *string `json:"c_nonce,omitempty"`

	// JSON integer denoting the lifetime in seconds of the c_nonce.
	CNonceExpiresIn *int        `json:"c_nonce_expires_in,omitempty"`
	Credential      interface{} `json:"credential"`

	// JSON string denoting the format of the issued Credential.
	Format string `json:"format"`
}

// JWTProof defines model for JWTProof.
type JWTProof struct {
	// REQUIRED. Signed JWT as proof of key possession.
	Jwt string `json:"jwt"`

	// REQUIRED. JSON String denoting the proof type. Currently the only supported proof type is 'jwt'.
	ProofType string `json:"proof_type"`
}

// Model for Pushed Authorization Response.
type PushedAuthorizationResponse struct {
	// A JSON number that represents the lifetime of the request URI in seconds as a positive integer. The request URI lifetime is at the discretion of the authorization server but will typically be relatively short (e.g., between 5 and 600 seconds).
	ExpiresIn int `json:"expires_in"`

	// The request URI corresponding to the authorization request posted. This URI is a single-use reference to the respective request data in the subsequent authorization request.
	RequestUri string `json:"request_uri"`
}

// OidcAuthorizeParams defines parameters for OidcAuthorize.
type OidcAuthorizeParams struct {
	// Value MUST be set to "code".
	ResponseType string `form:"response_type" json:"response_type"`

	// The client identifier.
	ClientId string `form:"client_id" json:"client_id"`

	// A challenge derived from the code verifier that is sent in the authorization request, to be verified against later.
	CodeChallenge string `form:"code_challenge" json:"code_challenge"`

	// A method that was used to derive code challenge.
	CodeChallengeMethod *string `form:"code_challenge_method,omitempty" json:"code_challenge_method,omitempty"`

	// The authorization server redirects the user-agent to the client's redirection endpoint previously established with the authorization server during the client registration process or when making the authorization request.
	RedirectUri *string `form:"redirect_uri,omitempty" json:"redirect_uri,omitempty"`

	// The scope of the access request.
	Scope *string `form:"scope,omitempty" json:"scope,omitempty"`

	// An opaque value used by the client to maintain state between the request and callback. The authorization server includes this value when redirecting the user-agent back to the client. The parameter SHOULD be used for preventing cross-site request forgery.
	State *string `form:"state,omitempty" json:"state,omitempty"`

	// The authorization_details conveys the details about the credentials the wallet wants to obtain. Multiple authorization_details can be used with type openid_credential to request authorization in case of multiple credentials.
	AuthorizationDetails *string `form:"authorization_details,omitempty" json:"authorization_details,omitempty"`

	// Wallet's OpenID Connect Issuer URL. The Issuer will use the discovery process to determine the wallet's capabilities and endpoints. RECOMMENDED in Dynamic Credential Request.
	WalletIssuer *string `form:"wallet_issuer,omitempty" json:"wallet_issuer,omitempty"`

	// An opaque user hint the wallet MAY use in subsequent callbacks to optimize the user's experience. RECOMMENDED in Dynamic Credential Request.
	UserHint *string `form:"user_hint,omitempty" json:"user_hint,omitempty"`

	// String value identifying a certain processing context at the credential issuer. A value for this parameter is typically passed in an issuance initiation request from the issuer to the wallet. This request parameter is used to pass the  op_state value back to the credential issuer. The issuer must take into account that op_state is not guaranteed to originate from this issuer, could be an attack.
	OpState string `form:"op_state" json:"op_state"`
}

// OidcCredentialJSONBody defines parameters for OidcCredential.
type OidcCredentialJSONBody = CredentialRequest

// OidcRedirectParams defines parameters for OidcRedirect.
type OidcRedirectParams struct {
	// auth code for issuer provider
	Code string `form:"code" json:"code"`

	// state
	State string `form:"state" json:"state"`
}

// OidcCredentialJSONRequestBody defines body for OidcCredential for application/json ContentType.
type OidcCredentialJSONRequestBody = OidcCredentialJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// OIDC Authorization Request
	// (GET /oidc/authorize)
	OidcAuthorize(ctx echo.Context, params OidcAuthorizeParams) error
	// OIDC Credential
	// (POST /oidc/credential)
	OidcCredential(ctx echo.Context) error
	// OIDC Pushed Authorization Request
	// (POST /oidc/par)
	OidcPushedAuthorizationRequest(ctx echo.Context) error
	// OIDC Pre-Authorized code flow handler
	// (POST /oidc/pre-authorized-code)
	OidcPreAuthorizedCode(ctx echo.Context) error
	// OIDC Redirect
	// (GET /oidc/redirect)
	OidcRedirect(ctx echo.Context, params OidcRedirectParams) error
	// OIDC Token Request
	// (POST /oidc/token)
	OidcToken(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// OidcAuthorize converts echo context to params.
func (w *ServerInterfaceWrapper) OidcAuthorize(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params OidcAuthorizeParams
	// ------------- Required query parameter "response_type" -------------

	err = runtime.BindQueryParameter("form", true, true, "response_type", ctx.QueryParams(), &params.ResponseType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter response_type: %s", err))
	}

	// ------------- Required query parameter "client_id" -------------

	err = runtime.BindQueryParameter("form", true, true, "client_id", ctx.QueryParams(), &params.ClientId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter client_id: %s", err))
	}

	// ------------- Required query parameter "code_challenge" -------------

	err = runtime.BindQueryParameter("form", true, true, "code_challenge", ctx.QueryParams(), &params.CodeChallenge)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code_challenge: %s", err))
	}

	// ------------- Optional query parameter "code_challenge_method" -------------

	err = runtime.BindQueryParameter("form", true, false, "code_challenge_method", ctx.QueryParams(), &params.CodeChallengeMethod)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code_challenge_method: %s", err))
	}

	// ------------- Optional query parameter "redirect_uri" -------------

	err = runtime.BindQueryParameter("form", true, false, "redirect_uri", ctx.QueryParams(), &params.RedirectUri)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter redirect_uri: %s", err))
	}

	// ------------- Optional query parameter "scope" -------------

	err = runtime.BindQueryParameter("form", true, false, "scope", ctx.QueryParams(), &params.Scope)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scope: %s", err))
	}

	// ------------- Optional query parameter "state" -------------

	err = runtime.BindQueryParameter("form", true, false, "state", ctx.QueryParams(), &params.State)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter state: %s", err))
	}

	// ------------- Optional query parameter "authorization_details" -------------

	err = runtime.BindQueryParameter("form", true, false, "authorization_details", ctx.QueryParams(), &params.AuthorizationDetails)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter authorization_details: %s", err))
	}

	// ------------- Optional query parameter "wallet_issuer" -------------

	err = runtime.BindQueryParameter("form", true, false, "wallet_issuer", ctx.QueryParams(), &params.WalletIssuer)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter wallet_issuer: %s", err))
	}

	// ------------- Optional query parameter "user_hint" -------------

	err = runtime.BindQueryParameter("form", true, false, "user_hint", ctx.QueryParams(), &params.UserHint)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter user_hint: %s", err))
	}

	// ------------- Required query parameter "op_state" -------------

	err = runtime.BindQueryParameter("form", true, true, "op_state", ctx.QueryParams(), &params.OpState)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter op_state: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcAuthorize(ctx, params)
	return err
}

// OidcCredential converts echo context to params.
func (w *ServerInterfaceWrapper) OidcCredential(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcCredential(ctx)
	return err
}

// OidcPushedAuthorizationRequest converts echo context to params.
func (w *ServerInterfaceWrapper) OidcPushedAuthorizationRequest(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcPushedAuthorizationRequest(ctx)
	return err
}

// OidcPreAuthorizedCode converts echo context to params.
func (w *ServerInterfaceWrapper) OidcPreAuthorizedCode(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcPreAuthorizedCode(ctx)
	return err
}

// OidcRedirect converts echo context to params.
func (w *ServerInterfaceWrapper) OidcRedirect(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params OidcRedirectParams
	// ------------- Required query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, true, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code: %s", err))
	}

	// ------------- Required query parameter "state" -------------

	err = runtime.BindQueryParameter("form", true, true, "state", ctx.QueryParams(), &params.State)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter state: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcRedirect(ctx, params)
	return err
}

// OidcToken converts echo context to params.
func (w *ServerInterfaceWrapper) OidcToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.OidcToken(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/oidc/authorize", wrapper.OidcAuthorize)
	router.POST(baseURL+"/oidc/credential", wrapper.OidcCredential)
	router.POST(baseURL+"/oidc/par", wrapper.OidcPushedAuthorizationRequest)
	router.POST(baseURL+"/oidc/pre-authorized-code", wrapper.OidcPreAuthorizedCode)
	router.GET(baseURL+"/oidc/redirect", wrapper.OidcRedirect)
	router.POST(baseURL+"/oidc/token", wrapper.OidcToken)

}
