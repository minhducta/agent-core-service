# Business Requirements Document (BRD) - Agent Core Service

## 1. Executive Summary
Agent Core Service (bot-context-service/agent-control-plane) serves as the centralized brain for managing AI Agents in a multi-agent system. It transitions agent state from local file systems to a resilient, scalable, and queryable PostgreSQL database equipped with pgvector for semantic search.

## 2. Business Objectives
- Centralize and manage Bot Identity, Skills, Memories, and Policies.
- Provide a resilient heartbeat and command/task distribution mechanism for active agents.
- Replace local `.md` configuration files with dynamic, RESTful + SSE (Server-Sent Events) push updates.
- Decouple transactional communication (chats/messages) from core agent context management.

## 3. Core Features
1. **Agent Identity Management**: Manage agent profiles (Name, Role, Vibe, Avatar).
2. **Skill & Policy Configuration**: Define what an agent can do (Skills) and what it is allowed/denied to do (Policies).
3. **Long-Term Memory**: Store agent memories as vector embeddings via `pgvector` for context recall.
4. **Task Orchestration (Todos)**: Manage task states, dependencies (`dependency_id`), and checklists.
5. **Heartbeat & Command Queue**: Monitor agent health and dispatch real-time commands (e.g., `reload_skills`, `new_todo`).

## 4. Architecture
- **Tech Stack**: Node.js/TypeScript, Fastify, Prisma, PostgreSQL 16 + pgvector.
- **Communication Pattern**: Hybrid (REST API for CRUD/Initial Bootstrap + SSE for real-time updates).
- **Authentication**: M2M API Key (Bearer Token) validated via Middleware to resolve `bot_id`.

## 5. Non-Functional Requirements
- **Resilience (Degraded Mode)**: Agents must operate using local bootstrap cache if the service goes down.
- **Scalability**: Capable of handling hundreds of concurrent bots pinging every 60s.
- **Security**: IP whitelisting + API keys per bot.
