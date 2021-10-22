package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	wxas "github.com/darrenparkinson/wxa-skills-go"
)

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

//TODO: Validate that this works ok.
func (app *application) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	encodedSignature := r.URL.Query().Get("signature")
	encodedCipher := r.URL.Query().Get("challenge")
	if encodedSignature == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "missing signature")
		return
	}
	if encodedCipher == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "missing message")
		return
	}
	decodedSignature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		app.errorLog.Println("error decoding signature:", err.Error())
		app.errorResponse(w, r, http.StatusBadRequest, "error decoding signature")
		return
	}
	if !verifySignature(app.config.Skill.Secret, encodedCipher, decodedSignature) {
		app.errorLog.Println("message has invalid signature and will not be processed")
		app.invalidSignatureResponse(w, r)
		return
	}
	decryptedChallenge, err := decryptMessage(app.config.Skill.PrivateKey, encodedCipher)
	if err != nil {
		app.errorLog.Println("message has invalid signature and will not be processed", err)
		app.errorResponse(w, r, http.StatusBadRequest, "unable to decrypt message")
		return
	}
	whr := wxas.WebexAssistantHealthResponse{
		Challenge: decryptedChallenge,
		Status:    "OK",
	}
	renderJSON(w, whr)
}

func (app *application) handleSkills(w http.ResponseWriter, r *http.Request) {

	var wr wxas.WebexAssistantRequest
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

	if wr.Signature == "" || wr.Message == "" {
		app.invalidRequestResponse(w, r)
		return
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(wr.Signature)
	if err != nil {
		app.errorLog.Println("error decoding signature:", err.Error())
		app.errorResponse(w, r, http.StatusBadRequest, "error decoding signature")
		return
	}
	if !verifySignature(app.config.Skill.Secret, wr.Message, decodedSignature) {
		app.errorLog.Println("message has invalid signature and will not be processed")
		app.invalidSignatureResponse(w, r)
		return
	}

	decryptedMessage, err := decryptMessage(app.config.Skill.PrivateKey, wr.Message)
	if err != nil {
		app.errorLog.Println("unable to decrypt message:", err)
		app.errorResponse(w, r, http.StatusBadRequest, "unable to decrypt message")
		return
	}
	var wam wxas.WebexAssistantMessage
	err = json.Unmarshal([]byte(decryptedMessage), &wam)
	if err != nil {
		app.errorLog.Println("error unmarshalling message", err)
		app.errorResponse(w, r, http.StatusBadRequest, "unable to unmarshal message")
	}
	shouldListen := false
	var text string
	if wam.Params.TargetDialogueState == "skill_intro" {
		text = "This is the echo skill.  Say something and I will echo it back."
		shouldListen = true
	} else {
		// TODO: this might be a string or []string I think, need to check.
		if len(wam.Text) > 0 {
			// text = wam.Text[0]
			text = wam.Text
		} else {
			text = "Hmm... I didn't get anything to echo"
		}
	}
	resp, err := buildResponse(text, wam.Challenge, shouldListen)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	renderJSON(w, resp)
}

func buildResponse(text string, challenge string, shouldListen bool) (wxas.WebexAssistantResponse, error) {
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
	wr.Challenge = challenge
	return wr, nil
}
