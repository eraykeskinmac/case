package docs

var ModelDefinitions = map[string]map[string]any{
	"Invoice": {
		"type": "object",
		"properties": map[string]any{
			"id": map[string]any{
				"type":    "integer",
				"example": 1,
			},
			"service_name": map[string]any{
				"type":    "string",
				"example": "DMP Service",
			},
			"invoice_number": map[string]any{
				"type":    "integer",
				"example": 1001,
			},
			"date": map[string]any{
				"type":    "string",
				"format":  "date-time",
				"example": "2024-03-16T00:00:00Z",
			},
			"amount": map[string]any{
				"type":    "number",
				"format":  "float",
				"example": 1500.50,
			},
			"status": map[string]any{
				"type":    "string",
				"enum":    []string{"Paid", "Pending", "Unpaid"},
				"example": "Pending",
			},
			"created_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
			"updated_at": map[string]any{
				"type":   "string",
				"format": "date-time",
			},
		},
		"required": []string{"service_name", "invoice_number", "date", "amount", "status"},
	},
	"InvoiceResponse": {
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":    "string",
				"example": "Operation successful",
			},
			"data": map[string]any{
				"$ref": "#/definitions/Invoice",
			},
		},
	},
	"InvoiceListResponse": {
		"type": "object",
		"properties": map[string]any{
			"data": map[string]any{
				"type": "array",
				"items": map[string]any{
					"$ref": "#/definitions/Invoice",
				},
			},
			"meta": map[string]any{
				"$ref": "#/definitions/MetaData",
			},
		},
	},
	"MetaData": {
		"type": "object",
		"properties": map[string]any{
			"total": map[string]any{
				"type":    "integer",
				"example": 100,
			},
			"page": map[string]any{
				"type":    "integer",
				"example": 1,
			},
			"limit": map[string]any{
				"type":    "integer",
				"example": 10,
			},
			"total_pages": map[string]any{
				"type":    "integer",
				"example": 10,
			},
			"sort_by": map[string]any{
				"type":    "string",
				"example": "created_at",
			},
			"sort_dir": map[string]any{
				"type":    "string",
				"example": "desc",
			},
		},
	},
	"ErrorResponse": {
		"type": "object",
		"properties": map[string]any{
			"message": map[string]any{
				"type":    "string",
				"example": "Error occurred",
			},
			"error": map[string]any{
				"type":    "string",
				"example": "Invalid input",
			},
		},
	},
}
