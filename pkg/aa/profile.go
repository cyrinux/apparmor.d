// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Warning: this is purposely not using a Yacc parser. Its only aim is to
// extract variables and attachments for apparmor.d profile

package aa

import (
	"regexp"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	regVariablesDef = regexp.MustCompile(`@{(.*)}\s*[+=]+\s*(.*)`)
	regVariablesRef = regexp.MustCompile(`@{([^{}]+)}`)

	// Tunables
	Tunables = map[string][]string{
		"bin":             {"/{usr/,}{s,}bin"},
		"lib":             {"/{usr/,}lib{,exec,32,64}"},
		"multiarch":       {"*-linux-gnu*"},
		"user_share_dirs": {"/home/*/.local/share"},
		"etc_ro":          {"/{usr/,}etc/"},
	}
)

type AppArmorProfile struct {
	Variables   map[string][]string
	Attachments []string
}

func NewAppArmorProfile() *AppArmorProfile {
	variables := make(map[string][]string)
	maps.Copy(variables, Tunables)
	return &AppArmorProfile{
		Variables:   variables,
		Attachments: []string{},
	}
}

// ParseVariables extract all variables from the profile
func (p *AppArmorProfile) ParseVariables(content string) {
	matches := regVariablesDef.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			key := match[1]
			values := match[2]
			if _, ok := p.Variables[key]; ok {
				p.Variables[key] = append(p.Variables[key], strings.Split(values, " ")...)
			} else {
				p.Variables[key] = strings.Split(values, " ")
			}
		}
	}
}

// resolve recursively resolves all variables references
func (p *AppArmorProfile) resolve(str string) []string {
	if strings.Contains(str, "@{") {
		vars := []string{}
		match := regVariablesRef.FindStringSubmatch(str)
		if len(match) > 1 {
			variable := match[0]
			varname := match[1]
			for _, value := range p.Variables[varname] {
				newVar := strings.ReplaceAll(str, variable, value)
				vars = append(vars, p.resolve(newVar)...)
			}
		} else {
			vars = append(vars, str)
		}
		return vars
	}
	return []string{str}
}

// ResolveAttachments resolve profile attachments defined in exec_path
func (p *AppArmorProfile) ResolveAttachments() {
	for _, exec := range p.Variables["exec_path"] {
		p.Attachments = append(p.Attachments, p.resolve(exec)...)
	}
}

// NestAttachments return a nested attachment string
func (p *AppArmorProfile) NestAttachments() string {
	if len(p.Attachments) == 0 {
		return ""
	} else if len(p.Attachments) == 1 {
		return p.Attachments[0]
	} else {
		res := []string{}
		for _, attachment := range p.Attachments {
			if strings.HasPrefix(attachment, "/") {
				res = append(res, attachment[1:])
			} else {
				res = append(res, attachment)
			}
		}
		return "/{" + strings.Join(res, ",") + "}"
	}
}
