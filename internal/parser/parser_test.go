package parser

import (
	"os"
	"testing"

	"github.com/masonictemple4/masonictempl/internal/dtos"
)

func TestSkipFrontmatter(t *testing.T) {
	fp := ""

	data, err := os.ReadFile(fp)
	if err != nil {
		t.Errorf("there was an error opening the file: %v", err)
	}

	newData, err := SkipFrontmatter(data)
	if err != nil {
		t.Errorf("there was an error skipping the frontmatter: %v", err)
	}

	t.Logf("New data: %s\n", string(newData))
}

func TestStandaloneParser(t *testing.T) {

	// TODO: Create some test blogs to go against and replace the second string
	// in our table below.
	var tests = []struct {
		name string
		fp   string
	}{
		{"Test parse with just frontmatter in the file", ""},
		{"Test with excess content after frontmatter", ""},
		// TODO: This test case could pass in the case that there is no additonal page content
		// after the frontmatter content and the parser reaches the end of file.
		// We'll probably want to check for both the begin and end precence otherwise just call
		// it invalid
		{"Test with valid open but no present close", ""},
		{"Test with no formatter open, but with a valid close.. So technically is it an open???", ""},
		{"Test with no fromatter precense at all.", ""},
		{"Test plain blog post inside of the frontmatter.", ""},
	}

	// TODO: Come back and refactor this.
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result dtos.BlogInput
			err := ParseFile(tt.fp, &result)
			if err != nil {
				t.Errorf("there was an error parsing the file: %v", err)
			}
			t.Logf("Test %d: %+v\n", i, result)
		})
	}

}

func TestParserObject(t *testing.T) {
	t.Run("Test just frontmatter in file with parser object", func(t *testing.T) {
		var result dtos.BlogInput
		fp := ""
		f, err := os.Open(fp)
		if err != nil {
			t.Errorf("there was an error opening the file: %v", err)
		}
		defer f.Close()
		ymlFormat := NewYamlFrontMatterFormat()
		psr := New(f, ymlFormat)
		err = psr.Parse(&result)
	})

	t.Run("Test with the parser object", func(t *testing.T) {
		var result dtos.BlogInput
		fp := ""
		f, err := os.Open(fp)
		if err != nil {
			t.Errorf("there was an error opening the file: %v", err)
		}
		defer f.Close()
		ymlFormat := NewYamlFrontMatterFormat()
		psr := New(f, ymlFormat)
		err = psr.Parse(&result)
	})

}
