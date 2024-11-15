// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"time"
)

// BadRequest Error
type BadRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// InternalServerError Internal server error
type InternalServerError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LinkItem defines model for LinkItem.
type LinkItem struct {
	AccessCount int        `json:"access_count"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	ExpireAt    time.Time  `json:"expire_at"`
	Id          string     `json:"id"`
	LastAccess  *time.Time `json:"last_access,omitempty"`
	ShortLink   string     `json:"short_link"`
	TargetUrl   string     `json:"target_url"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// LinkListResponse response
type LinkListResponse = []LinkItem

// NotFound not found
type NotFound struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Ok Success
type Ok struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RedirectResponse defines model for RedirectResponse.
type RedirectResponse = string

// ShortenerPostRequest request body
type ShortenerPostRequest struct {
	ExpireDays int    `json:"expire_days"`
	Url        string `json:"url"`
}

// ShortenerPostResponse response
type ShortenerPostResponse struct {
	ShortLink string `json:"short_link"`
}

// PostShortenerJSONRequestBody defines body for PostShortener for application/json ContentType.
type PostShortenerJSONRequestBody = ShortenerPostRequest
