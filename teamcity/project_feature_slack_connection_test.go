package teamcity_test

import (
	"testing"

	"github.com/cvbarros/go-teamcity/teamcity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectFeatureSlackConnection_CreateSlack(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureSlackConnection(project.ID, teamcity.ProjectFeatureSlackConnectionOptions{
		ClientId:     "abcd.1234",
		ClientSecret: "xyz",
		DisplayName:  "Notifier",
		Token:        "ABCD1234EFG",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	assert.Equal(t, "abcd.1234", createdFeature.Properties().Map()["clientId"])
	assert.Equal(t, "Notifier", createdFeature.Properties().Map()["displayName"])
	assert.Equal(t, "slackConnection", createdFeature.Properties().Map()["providerType"])
	assert.Empty(t, createdFeature.Properties().Map()["secure:clientSecret"])
	assert.Empty(t, createdFeature.Properties().Map()["secure:token"])
}

func TestProjectFeatureSlackConnection_Delete(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureSlackConnection(project.ID, teamcity.ProjectFeatureSlackConnectionOptions{
		ClientId:     "abcd.1234",
		ClientSecret: "xyz",
		DisplayName:  "Notifier",
		Token:        "ABCD1234EFG",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	err = service.Delete(createdFeature.ID())
	require.NoError(t, err)

	deletedFeature, err := service.GetByID(createdFeature.ID())
	assert.NotNil(t, err)
	assert.Nil(t, deletedFeature)
}

func TestProjectFeatureSlackConnection_Update(t *testing.T) {
	client := safeSetup(t)

	project := createTestProjectWithImplicitName(t, client)
	defer cleanUpProject(t, client, project.ID)

	service := client.ProjectFeatureService(project.ID)

	feature := teamcity.NewProjectFeatureSlackConnection(project.ID, teamcity.ProjectFeatureSlackConnectionOptions{
		ClientId:     "abcd.1234",
		ClientSecret: "xyz",
		DisplayName:  "Notifier",
		Token:        "ABCD1234EFG",
	})

	createdFeature, err := service.Create(feature)
	require.NoError(t, err)
	assert.NotEmpty(t, createdFeature.ID)

	var validate = func(t *testing.T, id string, expected teamcity.ProjectFeatureSlackConnectionOptions) {
		retrievedFeature, err := service.GetByID(id)
		require.NoError(t, err)
		slackConnection, ok := retrievedFeature.(*teamcity.ProjectFeatureSlackConnection)
		assert.True(t, ok)

		assert.Equal(t, expected.ClientId, slackConnection.Options.ClientId)
		assert.Equal(t, expected.DisplayName, slackConnection.Options.DisplayName)
		assert.Equal(t, "slackConnection", slackConnection.Options.ProviderType)
	}
	t.Log("Validating initial creation")
	validate(t, createdFeature.ID(), teamcity.ProjectFeatureSlackConnectionOptions{
		ClientId:    "abcd.1234",
		DisplayName: "Notifier",
	})

	// then let's toggle some things
	update := teamcity.ProjectFeatureSlackConnectionOptions{
		ClientId:     "1234.abcd",
		ClientSecret: "abc",
		DisplayName:  "Updated",
		Token:        "XYZ789ABCD",
	}
	t.Log("Validating update")
	existing, err := service.GetByID(createdFeature.ID())
	require.NoError(t, err)

	settings, ok := existing.(*teamcity.ProjectFeatureSlackConnection)
	assert.True(t, ok)

	settings.Options.ClientId = update.ClientId
	settings.Options.ClientSecret = update.ClientSecret
	settings.Options.DisplayName = update.DisplayName
	settings.Options.Token = update.Token

	updatedFeature, err := service.Update(settings)
	require.NoError(t, err)
	assert.NotEmpty(t, updatedFeature.ID)

	// sanity check since we're updating with the same ID
	assert.Equal(t, createdFeature.ID(), updatedFeature.ID())

	validate(t, updatedFeature.ID(), update)
}
