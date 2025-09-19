# New Topic and Stance Features

## Overview
The debate bot now supports:
1. **Any topic** - Users can provide their own debate topics
2. **Dynamic stance assignment** - Bot takes the opposite stance of the user
3. **Safety validation** - Topics are filtered for inappropriate content

## API Usage

### Request Format
```json
{
  "conversation_id": "optional-conversation-id",
  "message": "I think climate change is a hoax",
  "topic": "Climate change is real"
}
```

### Response
The bot will:
- Validate the topic for safety
- Determine the user's stance from their message
- Take the opposite stance (CON in this case)
- Respond with arguments against climate change being a hoax

## Safety Features

### Blocked Content
The system automatically blocks topics containing:
- Violence-related keywords
- Sexual content
- Hate speech
- Drug-related content
- Illegal activities

### Fallback Behavior
If a topic is deemed inappropriate, the system will:
- Use a predefined safe topic instead
- Assign a default stance
- Continue the conversation normally

## Examples

### Valid Topics
- "Remote work is better than office work"
- "Electric cars are the future"
- "Social media has negative effects on society"

### Invalid Topics (will be replaced with fallback)
- "Violence is good" → Replaced with safe topic
- "Drugs should be legal" → Replaced with safe topic
- "Racism is acceptable" → Replaced with safe topic

### Stance Detection
- "I agree that..." → Bot takes CON stance
- "I disagree with..." → Bot takes PRO stance
- "This is good/bad" → Bot takes opposite stance
- Unclear messages → Bot defaults to CON stance
