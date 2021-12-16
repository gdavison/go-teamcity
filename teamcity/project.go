package teamcity

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// Project is the model for project entities in TeamCity
type Project struct {
	Archived        *bool               `json:"archived,omitempty" xml:"archived"`
	Description     string              `json:"description,omitempty" xml:"description"`
	Href            string              `json:"href,omitempty" xml:"href"`
	ID              string              `json:"id,omitempty" xml:"id"`
	Name            string              `json:"name,omitempty" xml:"name"`
	Parameters      *Parameters         `json:"parameters,omitempty"`
	ParentProject   *ProjectReference   `json:"parentProject,omitempty"`
	ParentProjectID string              `json:"parentProjectId,omitempty" xml:"parentProjectId"`
	WebURL          string              `json:"webUrl,omitempty" xml:"webUrl"`
	BuildTypes      BuildTypeReferences `json:"buildTypes,omitempty" xml:"buildTypes"`
}

// ProjectReference contains basic information, usually enough to use as a type for relationships.
// In addition to that, TeamCity does not return the full detailed representation when creating objects, thus the need for a reference.
type ProjectReference struct {
	ID          string `json:"id,omitempty" xml:"id"`
	Name        string `json:"name,omitempty" xml:"name"`
	Description string `json:"description,omitempty" xml:"description"`
	Href        string `json:"href,omitempty" xml:"href"`
	WebURL      string `json:"webUrl,omitempty" xml:"webUrl"`
}

// ProjectService has operations for handling projects
type ProjectService struct {
	sling      *sling.Sling
	httpClient *http.Client
	restHelper *restHelper
}

//NewProject returns an instance of a Project. A non-empty name is required.
//Description can be an empty string and will be omitted.
//For creating a top-level project, pass empty to parentProjectId.
func NewProject(name string, description string, parentProjectID string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	var parent *ProjectReference
	if parentProjectID != "" {
		parent = &ProjectReference{
			ID: parentProjectID,
		}
	}
	return &Project{
		Name:            name,
		Description:     description,
		ParentProject:   parent,
		ParentProjectID: parentProjectID,
		Parameters:      NewParametersEmpty(),
	}, nil
}

//SetParentProject changes this Project instance's parent project
func (p *Project) SetParentProject(parentID string) {
	p.ParentProjectID = parentID
	p.ParentProject = &ProjectReference{
		ID: parentID,
	}
}

//ProjectReference converts a project instance to a ProjectReference
func (p *Project) ProjectReference() *ProjectReference {
	return &ProjectReference{
		ID:          p.ID,
		Description: p.Description,
		Name:        p.Name,
		WebURL:      p.WebURL,
		Href:        p.Href,
	}
}

func newProjectService(base *sling.Sling, client *http.Client) *ProjectService {
	sling := base.Path("projects/")
	return &ProjectService{
		sling:      sling,
		httpClient: client,
		restHelper: newRestHelper(client, sling),
	}
}

// Create creates a new project at root project level
func (s *ProjectService) Create(project *Project) (*Project, error) {
	var created ProjectReference
	err := s.restHelper.post("", project, &created, "project")
	if err != nil {
		return nil, err
	}

	//initial creation does not persist "description" or parameters, so in order to be consistent with the constructor, call an update after
	project.ID = created.ID
	updated, err := s.updateProject(project, true)

	if err != nil {
		return nil, err
	}

	return updated, nil
}

// GetByID Retrieves a project resource by ID
func (s *ProjectService) GetByID(id string) (*Project, error) {
	var out Project
	locator := LocatorID(id).String()
	err := s.restHelper.get(locator, &out, "project")
	if err != nil {
		return nil, err
	}

	//For now, filter all inherited parameters, until figuring out a proper way of exposing filtering options to the caller
	out.Parameters = out.Parameters.NonInherited()
	return &out, err
}

//GetByName returns a project by its name. There are no duplicate names in projects for TeamCity
func (s *ProjectService) GetByName(name string) (*Project, error) {
	var out Project

	err := s.restHelper.get(LocatorName(name).String(), &out, "project")
	if err != nil {
		return nil, err
	}

	//For now, filter all inherited parameters, until figuring out a proper way of exposing filtering options to the caller
	out.Parameters = out.Parameters.NonInherited()
	return &out, err
}

//Update changes the resource in-place for this project.
//TeamCity API does not support "PUT" on the whole project resource, so the only updateable field is "Description". Other field updates will be ignored.
//This method also updates Settings and Parameters, but this is not an atomic operation. If an error occurs, it will be returned to caller what was updated or not.
func (s *ProjectService) Update(project *Project) (*Project, error) {
	return s.updateProject(project, false)
}

//Delete - Deletes a project
func (s *ProjectService) Delete(id string) error {
	err := s.restHelper.deleteByIDWithSling(s.sling.New(), id, "project")
	return err
}

func (s *ProjectService) updateProject(project *Project, isCreate bool) (*Project, error) {
	_, err := s.restHelper.putTextPlain(project.ID+"/name", project.Name, "project name")
	if err != nil {
		return nil, err
	}

	_, err = s.restHelper.putTextPlain(project.ID+"/description", project.Description, "project description")
	if err != nil {
		return nil, err
	}

	//Update Parent
	if !isCreate {
		current, err := s.GetByID(project.ID)
		if err != nil {
			return nil, err
		}
		// Only perform update if there is a change.
		// Or else TeamCity will "copy" the project to the same parent project, altering it's name
		// For instance: "project" -> "project (1)"
		if (project.ParentProjectID != "" || project.ParentProject != nil) && current.ParentProjectID != project.ParentProjectID {
			var parent ProjectReference
			err = s.restHelper.put(project.ID+"/parentProject", project.ParentProject, &parent, "parent project")
			if err != nil {
				return nil, nil
			}
		}
	}

	//Update Parameters
	if project.Parameters.Count > 0 {
		var parameters *Parameters
		err = s.restHelper.put(project.ID+"/parameters", project.Parameters, &parameters, "project parameters")
		if err != nil {
			return nil, err
		}
	}
	out, err := s.GetByID(project.ID) //Refresh after update
	if err != nil {
		return nil, err
	}

	return out, nil
}
