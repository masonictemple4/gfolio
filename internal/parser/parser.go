package parser

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// For now this will default to yaml frontmatter
type Parser struct {
	input  *bufio.Reader
	output *bytes.Buffer

	Format *Format
}

func New(rdr io.Reader, fmt *Format) *Parser {
	p := &Parser{
		input:  bufio.NewReader(rdr),
		output: bytes.NewBuffer(nil),
		Format: fmt,
	}

	return p
}

func (p *Parser) Parse(result interface{}) error {

	shouldWrite := false

	for {
		line, err := p.input.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err != nil && err == io.EOF {
			format := NewYamlFrontMatterFormat()
			err = format.DecodeFn(p.output.Bytes(), &result)
			if err != nil {
				return err
			}
			return nil
		}

		line = bytes.TrimSuffix(line, []byte("\n"))

		if !shouldWrite && string(line) == p.Format.Begin {
			shouldWrite = true
			// Skip so we don't write the begin indicator to the output..
			continue
		}

		if shouldWrite && string(line) != p.Format.End {
			line = append(line, []byte("\n")...)
			_, err := p.output.Write(line)
			if err != nil {
				return err
			}
		}

		if string(line) == p.Format.End {
			break
		}
	}

	// Most of the time we will hit this case, chances are we won't
	// be publishing blogs with just the frontmatter and no content.
	if err := p.Format.DecodeFn(p.output.Bytes(), &result); err != nil {
		return err
	}

	return nil
}

func ParseFile[K any](path string, result K) error {
	fmFormat := NewYamlFrontMatterFormat()
	out := bytes.NewBuffer(nil)

	shouldWrite := false

	f, err := os.Open(path)
	if err != nil {
		println("could not open file")
		return err
	}

	defer f.Close()

	rdr := bufio.NewReader(f)

	for {
		line, err := rdr.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err != nil && err == io.EOF {
			err = fmFormat.DecodeFn(out.Bytes(), &result)
			if err != nil {
				println("eof cannot decode")
				return err
			}
			return nil
		}

		line = bytes.TrimSuffix(line, []byte("\n"))

		if !shouldWrite && string(line) == fmFormat.Begin {
			shouldWrite = true
			continue
		}

		if shouldWrite && string(line) != fmFormat.End {

			line = append(line, []byte("\n")...)
			_, err := out.Write(line)
			if err != nil {
				return err
			}
		}

		if string(line) == fmFormat.End {
			break
		}

	}

	// need to decode here as well because the file will most likely
	// have data beyond the frontmatter.

	if err := fmFormat.DecodeFn(out.Bytes(), &result); err != nil {
		println("could not decode with the frontmatter")
		return err
	}

	return nil
}

func SkipFrontmatter(data []byte) ([]byte, error) {
	fmFormat := NewYamlFrontMatterFormat()
	out := bytes.NewBuffer(nil)

	skipping := false

	rdr := bufio.NewReader(bytes.NewReader(data))

	for {
		line, err := rdr.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err != nil && err == io.EOF {
			return out.Bytes(), nil
		}

		line = bytes.TrimSuffix(line, []byte("\n"))

		if !skipping && string(line) == fmFormat.Begin {
			skipping = true
			continue
		}

		if skipping && string(line) != fmFormat.End {
			continue
		}

		if skipping && string(line) == fmFormat.End {
			skipping = false
			continue
		}

		line = append(line, []byte("\n")...)
		_, err = out.Write(line)
		if err != nil {
			return nil, err
		}

	}

	// need to decode here as well because the file will most likely
	// have data beyond the frontmatter.

	// return out.Bytes(), nil
}
