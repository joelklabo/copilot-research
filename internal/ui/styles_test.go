package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultStyles_Defined(t *testing.T) {
	styles := DefaultStyles()
	
	// All styles should be non-nil
	assert.NotNil(t, styles.TitleStyle)
	assert.NotNil(t, styles.SpinnerStyle)
	assert.NotNil(t, styles.MessageStyle)
	assert.NotNil(t, styles.ResultStyle)
	assert.NotNil(t, styles.ErrorStyle)
	assert.NotNil(t, styles.SuccessStyle)
}

func TestStyles_Render(t *testing.T) {
	styles := DefaultStyles()
	
	// Test that styles can render text
	title := styles.TitleStyle.Render("Title")
	assert.NotEmpty(t, title)
	assert.Contains(t, title, "Title")
	
	message := styles.MessageStyle.Render("Message")
	assert.NotEmpty(t, message)
	assert.Contains(t, message, "Message")
	
	result := styles.ResultStyle.Render("Result")
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "Result")
	
	errMsg := styles.ErrorStyle.Render("Error")
	assert.NotEmpty(t, errMsg)
	assert.Contains(t, errMsg, "Error")
	
	success := styles.SuccessStyle.Render("Success")
	assert.NotEmpty(t, success)
	assert.Contains(t, success, "Success")
}

func TestStyles_ConsistentApplication(t *testing.T) {
	styles := DefaultStyles()
	
	// Same text should produce same output
	text := "Test"
	output1 := styles.TitleStyle.Render(text)
	output2 := styles.TitleStyle.Render(text)
	assert.Equal(t, output1, output2)
}
