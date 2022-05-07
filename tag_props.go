package strmap

import (
	"strings"
)

type Meta map[string]string

type TagProps struct {
	Name    *string
	Default *[]string
	Ignore  bool
	all     Meta
}

func (t TagProps) All() Meta {
	m := make(map[string]string, len(t.all))
	for k, v := range t.all {
		m[k] = v
	}

	return m
}

type ParseTagProp = func(*TagProps, string) error

func parseDefaultTagProp(t *TagProps, v string) error {
	s := trimAll(strings.Split(v, ";"))
	t.Default = &s

	return nil
}

var fieldsParser = map[string]ParseTagProp{
	"default": parseDefaultTagProp,
}

func trimAll(ss []string) []string {
	r := make([]string, len(ss))
	for i, s := range ss {
		r[i] = strings.TrimSpace(s)
	}

	return r
}

func getTagKeyValue(s string) (string, string) {
	kv := trimAll(strings.SplitN(s, "=", 2))
	if len(kv) == 2 {
		return kv[0], kv[1]
	}

	return kv[0], ""
}

func parseTagProp(tp *TagProps, kv string) error {
	k, v := getTagKeyValue(kv)
	tp.all[k] = v
	fp, ok := fieldsParser[k]
	if !ok {
		return nil
	}

	return fp(tp, v)
}

func parseTagProps(sp string) (TagProps, error) {
	fp := TagProps{
		Name:    nil,
		Default: nil,
		Ignore:  false,
		all:     make(map[string]string),
	}

	ps := strings.Split(sp, ",")
	name := &ps[0]

	switch *name {
	case "-":
		fp.Ignore = true
		return fp, nil
	case "":
		name = nil
	default:
		fp.Name = name

	}

	for _, p := range ps[1:] {
		err := parseTagProp(&fp, p)
		if err != nil {
			return fp, err
		}
	}

	return fp, nil
}
