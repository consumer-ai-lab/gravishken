Thank you for the clarification. I'll update the README to reflect these four specific test types. Here's a revised version of the README:

# WCL Test Application

## Overview

This application has been developed for World Computer Limited (WCL) to conduct annual tests in our college. The system supports four types of tests:

1. Typing Test
2. PowerPoint Test
3. Excel Test
4. Word Test

The application is designed to provide a seamless testing experience with robust features to handle technical issues and ensure data persistence.

## Features

- **Multiple Test Types**: Conduct four different types of tests through a single platform.
- **Real-time Data Persistence**: Continuous saving of user progress to prevent data loss.
- **Interruption Handling**: Ability to resume tests from where the user left off in case of technical issues.
- **Accurate Time Tracking**: Server-side timer to ensure precise test duration.
- **Application Integration**: Tests for PowerPoint, Excel, and Word to assess proficiency in these tools.

## Test Modules

1. **Typing Test**: Measures typing speed and accuracy.
2. **PowerPoint Test**: Assesses skills in creating and manipulating presentations.
3. **Excel Test**: Evaluates proficiency in spreadsheet operations and formulas.
4. **Word Test**: Gauges competency in document creation and formatting.

## Technical Stack

- **Backend**: Go (Golang)
- **Frontend**: [Your frontend technology, e.g., React, Vue, etc.]
- **Real-time Communication**: WebSockets
- **Data Storage**: In-memory data store with persistence capabilities

## Key Components

1. **WebSocket Server**: Handles real-time communication between client and server.
2. **DataStore**: Manages user data and test progress.
3. **Test Modules**: Separate modules for each test type.
4. **User Interface**: Responsive design for seamless test-taking experience.

## Installation

[Provide installation steps here]

## Usage

[Provide usage instructions here]

## Development

To run the application in development mode:

```bash
BUILD_MODE=DEV go run main.go
```

For production:

```bash
BUILD_MODE=PROD go run main.go
```

## Contributing

[Provide guidelines for contributing to the project]

## License

[Specify the license under which this project is released]

## Acknowledgements

- World Computer Limited (WCL) for the opportunity to develop this application
- [Any other acknowledgements]

## Contact

For any queries regarding this application, please contact:

[Provide contact information]

---

We hope this application serves WCL and our college well in conducting efficient and reliable annual tests across typing, PowerPoint, Excel, and Word skills. Your feedback and suggestions for improvement are always welcome!