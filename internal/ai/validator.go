package ai

import (
	"encoding/json"
	"fmt"
	"github.com/jijiechen/dami-ultra/internal/business"
	"strings"
)

type ValidateOpenAIResponse struct {
	Valid bool `json:"valid,omitempty"`
	//RawConfiguration map[string]interface{} `json:"raw_configuration,omitempty"`
	RawConfiguration string `json:"raw_configuration,omitempty"`
	ErrorMessages    string `json:"error_messages,omitempty"`
}

func (o *OpenAI) ValidateKongConfiguration(kongConfigString string) (ValidateOpenAIResponse, error) {
	aiResp, err := o.CallAI(
		fmt.Sprintf(validatorSystemPromptTemplate, luaValidatorCode),
		[]business.Message{
			{Author: "system", Content: validatorAssistantPrompt},
			{Author: "user", Content: fmt.Sprintf(validatorUserPromptTemplate, kongConfigString)},
		})

	var respObj ValidateOpenAIResponse
	if err != nil {
		return respObj, err
	}

	fmt.Println(aiResp)
	aiResp = strings.ReplaceAll(aiResp, "```json", "")
	aiResp = strings.ReplaceAll(aiResp, "```", "")

	err = json.Unmarshal([]byte(aiResp), &respObj)
	return respObj, err
}

func (o *OpenAI) ExtractValidatedConfig(messages []business.Message) (string, error) {
	msgJsonString, err := json.Marshal(messages)
	if err != nil {
		return "", err
	}

	config, err := o.CallAISingle(fmt.Sprintf(validatorSummarizePromptTemplate, msgJsonString))
	if err != nil {
		return "", err
	}

	if config == "NOT_FOUND" {
		return "", fmt.Errorf("could not determine the last validated Kong configuration")
	}
	return config, nil
}

var validatorSystemPromptTemplate = `You are a schema validator for the Kong Gateway that validates user's input according to the rules defined by some Lua code.
Please read the validator code wrapped in the pair of <code></code>, use it to validate the user's input and generate output. The user's input should be in JSON format.

Important notes for generating output:
1. use the following JSON format to wrap your output: '{"valid": true, "raw_configuration": "<JSON string of the raw configuration when it's valid>", "error_messages": "<potential messages that laugh at the user>"}'.
2. DO NOT include any explanation or any code splitter in your output, you output MUST BE in valid JSON format, this is VERY IMPORTANT, otherwise your output will not be handled correctly.
3. when the given configuration IS VALID, please output 'true' using field 'valid' and attach the raw configuration in field 'raw_configuration'. 
4. when the given configuration IS NOT VALID, please output 'false' using field 'valid' and tell a polite joke to laugh at the user based on the error of the user's input.

Lua code of the validator:
<code>
%s
</code>
`

var validatorAssistantPrompt = `Please enter your Kong Gateway configuration in JSON format, and I will validate it for you.`

var validatorUserPromptTemplate = `Here is my Kong Gateway configuration in JSON format, it's wrapped in the pair of <code></code>:
<code>
%s
</code>`

var validatorSummarizePromptTemplate = `Read the conversation described in given JSON wrapped in the pair of <code></code>, extract and output the last validated configuration discussed in the conversation; ONLY output the JSON content extracted form the conversation and DO NOT output any explanation or additional words for coherence usage; if there is no such configuration found, output "NOT_FOUND", without quotes.

JSON describing conversation messages:
<code>
%s
</code>
`

var luaValidatorCode = `
local typedefs = require("kong.db.schema.typedefs")
local deprecation = require("kong.deprecation")


local kong_router_flavor = kong and kong.configuration and kong.configuration.router_flavor


local PATH_V1_DEPRECATION_MSG =
  "path_handling='v1' is deprecated and " ..
  (kong_router_flavor == "traditional" and
    "will be removed in future version, " or
    "will not work under 'expressions' or 'traditional_compatible' router_flavor, ") ..
  "please use path_handling='v0' instead"


local entity_checks = {
  { conditional = { if_field = "protocols",
                    if_match = { elements = { type = "string", not_one_of = { "grpcs", "https", "tls", "tls_passthrough" }}},
                    then_field = "snis",
                    then_match = { len_eq = 0 },
                    then_err = "'snis' can only be set when 'protocols' is 'grpcs', 'https', 'tls' or 'tls_passthrough'",
                  }
  },

  { custom_entity_check = {
    field_sources = { "path_handling" },
    fn = function(entity)
      if entity.path_handling == "v1" then
        deprecation(PATH_V1_DEPRECATION_MSG, { after = "3.0", })
      end

      return true
    end,
  }},
}

local snis_elements_type = typedefs.wildcard_host

if kong_router_flavor == "traditional" then
  snis_elements_type = typedefs.sni
end

local validate_route
if kong_router_flavor == "traditional_compatible" or kong_router_flavor == "expressions" then
  local ipairs = ipairs
  local tonumber = tonumber
  local re_match = ngx.re.match

  local router = require("resty.router.router")
  local transform = require("kong.router.transform")
  local get_schema = require("kong.router.atc").schema

  local is_null = transform.is_null
  local is_empty_field = transform.is_empty_field
  local amending_expression = transform.amending_expression

  local HTTP_PATH_SEGMENTS_PREFIX = "http.path.segments."
  local HTTP_PATH_SEGMENTS_SUFFIX_REG = [[^(0|[1-9]\d*)(_([1-9]\d*))?$]]

  validate_route = function(entity)
    local is_expression_empty =
      is_null(entity.expression)   -- expression is not a table

    local is_others_empty =
      is_empty_field(entity.snis) and
      is_empty_field(entity.sources) and
      is_empty_field(entity.destinations) and
      is_empty_field(entity.methods) and
      is_empty_field(entity.hosts) and
      is_empty_field(entity.paths) and
      is_empty_field(entity.headers)

    if is_expression_empty and is_others_empty then
      return true
    end

    if not is_expression_empty and not is_others_empty then
      return nil, "Router Expression failed validation: " ..
                  "cannot set 'expression' with " ..
                  "'methods', 'hosts', 'paths', 'headers', 'snis', 'sources' or 'destinations' " ..
                  "simultaneously"
    end

    local is_regex_priority_empty = is_null(entity.regex_priority) or
                                    entity.regex_priority == 0    -- default value 0 means 'no set'
    if not is_expression_empty and not is_regex_priority_empty then
      return nil, "Router Expression failed validation: " ..
                  "cannot set 'regex_priority' with 'expression' " ..
                  "simultaneously"
    end

    local is_priority_empty = is_null(entity.priority) or
                              entity.priority == 0    -- default value 0 means 'no set'
    if not is_others_empty and not is_priority_empty then
      return nil, "Router Expression failed validation: " ..
                  "cannot set 'priority' with " ..
                  "'methods', 'hosts', 'paths', 'headers', 'snis', 'sources' or 'destinations' " ..
                  "simultaneously"
    end

    local schema = get_schema(entity.protocols)
    local exp = amending_expression(entity)

    local fields, err = router.validate(schema, exp)
    if not fields then
      return nil, "Router Expression failed validation: " .. err
    end

    for _, f in ipairs(fields) do
      if f:find(HTTP_PATH_SEGMENTS_PREFIX, 1, true) then
        local suffix = f:sub(#HTTP_PATH_SEGMENTS_PREFIX + 1)
        local m = re_match(suffix, HTTP_PATH_SEGMENTS_SUFFIX_REG, "jo")

        if (suffix ~= "len") and
           (not m or (m[2] and tonumber(m[1]) >= tonumber(m[3]))) then
          return nil, "Router Expression failed validation: " ..
                      "illformed http.path.segments.* field"
        end
      end -- if f:find
    end -- for fields

    return true
  end

  table.insert(entity_checks,
    { custom_entity_check = {
      field_sources = { "id", "protocols",
                        "snis", "sources", "destinations",
                        "methods", "hosts", "paths", "headers",
                        "expression",
                        "regex_priority", "priority",
                      },
      run_with_missing_fields = true,
      fn = validate_route,
    } }
  )
end   -- if kong_router_flavor ~= "traditional"


local routes = {
    name         = "routes",
    primary_key  = { "id" },
    endpoint_key = "name",
    workspaceable = true,
    subschema_key = "protocols",

    fields = {
      { id             = typedefs.uuid, },
      { created_at     = typedefs.auto_timestamp_s },
      { updated_at     = typedefs.auto_timestamp_s },
      { name           = typedefs.utf8_name },

      { protocols      = { type     = "set",
                           description = "An array of the protocols this Route should allow.",
                           len_min  = 1,
                           required = true,
                           elements = typedefs.protocol,
                           mutually_exclusive_subsets = {
                             { "http", "https" },
                             { "tcp", "tls", "udp" },
                             { "tls_passthrough" },
                             { "grpc", "grpcs" },
                           },
                           default = { "http", "https" }, -- TODO: different default depending on service's scheme
                         }, },

      { https_redirect_status_code = { type = "integer",
                                       description = "The status code Kong responds with when all properties of a Route match except the protocol",
                                       one_of = { 426, 301, 302, 307, 308 },
                                       default = 426, required = true,
                                     }, },
      { strip_path     = { description = "When matching a Route via one of the paths, strip the matching prefix from the upstream request URL.", type = "boolean", required = true, default = true }, },
      { preserve_host  = { description = "When matching a Route via one of the hosts domain names, use the request Host header in the upstream request headers.", type = "boolean", required = true, default = false }, },
      { request_buffering  = { description = "Whether to enable request body buffering or not. With HTTP 1.1.", type = "boolean", required = true, default = true }, },
      { response_buffering  = { description = "Whether to enable response body buffering or not.", type = "boolean", required = true, default = true }, },

      { tags             = typedefs.tags },
      { service = { description = "The Service this Route is associated to. This is where the Route proxies traffic to.", type = "foreign", reference = "services" }, },

      { snis = { type = "set",
                 description = "A list of SNIs that match this Route.",
                 elements = snis_elements_type }, },
      { sources = typedefs.sources },
      { destinations = typedefs.destinations },

      { methods        = typedefs.methods },
      { hosts          = typedefs.hosts },
      { paths          = typedefs.router_paths },
      { headers = typedefs.headers {
        keys = typedefs.header_name {
          match_none = {
            {
              pattern = "^[Hh][Oo][Ss][Tt]$",
              err = "cannot contain 'host' header, which must be specified in the 'hosts' attribute",
            },
          },
        },
      } },

      { regex_priority = { description = "A number used to choose which route resolves a given request when several routes match it using regexes simultaneously.", type = "integer", default = 0 }, },
      { path_handling  = { description = "Controls how the Service path, Route path and requested path are combined when sending a request to the upstream.", type = "string", default = "v0", one_of = { "v0", "v1" }, }, },
    },  -- fields

    entity_checks = entity_checks,
} -- routes


if kong_router_flavor == "expressions" then

  local special_fields = {
    { expression = { description = "The route expression.", type = "string" }, },   -- not required now
    { priority = { description = "A number used to specify the matching order for expression routes. The higher the 'priority', the sooner an route will be evaluated. This field is ignored unless 'expression' field is set.", type = "integer", between = { 0, 2^46 - 1 }, required = true, default = 0 }, },
  }

  for _, v in ipairs(special_fields) do
    table.insert(routes.fields, v)
  end
end


return routes
`
