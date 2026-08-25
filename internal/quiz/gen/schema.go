package gen

type ResponseSchema struct {
	Name   string
	Schema map[string]any
}

const generationSchemaName = "quiz_generation"

const generationSchemaJSON = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["questions"],
  "properties": {
    "questions": {
      "type": "array",
      "minItems": 5,
      "maxItems": 10,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": [
          "text",
          "options",
          "correctIndex",
          "explanation",
          "evidence"
        ],
        "properties": {
          "text": {
            "type": "string",
            "minLength": 1,
            "description": "Question testing understanding of the source material"
          },
          "options": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
            "uniqueItems": true,
            "items": {
              "type": "string",
              "minLength": 1
            },
            "description": "Exactly four unique answer options"
          },
          "correctIndex": {
            "type": "integer",
            "minimum": 0,
            "maximum": 3,
            "description": "Zero-based index of the only correct option"
          },
          "explanation": {
            "type": "string",
            "minLength": 1,
            "description": "Short explanation of why the answer is correct"
          },
          "evidence": {
            "type": "string",
            "minLength": 1,
            "description": "Exact excerpt from SOURCE supporting the correct answer"
          }
        }
      }
    }
  }
}`
