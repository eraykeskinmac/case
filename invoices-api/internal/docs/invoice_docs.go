package docs

var InvoiceEndpoints = map[string]EndpointDoc{
	"GetInvoices": {
		Summary:     "List all invoices",
		Description: "Get all invoices with pagination, sorting and search functionality",
		Tags:        []string{"invoices"},
		Method:      "GET",
		Path:        "/v1/invoices",
		Parameters: []Parameter{
			{
				Name:        "page",
				In:          "query",
				Type:        "integer",
				Required:    false,
				Default:     "1",
				Description: "Page number",
			},
			{
				Name:        "limit",
				In:          "query",
				Type:        "integer",
				Required:    false,
				Default:     "10",
				Description: "Items per page",
			},
			{
				Name:        "search",
				In:          "query",
				Type:        "string",
				Required:    false,
				Description: "Search by service name",
			},
			{
				Name:        "sort_by",
				In:          "query",
				Type:        "string",
				Required:    false,
				Description: "Field to sort by",
			},
			{
				Name:        "sort_dir",
				In:          "query",
				Type:        "string",
				Required:    false,
				Default:     "asc",
				Description: "Sort direction (asc/desc)",
			},
		},
		Responses: map[int]Response{
			200: {
				Description: "Successful response",
				Schema:      "InvoiceListResponse",
			},
			500: {
				Description: "Internal server error",
				Schema:      "ErrorResponse",
			},
		},
	},
	"GetInvoiceByID": {
		Summary:     "Get invoice by ID",
		Description: "Get invoice details by its ID",
		Tags:        []string{"invoices"},
		Method:      "GET",
		Path:        "/v1/invoices/{id}",
		Parameters: []Parameter{
			{
				Name:        "id",
				In:          "path",
				Type:        "integer",
				Required:    true,
				Description: "Invoice ID",
			},
		},
		Responses: map[int]Response{
			200: {
				Description: "Successful response",
				Schema:      "InvoiceResponse",
			},
			404: {
				Description: "Invoice not found",
				Schema:      "ErrorResponse",
			},
		},
	},
	"CreateInvoice": {
		Summary:     "Create new invoice",
		Description: "Create a new invoice with the provided details",
		Tags:        []string{"invoices"},
		Method:      "POST",
		Path:        "/v1/invoices",
		Parameters: []Parameter{
			{
				Name:        "body",
				In:          "body",
				Type:        "object",
				Required:    true,
				Schema:      "Invoice",
				Description: "Invoice details",
			},
		},
		Responses: map[int]Response{
			201: {
				Description: "Invoice created successfully",
				Schema:      "InvoiceResponse",
			},
			400: {
				Description: "Invalid input",
				Schema:      "ErrorResponse",
			},
			500: {
				Description: "Internal server error",
				Schema:      "ErrorResponse",
			},
		},
	},

	"UpdateInvoice": {
		Summary:     "Update invoice",
		Description: "Update an existing invoice by ID",
		Tags:        []string{"invoices"},
		Method:      "PUT",
		Path:        "/v1/invoices/{id}",
		Parameters: []Parameter{
			{
				Name:        "id",
				In:          "path",
				Type:        "integer",
				Required:    true,
				Description: "Invoice ID",
			},
			{
				Name:        "body",
				In:          "body",
				Type:        "object",
				Required:    true,
				Schema:      "Invoice",
				Description: "Updated invoice details",
			},
		},
		Responses: map[int]Response{
			200: {
				Description: "Invoice updated successfully",
				Schema:      "InvoiceResponse",
			},
			400: {
				Description: "Invalid input",
				Schema:      "ErrorResponse",
			},
			404: {
				Description: "Invoice not found",
				Schema:      "ErrorResponse",
			},
			500: {
				Description: "Internal server error",
				Schema:      "ErrorResponse",
			},
		},
	},

	"DeleteInvoice": {
		Summary:     "Delete invoice",
		Description: "Delete an invoice by ID",
		Tags:        []string{"invoices"},
		Method:      "DELETE",
		Path:        "/v1/invoices/{id}",
		Parameters: []Parameter{
			{
				Name:        "id",
				In:          "path",
				Type:        "integer",
				Required:    true,
				Description: "Invoice ID",
			},
		},
		Responses: map[int]Response{
			200: {
				Description: "Invoice deleted successfully",
				Schema:      "InvoiceResponse",
			},
			404: {
				Description: "Invoice not found",
				Schema:      "ErrorResponse",
			},
			500: {
				Description: "Internal server error",
				Schema:      "ErrorResponse",
			},
		},
	},
}
