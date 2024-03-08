// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../../../docs/guide/zh-CN/api/error_code_generated.md

const (
	RuntimeDocker string = "docker"
	RuntimeK8S    string = "k8s"
)

const (
	ServerModelDebug   string = "debug"
	ServerModelTest    string = "test"
	ServerModelRelease string = "release"
)

// Common: basic errors.
// Code must start with 1xxxxx.
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = 100000

	// ErrInternalServerError - 500: Internal server error.
	ErrInternalServerError int = 100001

	ErrUnauthorized int = 100002

	// ErrBind - 400: Error occurred while binding the request body to the struct.
	ErrBind int = 100003

	// 请求参数无效
	ErrInvalidParams int = 100004

	// ErrTokenInvalid - 401: Token invalid.
	ErrExpiredToken int = 100005

	ErrDeny int = 100006

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound int = 100007
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)
