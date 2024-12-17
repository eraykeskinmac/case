package docs

var SwaggerInfo = struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
}{
	Title:       "Invoice Management API",
	Description: "API for managing invoices with CRUD operations",
	Version:     "1.0",
	Host:        "localhost:3000",
	BasePath:    "/api",
	Schemes:     []string{"http", "https"},
}
