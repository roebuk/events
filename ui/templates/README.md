# Templ Templates Guide

## Overview

This directory contains [templ](https://templ.guide) templates for the Firecrest Go application. Templ is a templating language for Go that provides type-safe HTML generation.

## Working with Templ

### File Structure

- Template files use the `.templ` extension
- Each file must declare a package (e.g., `package templates`)
- Components are defined using the `templ` keyword

### Basic Template Syntax

```templ
templ ComponentName(param1 string, param2 int) {
    <div>
        <h1>{ param1 }</h1>
        <p>Count: { fmt.Sprintf("%d", param2) }</p>
    </div>
}
```

### Using Children Content

Use `{ children... }` to accept nested content:

```templ
templ Layout(title string) {
    <html>
        <head><title>{ title }</title></head>
        <body>
            { children... }
        </body>
    </html>
}

// Usage
templ Page() {
    @Layout("My Page") {
        <h1>Content goes here</h1>
    }
}
```

### Conditional Rendering

```templ
templ ShowMessage(show bool, message string) {
    if show {
        <p>{ message }</p>
    }
}
```

### Loops

```templ
templ List(items []string) {
    <ul>
        for _, item := range items {
            <li>{ item }</li>
        }
    </ul>
}
```

### Component Composition

```templ
templ Page(meta templ.Component) {
    <head>
        if meta != nil {
            @meta
        }
    </head>
}
```

## Generating Templates

After modifying `.templ` files, you must regenerate the Go code.

### Prerequisites

Install templ CLI:

```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### Generate Command

From the project root:

```bash
templ generate
```

Or from this directory:

```bash
templ generate ./ui/templates
```

### Watch Mode (Development)

Automatically regenerate on file changes:

```bash
templ generate --watch
```

### Build Integration

Add to your build process:

```bash
templ generate && go build ./...
```

## Generated Files

- Templ generates `*_templ.go` files alongside your `.templ` files
- These files are auto-generated - **do not edit them manually**
- Add `*_templ.go` to `.gitignore` or commit them (project preference)

## Best Practices

1. **Keep components small and focused** - One responsibility per component
2. **Use parameters for dynamic content** - Avoid global state
3. **Leverage type safety** - Use Go types for parameters
4. **Component composition** - Build complex UIs from simple components
5. **Naming conventions** - Use PascalCase for component names
6. **Generate before committing** - Ensure generated code is up to date

## IDE Support

- **VS Code**: Install the [templ extension](https://marketplace.visualstudio.com/items?itemName=a-h.templ)
- **GoLand/IntelliJ**: Plugin available in marketplace
- **Neovim**: LSP support via templ-lsp

## Common Patterns

### Layout with Meta Tags

```templ
templ Layout(title string, meta templ.Component) {
    <html>
        <head>
            <title>{ title }</title>
            if meta != nil {
                @meta
            }
        </head>
        <body>{ children... }</body>
    </html>
}

templ MetaTags(description, keywords string) {
    <meta name="description" content={ description }/>
    <meta name="keywords" content={ keywords }/>
}
```

### Using in HTTP Handlers

```go
func handler(w http.ResponseWriter, r *http.Request) {
    component := templates.Page("Hello World")
    component.Render(r.Context(), w)
}
```

## Resources

- [Official Templ Documentation](https://templ.guide)
- [Templ GitHub Repository](https://github.com/a-h/templ)
- [Syntax Reference](https://templ.guide/syntax-and-usage/template-syntax)
