export interface Metric {
  id: string;
  interactionId?: string;
  provider: string;
  model: string;
  totalLatencyMs: number;
  ttftMs: number;
  promptTokens: number;
  completionTokens: number;
  totalTokens: number;
  finishReason: string;
  timestamp: string;
}

export interface Message {
  id: string;
  sessionId: string;
  agentId?: string;
  role: 'user' | 'assistant' | 'system' | 'tool';
  sender: string;
  content: string;
  telemetry?: Metric;
  createdAt: string;
}

export interface Session {
  id: string;
  projectId: string;
  title: string;
  agentIds: string[];
  status: 'active' | 'archived';
  createdAt: string;
  updatedAt: string;
}

export interface Project {
  id: string;
  name: string;
  description: string;
  defaultModel?: string;
  metadata?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export interface Agent {
  id: string;
  name: string;
  role: string;
  description: string;
  systemPrompt: string;
  provider: string;
  model: string;
  temperature: number;
  maxTokens: number;
  tools: string[];
  avatar?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Connector {
  id: string;
  name: string;
  type: 'mcp' | 'skill' | 'plugin';
  description: string;
  status: 'connected' | 'disconnected' | 'disabled';
  endpoint?: string;
  tools: string[];
  config?: Record<string, string>;
  createdAt: string;
  updatedAt: string;
}

export type ViewMode = 'gui' | 'cli';
export type NavTab = 'sessions' | 'projects' | 'connectors';
