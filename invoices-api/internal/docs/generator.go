package docs

import (
	"fmt"
	"strings"
)

type EndpointDoc struct {
	Summary     string
	Description string
	Tags        []string
	Method      string
	Path        string
	Parameters  []Parameter
	Responses   map[int]Response
}

type Parameter struct {
	Name        string
	In          string
	Type        string
	Required    bool
	Default     string
	Schema      string
	Description string
}

type Response struct {
	Description string
	Schema      string
}

func GenerateSwaggerSpec() map[string]any {
	paths := make(map[string]any)

	// Her endpoint için path ve method map'i oluştur
	for _, endpoint := range InvoiceEndpoints {
		// Eğer path daha önce oluşturulmadıysa
		if _, exists := paths[endpoint.Path]; !exists {
			paths[endpoint.Path] = make(map[string]any)
		}

		// HTTP metodu için detayları ekle (lowercase olmalı)
		method := strings.ToLower(endpoint.Method)
		pathMap := paths[endpoint.Path].(map[string]any)
		pathMap[method] = map[string]any{
			"tags":        endpoint.Tags,
			"summary":     endpoint.Summary,
			"description": endpoint.Description,
			"parameters":  generateParametersSpec(endpoint.Parameters),
			"responses":   generateResponsesSpec(endpoint.Responses),
			"produces":    []string{"application/json"},
			"consumes":    []string{"application/json"},
		}
	}

	return map[string]any{
		"swagger": "2.0",
		"info": map[string]any{
			"title":       SwaggerInfo.Title,
			"description": SwaggerInfo.Description,
			"version":     SwaggerInfo.Version,
		},
		"host":        SwaggerInfo.Host,
		"basePath":    SwaggerInfo.BasePath,
		"schemes":     SwaggerInfo.Schemes,
		"paths":       paths,
		"definitions": ModelDefinitions,
	}
}

func generateParametersSpec(params []Parameter) []map[string]any {
	result := make([]map[string]any, 0)
	for _, param := range params {
		paramSpec := map[string]any{
			"name":        param.Name,
			"in":          param.In,
			"required":    param.Required,
			"description": param.Description,
		}

		if param.In == "body" {
			paramSpec["schema"] = map[string]any{
				"$ref": fmt.Sprintf("#/definitions/%s", param.Schema),
			}
		} else {
			paramSpec["type"] = param.Type
			if param.Default != "" {
				paramSpec["default"] = param.Default
			}
		}

		result = append(result, paramSpec)
	}
	return result
}

func generateResponsesSpec(responses map[int]Response) map[string]any {
	result := make(map[string]any)
	for code, response := range responses {
		result[fmt.Sprint(code)] = map[string]any{
			"description": response.Description,
			"schema": map[string]any{
				"$ref": fmt.Sprintf("#/definitions/%s", response.Schema),
			},
		}
	}
	return result
}
