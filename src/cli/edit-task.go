package main

import (
	"errors"
	"fmt"
	"sort"

	"github.com/charmbracelet/huh"
)

func BuildTaskEditForm(task Task, allTags Tags, allProjects []ProjectRow) *huh.Form {
	priority := fmt.Sprint(int(task.priority))
	tagOptions := []huh.Option[string]{}

	tagTexts := []string{}
	for _, text := range allTags {
		tagTexts = append(tagTexts, text)
	}

	sort.Strings(tagTexts)

	for _, text := range tagTexts {
		for id, tagText := range allTags {
			if text == tagText {
				isSelected := false
				for tagId := range task.tags {
					if tagId == id {
						isSelected = true
						break
					}
				}
				tagOptions = append(tagOptions, huh.NewOption(tagText, fmt.Sprint(id)).Selected(isSelected))
			}
		}
	}

	projectOptions := []huh.Option[string]{}
	for _, project := range allProjects {
		projectOptions = append(projectOptions, huh.NewOption(project.title, fmt.Sprint(project.id)))
	}

	projectOptions = append(projectOptions, huh.NewOption("No project", "0"))

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Description").
				Key("description").
				Value(&task.description).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Description cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Priority").
				Key("priority").
				Value(&priority).
				Options(
					huh.NewOption("Low", "0"),
					huh.NewOption("Medium", "1"),
					huh.NewOption("High", "2"),
					huh.NewOption("Critical", "3"),
				),
			huh.NewSelect[string]().
				Title("Project").
				Key("project").
				Value(&task.projectId).
				Options(projectOptions...),
			huh.NewMultiSelect[string]().
				Title("Tags").
				Key("tags").
				Options(tagOptions...).
				Limit(3),
		).WithTheme(&*huh.ThemeCatppuccin()),
	)
}
