package main

import (
	"errors"
	"net/http"

	wxas "github.com/darrenparkinson/wxa-skills-go"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) handleSkills(w http.ResponseWriter, r *http.Request) {
	var wr wxas.WebexAssistantMessage
	err := decodeJSONBody(w, r, &wr)
	if err != nil {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			app.errorLog.Println(err.Error())
			app.errorResponse(w, r, mr.status, mr.msg)
		} else {
			app.errorLog.Println(err.Error())
			app.invalidRequestResponse(w, r)
		}
		return
	}

	shouldListen := false
	var text string
	if wr.Params.TargetDialogueState == "skill_intro" {
		text = "This is the echo skill.  Say something and I will echo it back."
		shouldListen = true
	} else {
		// TODO: this might be a string or []string I think, need to check.
		if len(wr.Text) > 0 {
			// text = wr.Text[0]
			text = wr.Text
		} else {
			text = "Hmm... I didn't get anything to echo"
		}
	}
	resp, err := buildResponse(text, shouldListen)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	renderJSON(w, resp)
}

func buildResponse(text string, shouldListen bool) (wxas.WebexAssistantResponse, error) {
	var wr wxas.WebexAssistantResponse
	var sleepOrListen wxas.DirectiveName
	if shouldListen {
		sleepOrListen = wxas.DirectiveNameListen
	} else {
		sleepOrListen = wxas.DirectiveNameSleep
	}
	wr.Directives = []wxas.WebexAssistantDirective{
		{Name: wxas.DirectiveNameReply, Type: wxas.DirectiveTypeView, Payload: wxas.Payload{Text: &text}},
		{Name: wxas.DirectiveNameSpeak, Type: wxas.DirectiveTypeAction, Payload: wxas.Payload{Text: &text}},
		{Name: sleepOrListen, Type: wxas.DirectiveTypeAction},
	}
	return wr, nil
}
