package ai

import (
	"encoding/json"
	"fmt"
	"github.com/jijiechen/dami-ultra/internal/business"
	"strings"
)

type UserOperation string

var OperationNone UserOperation = ""
var OperationCheckValidity UserOperation = "check-validity"
var OperationApplyConfigYes UserOperation = "apply-config-yes"
var OperationApplyConfigNo UserOperation = "apply-config-no"
var OperationShowHelp UserOperation = "show-help"

func (o *OpenAI) GetOperation(messages []business.Message) (UserOperation, error) {

	var promptMessages []business.Message

	msgLen := len(messages)
	if msgLen == 1 {
		promptMessages = append(promptMessages, business.Message{
			Author:  "system",
			Content: "Hi, I am Dami, your assistant. I am here to help you with your Kong configuration. Please provide me with the configuration you want to apply.",
		})
	} else {
		promptMessages = append(promptMessages, messages[msgLen-2])
	}
	promptMessages = append(promptMessages, messages[msgLen-1])
	promptMsgJson, err := json.Marshal(promptMessages)
	if err != nil {
		return "", err
	}

	operationResp, err := o.CallAISingle(fmt.Sprintf(operationPromptTemplate, promptMsgJson))
	if err != nil {
		return OperationNone, err
	}

	operationResp = strings.TrimSpace(operationResp)
	switch operationResp {
	case "check-validity":
		return OperationCheckValidity, nil
	case "apply-config-yes":
		return OperationApplyConfigYes, nil
	case "apply-config-no":
		return OperationApplyConfigNo, nil
	case "show-help":
		return OperationShowHelp, nil
	}
	return OperationNone, nil
}

var operationPromptTemplate = `Read the conversation attached and generate which operation matches the user's intention the best according to the conversation. 
The conversation is given in JSON wrapped in the pair of <code></code> and is attached in the end of this message.

The following table in markdown format shows the operations the user may be asking for:
| Operation Name | Description |
|-------|------|
| check-validity | matches when the content given by the user is a valid JSON string |
| apply-config-yes | the user wants to apply a configuration |
| abort-config-no | the user does not want to apply a configuration |
| show-help | the user asks for help |

Important notes for generating output:
1. if the user's intent matched one of supported operations, please ONLY output the matched operation name from the table in the first column, without any other explanation or decoration.
2. if the user's intent did not match any supported operations, please output exactly, without quotes: "Operation Not Supported"

JSON describing conversation messages:
<code>
%s
</code>
`
