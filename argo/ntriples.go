package argo

import (
    "fmt"
    "github.com/kierdavis/goutil"
    "io"
    "os"
    "strings"
)

var (
    Ignore = iota
    Print
    Error
    Aggregate
)

var Whitespace = " \t\r\n"

type NTriplesParseError struct {
    Message string
    Line int
}

func (err *NTriplesParseError) Error() (msg string) {
    return fmt.Sprintf("Error when parsing line %d: %s", err.Line, err.Message)
}

type NTriplesParseErrors []*NTriplesParseError

func (err NTriplesParseErrors) Error() (msg string) {
    return "Multiple errors occurred during parsing"
}

func WriteNTriples(graph *Graph, w io.Writer) (err error) {
    for triple := range graph.IterTriples() {
        _, err = fmt.Fprintln(w, triple.String())
        if err != nil {
            return err
        }
    }
}

func ReadNTriples(graph *Graph, r io.Reader, parseErrors int) (err error) {
    lineChan, errChan := util.IterLines(r)
    lineno := 0
    
    var errors NTriplesParseErrors
    
    for line := range lineChan {
        lineno++
        err = nil
        
        line = strings.Trim(line, Whitespace)
        if line == "" || line[0] == '#' {
            continue
        }
        
        for {
            line, subj, err := readTerm(line, lineno)
            if err != nil {
                break
            }
            
            line, pred, err := readTerm(line, lineno)
            if err != nil {
                break
            }
            
            line, obj, err := readTerm(line, lineno)
            if err != nil {
                break
            }
            
            ctx := nil
            
            if line[0] != '.' {
                line, ctx, err = readTerm(line, lineno)
                if err != nil {
                    break
                }
                
                if line[0] != '.' {
                    err = NTriplesParseError{"Unterminated line (lines must end with a period - '.')", lineno}
                    break
                }
            }
            
            graph.AddQuad(subj, pred, obj, ctx)
            break
        }
        
        if err != nil {
            switch parseErrors {
            case Ignore:
                /* Do nothing except drop the triple */
            
            case Print:
                fmt.Fprintf(os.Stderr, "NTriples %s\n", err.Error())
            
            case Error:
                return err
            
            case Aggregate:
                errors = append(errors, err)
            }
        }
    }
    
    err = <-errChan
    if err != nil {
        return err
    }
    
    if errors != nil {
        return errors
    }
    
    return nil
}

func readTerm(line string, lineno int) (remainder string, term Term, err error) {
    line = strings.TrimLeft(line, Whitespace)
    
    // Resource
    if line[0] == '<' {
        end := strings.Index(line, ">")
        if end < 0 {
            return "", nil, NTriplesParseError{"Unterminated resource URI (<...>); '>' character is required", lineno}
        }
        
        return line[end+1:], NewResource(line[1:end]), nil
    
    // Blank node
    } else if line[:2] == "_:" {
        end := strings.IndexAny(line, Whitespace)
        if end < 0 {
            return "", nil, NTriplesParseError{"Unterminated blank node (_:...); delimiting whitespace is required", lineno}
        }
        
        return line[end+1:], NewBlank(line[2:end]), nil
    
    // Literal
    } else if line[0] == '"' {
        end := strings.Index(line, '"')
        if end < 0 {
            return "", nil, NTriplesParseError{"Unterminated literal (\"...\"); '\"' character is required", lineno}
        }
        
        text := line[1:end]
        line = line[end+1:]
        
        if line[0] == '@' {
            end = strings.IndexAny(line, Whitespace)
            if end < 0 {
                return "", nil, NTriplesParseError{"Unterminated language identifier (@...); delimiting whitespace is required", lineno}
            }
            
            return line[end+1:], NewLiteralWithLanguage(text, line[1:end]), nil
        
        } else if strings.HasPrefix(line, "^^") {
            remainder, datatype, err := readTerm(line[2:])
            if err != nil {
                return "", nil, err
            }
            
            return remainder, NewLiteralWithDatatype(text, datatype), nil
        }
        
        return line, NewLiteral(text), nil
    }
    
    return "", nil, NTriplesParseError{"Invalid term syntax", lineno}
}
