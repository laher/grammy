package main

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

// Ensure the parser can parse strings into Statement ASTs.
func TestParser_ParseStatement(t *testing.T) {
	var tests = []struct {
		s     string
		stmts []*Statement
		err   string
	}{
		// Single statement
		{
			s: `ACTOR me`,
			stmts: []*Statement{
				{
					Type: DefinitionStatement,
					Definition: &Definition{
						Name:  "me",
						Type:  ACTOR,
						Alias: "",
					},
				},
			},
		},
		{
			s: `ACTOR me
			PARTICIPANT world as w
			me -> w : hello
			`,
			stmts: []*Statement{
				{
					Type: DefinitionStatement,
					Definition: &Definition{
						Name:  "me",
						Type:  ACTOR,
						Alias: "",
					},
				},
				{
					Type: DefinitionStatement,
					Definition: &Definition{
						Name:  "world",
						Type:  PARTICIPANT,
						Alias: "w",
					},
				},
				{
					Type: MessageStatement,
					Message: &Message{
						From:        "me",
						To:          "w",
						Arrow:       RARR,
						Description: "hello",
					},
				},
			},
		},
		/*
			// Errors
			{s: `foo`, err: `found "foo", expected ACTOR`},
			{s: `ACTOR !`, err: `found "!", expected field`},
			{s: `ACTOR field xxx`, err: `found "xxx", expected PARTICIPANT`},
			{s: `ACTOR field PARTICIPANT *`, err: `found "*", expected table name`},
		*/
	}

	for i, tt := range tests {
		p := NewParser(strings.NewReader(tt.s))
		ss := []*Statement{}
		j := 0
		for {
			stmt, err := p.Parse()
			if err == io.EOF {
				if len(tt.stmts) != len(ss) {
					t.Errorf("Wrong number of lines (%d != %d)", len(tt.stmts), len(ss))
				}
				break
			}
			if len(tt.stmts) < i {
				t.Errorf("Wrong number of lines (too many)")
			}
			if err != nil {
				t.Fatalf("Failed with error: %v", err)
			}
			if !reflect.DeepEqual(tt.err, errstring(err)) {
				t.Errorf("%d. %q: error mismatch:\n  exp=%s\n  got=%s\n\n", i, tt.s, tt.err, err)
			} else if tt.err == "" {
				exp := tt.stmts[j]
				if exp.Type != stmt.Type {
					t.Errorf("%d.%d. %q\n\ndefinition type mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, j, tt.s, exp.Type, stmt.Type)
				}
				if !reflect.DeepEqual(exp.Definition, stmt.Definition) {
					t.Errorf("%d.%d. %q\n\ndefinition statement mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, j, tt.s, exp.Definition, stmt.Definition)
				}
				if !reflect.DeepEqual(exp.Message, stmt.Message) {
					t.Errorf("%d.%d. %q\n\nmessage statement mismatch:\n\nexp=%#v\n\ngot=%#v\n\n", i, j, tt.s, exp.Message, stmt.Message)
				}

			}
			ss = append(ss, stmt)
			j++
		}
	}
}

// errstring returns the string representation of an error.
func errstring(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
