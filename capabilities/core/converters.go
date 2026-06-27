package core

import (
	skillcore "github.com/plexusone/assistantkit/skills/core"
	validationcore "github.com/plexusone/assistantkit/validation/core"
)

// ToSkill converts a Capability to a Skill for human-guided use.
// Skills are invoked via slash commands or triggers and provide
// interactive guidance to the user.
func (c *Capability) ToSkill() *skillcore.Skill {
	skill := skillcore.NewSkill(c.Name, c.Description)
	skill.Instructions = c.Instructions
	skill.Triggers = c.Triggers
	skill.Dependencies = c.Dependencies
	return skill
}

// ToCheck converts a Capability to a ValidationArea Check for
// autonomous QA validation. Checks are executed by validation
// agents and report GO/NO-GO status.
func (c *Capability) ToCheck() validationcore.Check {
	command := c.Command
	if len(c.CommandArgs) > 0 {
		for _, arg := range c.CommandArgs {
			command += " " + arg
		}
	}

	return validationcore.Check{
		Name:        c.Name,
		Description: c.Description,
		Command:     command,
		Pattern:     c.FailurePattern,
		Required:    c.Required,
	}
}

// CapabilitySet represents a collection of related capabilities.
// This allows grouping capabilities by domain (e.g., all Go checks).
type CapabilitySet struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Language     string        `json:"language,omitempty"`
	Capabilities []*Capability `json:"capabilities"`
}

// NewCapabilitySet creates a new CapabilitySet.
func NewCapabilitySet(name, description string) *CapabilitySet {
	return &CapabilitySet{
		Name:        name,
		Description: description,
	}
}

// Add adds a capability to the set.
func (cs *CapabilitySet) Add(c *Capability) *CapabilitySet {
	cs.Capabilities = append(cs.Capabilities, c)
	return cs
}

// ToSkills converts all capabilities to Skills.
func (cs *CapabilitySet) ToSkills() []*skillcore.Skill {
	skills := make([]*skillcore.Skill, len(cs.Capabilities))
	for i, c := range cs.Capabilities {
		skills[i] = c.ToSkill()
	}
	return skills
}

// ToChecks converts all capabilities to ValidationArea Checks.
func (cs *CapabilitySet) ToChecks() []validationcore.Check {
	checks := make([]validationcore.Check, len(cs.Capabilities))
	for i, c := range cs.Capabilities {
		checks[i] = c.ToCheck()
	}
	return checks
}

// ToValidationArea creates a ValidationArea from the capability set.
// This is useful for creating a QA validation area that composes
// multiple capabilities.
func (cs *CapabilitySet) ToValidationArea(signOffCriteria, instructions string) *validationcore.ValidationArea {
	va := validationcore.NewValidationArea(cs.Name, cs.Description)
	va.SignOffCriteria = signOffCriteria
	va.Instructions = instructions

	// Add all checks
	for _, c := range cs.Capabilities {
		va.AddCheck(c.ToCheck())
	}

	// Collect unique dependencies
	deps := make(map[string]bool)
	for _, c := range cs.Capabilities {
		for _, dep := range c.Dependencies {
			deps[dep] = true
		}
	}
	for dep := range deps {
		va.AddDependency(dep)
	}

	return va
}
