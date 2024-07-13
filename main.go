package main

import "bitbucket.org/dyfrag-internal/mass-media-core/pkg/cli"

// @title SE Project
// @version 1.0
// @description Software Engineering Course Project
// @contact.name Mahdieh Moghiseh
// @contact.email mahdiehmoghiseh81@gmail.com
// @securityDefinitions.apiKey BearerAuth
// @type apiKey
// @name Authorization
// @in header
// @externalDocs.description OpenAPI
func main() {
	cli.Execute()
}
