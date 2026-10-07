// apidoc emits the actual route inventory as an OpenAPI document.
package main

import (
	"chrisfoong/chaum-work-management-backend/internal/contract"
	"chrisfoong/chaum-work-management-backend/internal/server"
	"chrisfoong/chaum-work-management-backend/internal/work"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	noop := func(c *gin.Context) { c.Next() }
	r := server.New(server.Deps{WebAuth: noop, WorkerAuth: noop})
	contract.RegisterRoutes(r.Web, contract.NewHandler(nil))
	work.Register(r.Web, r.Worker, &work.Service{})
	st := work.NewStorage("", "", "")
	for _, g := range []*gin.RouterGroup{r.Web, r.Worker} {
		g.POST("/files", st.Upload)
		g.GET("/files", st.Read(&work.Service{}))
	}
	paths := map[string]any{}
	for _, route := range r.Engine.Routes() {
		parts := strings.Split(route.Path, "/")
		var params []map[string]any
		for i, part := range parts {
			if strings.HasPrefix(part, ":") {
				name := strings.TrimPrefix(part, ":")
				parts[i] = "{" + name + "}"
				params = append(params, map[string]any{"name": name, "in": "path", "required": true, "schema": map[string]string{"type": "string", "format": "uuid"}})
			}
		}
		path := strings.Join(parts, "/")
		ops, ok := paths[path].(map[string]any)
		if !ok {
			ops = map[string]any{}
			paths[path] = ops
		}
		op := map[string]any{"summary": route.Method + " " + path, "responses": map[string]any{"200": map[string]string{"description": "Successful operation; see docs/API.md for payloads"}, "400": map[string]string{"description": "Invalid input"}, "401": map[string]string{"description": "Unverified, inactive or unknown LINE identity"}, "403": map[string]string{"description": "Role or ownership denied"}, "409": map[string]string{"description": "State conflict"}}}
		if strings.Contains(path, "/contracts/confirm") {
			op["responses"].(map[string]any)["201"] = map[string]string{"description": "Contract and scope created atomically"}
		}
		if len(params) > 0 {
			op["parameters"] = params
		}
		describeRequest(route.Method, path, op)
		if strings.HasPrefix(path, "/api/") {
			op["security"] = []map[string][]string{{"lineIDToken": {}}}
		}
		ops[strings.ToLower(route.Method)] = op
	}
	doc := map[string]any{"openapi": "3.0.3", "info": map[string]string{"title": "Chaum MVP API", "version": "0.1.0", "description": "Implemented routes only. JSON schemas and required payloads are in docs/API.md. Integration against deployed Supabase/LINE requires configuration."}, "paths": paths, "components": map[string]any{"securitySchemes": map[string]any{"lineIDToken": map[string]string{"type": "http", "scheme": "bearer", "bearerFormat": "LINE ID token"}}}}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if e := enc.Encode(doc); e != nil {
		panic(e)
	}
}
