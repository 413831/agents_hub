import { Project, Session, Message, Connector, Agent, Metric } from '../types';

// Import Wails generated bindings
import * as HubHandler from '../../wailsjs/go/wails/HubHandler';

export const HubApi = {
  async listProjects(): Promise<Project[]> {
    try {
      const res = await HubHandler.ListProjects();
      return (res || []) as unknown as Project[];
    } catch (err) {
      console.warn("Wails IPC listProjects failed, using fallback:", err);
      return [
        {
          id: 'proj-default-1',
          name: 'AI Agent Laboratory',
          description: 'Experimentación multi-proveedor con OpenAI, Gemini y Llama local.',
          defaultModel: 'gemini-1.5-flash',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }
      ];
    }
  },

  async createProject(name: string, description: string): Promise<Project> {
    try {
      const res = await HubHandler.CreateProject(name, description);
      return res as unknown as Project;
    } catch (err) {
      console.error("CreateProject error:", err);
      throw err;
    }
  },

  async listSessions(projectId: string): Promise<Session[]> {
    try {
      const res = await HubHandler.ListSessions(projectId);
      return (res || []) as unknown as Session[];
    } catch (err) {
      console.warn("Wails IPC listSessions failed, using fallback:", err);
      return [
        {
          id: 'sess-1',
          projectId: projectId || 'proj-default-1',
          title: 'Arquitectura DDD y Clean Architecture',
          agentIds: ['agent-arch-1'],
          status: 'active',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        }
      ];
    }
  },

  async createSession(projectId: string, title: string, agentIds: string[] = []): Promise<Session> {
    try {
      const res = await HubHandler.CreateSession(projectId, title, agentIds);
      return res as unknown as Session;
    } catch (err) {
      console.error("CreateSession error:", err);
      throw err;
    }
  },

  async getSessionMessages(sessionId: string): Promise<Message[]> {
    try {
      const res = await HubHandler.GetSessionMessages(sessionId);
      return (res || []) as unknown as Message[];
    } catch (err) {
      console.warn("Wails IPC getSessionMessages failed:", err);
      return [];
    }
  },

  async sendMessage(sessionId: string, agentId: string, sender: string, content: string): Promise<Message> {
    try {
      const res = await HubHandler.SendMessage(sessionId, agentId, sender, content);
      return res as unknown as Message;
    } catch (err) {
      console.error("SendMessage error:", err);
      throw err;
    }
  },

  async listConnectors(): Promise<Connector[]> {
    try {
      const res = await HubHandler.ListConnectors();
      return (res || []) as unknown as Connector[];
    } catch (err) {
      console.warn("Wails IPC listConnectors failed:", err);
      return [];
    }
  },

  async toggleConnector(connectorId: string): Promise<Connector> {
    try {
      const res = await HubHandler.ToggleConnector(connectorId);
      return res as unknown as Connector;
    } catch (err) {
      console.error("ToggleConnector error:", err);
      throw err;
    }
  },

  async listAgents(): Promise<Agent[]> {
    try {
      const res = await HubHandler.ListAgents();
      return (res || []) as unknown as Agent[];
    } catch (err) {
      console.warn("Wails IPC listAgents failed:", err);
      return [];
    }
  },

  async getMetrics(limit: number = 20): Promise<Metric[]> {
    try {
      const res = await HubHandler.GetMetrics(limit);
      return (res || []) as unknown as Metric[];
    } catch (err) {
      console.warn("Wails IPC getMetrics failed:", err);
      return [];
    }
  }
};
