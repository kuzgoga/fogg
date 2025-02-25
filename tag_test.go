package main

import "testing"

func TestTagMethods(t *testing.T) {
	tag := Tag{
		name:    "example",
		value:   "value",
		params:  map[string]TagParam{"param1": {Value: "value1"}},
		options: []string{"option1", "option2"},
	}

	if tag.Name() != "example" {
		t.Errorf("expected name to be 'example', got %s", tag.Name())
	}

	if !tag.HasOption("option1") {
		t.Errorf("expected to have option 'option1'")
	}

	if tag.HasOption("option3") {
		t.Errorf("expected not to have option 'option3'")
	}

	if !tag.HasParam("param1") {
		t.Errorf("expected to have param 'param1'")
	}

	if tag.HasParam("param2") {
		t.Errorf("expected not to have param 'param2'")
	}

	if param := tag.GetParam("param1"); param == nil || param.Value != "value1" {
		t.Errorf("expected to get param 'param1' with value 'value1'")
	}

	if param := tag.GetParam("param2"); param != nil {
		t.Errorf("expected to get nil for param 'param2'")
	}

	if value := tag.GetParamOr("param1", "default"); value != "value1" {
		t.Errorf("expected to get 'value1' for param 'param1'")
	}

	if value := tag.GetParamOr("param2", "default"); value != "default" {
		t.Errorf("expected to get 'default' for param 'param2'")
	}

	if tag.GetValue() != "value" {
		t.Errorf("expected value to be 'value', got %s", tag.GetValue())
	}
}
