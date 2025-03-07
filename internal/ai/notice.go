package ai

import "fmt"

func (o *OpenAI) OperationNotUnderstood() (string, error) {
	return o.CallAISingle(dontUnderstandPrompt)
}

func (o *OpenAI) ShowHelpMessage(websiteURL string) (string, error) {
	return o.CallAISingle(fmt.Sprintf(showHelpPrompt, websiteURL))
}

func (o *OpenAI) AskIfApplyConfig() (string, error) {
	return o.CallAISingle(configValidatedPrompt)
}

func (o *OpenAI) ShowConfigDiscarded() (string, error) {
	return o.CallAISingle(configDiscardedPrompt)
}

var dontUnderstandPrompt = `generate a short humorous sentence expressing "I didn't understand what do you mean, we only support validating JSON configuration for the Kong Gateway and if you don't know how to write, you can ask for help"`
var showHelpPrompt = `generate a short sentence to guide a user asking for help on writing Kong Gateway configuration to the website at "%s", make sure the link is included in the output text, behave like a wise old man.`
var configDiscardedPrompt = `generate a short sentence to tell the configuration is discarded and the gateway won't be changed.`
var configValidatedPrompt = `generate a short compliment to a user that they've successfully written a valid configuration and ask them if they would like to apply it now`
