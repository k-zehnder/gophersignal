// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/articles": {
            "get": {
                "description": "Retrieve a list of articles from the database. Optional query parameters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Articles"
                ],
                "summary": "Get articles",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "Filter by flagged status",
                        "name": "flagged",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter by dead status",
                        "name": "dead",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Filter by duplicate status",
                        "name": "dupe",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 30,
                        "description": "Limit the number of articles returned (default is 30)",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Offset for pagination (default is 0)",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Articles data",
                        "schema": {
                            "$ref": "#/definitions/models.ArticlesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Article": {
            "type": "object",
            "properties": {
                "comment_count": {
                    "description": "Number of comments from Hacker News or similar.",
                    "type": "integer"
                },
                "comment_link": {
                    "description": "Link to the comment thread (if any).",
                    "type": "string"
                },
                "content": {
                    "description": "Full content of the article.",
                    "type": "string"
                },
                "created_at": {
                    "description": "Timestamp when the article was created.",
                    "type": "string"
                },
                "dead": {
                    "description": "Whether the article is dead.",
                    "type": "boolean"
                },
                "dupe": {
                    "description": "Whether the article is marked as duplicate.",
                    "type": "boolean"
                },
                "flagged": {
                    "description": "Whether the article is flagged.",
                    "type": "boolean"
                },
                "id": {
                    "description": "Unique identifier for the article.",
                    "type": "integer"
                },
                "link": {
                    "description": "URL link to the article.",
                    "type": "string"
                },
                "source": {
                    "description": "Source from where the article was fetched.",
                    "type": "string"
                },
                "summary": {
                    "description": "Summary of the article, nullable.",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the article.",
                    "type": "string"
                },
                "updated_at": {
                    "description": "Timestamp when the article was last updated.",
                    "type": "string"
                },
                "upvotes": {
                    "description": "Upvote count from Hacker News or similar.",
                    "type": "integer"
                }
            }
        },
        "models.ArticlesResponse": {
            "type": "object",
            "properties": {
                "articles": {
                    "description": "List of articles.",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Article"
                    }
                },
                "code": {
                    "description": "HTTP status code.",
                    "type": "integer"
                },
                "status": {
                    "description": "Status message.",
                    "type": "string"
                },
                "total_count": {
                    "description": "Total count of articles.",
                    "type": "integer"
                }
            }
        },
        "models.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "HTTP status code.",
                    "type": "integer"
                },
                "message": {
                    "description": "Detailed error message.",
                    "type": "string"
                },
                "status": {
                    "description": "Error status message.",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "GopherSignal API",
	Description:      "API server for the GopherSignal application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
