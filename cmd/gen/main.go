package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

type packetField struct {
	Name   string
	Type   string
	IsJson bool
}
type PacketInfo struct {
	CommandArgs string
	PacketType  string
	PacketId    int64
	HasJson     bool
	Fields      []packetField
}

var packetType = flag.String("type", "", "struct name to generate methods for")
var packetId = flag.String("id", "", "packet id")

func main() {
	flag.Parse()
	if *packetType == "" {
		log.Fatalln("--type is required")
	}

	if *packetId == "" {
		log.Fatalln("--id is required")
	}

	id, err := strconv.ParseInt(*packetId, 0, 64)
	if err != nil {
		log.Fatalln("cannot parse packet id:", err)
	}

	directory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file := os.Getenv("GOFILE")

	log.Printf("Processing %v (%v) in %v/%v\n", *packetType, id, directory, file)

	fset := token.NewFileSet()

	conf := &packages.Config{
		Mode:  packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Fset:  fset,
		Tests: false,
	}

	pkgs, err := packages.Load(conf, "github.com/iCiaran/mc-server/internal/packets")
	if err != nil {
		log.Fatal(err)
	}

	scope := pkgs[0].Types.Scope()
	obj := scope.Lookup(*packetType)
	if obj == nil {
		log.Fatalf("no type found for %v", *packetType)
	}

	named, ok := obj.Type().(*types.Named)
	if !ok {
		log.Fatalf("%v is not a named type", *packetType)
	}

	typeStruct, ok := named.Underlying().(*types.Struct)
	if !ok {
		log.Fatalf("%s is not a struct", *packetType)
	}

	packetInfo := PacketInfo{
		strings.Join(os.Args[1:], " "),
		*packetType,
		id,
		false,
		[]packetField{},
	}

	for i := 0; i < typeStruct.NumFields(); i++ {
		field := typeStruct.Field(i)
		tag := typeStruct.Tag(i)

		isJson := strings.Contains(tag, "json")

		packetInfo.Fields = append(packetInfo.Fields, packetField{
			Name:   field.Name(),
			Type:   types.TypeString(field.Type(), func(*types.Package) string { return "" }),
			IsJson: isJson,
		})

		if isJson {
			packetInfo.HasJson = true

		}
	}

	_, filename, _, _ := runtime.Caller(0)
	tmpl := template.Must(template.New("gen.tmpl").ParseFiles(filepath.Dir(filename) + "/gen.tmpl"))

	buf := &bytes.Buffer{}

	err = tmpl.Execute(buf, packetInfo)

	if err != nil {
		log.Fatalf("Failed to execute template for %v: %v", *packetType, err)
	}

	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		log.Fatalf("Failed to format %v: %v", *packetType, err)
	}

	outputFile := fmt.Sprintf("%s_gen.go", *packetType)
	if err := os.WriteFile(outputFile, formatted, 0644); err != nil {
		log.Fatalf("Failed to write %v: %v", outputFile, err)
	}
}
