package teamcity

// ProjectFeatureSlackConnectionOptions holds all properties for the versioned settings project feature.
type ProjectFeatureSlackConnectionOptions struct {
	ClientId     string
	ClientSecret string `json:"secure:clientSecret"`
	DisplayName  string
	ProviderType string
	Token        string `json:"secure:token"`
}

// ProjectFeatureSlackConnection represents the versioned settings feature for a project.
// Can be used to configure https://confluence.jetbrains.com/display/TCD10/Storing+Project+Settings+in+Version+Control.
type ProjectFeatureSlackConnection struct {
	id        string
	projectID string

	Options ProjectFeatureSlackConnectionOptions
}

// NewProjectFeatureSlackConnection creates a new Versioned Settings project feature.
func NewProjectFeatureSlackConnection(projectID string, options ProjectFeatureSlackConnectionOptions) *ProjectFeatureSlackConnection {
	return &ProjectFeatureSlackConnection{
		projectID: projectID,
		Options:   options,
	}
}

// ID returns the ID of this project feature.
func (f *ProjectFeatureSlackConnection) ID() string {
	return f.id
}

// SetID sets the ID of this project feature.
func (f *ProjectFeatureSlackConnection) SetID(value string) {
	f.id = value
}

// Type represents the type of the project feature as a string.
func (f *ProjectFeatureSlackConnection) Type() string {
	return "OAuthProvider"
}

// ProjectID represents the ID of the project the project feature is assigned to.
func (f *ProjectFeatureSlackConnection) ProjectID() string {
	return f.projectID
}

// SetProjectID sets the ID of the project the project feature is assigned to.
func (f *ProjectFeatureSlackConnection) SetProjectID(value string) {
	f.projectID = value
}

// Properties returns all properties for the versioned settings project feature.
func (f *ProjectFeatureSlackConnection) Properties() *Properties {
	props := NewProperties(
		NewProperty("clientId", string(f.Options.ClientId)),
		NewProperty("secure:clientSecret", f.Options.ClientSecret),
		NewProperty("displayName", string(f.Options.DisplayName)),
		NewProperty("providerType", "slackConnection"),
		NewProperty("secure:token", f.Options.Token),
	)

	return props
}

func loadProjectFeatureSlackConnection(projectID string, feature projectFeatureJSON) (ProjectFeature, error) {
	settings := &ProjectFeatureSlackConnection{
		id:        feature.ID,
		projectID: projectID,
		Options:   ProjectFeatureSlackConnectionOptions{},
	}

	if clientId, ok := feature.Properties.GetOk("clientId"); ok {
		settings.Options.ClientId = clientId
	}

	if displayName, ok := feature.Properties.GetOk("displayName"); ok {
		settings.Options.DisplayName = displayName
	}

	if providerType, ok := feature.Properties.GetOk("providerType"); ok {
		settings.Options.ProviderType = providerType
	}

	return settings, nil
}
