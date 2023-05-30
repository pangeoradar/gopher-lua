package ast

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestNumberExprJSON(t *testing.T) {
	var expr Expr
	expr = &NumberExpr{Value: "10"}

	data, err := json.Marshal(expr)
	if err != nil {
		t.Error(err)
	}

	expr1, err := unmarshalExpr(data)
	if err != nil {
		t.Fatal(err)
	}
	if ne, ok := expr1.(*NumberExpr); !ok || ne.Value != "10" {
		t.Log("should be number expr with value equal 10")
		t.Fail()
	}
}

func TestRulePrint(t *testing.T) {
	p := json.RawMessage(`[{"names":["detection_windows"],"exprs":[{"value":"5s","_type":"string_expr"}],"_type":"local_assign_stmt"},{"names":["risk_score"],"exprs":[{"value":"7","_type":"number_expr"}],"_type":"local_assign_stmt"},{"names":["create_incident"],"exprs":[{"_type":"true_expr"}],"_type":"local_assign_stmt"},{"names":["assign_to_customer"],"exprs":[{"_type":"false_expr"}],"_type":"local_assign_stmt"},{"names":["grouped_by"],"exprs":[{"fields":[{"key":null,"value":{"value":"initiator.host.ip","_type":"string_expr"}},{"key":null,"value":{"value":"target.host.ip","_type":"string_expr"}},{"key":null,"value":{"value":"target.threat.name","_type":"string_expr"}}],"_type":"table_expr"}],"_type":"local_assign_stmt"},{"names":["aggregated_by"],"exprs":[{"fields":[{"key":null,"value":{"value":"target.threat.category","_type":"string_expr"}}],"_type":"table_expr"}],"_type":"local_assign_stmt"},{"names":["grouped_time_field"],"exprs":[{"value":"@timestamp","_type":"string_expr"}],"_type":"local_assign_stmt"},{"names":["template"],"exprs":[{"value":"\"Результат анализа.\n\nС узла {{ .First.initiator.host.ip | join \", \" }} была произведена попытка сканирования на уязвимость {{ .First.target.threat.name }} на узел {{ .First.target.host.ip | join \", \" }}.\"","_type":"string_expr"}],"_type":"local_assign_stmt"},{"names":["wl_src_ip"],"exprs":[{"fields":[],"_type":"table_expr"}],"_type":"local_assign_stmt"},{"name":{"func":{"value":"on_logline","_type":"ident_expr"},"receiver":null,"method":""},"func":{"par_list":{"has_vargs":false,"names":["logline"]},"stmts":[{"expr":{"func":null,"receiver":{"value":"logline","_type":"ident_expr"},"method":"get","args":[{"value":"smth","_type":"string_expr"}],"adjust_ret":false,"_type":"func_call_expr"},"_type":"func_call_stmt"},{"condition":{"expr":{"func":{"value":"contains","_type":"ident_expr"},"receiver":null,"method":"","args":[{"value":"wl_src_ip","_type":"ident_expr"},{"func":null,"receiver":{"value":"logline","_type":"ident_expr"},"method":"get","args":[{"value":"initiator.host.ip","_type":"string_expr"},{"value":"","_type":"string_expr"}],"adjust_ret":false,"_type":"func_call_expr"}],"adjust_ret":false,"_type":"func_call_expr"},"_type":"unary_not_op_expr"},"then":[{"expr":{"func":null,"receiver":{"value":"grouper1","_type":"ident_expr"},"method":"feed","args":[{"value":"logline","_type":"ident_expr"}],"adjust_ret":false,"_type":"func_call_expr"},"_type":"func_call_stmt"}],"else":null,"_type":"if_stmt"}],"_type":"function_expr"},"_type":"func_def_stmt"},{"name":{"func":{"value":"on_grouped","_type":"ident_expr"},"receiver":null,"method":""},"func":{"par_list":{"has_vargs":false,"names":["grouped"]},"stmts":[{"expr":{"func":{"value":"log","_type":"ident_expr"},"receiver":null,"method":"","args":[{"lhs":{"value":"grouped.aggregatedData.aggregated.total: ","_type":"string_expr"},"rhs":{"object":{"object":{"object":{"value":"grouped","_type":"ident_expr"},"key":{"value":"aggregatedData","_type":"string_expr"},"_type":"attr_get_expr"},"key":{"value":"aggregated","_type":"string_expr"},"_type":"attr_get_expr"},"key":{"value":"total","_type":"string_expr"},"_type":"attr_get_expr"},"_type":"string_concat_op_expr"}],"adjust_ret":false,"_type":"func_call_expr"},"_type":"func_call_stmt"},{"lhs":[{"value":"logline","_type":"ident_expr"}],"rhs":[{"object":{"object":{"object":{"value":"grouped","_type":"ident_expr"},"key":{"value":"aggregatedData","_type":"string_expr"},"_type":"attr_get_expr"},"key":{"value":"loglines","_type":"string_expr"},"_type":"attr_get_expr"},"key":{"value":"1","_type":"number_expr"},"_type":"attr_get_expr"}],"_type":"assign_stmt"},{"lhs":[{"value":"asset","_type":"ident_expr"}],"rhs":[{"func":{"value":"get_fields_value","_type":"ident_expr"},"receiver":null,"method":"","args":[{"value":"logline","_type":"ident_expr"},{"fields":[{"key":null,"value":{"value":"target.host.ip.0","_type":"string_expr"}},{"key":null,"value":{"value":"target.host.fqdn.0","_type":"string_expr"}},{"key":null,"value":{"value":"target.host.hostname.0","_type":"string_expr"}}],"_type":"table_expr"}],"adjust_ret":false,"_type":"func_call_expr"}],"_type":"assign_stmt"},{"lhs":[{"value":"meta","_type":"ident_expr"}],"rhs":[{"fields":[],"_type":"table_expr"}],"_type":"assign_stmt"},{"lhs":[{"value":"incident_identifier","_type":"ident_expr"}],"rhs":[{"func":{"value":"get_field_value","_type":"ident_expr"},"receiver":null,"method":"","args":[{"value":"logline","_type":"ident_expr"},{"value":"target.threat.category","_type":"string_expr"},{"value":"","_type":"string_expr"}],"adjust_ret":false,"_type":"func_call_expr"}],"_type":"assign_stmt"},{"expr":{"func":{"value":"alert","_type":"ident_expr"},"receiver":null,"method":"","args":[{"fields":[{"key":{"value":"template","_type":"string_expr"},"value":{"value":"template","_type":"ident_expr"}},{"key":{"value":"risk_level","_type":"string_expr"},"value":{"value":"risk_score","_type":"ident_expr"}},{"key":{"value":"asset_ip","_type":"string_expr"},"value":{"object":{"value":"asset","_type":"ident_expr"},"key":{"value":"1","_type":"number_expr"},"_type":"attr_get_expr"}},{"key":{"value":"asset_hostname","_type":"string_expr"},"value":{"object":{"value":"asset","_type":"ident_expr"},"key":{"value":"2","_type":"number_expr"},"_type":"attr_get_expr"}},{"key":{"value":"asset_fqdn","_type":"string_expr"},"value":{"object":{"value":"asset","_type":"ident_expr"},"key":{"value":"3","_type":"number_expr"},"_type":"attr_get_expr"}},{"key":{"value":"asset_mac","_type":"string_expr"},"value":{"value":"","_type":"string_expr"}},{"key":{"value":"create_incident","_type":"string_expr"},"value":{"value":"create_incident","_type":"ident_expr"}},{"key":{"value":"assign_to_customer","_type":"string_expr"},"value":{"value":"assign_to_customer","_type":"ident_expr"}},{"key":{"value":"logs","_type":"string_expr"},"value":{"object":{"object":{"value":"grouped","_type":"ident_expr"},"key":{"value":"aggregatedData","_type":"string_expr"},"_type":"attr_get_expr"},"key":{"value":"loglines","_type":"string_expr"},"_type":"attr_get_expr"}},{"key":{"value":"meta","_type":"string_expr"},"value":{"value":"meta","_type":"ident_expr"}},{"key":{"value":"incident_identifier","_type":"string_expr"},"value":{"value":"incident_identifier","_type":"ident_expr"}}],"_type":"table_expr"}],"adjust_ret":false,"_type":"func_call_expr"},"_type":"func_call_stmt"},{"expr":{"func":null,"receiver":{"value":"grouper1","_type":"ident_expr"},"method":"clear","args":[],"adjust_ret":false,"_type":"func_call_expr"},"_type":"func_call_stmt"}],"_type":"function_expr"},"_type":"func_def_stmt"},{"lhs":[{"value":"grouper1","_type":"ident_expr"}],"rhs":[{"func":{"object":{"value":"grouper","_type":"ident_expr"},"key":{"value":"new","_type":"string_expr"},"_type":"attr_get_expr"},"receiver":null,"method":"","args":[{"value":"grouped_by","_type":"ident_expr"},{"value":"aggregated_by","_type":"ident_expr"},{"value":"grouped_time_field","_type":"ident_expr"},{"value":"detection_windows","_type":"ident_expr"},{"value":"on_grouped","_type":"ident_expr"}],"adjust_ret":false,"_type":"func_call_expr"}],"_type":"assign_stmt"}]`)
	r, err := ParseRule(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(PrintRule(r))
}

func TestNestedLogicalOp(t *testing.T) {
	p := json.RawMessage(`[
  {
    "name": {
      "func": {
        "value": "on_logline",
        "_type": "ident_expr"
      },
      "receiver": null,
      "method": ""
    },
    "func": {
      "par_list": {
        "has_vargs": false,
        "names": [
          "logline"
        ]
      },
      "stmts": [
        {
          "condition": {
            "operator": "and",
            "lhs": {
              "operator": "and",
              "lhs": {
                "operator": "and",
                "lhs": {
                  "operator": "==",
                  "lhs": {
                    "func": null,
                    "receiver": {
                      "value": "logline",
                      "_type": "ident_expr"
                    },
                    "method": "get",
                    "args": [
                      {
                        "value": "action",
                        "_type": "string_expr"
                      },
                      {
                        "value": "",
                        "_type": "string_expr"
                      }
                    ],
                    "adjust_ret": false,
                    "_type": "func_call_expr"
                  },
                  "rhs": {
                    "value": "1234",
                    "_type": "string_expr"
                  },
                  "_type": "relational_op_expr"
                },
                "rhs": {
                  "operator": "==",
                  "lhs": {
                    "func": null,
                    "receiver": {
                      "value": "logline",
                      "_type": "ident_expr"
                    },
                    "method": "get",
                    "args": [
                      {
                        "value": "event.category",
                        "_type": "string_expr"
                      },
                      {
                        "value": "",
                        "_type": "string_expr"
                      }
                    ],
                    "adjust_ret": false,
                    "_type": "func_call_expr"
                  },
                  "rhs": {
                    "value": "test",
                    "_type": "string_expr"
                  },
                  "_type": "relational_op_expr"
                },
                "_type": "logical_op_expr"
              },
              "rhs": {
                "operator": "==",
                "lhs": {
                  "func": null,
                  "receiver": {
                    "value": "logline",
                    "_type": "ident_expr"
                  },
                  "method": "get",
                  "args": [
                    {
                      "value": "event.category",
                      "_type": "string_expr"
                    },
                    {
                      "value": "",
                      "_type": "string_expr"
                    }
                  ],
                  "adjust_ret": false,
                  "_type": "func_call_expr"
                },
                "rhs": {
                  "value": "test",
                  "_type": "string_expr"
                },
                "_type": "relational_op_expr"
              },
              "_type": "logical_op_expr"
            },
            "rhs": {
              "operator": "==",
              "lhs": {
                "func": null,
                "receiver": {
                  "func": null,
                  "receiver": {
                    "value": "logline",
                    "_type": "ident_expr"
                  },
                  "method": "get",
                  "args": [
                    {
                      "value": "action",
                      "_type": "string_expr"
                    },
                    {
                      "value": "",
                      "_type": "string_expr"
                    }
                  ],
                  "adjust_ret": false,
                  "_type": "func_call_expr"
                },
                "method": "lower",
                "args": null,
                "adjust_ret": false,
                "_type": "func_call_expr"
              },
              "rhs": {
                "func": null,
                "receiver": {
                  "value": "1234",
                  "_type": "string_expr"
                },
                "method": "lower",
                "args": null,
                "adjust_ret": false,
                "_type": "func_call_expr"
              },
              "_type": "relational_op_expr"
            },
            "_type": "logical_op_expr"
          },
          "then": [
            {
              "expr": {
                "func": null,
                "receiver": {
                  "value": "grouper1",
                  "_type": "ident_expr"
                },
                "method": "feed",
                "args": [
                  {
                    "value": "logline",
                    "_type": "ident_expr"
                  }
                ],
                "adjust_ret": false,
                "_type": "func_call_expr"
              },
              "_type": "func_call_stmt"
            }
          ],
          "else": null,
          "_type": "if_stmt"
        }
      ],
      "_type": "function_expr"
    },
    "_type": "func_def_stmt"
  },
  {
    "name": {
      "func": {
        "value": "on_grouped",
        "_type": "ident_expr"
      },
      "receiver": null,
      "method": ""
    },
    "func": {
      "par_list": {
        "has_vargs": false,
        "names": [
          "grouped"
        ]
      },
      "stmts": [
        {
          "names": [
            "logline"
          ],
          "exprs": [
            {
              "object": {
                "object": {
                  "object": {
                    "value": "grouped",
                    "_type": "ident_expr"
                  },
                  "key": {
                    "value": "aggregatedData",
                    "_type": "string_expr"
                  },
                  "_type": "attr_get_expr"
                },
                "key": {
                  "value": "loglines",
                  "_type": "string_expr"
                },
                "_type": "attr_get_expr"
              },
              "key": {
                "value": "1",
                "_type": "number_expr"
              },
              "_type": "attr_get_expr"
            }
          ],
          "_type": "local_assign_stmt"
        },
        {
          "condition": {
            "operator": "\u003e=",
            "lhs": {
              "object": {
                "object": {
                  "object": {
                    "value": "grouped",
                    "_type": "ident_expr"
                  },
                  "key": {
                    "value": "aggregatedData",
                    "_type": "string_expr"
                  },
                  "_type": "attr_get_expr"
                },
                "key": {
                  "value": "aggregated",
                  "_type": "string_expr"
                },
                "_type": "attr_get_expr"
              },
              "key": {
                "value": "total",
                "_type": "string_expr"
              },
              "_type": "attr_get_expr"
            },
            "rhs": {
              "value": "1",
              "_type": "number_expr"
            },
            "_type": "relational_op_expr"
          },
          "then": [
            {
              "expr": {
                "func": {
                  "value": "alert",
                  "_type": "ident_expr"
                },
                "receiver": null,
                "method": "",
                "args": [
                  {
                    "fields": [
                      {
                        "key": {
                          "value": "template",
                          "_type": "string_expr"
                        },
                        "value": {
                          "value": "",
                          "_type": "string_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "risk_level",
                          "_type": "string_expr"
                        },
                        "value": {
                          "value": "0.0",
                          "_type": "number_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "asset_ip",
                          "_type": "string_expr"
                        },
                        "value": {
                          "func": null,
                          "receiver": {
                            "value": "logline",
                            "_type": "ident_expr"
                          },
                          "method": "get_asset_data",
                          "args": [
                            {
                              "value": "",
                              "_type": "string_expr"
                            }
                          ],
                          "adjust_ret": false,
                          "_type": "func_call_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "asset_hostname",
                          "_type": "string_expr"
                        },
                        "value": {
                          "func": null,
                          "receiver": {
                            "value": "logline",
                            "_type": "ident_expr"
                          },
                          "method": "get_asset_data",
                          "args": [
                            {
                              "value": "",
                              "_type": "string_expr"
                            }
                          ],
                          "adjust_ret": false,
                          "_type": "func_call_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "asset_fqdn",
                          "_type": "string_expr"
                        },
                        "value": {
                          "func": null,
                          "receiver": {
                            "value": "logline",
                            "_type": "ident_expr"
                          },
                          "method": "get_asset_data",
                          "args": [
                            {
                              "value": "",
                              "_type": "string_expr"
                            }
                          ],
                          "adjust_ret": false,
                          "_type": "func_call_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "asset_mac",
                          "_type": "string_expr"
                        },
                        "value": {
                          "func": null,
                          "receiver": {
                            "value": "logline",
                            "_type": "ident_expr"
                          },
                          "method": "get_asset_data",
                          "args": [
                            {
                              "value": "",
                              "_type": "string_expr"
                            }
                          ],
                          "adjust_ret": false,
                          "_type": "func_call_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "create_incident",
                          "_type": "string_expr"
                        },
                        "value": {
                          "_type": "false_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "assign_to_customer",
                          "_type": "string_expr"
                        },
                        "value": {
                          "_type": "false_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "incident_identifier",
                          "_type": "string_expr"
                        },
                        "value": {
                          "value": "",
                          "_type": "string_expr"
                        }
                      },
                      {
                        "key": {
                          "value": "logs",
                          "_type": "string_expr"
                        },
                        "value": {
                          "object": {
                            "object": {
                              "value": "grouped",
                              "_type": "ident_expr"
                            },
                            "key": {
                              "value": "aggregatedData",
                              "_type": "string_expr"
                            },
                            "_type": "attr_get_expr"
                          },
                          "key": {
                            "value": "loglines",
                            "_type": "string_expr"
                          },
                          "_type": "attr_get_expr"
                        }
                      }
                    ],
                    "_type": "table_expr"
                  }
                ],
                "adjust_ret": false,
                "_type": "func_call_expr"
              },
              "_type": "func_call_stmt"
            }
          ],
          "else": null,
          "_type": "if_stmt"
        }
      ],
      "_type": "function_expr"
    },
    "_type": "func_def_stmt"
  },
  {
    "lhs": [
      {
        "value": "grouper1",
        "_type": "ident_expr"
      }
    ],
    "rhs": [
      {
        "func": {
          "object": {
            "value": "grouper",
            "_type": "ident_expr"
          },
          "key": {
            "value": "new",
            "_type": "string_expr"
          },
          "_type": "attr_get_expr"
        },
        "receiver": null,
        "method": "",
        "args": [
          {
            "fields": [
              {
                "key": null,
                "value": {
                  "value": "target.host.ip",
                  "_type": "string_expr"
                }
              },
              {
                "key": null,
                "value": {
                  "value": "target.host.hostname",
                  "_type": "string_expr"
                }
              }
            ],
            "_type": "table_expr"
          },
          {
            "fields": [
              {
                "key": null,
                "value": {
                  "value": "target.host.ip",
                  "_type": "string_expr"
                }
              },
              {
                "key": null,
                "value": {
                  "value": "target.host.hostname",
                  "_type": "string_expr"
                }
              }
            ],
            "_type": "table_expr"
          },
          {
            "value": "@timestamp,RFC3339Nano",
            "_type": "string_expr"
          },
          {
            "value": "5m",
            "_type": "string_expr"
          },
          {
            "value": "on_grouped",
            "_type": "ident_expr"
          }
        ],
        "adjust_ret": false,
        "_type": "func_call_expr"
      }
    ],
    "_type": "assign_stmt"
  }
]`)
	r, err := ParseRule(p)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(PrintRule(r))
}
