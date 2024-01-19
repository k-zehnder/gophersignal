package docs

import "github.com/swaggo/swag"

const docTemplate = `
openapi: 3.0.0
info:
  description: "An Article service API in Go using Gin framework"
  title: "Article Service API"
  version: "1.0"
servers:
  - url: http://localhost:8080/api/v1
tags:
  - name: "Article"
    description: "Operations related to articles"
paths:
  /articles:
    get:
      tags:
        - "Article"
      summary: "List Articles"
      description: "Retrieves a list of all articles"
      responses:
        '200':
          description: "An array of articles"
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/response.ArticleResponse"
    post:
      tags:
        - "Article"
      summary: "Create Article"
      description: "Creates a new article"
      requestBody:
        description: "Article data"
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/request.CreateArticleRequest"
      responses:
        '201':
          description: "Article created successfully"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/response.Response"
        '400':
          description: "Invalid request format"
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/response.Response"
components:
  schemas:
    request.CreateArticleRequest:
      type: object
      properties:
        content:
          type: string
        title:
          type: string
        link:
          type: string
        summary:
          type: string
        source:
          type: string
    response.ArticleResponse:
      type: object
      properties:
        id:
          type: integer
        title:
          type: string
        content:
          type: string
        link:
          type: string
        summary:
          type: string
        source:
          type: string
        createdAt:
          type: string
        updatedAt:
          type: string
        is_on_homepage:
          type: boolean
    response.Response:
      type: object
      properties:
        code:
          type: integer
        data:
          type: object
        status:
          type: string
`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Article Service API",
	Description:      "An Article service API in Go using Gin framework",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
