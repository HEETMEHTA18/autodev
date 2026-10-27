package cmd

// ── Next.js Boilerplate Templates ──────────────────────────────────────────

const nextJsPackageJson = `{
  "name": "autodev-nextjs-app",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "next": "^14.2.4",
    "lucide-react": "^0.395.0"
  },
  "devDependencies": {
    "typescript": "^5.4.5",
    "@types/node": "^20.14.2",
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "postcss": "^8.4.38",
    "autoprefixer": "^10.4.19",
    "tailwindcss": "^3.4.4",
    "eslint": "^8.57.0",
    "eslint-config-next": "14.2.4",
    "prettier": "^3.3.1"
  }
}`

const nextJsTsConfig = `{
  "compilerOptions": {
    "target": "es5",
    "lib": ["dom", "dom.iterable", "esnext"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}`

const nextJsConfig = `/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone', // Optimized for Docker containerization
};

module.exports = nextConfig;
`

const nextJsTailwindConfig = `/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "var(--background)",
        foreground: "var(--foreground)",
      },
    },
  },
  plugins: [],
};
`

const nextJsPostcssConfig = `module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`

const nextJsLayout = `import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AutoDev NextJS Starter",
  description: "Production ready Next.js starter created by AutoDev CLI",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="bg-slate-955 text-slate-100 antialiased">
        {children}
      </body>
    </html>
  );
}
`

const nextJsPage = `import { Terminal, Shield, Cpu, ExternalLink } from "lucide-react";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center justify-center p-6 bg-gradient-to-b from-slate-950 via-slate-900 to-indigo-950">
      <div className="max-w-4xl w-full text-center">
        <div className="inline-flex items-center gap-2 px-3 py-1 border border-amber-500/30 bg-amber-500/10 text-amber-400 text-xs font-bold uppercase tracking-wider rounded-full mb-8">
          <Terminal className="w-3.5 h-3.5" /> Next.js Production Boilerplate
        </div>

        <h1 className="text-5xl md:text-7xl font-extrabold text-transparent bg-clip-text bg-gradient-to-r from-amber-400 via-emerald-400 to-cyan-400 tracking-tight mb-6">
          Ready for Production.
        </h1>

        <p className="text-lg md:text-xl text-slate-400 max-w-2xl mx-auto mb-12">
          This Next.js application has been bootstrapped with TS, Tailwind, multi-stage Docker builds, and CI/CD GitHub workflows.
        </p>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-left mb-16">
          <div className="bg-slate-900/50 backdrop-blur-md border border-slate-805 rounded-xl p-6 hover:border-slate-700 transition-all">
            <Shield className="w-8 h-8 text-emerald-400 mb-4" />
            <h3 className="text-lg font-bold text-white mb-2">Dockerized</h3>
            <p className="text-sm text-slate-400">Multi-stage standalone Docker configuration ready to deploy to K8s or ECS.</p>
          </div>
          <div className="bg-slate-900/50 backdrop-blur-md border border-slate-805 rounded-xl p-6 hover:border-slate-700 transition-all">
            <Cpu className="w-8 h-8 text-amber-400 mb-4" />
            <h3 className="text-lg font-bold text-white mb-2">GitHub Actions</h3>
            <p className="text-sm text-slate-400">Pre-configured testing, building, and security analysis workflow pipelines.</p>
          </div>
          <div className="bg-slate-900/50 backdrop-blur-md border border-slate-805 rounded-xl p-6 hover:border-slate-700 transition-all">
            <Terminal className="w-8 h-8 text-cyan-400 mb-4" />
            <h3 className="text-lg font-bold text-white mb-2">SEO Optimized</h3>
            <p className="text-sm text-slate-400">Layout architecture with built-in metadata support and Google-friendly tags.</p>
          </div>
        </div>

        <div className="flex flex-wrap justify-center gap-4">
          <a
            href="https://autodevs.dev"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-1.5 bg-gradient-to-r from-amber-505 to-emerald-500 text-slate-950 font-bold px-6 py-3 rounded-lg hover:opacity-90 active:scale-95 transition-all shadow-lg"
          >
            Explore AutoDev <ExternalLink className="w-4 h-4" />
          </a>
        </div>
      </div>
    </main>
  );
}
`

const nextJsGlobalsCss = `@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --background: #020617;
  --foreground: #f8fafc;
}
`

const nextJsDockerfile = `# Multi-stage Build for standalone Next.js deployment
FROM node:18-alpine AS base

# Install dependencies only when needed
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# Install dependencies based on the preferred package manager
COPY package.json ./
RUN npm install

# Rebuild the source code only when needed
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# Disable telemetry during build
ENV NEXT_TELEMETRY_DISABLED 1
RUN npm run build

# Production image, copy all the files and run next
FROM base AS runner
WORKDIR /app

ENV NODE_ENV production
ENV NEXT_TELEMETRY_DISABLED 1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public

# Set the correct permission for prerender cache
RUN mkdir .next
RUN chown nextjs:nodejs .next

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

# server.js is created by next build from the standalone output
# https://nextjs.org/docs/pages/api-reference/next-config-js/output
CMD ["node", "server.js"]
`

const nextJsGithubAction = `name: NextJS CI/CD

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v4

    - name: Use Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 18.x
        cache: 'npm'

    - name: Install dependencies
      run: npm ci

    - name: Lint project
      run: npm run lint

    - name: Build standalone app
      run: npm run build
`

const nextJsReadme = `# AutoDev NextJS Production Starter

This project is bootstrapped with **Next.js** + **TypeScript** + **Tailwind CSS** + **Docker** + **GitHub Actions CI/CD** via AutoDev.

## Getting Started

First, run the development server:

__BT____BT____BT__bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
__BT____BT____BT__

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Docker Deployment

To build and run the Docker container locally:

__BT____BT____BT__bash
docker build -t nextjs-autodev-app .
docker run -p 3000:3000 nextjs-autodev-app
__BT____BT____BT__

## Features Added
* **TypeScript** for strict type checking
* **Tailwind CSS** for layout utility styles
* **ESLint** & **Prettier** pre-configured for code quality
* **Standalone Docker build** mapping files for high performance and low image size
* **GitHub Actions pipeline** for automated testing and builds on push/pull requests
`

// ── AI Chatbot (Gemini) Templates ──────────────────────────────────────────

const aiChatbotPackageJson = `{
  "name": "autodev-ai-chatbot",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev:server": "node server.js",
    "dev:client": "vite",
    "dev": "concurrently \"npm run dev:server\" \"npm run dev:client\"",
    "build": "vite build",
    "start": "node server.js"
  },
  "dependencies": {
    "@google/genai": "^0.2.0",
    "cors": "^2.8.5",
    "dotenv": "^16.4.5",
    "express": "^4.19.2",
    "lucide-react": "^0.395.0",
    "react": "^18.3.1",
    "react-dom": "^18.3.1"
  },
  "devDependencies": {
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.1",
    "autoprefixer": "^10.4.19",
    "concurrently": "^8.2.2",
    "postcss": "^8.4.39",
    "tailwindcss": "^3.4.4",
    "typescript": "^5.2.2",
    "vite": "^5.3.1"
  }
}`

const aiChatbotServer = `import express from 'express';
import cors from 'cors';
import dotenv from 'dotenv';
import { GoogleGenAI } from '@google/genai';
import path from 'path';
import { fileURLToPath } from 'url';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 5000;

app.use(cors());
app.use(express.json());

// Initialize Gemini API client
const apiKey = process.env.GEMINI_API_KEY;
let ai = null;

if (apiKey) {
  ai = new GoogleGenAI({ apiKey });
  console.log("Gemini API Client successfully initialized.");
} else {
  console.warn("! Warning: GEMINI_API_KEY environment variable is missing.");
}

// API Chat Endpoint
app.post('/api/chat', async (req, res) => {
  const { message } = req.body;

  if (!message) {
    return res.status(400).json({ error: "Message is required" });
  }

  if (!ai) {
    return res.json({ 
      reply: "Gemini API key is missing. Please set your GEMINI_API_KEY environment variable to start chatting! You can generate keys in Google AI Studio." 
    });
  }

  try {
    // Generate text content using Gemini
    const response = await ai.models.generateContent({
      model: 'gemini-2.5-flash',
      contents: message,
    });

    const reply = response.text || "Sorry, I couldn't generate a response.";
    res.json({ reply });
  } catch (error) {
    console.error("Gemini Generation Error:", error);
    res.status(500).json({ error: "Failed to communicate with AI model", details: error.message });
  }
});

// Serve static assets in production
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

if (process.env.NODE_ENV === 'production') {
  app.use(express.static(path.join(__dirname, 'dist')));
  app.get('*', (req, res) => {
    res.sendFile(path.resolve(__dirname, 'dist', 'index.html'));
  });
}

app.listen(PORT, () => {
  console.log(__BT__Server running in ${process.env.NODE_ENV || 'development'} mode on port ${PORT}__BT__);
});
`

const aiChatbotIndexHtml = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>AutoDev AI Chatbot</title>
  </head>
  <body class="bg-slate-950 text-white">
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>`

const aiChatbotApp = `import { useState, useRef, useEffect } from 'react';
import { Send, Bot, User, Sparkles, Key } from 'lucide-react';

interface Message {
  id: string;
  sender: 'ai' | 'user';
  text: string;
}

function App() {
  const [messages, setMessages] = useState<Message[]>([
    { id: '1', sender: 'ai', text: 'Hello! I am a Gemini-powered AI chatbot initialized by AutoDev. How can I help you build today?' }
  ]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages, loading]);

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || loading) return;

    const userMessage: Message = {
      id: Date.now().toString(),
      sender: 'user',
      text: input.trim()
    };

    setMessages((prev) => [...prev, userMessage]);
    setInput('');
    setLoading(true);

    try {
      const response = await fetch('/api/chat', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: userMessage.text })
      });
      const data = await response.json();
      
      setMessages((prev) => [...prev, {
        id: (Date.now() + 1).toString(),
        sender: 'ai',
        text: data.reply || data.error || 'Server error occurred'
      }]);
    } catch (err) {
      setMessages((prev) => [...prev, {
        id: (Date.now() + 1).toString(),
        sender: 'ai',
        text: 'Failed to connect to local server backend. Make sure the Node server is running on port 5000.'
      }]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-slate-955 text-slate-100 font-sans">
      {/* Header */}
      <header className="border-b border-slate-800 bg-slate-900/60 backdrop-blur-md px-6 py-4 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-amber-400 to-emerald-400 flex items-center justify-center shadow-lg shadow-emerald-500/10">
            <Bot className="w-5 h-5 text-slate-950 font-black" />
          </div>
          <div>
            <h1 className="font-extrabold text-white tracking-tight flex items-center gap-1.5 text-lg">
              AutoDev AI Chatbot <Sparkles className="w-4 h-4 text-amber-400 fill-amber-400" />
            </h1>
            <p className="text-xs text-slate-400">Powered by Gemini 2.5 Flash</p>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <span className="hidden md:inline-flex items-center gap-1.5 px-3 py-1 border border-slate-700 bg-slate-800/50 rounded-full text-xs font-mono text-slate-400">
            <Key className="w-3.5 h-3.5" /> Set GEMINI_API_KEY
          </span>
        </div>
      </header>

      {/* Main chat window */}
      <div className="flex-1 overflow-y-auto max-w-4xl w-full mx-auto p-6 space-y-6">
        {messages.map((m) => (
          <div
            key={m.id}
            className={__BT__flex gap-4 max-w-[85%] ${m.sender === 'user' ? 'ml-auto flex-row-reverse' : ''}__BT__}
          >
            <div className={__BT__w-8 h-8 rounded-lg flex items-center justify-center shrink-0 shadow-md ${
              m.sender === 'user' 
                ? 'bg-emerald-500 text-slate-950' 
                : 'bg-slate-800 text-amber-400 border border-slate-700'
            }__BT__}>
              {m.sender === 'user' ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
            </div>
            
            <div className={__BT__rounded-2xl px-5 py-3.5 text-sm leading-relaxed border shadow-xl ${
              m.sender === 'user'
                ? 'bg-emerald-600/10 border-emerald-500/30 text-white'
                : 'bg-slate-900 border-slate-805 text-slate-200'
            }__BT__}>
              <p className="whitespace-pre-wrap">{m.text}</p>
            </div>
          </div>
        ))}

        {loading && (
          <div className="flex gap-4 max-w-[80%]">
            <div className="w-8 h-8 rounded-lg bg-slate-800 border border-slate-700 text-amber-400 flex items-center justify-center shrink-0">
              <Bot className="w-4 h-4" />
            </div>
            <div className="bg-slate-900 border border-slate-805 rounded-2xl px-5 py-3.5 text-sm flex items-center gap-1">
              <span className="w-1.5 h-1.5 bg-amber-400 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
              <span className="w-1.5 h-1.5 bg-amber-400 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
              <span className="w-1.5 h-1.5 bg-amber-400 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input section */}
      <footer className="border-t border-slate-900 bg-slate-955/80 backdrop-blur-md px-6 py-4 sticky bottom-0 z-10">
        <form onSubmit={handleSend} className="max-w-4xl w-full mx-auto flex gap-3">
          <input
            type="text"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Type your prompt here..."
            className="flex-1 bg-slate-900 border border-slate-805 rounded-xl px-5 py-3 text-sm focus:outline-none focus:border-amber-500 focus:ring-1 focus:ring-amber-500 transition-all text-white placeholder-slate-500"
          />
          <button
            type="submit"
            disabled={!input.trim() || loading}
            className="px-5 py-3 bg-gradient-to-r from-amber-500 to-emerald-500 text-slate-950 rounded-xl font-bold hover:opacity-90 active:scale-95 disabled:opacity-40 disabled:pointer-events-none transition-all flex items-center justify-center"
          >
            <Send className="w-4 h-4" />
          </button>
        </form>
      </footer>
    </div>
  );
}

export default App;
`

const aiChatbotDockerfile = `FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:18-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY package.json ./
RUN npm install --only=production
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/server.js ./server.js

EXPOSE 5000
CMD ["node", "server.js"]
`

const aiChatbotReadme = `# AutoDev AI Chatbot (Gemini) Boilerplate

This template sets up a beautiful chatbot frontend (Vite + React) linked to a fast Node.js backend integrating Google's **Gemini 2.5 Flash** model via the official __BT__@google/genai__BT__ library.

## · Configuration

Generate a free Gemini API key in [Google AI Studio](https://aistudio.google.com/).

Create a __BT__.env__BT__ file in the root directory:

__BT____BT____BT__env
GEMINI_API_KEY=your_gemini_api_key_here
PORT=5000
__BT____BT____BT__

## Running the App

Run both the server and the frontend client concurrently:

__BT____BT____BT__bash
npm install
npm run dev
__BT____BT____BT__

Open [http://localhost:5173](http://localhost:5173) to start testing.

## Running inside Docker

Build and spin up the Docker container:

__BT____BT____BT__bash
docker build -t autodev-gemini-bot .
docker run -p 5000:5000 -e GEMINI_API_KEY="your_api_key" autodev-gemini-bot
__BT____BT____BT__
`

const aiChatbotViteConfig = `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:5000',
        changeOrigin: true,
      },
    },
  },
})`

// ── MERN Stack Templates ──────────────────────────────────────────

const mernDockerCompose = `version: '3.8'

services:
  database:
    image: mongo:6.0
    container_name: autodev-mern-db
    ports:
      - "27017:27017"
    volumes:
      - mongo-data:/data/db

  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    container_name: autodev-mern-server
    ports:
      - "5000:5000"
    environment:
      - PORT=5000
      - MONGO_URI=mongodb://database:27017/autodev-mern
    depends_on:
      - database

  client:
    build:
      context: ./client
      dockerfile: Dockerfile
    container_name: autodev-mern-client
    ports:
      - "5173:80"
    depends_on:
      - server

volumes:
  mongo-data:
`

const mernServerPackageJson = `{
  "name": "autodev-mern-server",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "start": "node server.js",
    "dev": "nodemon server.js"
  },
  "dependencies": {
    "cors": "^2.8.5",
    "dotenv": "^16.4.5",
    "express": "^4.19.2",
    "mongoose": "^8.4.1"
  },
  "devDependencies": {
    "nodemon": "^3.1.2"
  }
}`

const mernServerJs = `import express from 'express';
import mongoose from 'mongoose';
import cors from 'cors';
import dotenv from 'dotenv';

dotenv.config();

const app = express();
const PORT = process.env.PORT || 5000;
const MONGO_URI = process.env.MONGO_URI || 'mongodb://localhost:27017/mern-app';

app.use(cors());
app.use(express.json());

// Connection to MongoDB
mongoose.connect(MONGO_URI)
  .then(() => console.log('MongoDB successfully connected.'))
  .catch(err => console.error('✗ Failed to connect to MongoDB:', err));

// Test Schema
const ItemSchema = new mongoose.Schema({
  name: { type: String, required: true },
  createdAt: { type: Date, default: Date.now }
});
const Item = mongoose.model('Item', ItemSchema);

// API Endpoints
app.get('/api/items', async (req, res) => {
  try {
    const items = await Item.find().sort({ createdAt: -1 });
    res.json(items);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
});

app.post('/api/items', async (req, res) => {
  try {
    const newItem = new Item({ name: req.body.name });
    const savedItem = await newItem.save();
    res.status(201).json(savedItem);
  } catch (err) {
    res.status(400).json({ error: err.message });
  }
});

app.listen(PORT, () => {
  console.log(__BT__Express server running on port ${PORT}__BT__);
});
`

const mernServerDockerfile = `FROM node:18-alpine
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
EXPOSE 5000
CMD ["npm", "start"]
`

const mernClientPackageJson = `{
  "name": "autodev-mern-client",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "lucide-react": "^0.395.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.1",
    "autoprefixer": "^10.4.19",
    "postcss": "^8.4.39",
    "tailwindcss": "^3.4.4",
    "typescript": "^5.2.2",
    "vite": "^5.3.1"
  }
}`

const mernClientApp = `import { useState, useEffect } from 'react';
import { Database, Plus, Database as DbIcon } from 'lucide-react';

interface DBItem {
  _id: string;
  name: string;
}

function App() {
  const [items, setItems] = useState<DBItem[]>([]);
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const res = await fetch('http://localhost:5000/api/items');
      const data = await res.json();
      setItems(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      const res = await fetch('http://localhost:5000/api/items', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name })
      });
      const newItem = await res.json();
      setItems((prev) => [newItem, ...prev]);
      setName('');
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div className="min-h-screen bg-slate-955 text-slate-100 flex flex-col justify-between">
      <header className="border-b border-slate-900 bg-slate-900/40 px-8 py-5 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Database className="w-6 h-6 text-emerald-400" />
          <span className="font-extrabold tracking-tight text-white">AUTODEV MERN STACK</span>
        </div>
      </header>

      <main className="flex-1 max-w-2xl w-full mx-auto p-6">
        <div className="bg-slate-900/50 border border-slate-800 rounded-2xl p-8 mb-8 text-center shadow-xl">
          <h2 className="text-3xl font-black text-white mb-2">MongoDB + Express + React + Node</h2>
          <p className="text-sm text-slate-400">Full stack containerized development environment ready for code.</p>
        </div>

        <form onSubmit={handleAdd} className="flex gap-2 mb-8">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Add dynamic DB item..."
            className="flex-1 bg-slate-900 border border-slate-805 rounded-xl px-5 py-3 text-sm focus:outline-none focus:border-emerald-400 text-white"
          />
          <button type="submit" className="bg-emerald-500 hover:bg-emerald-600 active:scale-95 text-slate-950 font-bold px-6 py-3 rounded-xl flex items-center justify-center transition-all">
            <Plus className="w-5 h-5" />
          </button>
        </form>

        <div className="space-y-3">
          <h3 className="font-bold text-slate-400 text-sm tracking-wider uppercase">Database Items</h3>
          {loading ? (
            <p className="text-sm text-slate-500">Loading db items...</p>
          ) : items.length === 0 ? (
            <p className="text-sm text-slate-500">No items in DB. Try adding some above!</p>
          ) : (
            items.map((item) => (
              <div key={item._id} className="bg-slate-900 border border-slate-805 px-5 py-4 rounded-xl flex justify-between items-center">
                <span className="text-sm text-slate-200 font-medium">{item.name}</span>
                <span className="text-xs text-slate-500 font-mono">{item._id}</span>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}

export default App;
`

const mernClientDockerfile = `FROM node:18-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:stable-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`

const mernReadme = `# AutoDev MERN Stack (Containerized)

This template boots a full-stack container environment containing a **MongoDB** database, an **Express Node.js** API backend, and a **Vite + React** client frontend.

## Launching the Environment

To compile, link and start all the containers, use Docker Compose in one command:

__BT____BT____BT__bash
docker-compose up --build
__BT____BT____BT__

* Frontend Dashboard: [http://localhost:5173](http://localhost:5173)
* Express API Server: [http://localhost:5000](http://localhost:5000)
* MongoDB endpoint: mongodb://localhost:27017

## Directory Structure
* **/client**: Vite + React + TS dashboard client
* **/server**: Express database API server endpoint
* **/docker-compose.yml**: Container deployment links orchestrator
`

// ── Flutter templates ──────────────────────────────────────────────

const flutterPubspec = `name: autodev_flutter_app
description: A clean Flutter UI created by AutoDev.
version: 1.0.0+1
publish_to: 'none'

environment:
  sdk: '>=3.0.0 <4.0.0'

dependencies:
  flutter:
    sdk: flutter
  cupertino_icons: ^1.0.5

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^2.0.0

flutter:
  uses-material-design: true
`

const flutterMainDart = `import 'package:flutter/material';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'AutoDev Flutter App',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.amber, brightness: Brightness.dark),
        useMaterial3: true,
      ),
      home: const HomeScreen(),
      debugShowCheckedModeBanner: false,
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _counter = 0;

  void _increment() {
    setState(() {
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('AutoDev Flutter'),
        centerTitle: true,
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const Text(
              'Flutter Boilerplate successfully created!',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 12),
            const Text(
              'Interactions Counter:',
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                    color: Colors.amber,
                    fontWeight: FontWeight.w900,
                  ),
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _increment,
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ),
    );
  }
}
`

const flutterDockerfile = `# Build Flutter Web App and Serve with Nginx
FROM ubuntu:20.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y curl git unzip xz-utils zip libglu1-mesa wget

WORKDIR /usr/local
RUN git clone https://github.com/flutter/flutter.git -b stable

ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:\${PATH}"
RUN flutter doctor -v

WORKDIR /app
COPY . .
RUN flutter build web

FROM nginx:alpine
COPY --from=builder /app/build/web /usr/share/nginx/html
EXPOSE 80
`

const flutterGithubAction = `name: Flutter Web Build CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Java JDK
      uses: actions/setup-java@v3
      with:
        distribution: 'zulu'
        java-version: '11'
        
    - name: Setup Flutter
      uses: subosito/flutter-action@v2
      with:
        flutter-version: '3.x'
        channel: 'stable'
        
    - name: Install dependencies
      run: flutter pub get
      
    - name: Build Web Bundle
      run: flutter build web
`

const flutterReadme = `# AutoDev Flutter Starter

This is a Flutter mobile and web environment starter generated automatically by AutoDev CLI, incorporating clean project architecture structures, a Docker build workflow, and a pre-set GitHub Actions pipeline.

## Commands

First, retrieve the package dependencies:

__BT____BT____BT__bash
flutter pub get
__BT____BT____BT__

Run the application:

__BT____BT____BT__bash
flutter run
__BT____BT____BT__

## Web Container Deployment

To build and compile the Flutter web application using multi-stage Docker and serve it via Nginx:

__BT____BT____BT__bash
docker build -t flutter-web-autodev .
docker run -p 8080:80 flutter-web-autodev
__BT____BT____BT__
`

const mernRootPackageJson = `{
  "name": "autodev-mern-stack",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev:server": "cd server && npm run dev",
    "dev:client": "cd client && npm run dev",
    "dev": "concurrently \"npm run dev:server\" \"npm run dev:client\"",
    "build": "cd client && npm run build",
    "start": "cd server && npm start"
  },
  "devDependencies": {
    "concurrently": "^8.2.2"
  }
}`

