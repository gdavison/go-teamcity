package teamcity

import (
	"encoding/json"
)

// FeatureGolang represents a golang build feature. Implements BuildFeature interface
type FeatureGolang struct {
	id          string
	disabled    bool
	buildTypeID string

	properties *Properties
}

// NewFeatureGolang returns a new instance of the FeatureGolang struct
func NewFeatureGolang() *FeatureGolang {
	return &FeatureGolang{
		properties: NewProperties(),
	}
}

// ID returns the ID for this instance.
func (f *FeatureGolang) ID() string {
	return f.id
}

// SetID sets the ID for this instance.
func (f *FeatureGolang) SetID(value string) {
	f.id = value
}

// Type returns the "golang", the keyed-type for this build feature instance
func (f *FeatureGolang) Type() string {
	return "golang"
}

// Disabled returns whether this build feature is disabled or not.
func (f *FeatureGolang) Disabled() bool {
	return f.disabled
}

// SetDisabled sets whether this build feature is disabled or not.
func (f *FeatureGolang) SetDisabled(value bool) {
	f.disabled = value
}

// BuildTypeID is a getter for the Build Type ID associated with this build feature.
func (f *FeatureGolang) BuildTypeID() string {
	return f.buildTypeID
}

// SetBuildTypeID is a setter for the Build Type ID associated with this build feature.
func (f *FeatureGolang) SetBuildTypeID(value string) {
	f.buildTypeID = value
}

// Properties returns a *Properties instance representing a serializable collection to be used.
func (f *FeatureGolang) Properties() *Properties {
	return f.properties
}

// MarshalJSON implements JSON serialization for FeatureGolang
func (f *FeatureGolang) MarshalJSON() ([]byte, error) {
	out := &buildFeatureJSON{
		ID:         f.id,
		Disabled:   NewBool(f.disabled),
		Properties: f.properties,
		Inherited:  NewFalse(),
		Type:       f.Type(),
	}

	// this is the only value and has to be set to this - no no point making it user configurable
	out.Properties.AddOrReplaceValue("test.format", "json")
	return json.Marshal(out)
}

// UnmarshalJSON implements JSON deserialization for FeatureGolang
func (f *FeatureGolang) UnmarshalJSON(data []byte) error {
	var aux buildFeatureJSON
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	f.id = aux.ID

	disabled := aux.Disabled
	if disabled == nil {
		disabled = NewFalse()
	}
	f.disabled = *disabled
	f.properties = NewProperties(aux.Properties.Items...)

	return nil
}
