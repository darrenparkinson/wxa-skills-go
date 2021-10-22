// Copyright 2021 Darren Parkinson

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wxas

import (
	"encoding/json"
)

// DirectiveName provides a strongly typed enum for directive names
type DirectiveName struct {
	slug string
}

// String implements the stringer interface
func (n DirectiveName) String() string {
	return n.slug
}

// MarshalJSON implements the Marshaler interface
func (n DirectiveName) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.slug)
}

// UnmarshalJSON implements the Unmarshaler interface
func (n *DirectiveName) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	n.slug = s
	return nil
}

// DirectiveType provides a strongly typed enum for directive types
type DirectiveType struct {
	slug string
}

// String implements the stringer interface
func (t DirectiveType) String() string {
	return t.slug
}

// MarshalJSON implements the Marshaler interface
func (t DirectiveType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.slug)
}

// UnmarshalJSON implements the Unmarshaler interface
func (t *DirectiveType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t.slug = s
	return nil
}

var (
	// DirectiveNameReply provides the reply directive name
	DirectiveNameReply = DirectiveName{"reply"}
	// DirectiveNameSpeak provides the speak directive name
	DirectiveNameSpeak = DirectiveName{"speak"}
	// DirectiveNameListen provides the listen directive name
	DirectiveNameListen = DirectiveName{"listen"}
	// DirectiveNameSleep provides the sleep directive name
	DirectiveNameSleep = DirectiveName{"sleep"}
	// DirectiveNameUIHint provides the ui-hint directive name
	DirectiveNameUIHint = DirectiveName{"ui-hint"}
	// DirectiveNameDisplayWebView provides the display-web-view directive name
	DirectiveNameDisplayWebView = DirectiveName{"display-web-view"}
	// DirectiveNameClearWebView provides the clear-web-view directive name
	DirectiveNameClearWebView = DirectiveName{"clear-web-view"}
	// DirectiveNameAssistantEvent provides the assistant-event directive name
	DirectiveNameAssistantEvent = DirectiveName{"assistant-event"}
	// DirectiveTypeView provides the view directive name
	DirectiveTypeView = DirectiveType{"view"}
	// DirectiveTypeAction provides the action directive name
	DirectiveTypeAction = DirectiveType{"action"}
)

// WebexAssistantHealthResponse is the response required for Webex Assistant Health Checks on our skill.
type WebexAssistantHealthResponse struct {
	Challenge string `json:"challenge,omitempty"`
	Status    string `json:"status,omitempty"`
}

// WebexAssistantRequest is the encrypted request we receive from Webex Assistant.
type WebexAssistantRequest struct {
	Signature string
	Message   string
}

// WebexAssistantMessage is the message we receive in the encrypted request from Webex Assistant.
type WebexAssistantMessage struct {
	// Text      []string `json:"text"`
	Text      string  `json:"text"`
	Context   Context `json:"context"`
	Params    Params  `json:"params"`
	Frame     Frame   `json:"frame"`
	History   History `json:"history"`
	Challenge string  `json:"challenge"`
}

// Frame contains information that needs to be preserved during multiple continuous interactions with the skill.
type Frame struct {
	//TODO: Find details about the Frame item
}

// History contains the history of the conversation in a multi-turn interaction.
type History struct {
	//TODO: Find details about the History item
}

// Params Contains information like time_zone, timestamp of the query, language, etc...
// One particular field here is target_dialogue_state this can be used to tell us what the user
// intended to do. In this particular case, if the field is equal to skill_intro, we need to return
// an introductory message from the skill. TODO:Add missing fields.
type Params struct {
	TargetDialogueState string      `json:"target_dialogue_state,omitempty"` // Possible values: "skill_intro", "TODO:?what else?"
	TimeZone            string      `json:"time_zone,omitempty"`
	Timestamp           int64       `json:"timestamp,omitempty"`
	Language            string      `json:"language,omitempty"`
	Locale              string      `json:"locale,omitempty"`
	DynamicResource     interface{} `json:"dynamic_resource,omitempty"`
	AllowedIntents      interface{} `json:"allowed_intents,omitempty"`
}

// Context contains some information about how the user is making the request.
type Context struct {
	OrgID               *string  `json:"orgId,omitempty"`               // The org id of the user making the request
	UserID              *string  `json:"userId,omitempty"`              // The id of the user making the request
	UserType            *string  `json:"userType,omitempty"`            // "user",  # The user type
	SupportedDirectives []string `json:"supportedDirectives,omitempty"` // The list of directives supported by the client making the request, e.g. sleep, listen, reply, speak
	DeveloperDeviceID   *string  `json:"developerDeviceId,omitempty"`
}

// WebexAssistantResponse is what we can send back to webex assistant.
type WebexAssistantResponse struct {
	Directives []WebexAssistantDirective `json:"directives"`
	Challenge  string                    `json:"challenge"`
}

// WebexAssistantDirective is the response we need to sent back to Webex Assistant
type WebexAssistantDirective struct {
	Name    DirectiveName `json:"name"` // reply, speak, listen, ui-hint, display-web-view, clear-web-view, assistant-event
	Type    DirectiveType `json:"type"` // view, action
	Payload Payload       `json:"payload"`
}

// Payload is the payload for the WebexAssistantDirective
type Payload struct {
	Text               *string           `json:"text,omitempty"` // may be []string?
	Delay              *int              `json:"delay,omitempty"`
	Prompt             *string           `json:"prompt,omitempty"`
	DisplayImmediately *bool             `json:"displayImmediately,omitempty"`
	Title              *string           `json:"title,omitempty"`
	URL                *string           `json:"url,omitempty"`
	Payload            map[string]string `json:"payload,omitempty"`
}
