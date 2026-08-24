package llm

const quizSchemaJSON = `{
  "type": "object",
  "additionalProperties": false,
  "required": [
    "title",
    "topic",
    "difficulty",
    "questions"
  ],
  "properties": {
    "title": {
      "type": "string",
      "minLength": 1
    },
    "topic": {
      "type": "string",
      "minLength": 1
    },
    "difficulty": {
      "type": "string",
      "enum": [
        "easy",
        "medium",
        "hard"
      ]
    },
    "questions": {
      "type": "array",
      "minItems": 5,
      "maxItems": 5,
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": [
          "question",
          "options",
          "correctIndex",
          "explanation",
          "evidence"
        ],
        "properties": {
          "question": {
            "type": "string",
            "minLength": 1
          },
          "options": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
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
