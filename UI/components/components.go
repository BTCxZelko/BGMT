package components

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// CreateLabeledTextInput creates and returns a new container with a label and a text input widget
func CreateLabeledTextInput(placeholder string, onChange func(string)) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder(placeholder)
	input.OnChanged = onChange

	// Combine the label and text input in a vertical box layout
	return container.NewVBox(
		input,
	)
}

// CreateLabeledTextInputInt creates a labeled text input that accepts integers and passes a string representation of integers within the specified range to the onChange function.
func CreateLabeledTextInputInt(window fyne.Window, placeholder string, onChange func(s string)) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder(placeholder)

	input.OnChanged = func(inputValue string) {
		if inputValue == "" {
			// If the input is empty, do nothing (or reset to a default state if desired)
			dialog.ShowError(errors.New("must be a number between -255 and 255"), window)
		}

		// Try to convert the inputValue to an integer
		num, err := strconv.Atoi(inputValue)
		if err != nil || num < -255 || num > 255 {
			// Input is not a valid number or out of bounds; show an error dialog
			dialog.ShowError(errors.New("must be a number between -255 and 255"), window)
		} else {
			// Input is valid; convert the number back to a string and call the onChange function
			onChange(strconv.Itoa(num))
		}
	}

	// Combine the label and text input in a vertical box layout
	return container.NewVBox(
		input,
	)
}

func ShowFileLoadDialog(window fyne.Window, expectedExtension string, onValidFileSelected func(filePath string)) {
	dialog.ShowFileOpen(func(file fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, window)
			return
		}
		if file == nil {
			// User cancelled the dialog
			return
		}
		// Check if the file has the appropriate extension
		if filepath.Ext(file.URI().Path()) == expectedExtension {
			// Execute the callback with the selected file path
			onValidFileSelected(file.URI().Path())
		} else {
			dialog.ShowInformation("Invalid File", "Please select a file with a "+expectedExtension+" extension", window)
		}
	}, window)
}

// CreateLoadModButton creates and returns a button for loading and copying a mod file.
func CreateLoadModButton(window fyne.Window, expectedExtension string, postProcess func()) *widget.Button {
	btn := widget.NewButton("Load Mod File", func() {
		ShowFileLoadDialog(window, expectedExtension, func(filePath string) {
			// File is valid, now copy to working.tmp
			err := CopyFile(filePath, "working.tmp")
			if err != nil {
				dialog.ShowError(err, window)
			} else {
				dialog.ShowInformation("Success", "File loaded and copied to working.tmp", window)
				if postProcess != nil {
					postProcess() // Execute any post-processing function if provided
				}
			}
		})
	})

	return btn
}

// CopyFile copies the contents of the src file to the dst file.
func CopyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, input, 0644)
}

func MakeNewFile(templatePath, projectDir string, window fyne.Window) {
	dst := filepath.Join(projectDir, "working.tmp")
	copyErr := CopyFile(templatePath, dst)
	if copyErr != nil {
		dialog.ShowError(copyErr, window)
		return
	}
	dialog.ShowInformation("Success", "New project file created", window)
}

// InsertUserInputs replaces placeholders in a file with user inputs.
func InsertUserInputs(filePath string, inputs *UserInputs) error {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Convert the content to a string for replacement
	newContent := string(content)

	// Replace each placeholder with the corresponding value from UserInputs
	newContent = strings.ReplaceAll(newContent, "$variable", inputs.Variable)
	newContent = strings.ReplaceAll(newContent, "$type", inputs.Type)
	newContent = strings.ReplaceAll(newContent, "$value", inputs.Value)
	newContent = strings.ReplaceAll(newContent, "$var", inputs.Variable2)
	newContent = strings.ReplaceAll(newContent, "$typ", inputs.Type2)
	newContent = strings.ReplaceAll(newContent, "$val", inputs.Value2)
	newContent = strings.ReplaceAll(newContent, "$creatureID", inputs.CreatureID)
	newContent = strings.ReplaceAll(newContent, "$dialogueID", inputs.DialogueID)

	// Write the updated content back to the file
	return ioutil.WriteFile(filePath, []byte(newContent), 0644)
}

type UserInputs struct {
	Variable   string
	Type       string
	Value      string
	Variable2  string
	Type2      string
	Value2     string
	DialogueID string
	CreatureID string
}
