# Next.js Frontend with Firebase Authentication

This is the frontend application built with Next.js and Firebase Authentication.

## Setup Firebase

1. Create a new Firebase project at [Firebase Console](https://console.firebase.google.com)
2. Enable Authentication in your Firebase project:
   - Go to Authentication > Get Started
   - Enable Google Authentication provider
   - Add any other authentication providers you want to use

3. Get your Firebase configuration:
   - Go to Project Settings (⚙️)
   - Scroll down to "Your apps" section
   - Click the web icon (</>)
   - Register your app with a nickname
   - Copy the configuration object

4. Set up environment variables:
   Create or update your `.env` file with your Firebase configuration:

   ```env
   NEXT_PUBLIC_FIREBASE_API_KEY=your_api_key
   NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your_auth_domain
   NEXT_PUBLIC_FIREBASE_PROJECT_ID=your_project_id
   NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=your_storage_bucket
   NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=your_messaging_sender_id
   NEXT_PUBLIC_FIREBASE_APP_ID=your_app_id
   ```

## Getting Started

First, install the dependencies:

```bash
bun install
```

Then, run the development server:

```bash
bun run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Authentication Features

- Sign in with Google and Email/Password
- Protected routes under /(protected)/*
- Authentication status management
- Automatic redirection for unauthenticated users
- Password reset functionality
- Email verification
- Two-factor authentication with login verification
- Logout handling

## Project Structure

```
src/
├── app/                      # Next.js app router pages
│   ├── (protected)/         # Protected routes
│   ├── login/              # Authentication pages
│   ├── register/
│   ├── forgot-password/
│   ├── reset-password/
│   ├── verify-email/
│   └── verify-login/       # 2FA verification
├── components/
│   ├── auth/              # Authentication components
│   ├── layout/            # Layout components
│   ├── providers/         # Context providers
│   ├── rbac/             # Role-based access control
│   ├── tenants/          # Tenant management
│   └── ui/               # Shared UI components
├── contexts/
│   └── AuthContext.tsx    # Authentication context
├── hooks/                 # Custom React hooks
├── lib/
│   ├── apiClient.ts      # API client configuration
│   ├── firebase.ts       # Firebase configuration
│   └── utils.ts          # Utility functions
├── services/             # API services
└── types/                # TypeScript type definitions
```

## Features

- Role-Based Access Control (RBAC)
- Multi-tenant support
- System admin functionality
- Staff management
- User profile management
- Breadcrumb navigation
- Responsive UI components with shadcn/ui

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Set up your Firebase configuration in `.env`
4. Install dependencies:
   ```bash
   bun install
   ```
5. Run the development server:
   ```bash
   bun run dev
   ```
6. Make your changes
7. Test thoroughly
8. Commit your changes (`git commit -m 'Add some amazing feature'`)
9. Push to the branch (`git push origin feature/amazing-feature`)
10. Open a Pull Request

## Development Guidelines

- Keep components small and focused
- Follow the established project structure
- Use TypeScript for type safety
- Add appropriate documentation for new features
- Ensure code passes linting and type checking
- Test your changes across different browsers
