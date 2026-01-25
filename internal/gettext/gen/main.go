package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Message struct {
	ID       string
	IDPlural string
	File     string
	Line     int
	IsPlural bool
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: go run main.go [output_file]")
		fmt.Println("Scans Go files for gettext.Get and gettext.GetN calls and generates a POT file")
		fmt.Println("Default output: messages.pot")
		os.Exit(0)
	}

	outputFile := "messages.pot"
	if len(os.Args) > 1 {
		outputFile = os.Args[1]
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding project root: %v\n", err)
		os.Exit(1)
	}

	messages, err := scanGoFiles(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning Go files: %v\n", err)
		os.Exit(1)
	}

	if len(messages) == 0 {
		fmt.Println("No translatable messages found")
		return
	}

	err = generatePOTFile(outputFile, messages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating POT file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s with %d messages\n", outputFile, len(messages))
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("go.mod not found in current directory or any parent directory")
}

func scanGoFiles(root string) ([]Message, error) {
	var messages []Message
	messageMap := make(map[string]Message)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || name == "node_modules" {
				return filepath.SkipDir
			}
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileMessages, err := extractMessages(path, root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: error processing %s: %v\n", path, err)
			return nil
		}

		for _, msg := range fileMessages {
			key := msg.ID
			if msg.IsPlural {
				key = msg.ID + "|" + msg.IDPlural
			}

			if _, exists := messageMap[key]; !exists {
				messageMap[key] = msg
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, msg := range messageMap {
		messages = append(messages, msg)
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].ID < messages[j].ID
	})

	return messages, nil
}

func extractMessages(filePath, projectRoot string) ([]Message, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var messages []Message
	relPath, _ := filepath.Rel(projectRoot, filePath)

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if ident.Name != "gettext" {
			return true
		}

		pos := fset.Position(call.Pos())
		line := pos.Line

		switch sel.Sel.Name {
		case "Get":
			if len(call.Args) < 1 {
				return true
			}

			if msgid := extractStringLiteral(call.Args[0]); msgid != "" {
				messages = append(messages, Message{
					ID:       msgid,
					File:     relPath,
					Line:     line,
					IsPlural: false,
				})
			}

		case "GetN":
			if len(call.Args) < 3 {
				return true
			}

			msgid := extractStringLiteral(call.Args[0])
			msgidPlural := extractStringLiteral(call.Args[1])

			if msgid != "" && msgidPlural != "" {
				messages = append(messages, Message{
					ID:       msgid,
					IDPlural: msgidPlural,
					File:     relPath,
					Line:     line,
					IsPlural: true,
				})
			}
		}

		return true
	})

	return messages, nil
}

func extractStringLiteral(expr ast.Expr) string {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		str := lit.Value
		if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
			str = str[1 : len(str)-1]
			str = strings.ReplaceAll(str, `\"`, `"`)
			str = strings.ReplaceAll(str, `\\`, `\`)
			str = strings.ReplaceAll(str, `\n`, "\n")
			str = strings.ReplaceAll(str, `\t`, "\t")
			return str
		}
	}
	return ""
}

func generatePOTFile(filename string, messages []Message) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	version := "PACKAGE VERSION"
	cmd := exec.Command("git", "describe", "--tags", "--long", "--abbrev=7")
	if output, err := cmd.Output(); err == nil {
		version = strings.TrimSpace(string(output))
	}

	fmt.Fprintln(file, "# SOME DESCRIPTIVE TITLE.")
	fmt.Fprintln(file, "# Copyright (C) YEAR THE PACKAGE'S COPYRIGHT HOLDER")
	fmt.Fprintln(file, "# This file is distributed under the same license as the PACKAGE package.")
	fmt.Fprintln(file, "# FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.")
	fmt.Fprintln(file, "#")
	fmt.Fprintln(file, "#, fuzzy")
	fmt.Fprintln(file, "msgid \"\"")
	fmt.Fprintln(file, "msgstr \"\"")
	fmt.Fprintln(file, "\"Project-Id-Version: "+version+"\\n\"")
	fmt.Fprintln(file, "\"Report-Msgid-Bugs-To: \\n\"")
	fmt.Fprintf(file, "\"POT-Creation-Date: %s\\n\"\n", time.Now().Format("2006-01-02 15:04-0700"))
	fmt.Fprintln(file, "\"PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE\\n\"")
	fmt.Fprintln(file, "\"Last-Translator: FULL NAME <EMAIL@ADDRESS>\\n\"")
	fmt.Fprintln(file, "\"Language-Team: LANGUAGE <LL@li.org>\\n\"")
	fmt.Fprintln(file, "\"Language: \\n\"")
	fmt.Fprintln(file, "\"MIME-Version: 1.0\\n\"")
	fmt.Fprintln(file, "\"Content-Type: text/plain; charset=UTF-8\\n\"")
	fmt.Fprintln(file, "\"Content-Transfer-Encoding: 8bit\\n\"")
	fmt.Fprintln(file)

	// Write messages
	for _, msg := range messages {
		fmt.Fprintf(file, "#: %s:%d\n", msg.File, msg.Line)

		if msg.IsPlural {
			fmt.Fprintf(file, "msgid %s\n", quotePOString(msg.ID))
			fmt.Fprintf(file, "msgid_plural %s\n", quotePOString(msg.IDPlural))
			fmt.Fprintln(file, "msgstr[0] \"\"")
			fmt.Fprintln(file, "msgstr[1] \"\"")
		} else {
			fmt.Fprintf(file, "msgid %s\n", quotePOString(msg.ID))
			fmt.Fprintln(file, "msgstr \"\"")
		}
		fmt.Fprintln(file)
	}

	return nil
}

func quotePOString(s string) string {
	if strings.Contains(s, "\n") {
		lines := strings.Split(s, "\n")
		if len(lines) == 1 {
			return fmt.Sprintf("\"%s\"", escapeForPO(s))
		}

		result := "\"\"\n"
		for i, line := range lines {
			if i == len(lines)-1 && line == "" {
				break
			}
			if i == len(lines)-1 {
				result += fmt.Sprintf("\"%s\"", escapeForPO(line))
			} else {
				result += fmt.Sprintf("\"%s\\n\"\n", escapeForPO(line))
			}
		}
		return result
	}

	return fmt.Sprintf("\"%s\"", escapeForPO(s))
}

func escapeForPO(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\t", "\\t")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
