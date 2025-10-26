Grubzo

Grubzo is a modern web application for cafeteria and canteen food ordering. Users can browse menu items, place orders, and employees can manage and deliver those orders efficiently.

Features

- User registration & authentication (Email/Password + OAuth with Google/GitHub)
- Tenant management for cafeteria owners
- Employee management under tenant accounts
- Menu management (create, update, list items)
- Order creation and tracking
- File upload & download

Configuration
All application configuration (database, OAuth credentials, JWT secret, etc.) is handled via:

internal/config/config.go

Please go through config.go to set up your environment variables or default values before running the application.

Installation

1. Clone the repository:
   git clone https://github.com/rahulreddy-001/grubzo.git
   cd grubzo
2. Install backend dependencies:
   go mod tidy
3. Navigate to frontend folder and install dependencies:
   cd frontend
   npm install

Development Run
Start both backend and frontend in development mode:

make watch

Note: make watch will handle running the Go server and starting the frontend automatically.
