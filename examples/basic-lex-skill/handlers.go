package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/lexruntimeservice"
	owm "github.com/briandowns/openweathermap"
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

	// Now process the message
	fmt.Printf("%+v\n", wam)
	resp, err := app.buildResponse(wam)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	renderJSON(w, resp)
}

func (app *application) buildResponse(wam wxas.WebexAssistantMessage) (wxas.WebexAssistantResponse, error) {
	var wr wxas.WebexAssistantResponse
	if wam.Params.TargetDialogueState == "skill_intro" {
		return buildSkillIntroResponse(wam), nil
	}
	lr, err := app.lex.PostText(&lexruntimeservice.PostTextInput{
		BotAlias:  &app.config.Lex.Alias,
		BotName:   &app.config.Lex.BotName,
		InputText: &wam.Text,
		UserId:    wam.Context.UserID,
	})
	if err != nil {
		app.errorLog.Printf("error communicating with lex: %s", err)
		return wr, err
	}
	city, text := "", ""
	listenOrSleep := wxas.DirectiveNameSleep
	if lr.DialogState != nil && lr.IntentName != nil {
		if *lr.IntentName != "CityWeather" {
			text = "That isn't a skill I have just yet."
		} else {
			switch *lr.DialogState {
			case "ElicitSlot":
				text = *lr.Message
				listenOrSleep = wxas.DirectiveNameListen
			case "ReadyForFulfillment":
				city = *lr.Slots["city"]
				w, err := owm.NewCurrent("C", "en", app.config.OpenWeatherMap.APIKey)
				if err != nil {
					app.errorLog.Println("error retrieving weather information")
					text = "Sorry, there was an error retrieving weather information."
					break
				}
				w.CurrentByName(city)
				text = fmt.Sprintf("The current weather in %s shows %s, with a low of %2.0f and a high of %2.0f.", w.Name, w.Weather[0].Description, w.Main.TempMin, w.Main.TempMax)
				listenOrSleep = wxas.DirectiveNameSleep
			}
		}
	}
	if text == "" {
		text = "Sorry, I have nothing to say to that."
	}

	fmt.Printf("%+v\n", lr)
	wr.Directives = []wxas.WebexAssistantDirective{
		{Name: wxas.DirectiveNameReply, Type: wxas.DirectiveTypeView, Payload: wxas.Payload{Text: &text}},
		{Name: wxas.DirectiveNameSpeak, Type: wxas.DirectiveTypeAction, Payload: wxas.Payload{Text: &text}},
		{Name: listenOrSleep, Type: wxas.DirectiveTypeAction},
	}
	wr.Challenge = wam.Challenge
	return wr, nil
}

func buildSkillIntroResponse(wam wxas.WebexAssistantMessage) wxas.WebexAssistantResponse {
	var wr wxas.WebexAssistantResponse
	text := "Sorry, I didn't catch what you said."
	wr.Directives = []wxas.WebexAssistantDirective{
		{Name: wxas.DirectiveNameReply, Type: wxas.DirectiveTypeView, Payload: wxas.Payload{Text: &text}},
		{Name: wxas.DirectiveNameSpeak, Type: wxas.DirectiveTypeAction, Payload: wxas.Payload{Text: &text}},
		{Name: wxas.DirectiveNameListen, Type: wxas.DirectiveTypeAction},
	}
	wr.Challenge = wam.Challenge
	return wr
}
