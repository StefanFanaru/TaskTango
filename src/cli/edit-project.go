package main

import (
	"errors"

	"github.com/charmbracelet/huh"
)

func BuildProjectEditForm(project ProjectRow) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Title").
				Key("title").
				Value(&project.title).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Title cannot be empty")
					}
					return nil
				}),
		).WithTheme(&*huh.ThemeCatppuccin()),
	)
}
