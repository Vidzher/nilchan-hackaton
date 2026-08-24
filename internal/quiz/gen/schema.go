package gen

const generationSchemaName = "learning_backlog_quiz"

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
            "minLength": 1
          },
          "options": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
            "uniqueItems": true,
            "items": {
              "type": "string",
              "minLength": 1
            }
          },
          "correctIndex": {
            "type": "integer",
            "minimum": 0,
            "maximum": 3
          },
          "explanation": {
            "type": "string",
            "minLength": 1
          },
          "evidence": {
            "type": "string",
            "minLength": 1
          }
        }
      }
    }
  }
}`
