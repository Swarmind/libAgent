package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	shodanclient "github.com/shadowscatcher/shodan"
	"github.com/shadowscatcher/shodan/search"

	"github.com/Swarmind/libagent/internal/tools"
	"github.com/Swarmind/libagent/pkg/config"

	"github.com/tmc/langchaingo/llms"
)

var ShodanToolDefinition = llms.FunctionDefinition{
	Name: "shodan",
	Description: `Shodan internet intelligence platform API.
	Returns raw string json of collected data about matching resources.`,
	Parameters: map[string]any{
		"type": "object",
		"properties": map[string]any{
			"params": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "object",
						"description": "The root object for constructing a Shodan search query.",
						"properties": map[string]any{
							"Text": map[string]any{
								"type":        "string",
								"description": "Raw query text. You can use only this field if you want to. Overrides any other filter",
							},
							"all": map[string]any{
								"type":        "string",
								"description": "Filter by any text",
							},
							"ip": map[string]any{
								"type":        "string",
								"description": "IP-address",
							},
							"after": map[string]any{
								"type":        "string",
								"description": "Only show results that were collected after the given date (dd/mm/yyyy).",
							},
							"asn": map[string]any{
								"type":        "string",
								"description": "The Autonomous System Number that identifies the network the device is on; ex: \"AS15169\"",
							},
							"before": map[string]any{
								"type":        "string",
								"description": "Only show results that were collected before the given date (dd/mm/yyyy).",
							},
							"city": map[string]any{
								"type":        "string",
								"description": "Show results that are located in the given city.",
							},
							"country": map[string]any{
								"type":        "string",
								"description": "Show results that are located within the given country.",
							},
							"cpe": map[string]any{
								"type":        "string",
								"description": "Common platform enumeration",
							},
							"device": map[string]any{
								"type":        "string",
								"description": "Device type; ex: \"printer\", \"router\"",
							},
							"geo": map[string]any{
								"type":        "string",
								"description": "There are 2 modes to the geo filter: radius and bounding box. To limit results based on a radius around a pair of latitude/longitude, provide 3 parameters; ex: geo:50,50,100. If you want to find all results within a bounding box, supply the top left and bottom right coordinates for the region; ex: geo:10,10,50,50.",
							},
							"hostname": map[string]any{
								"type":        "string",
								"description": "Search for hosts that contain the given value in their hostname",
							},
							"isp": map[string]any{
								"type":        "string",
								"description": "Find devices based on the upstream owner of the IP netblock",
							},
							"link": map[string]any{
								"type": "string",
								"description": "Find devices depending on their connection to the Internet\n" +
									"Available values: \"Ethernet or modem\", \"generic tunnel or VPN\", DSL, \"IPIP or SIT\", SLIP, \"IPSec or GRE\", VLAN, \"jumbo Ethernet\", Google, GIF, PPTP, loopback, \"AX.25 radio modem\"",
							},
							"net": map[string]any{
								"type":        "string",
								"description": "Search by netblock using CIDR notation; ex: net:69.84.207.0/24",
							},
							"org": map[string]any{
								"type":        "string",
								"description": "Find devices based on the owner of the IP netblock.",
							},
							"os": map[string]any{
								"type":        "string",
								"description": "Filter results based on the operating system of the device.",
							},
							"postal": map[string]any{
								"type":        "string",
								"description": "Search by postal code.",
							},
							"product": map[string]any{
								"type":        "string",
								"description": "Filter using the name of the software/product; ex: product:Apache",
							},
							"state": map[string]any{
								"type":        "string",
								"description": "Search for devices based on the state/region they are located in",
							},
							"version": map[string]any{
								"type":        "string",
								"description": "Filter the results to include only products of the given version; ex: product:apache version:1.3.37",
							},
							// Shodan has this key duplicate, let llm use more advanced search
							// "ssl": map[string]any{
							// "type":        "string",
							// "description": "Search all SSL data",
							// },
							"port": map[string]any{
								"type":        "integer",
								"description": "Find devices based on the services/ports that are publicly exposed on the Internet",
							},
							"hash": map[string]any{
								"type":        "integer",
								"description": "Hash of the \"data\" property",
							},
							"has_ipv6": map[string]any{
								"type":        "boolean",
								"description": "If \"true\" only show results that were discovered on IPv6",
							},
							"has_screenshot": map[string]any{
								"type":        "boolean",
								"description": "If \"true\" only show results that have a screenshot available",
							},
							"has_ssl": map[string]any{
								"type":        "boolean",
								"description": "If \"true\" only show results that have SSL",
							},
							"has_vuln": map[string]any{
								"type":        "boolean",
								"description": "If \"true\" only show results that have vulnerabilities. Enterpise only.",
							},
							"region": map[string]any{
								"type":        "string",
								"description": "Region code",
							},
							"tag": map[string]any{
								"type":        "string",
								"description": "Host tag. Enterprise only.",
							},
							"scan": map[string]any{
								"type":        "string",
								"description": "Signature unknown",
							},
							"vuln": map[string]any{
								"type":        "string",
								"description": "Filter by vulnerability. Only available to academic users or Small Business API subscription and higher.",
							},
							"ssl": map[string]any{
								"type":        "object",
								"description": "SSL options",
								"properties": map[string]any{
									"chain_count": map[string]any{
										"type":        "integer",
										"description": "Number of certificates in the chain",
									},
									"alpn": map[string]any{
										"type":        "string",
										"description": "Application layer protocols such as HTTP/2 (\"h2\")",
									},
									"version": map[string]any{
										"type":        "string",
										"description": "Possible values: SSLv2, SSLv3, TLSv1, TLSv1.1, TLSv1.2",
									},
									"cert": map[string]any{
										"type":        "object",
										"description": "Various certificate options",
										"properties": map[string]any{
											"alg": map[string]any{
												"type":        "string",
												"description": "Certificate algorithm",
											},
											"expired": map[string]any{
												"type":        "boolean",
												"description": "Whether the SSL certificate is expired or not",
											},
											"extension": map[string]any{
												"type":        "string",
												"description": "Names of extensions in the certificate",
											},
											"serial": map[string]any{
												"type":        "string",
												"description": "Serial number as string",
											},
											"fingerprint": map[string]any{
												"type":        "string",
												"description": "SHA-1 fingerprint",
											},
											"issuer": map[string]any{
												"type":        "object",
												"description": "Cert issuer options",
												"properties": map[string]any{
													"cn": map[string]any{
														"type":        "string",
														"description": "Common name",
													},
												},
											},
											"subject": map[string]any{
												"type":        "object",
												"description": "Cert subject options",
												"properties": map[string]any{
													"cn": map[string]any{
														"type":        "string",
														"description": "Common name",
													},
												},
											},
											"pubkey": map[string]any{
												"type": "object",
												"properties": map[string]any{
													"bits": map[string]any{
														"type":        "integer",
														"description": "Number of bits in the public key",
													},
													"type": map[string]any{
														"type":        "string",
														"description": "Public key type",
													},
												},
											},
										},
									},
									"cipher": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"version": map[string]any{
												"type":        "string",
												"description": "SSL version of the preferred cipher",
											},
											"bits": map[string]any{
												"type":        "integer",
												"description": "Number of bits in the preferred cipher",
											},
											"name": map[string]any{
												"type":        "string",
												"description": "Name of the preferred cipher",
											},
										},
									},
								},
							},
							"bitcoin": map[string]any{
								"type":        "object",
								"description": "Bitcoin server options",
								"properties": map[string]any{
									"ip": map[string]any{
										"type":        "string",
										"description": "Find Bitcoin servers that had the given IP in their list of peers",
									},
									"ip_count": map[string]any{
										"type":        "integer",
										"description": "Find Bitcoin servers that return the given number of IPs in the list of peers",
									},
									"port": map[string]any{
										"type":        "integer",
										"description": "Find Bitcoin servers that had IPs with the given port in their list of peers",
									},
									"version": map[string]any{
										"type":        "string",
										"description": "Filter results based on the Bitcoin protocol version",
									},
								},
							},
							"telnet": map[string]any{
								"type":        "object",
								"description": "Telnet server options",
								"properties": map[string]any{
									"option": map[string]any{
										"type":        "string",
										"description": "Search all the options",
									},
									"do": map[string]any{
										"type":        "string",
										"description": "The server requests the client to support these options",
									},
									"dont": map[string]any{
										"type":        "string",
										"description": "The server requests the client to not support these options",
									},
									"will": map[string]any{
										"type":        "string",
										"description": "The server supports these options",
									},
									"wont": map[string]any{
										"type":        "string",
										"description": "The server doesnt support these options",
									},
								},
							},
							"http": map[string]any{
								"type":        "object",
								"description": "HTTP server options",
								"properties": map[string]any{
									"component": map[string]any{
										"type":        "string",
										"description": "Name of web technology used on the website",
									},
									"component_category": map[string]any{
										"type":        "string",
										"description": "Category of web components used on the website",
									},
									"html": map[string]any{
										"type":        "string",
										"description": "Search the HTML of the website for the given value.",
									},
									"title": map[string]any{
										"type":        "string",
										"description": "Search the title of the website",
									},
									"status": map[string]any{
										"type":        "integer",
										"description": "Response status code",
									},
									"html_hash": map[string]any{
										"type":        "integer",
										"description": "Hash of the website HTML",
									},
									"favicon": map[string]any{
										"type": "object",
										"properties": map[string]any{
											"hash": map[string]any{
												"type":        "integer",
												"description": "Hash of website favicon.ico file",
											},
										},
									},
									"robots_hash": map[string]any{
										"type":        "integer",
										"description": "Hash of website robots.txt file",
									},
									"securitytxt": map[string]any{
										"type":        "string",
										"description": "Search in contents of website's security.txt",
									},
									"waf": map[string]any{
										"type":        "string",
										"description": "Search by Web Application Firewall vendor/name",
									},
								},
							},
							"ntp": map[string]any{
								"type":        "object",
								"description": "NTP server options",
								"properties": map[string]any{
									"ip": map[string]any{
										"type":        "string",
										"description": "Find NTP servers that had the given IP in their monlist.",
									},
									"ip_count": map[string]any{
										"type":        "integer",
										"description": "Find NTP servers that return the given number of IPs in the initial monlist response.",
									},
									"port": map[string]any{
										"type":        "integer",
										"description": "Find NTP servers that had IPs with the given port in their monlist.",
									},
									"more": map[string]any{
										"type":        "boolean",
										"description": "Whether or not more IPs were available for the given NTP server.",
									},
								},
							},
							"screenshot": map[string]any{
								"type":        "object",
								"description": "Screenshot options",
								"properties": map[string]any{
									"label": map[string]any{
										"type":        "string",
										"description": "Label of screenshot (kind of tag, like \"login\", \"windows\")",
									},
								},
							},
							"shodan": map[string]any{
								"type":        "object",
								"description": "Shodan module options",
								"properties": map[string]any{
									"module": map[string]any{
										"type":        "string",
										"description": "Filter by shodan crawler module",
									},
								},
							},
							"snmp": map[string]any{
								"type":        "object",
								"description": "SNMP server options",
								"properties": map[string]any{
									"contact": map[string]any{
										"type":        "string",
										"description": "SNMP contact address",
									},
									"name": map[string]any{
										"type":        "string",
										"description": "SNMP server name",
									},
									"location": map[string]any{
										"type":        "string",
										"description": "Location",
									},
								},
							},
							"ssh": map[string]any{
								"type":        "object",
								"description": "SSH server options",
								"properties": map[string]any{
									"hassh": map[string]any{
										"type":        "string",
										"description": "HASSH Md5 fingerprint hash",
									},
									"type": map[string]any{
										"type":        "string",
										"description": "Type",
									},
								},
							},
						},
					},
					"facets": map[string]any{
						"type":        "array",
						"description": "A comma-separated list of properties to get summary information on. Property names can also be in the format of 'property:count'.",
						"items": map[string]any{
							"type": "string",
						},
					},
					"minify": map[string]any{
						"type":        "boolean",
						"description": "Whether or not to truncate some of the larger fields (default: true)",
					},
					"page": map[string]any{
						"type":        "integer",
						"description": "The page number to page through results 100 at a time (default: 1)",
						"minimum":     1,
					},
					"offset": map[string]any{
						"type":        "integer",
						"description": "Offset is an alternative to Page for pagination.",
						"minimum":     0,
					},
				},
			},
		},
	},
}

type ShodanTool struct {
	APIKey string
}

func (s ShodanTool) Call(ctx context.Context, input string) (string, error) {
	params, err := s.UnmarshalFromString(input)
	if err != nil {
		return "", fmt.Errorf("unmarshal from string")
	}

	client, _ := shodanclient.GetClient(
		s.APIKey,
		http.DefaultClient,

		// See https://github.com/shadowscatcher/shodan README.md:
		// "The client can be configured to automatically make one second pause
		// between requests (this interval required by Shodan's API terms of service)"
		true,
	)
	result, err := client.Search(ctx, params)
	if err != nil {
		return "", fmt.Errorf("shodan search failed: %w", err)
	}

	out, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(out), nil
}

func (s ShodanTool) UnmarshalFromString(input string) (search.Params, error) {
	dataMap := map[string]map[string]any{}
	if err := json.Unmarshal([]byte(input), &dataMap); err != nil {
		return search.Params{}, fmt.Errorf("unmarshal to map: %w", err)
	}
	if _, ok := dataMap["params"]; !ok {
		return search.Params{}, fmt.Errorf("no 'params' key")
	}

	params := search.Params{}
	err := unmarshalFromMap(dataMap["params"], &params)
	if err != nil {
		return search.Params{}, fmt.Errorf("unmarshal to shodan params: %w", err)
	}

	return params, nil
}

func init() {
	globalToolsRegistry = append(globalToolsRegistry,
		func(ctx context.Context, cfg config.Config) (*tools.ToolData, error) {
			if cfg.ShodanDisable {
				return nil, nil
			}
			if cfg.ShodanAPIKey == "" {
				return nil, fmt.Errorf("shodan api key empty")
			}

			shodanTool := ShodanTool{
				APIKey: cfg.ShodanAPIKey,
			}

			definition := ShodanToolDefinition

			return &tools.ToolData{
				Definition: definition,
				Call:       shodanTool.Call,
			}, nil
		},
	)
}

func unmarshalFromMap(dataMap map[string]any, targetStruct any) error {
	rv := reflect.ValueOf(targetStruct)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer to a struct")
	}
	rv = rv.Elem()

	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct, got %s", rv.Kind())
	}

	return unmarshalStruct(dataMap, rv)
}

func unmarshalStruct(dataMap map[string]any, structValue reflect.Value) error {

	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldType := structType.Field(i)
		fieldValue := structValue.Field(i)

		jsonKey := fieldType.Tag.Get("shodan_search")

		if jsonKey == "" && fieldType.Name == "Text" {
			jsonKey = "Text"
		}

		if !fieldValue.CanSet() {
			continue
		}

		// deal with params not shodan_search tagged fields
		if jsonKey == "" {
			jsonKey = strings.ToLower(fieldType.Name)
		}

		// deal with Shodan key duplicate, falling back to more advanced SSLOpts type
		if jsonKey == "ssl" && fieldType.Name == "SSL" {
			continue
		}

		mapValue, found := dataMap[jsonKey]
		if !found {
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			nestedMap, ok := mapValue.(map[string]any)
			if !ok {
				fmt.Println(mapValue, jsonKey, dataMap)
				return fmt.Errorf("expected map for field '%s', but got type %T", fieldType.Name, mapValue)
			}

			if err := unmarshalStruct(nestedMap, fieldValue); err != nil {
				return err
			}
			continue
		}

		switch fieldValue.Kind() {
		case reflect.String:
			if strVal, ok := mapValue.(string); ok {
				fieldValue.SetString(strVal)
			} else {
				return fmt.Errorf("field '%s' expected string, got %T", fieldType.Name, mapValue)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if floatVal, ok := mapValue.(float64); ok {
				fieldValue.SetInt(int64(floatVal))
			} else {
				return fmt.Errorf("field '%s' expected number, got %T", fieldType.Name, mapValue)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if floatVal, ok := mapValue.(float64); ok {
				if floatVal < 0 {
					return fmt.Errorf("field '%s' is uint but received negative value", fieldType.Name)
				}
				fieldValue.SetUint(uint64(floatVal))
			} else {
				return fmt.Errorf("field '%s' expected number (uint), got %T", fieldType.Name, mapValue)
			}
		case reflect.Bool:
			if boolVal, ok := mapValue.(bool); ok {
				fieldValue.SetBool(boolVal)
			} else {
				return fmt.Errorf("field '%s' expected boolean, got %T", fieldType.Name, mapValue)
			}
		default:
		}
	}

	return nil
}
