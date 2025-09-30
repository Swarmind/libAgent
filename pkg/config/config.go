package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const EnvPrefixKey = "LIBAGENT_ENV_PREFIX"

type Config struct {
	AIURL              string             `env:"AI_URL"`
	AIToken            string             `env:"AI_TOKEN"`
	Model              string             `env:"MODEL"`
	DefaultCallOptions DefaultCallOptions `env:"AI_DEFAULT_CALL_OPTION"`

	ReWOODisable            bool               `env:"REWOO_DISABLE"`
	RewOODefaultCallOptions DefaultCallOptions `env:"REWOO_DEFAULT_CALL_OPTION"`

	SemanticSearchDisable        bool   `env:"SEMANTIC_SEARCH_DISABLE"`
	SemanticSearchAIURL          string `env:"AI_URL,SEMANTIC_SEARCH_AI_URL"`
	SemanticSearchAIToken        string `env:"AI_TOKEN,SEMANTIC_SEARCH_AI_TOKEN"`
	SemanticSearchDBConnection   string `env:"SEMANTIC_SEARCH_DB_CONNECTION"`
	SemanticSearchEmbeddingModel string `env:"SEMANTIC_SEARCH_EMBEDDING_MODEL"`
	SemanticSearchMaxResults     int    `env:"SEMANTIC_SEARCH_MAX_RESULTS"`

	DDGSearchDisable    bool   `env:"DDG_SEARCH_DISABLE"`
	DDGSearchUserAgent  string `env:"DDG_SEARCH_USER_AGENT"`
	DDGSearchMaxResults int    `env:"DDG_SEARCH_MAX_RESULTS"`

	WebReaderDisable bool `env:"WEB_READER_DISABLE"`

	NmapDisable bool `env:"NMAP_DISABLE"`
	MsfDisable  bool `env:"MSF_DISABLE"`

	CommandExecutorDisable  bool              `env:"COMMAND_EXECUTOR_DISABLE"`
	CommandExecutorCommands map[string]string `env:"COMMAND_EXECUTOR_CMD_*"`
}

// See tmc/langchaingo/llms/options.go
type DefaultCallOptions struct {
	// Model is the model to use.
	Model *string `env:"MODEL"`
	// CandidateCount is the number of response candidates to generate.
	CandidateCount *int `env:"CANDIDATE_COUNT"`
	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens *int `env:"MAX_TOKENS"`
	// Temperature is the temperature for sampling, between 0 and 1.
	Temperature *float64 `env:"TEMPERATURE"`
	// StopWords is a list of words to stop on.
	StopWords *[]string `env:"STOP_WORDS"`
	// TopK is the number of tokens to consider for top-k sampling.
	TopK *int `env:"TOP_K"`
	// TopP is the cumulative probability for top-p sampling.
	TopP *float64 `env:"TOP_P"`
	// Seed is a seed for deterministic sampling.
	Seed *int `env:"SEED"`
	// MinLength is the minimum length of the generated text.
	MinLength *int `env:"MIN_LENGTH"`
	// MaxLength is the maximum length of the generated text.
	MaxLength *int `env:"MAX_LENGTH"`
	// N is how many chat completion choices to generate for each input message.
	N *int `env:"N"`
	// RepetitionPenalty is the repetition penalty for sampling.
	RepetitionPenalty *float64 `env:"REPETITION_PENALTY"`
	// FrequencyPenalty is the frequency penalty for sampling.
	FrequencyPenalty *float64 `env:"FREQUENCY_PENALTY"`
	// PresencePenalty is the presence penalty for sampling.
	PresencePenalty *float64 `env:"PRESENCE_PENALTY"`
	// JSONMode is a flag to enable JSON mode.
	JSONMode *bool `env:"JSON"`
	// ResponseMIMEType MIME type of the generated candidate text.
	// Supported MIME types are: text/plain: (default) Text output.
	// application/env: JSON response in the response candidates.
	ResponseMIMEType *string `env:"RESPONSE_MIME_TYPE"`
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg(".env file load")
	}

	envPrefix := os.Getenv(EnvPrefixKey)
	if envPrefix != "" && strings.HasSuffix(envPrefix, "_") {
		envPrefix = strings.TrimSuffix(envPrefix, "_")
	}

	val := reflect.ValueOf(&cfg).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if err := processField(fieldVal, field, envPrefix); err != nil {
			return Config{}, err
		}
	}

	if cfg.AIURL == "" {
		return cfg, fmt.Errorf("empty AI URL")
	}
	if cfg.AIToken == "" {
		return cfg, fmt.Errorf("empty AI Token")
	}
	if cfg.Model == "" {
		return cfg, fmt.Errorf("empty model")
	}

	return cfg, nil
}

func processField(fieldVal reflect.Value, field reflect.StructField, currentPrefix string) error {
	envTag := field.Tag.Get("env")
	if envTag == "" {
		return nil
	}

	if currentPrefix != "" {
		envTag = strings.Join([]string{currentPrefix, envTag}, "_")
	}

	fieldType := field.Type
	fieldKind := fieldType.Kind()

	if fieldKind == reflect.Struct {
		for i := 0; i < fieldType.NumField(); i++ {
			nestedFieldVal := fieldVal.Field(i)
			nestedField := fieldType.Field(i)
			if err := processField(
				nestedFieldVal,
				nestedField,
				envTag,
			); err != nil {
				return err
			}
		}
		return nil
	}

	envKeys := strings.Split(envTag, ",")
	var envValue string

	for _, key := range envKeys {
		key = strings.TrimSpace(key)
		envName := key
		if value, ok := os.LookupEnv(envName); ok && value != "" {
			envValue = value
			break
		}
	}

	if strings.HasSuffix(envTag, "_*") {
		return processWildcardField(fieldVal, field, envTag)
	}

	if envValue == "" {
		return nil
	}
	return setFieldValue(fieldVal, envValue, field.Name, envTag)
}

func processWildcardField(fieldVal reflect.Value, field reflect.StructField, envTag string) error {
	keyPrefix := strings.TrimSuffix(envTag, "_*")

	envMap := map[string]string{}
	envArrayString := map[int]string{}
	var mapKeyType reflect.Kind

	fieldKind := fieldVal.Kind()

	switch fieldKind {
	case reflect.Map:
		mapType := field.Type
		mapKeyType = mapType.Key().Kind()
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.MakeMap(mapType))
		}
	case reflect.Slice:
		if fieldVal.IsNil() {
			fieldVal.Set(reflect.MakeSlice(field.Type, 0, 0))
		}
	default:
		return fmt.Errorf("field '%s' with wildcard tag is not a map or slice", field.Name)
	}

	regex := fmt.Sprintf("^%s_(.+)$", regexp.QuoteMeta(keyPrefix))
	r := regexp.MustCompile(regex)

	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]

		if matches := r.FindStringSubmatch(key); len(matches) == 2 {
			suffix := matches[1]
			switch fieldKind {
			case reflect.Map:
				envMap[suffix] = value
			case reflect.Slice:
				if index, err := strconv.Atoi(suffix); err == nil {
					envArrayString[index] = value
				} else if mapKeyType == reflect.String {
					envMap[suffix] = value
				}
			}
		}
	}

	switch fieldKind {
	case reflect.Map:
		mapKeyType := field.Type.Key().Kind()
		mapValueType := field.Type.Elem()
		for k, v := range envMap {
			mapValue := reflect.New(mapValueType).Elem()
			if err := setFieldValue(mapValue, v, fmt.Sprintf("%s[%s]", fieldVal.Type().String(), k), envTag); err != nil {
				return err
			}
			keyVal := reflect.New(field.Type.Key()).Elem()
			if mapKeyType == reflect.String {
				keyVal.SetString(k)
			} else {
				return fmt.Errorf("unsupported map key type for wildcard '%s': %v", envTag, mapKeyType)
			}
			fieldVal.SetMapIndex(keyVal, mapValue)
		}
		return nil
	case reflect.Slice:
		sliceElemType := field.Type.Elem()
		if len(envArrayString) > 0 {
			sliceVal := reflect.MakeSlice(field.Type, len(envArrayString), len(envArrayString))
			for i, v := range envArrayString {
				elemVal := sliceVal.Index(i)
				if err := setFieldValue(elemVal, v, fmt.Sprintf("%s[%d]", fieldVal.Type().String(), i), envTag); err != nil {
					return err
				}
			}
			fieldVal.Set(sliceVal)
		} else if len(envMap) > 0 {
			sliceVal := reflect.MakeSlice(field.Type, 0, len(envMap))
			for _, v := range envMap {
				elemVal := reflect.New(sliceElemType).Elem()
				if err := setFieldValue(elemVal, v, fmt.Sprintf("%s[]", fieldVal.Type().String()), envTag); err != nil {
					return err
				}
				sliceVal = reflect.Append(sliceVal, elemVal)
			}
			fieldVal.Set(sliceVal)
		}
		return nil
	}

	return nil
}

func setFieldValue(fieldVal reflect.Value, envValue string, fieldName, envTag string) error {
	fieldType := fieldVal.Type()
	fieldKind := fieldType.Kind()

	if fieldKind == reflect.Ptr {
		if fieldVal.IsNil() {
			underlyingType := fieldType.Elem()
			ptr := reflect.New(underlyingType)
			fieldVal.Set(ptr)
			fieldVal = fieldVal.Elem()
		} else {
			fieldVal = fieldVal.Elem()
		}
		fieldType = fieldVal.Type()
		fieldKind = fieldType.Kind()
	}

	switch fieldKind {
	case reflect.String:
		fieldVal.SetString(envValue)
	case reflect.Int:
		intValue, err := strconv.Atoi(envValue)
		if err != nil {
			return fmt.Errorf("failed to parse int for field '%s' from env '%s': %w", fieldName, envTag, err)
		}
		fieldVal.SetInt(int64(intValue))
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(strings.ToLower(envValue))
		if err != nil {
			return fmt.Errorf("failed to parse bool for field '%s' from env '%s': %w", fieldName, envTag, err)
		}
		fieldVal.SetBool(boolValue)
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float64 for field '%s' from env '%s': %w", fieldName, envTag, err)
		}
		fieldVal.SetFloat(floatValue)
	case reflect.Float32:
		floatValue, err := strconv.ParseFloat(envValue, 32)
		if err != nil {
			return fmt.Errorf("failed to parse float32 for field '%s' from env '%s': %w", fieldName, envTag, err)
		}
		fieldVal.SetFloat(floatValue)
	case reflect.Slice:
		envSplit := strings.Split(envValue, ",")
		sliceType := fieldVal.Type()
		sliceVal := reflect.MakeSlice(sliceType, len(envSplit), len(envSplit))

		for i, v := range envSplit {
			elemVal := sliceVal.Index(i)
			if err := setFieldValue(elemVal, v, fmt.Sprintf("%s[%d]", elemVal.Type().String(), i), envTag); err != nil {
				return fmt.Errorf("failed to set slice value for field '%s' from env '%s': %w", fieldName, envTag, err)
			}
		}
		fieldVal.Set(sliceVal)
	default:
		return fmt.Errorf("unsupported type %s for field '%s' with env tag '%s'", fieldVal.Kind().String(), fieldName, envTag)
	}
	return nil
}
