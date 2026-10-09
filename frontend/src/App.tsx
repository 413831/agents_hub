import React, { useState, useEffect, useCallback } from 'react';
import { Sidebar } from './components/Sidebar';
import { GuiView } from './components/GuiView';
import { CliView } from './components/CliView';
import { Project, Session, Message, Connector, Agent, ViewMode, NavTab } from './types';
import { HubApi } from './services/api';

// Wails Runtime for events
declare global {
  interface Window {
    runtime?: {
      EventsOn: (eventName: string, callback: (data: any) => void) => () => void;
    };
  }
}

export const App: React.FC = () => {
  const [viewMode, setViewMode] = useState<ViewMode>('gui');
  const [currentTab, setCurrentTab] = useState<NavTab>('sessions');

  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);

  const [sessions, setSessions] = useState<Session[]>([]);
  const [selectedSession, setSelectedSession] = useState<Session | null>(null);

  const [messages, setMessages] = useState<Message[]>([]);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [connectors, setConnectors] = useState<Connector[]>([]);

  // Toggle GUI / CLI
  const handleToggleViewMode = useCallback(() => {
    setViewMode((prev) => (prev === 'gui' ? 'cli' : 'gui'));
  }, []);

  // Shortcut Ctrl+M or F2 to toggle dual-view
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey && e.key.toLowerCase() === 'm') || e.key === 'F2') {
        e.preventDefault();
        handleToggleViewMode();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleToggleViewMode]);

  // Initial Data Fetch
  useEffect(() => {
    const init = async () => {
      try {
        const [projList, agentList, connList] = await Promise.all([
          HubApi.listProjects(),
          HubApi.listAgents(),
          HubApi.listConnectors(),
        ]);
        setProjects(projList);
        setAgents(agentList);
        setConnectors(connList);

        if (projList.length > 0) {
          setSelectedProject(projList[0]);
        }
      } catch (err) {
        console.error("Error initializing Hub data:", err);
      }
    };
    init();
  }, []);

  // Fetch Sessions when selectedProject changes
  useEffect(() => {
    if (!selectedProject) return;
    const fetchSessions = async () => {
      try {
        const sessList = await HubApi.listSessions(selectedProject.id);
        setSessions(sessList);
        if (sessList.length > 0) {
          setSelectedSession(sessList[0]);
        } else {
          setSelectedSession(null);
          setMessages([]);
        }
      } catch (err) {
        console.error("Error fetching sessions:", err);
      }
    };
    fetchSessions();
  }, [selectedProject?.id]);

  // Fetch Messages when selectedSession changes
  useEffect(() => {
    if (!selectedSession) {
      setMessages([]);
      return;
    }
    const fetchMessages = async () => {
      try {
        const msgList = await HubApi.getSessionMessages(selectedSession.id);
        setMessages(msgList);
      } catch (err) {
        console.error("Error fetching messages:", err);
      }
    };
    fetchMessages();
  }, [selectedSession?.id]);

  // Wails Event Listeners
  useEffect(() => {
    if (window.runtime?.EventsOn) {
      const unsubMsg = window.runtime.EventsOn('message:received', (newMsg: Message) => {
        setMessages((prev) => {
          if (prev.some((m) => m.id === newMsg.id)) return prev;
          return [...prev, newMsg];
        });
      });
      return () => {
        unsubMsg();
      };
    }
  }, []);

  // Handlers
  const handleSendMessage = async (content: string, agentId: string) => {
    if (!selectedSession) return;

    // Optimistic user message addition
    const tempUserMsg: Message = {
      id: `usr-tmp-${Date.now()}`,
      sessionId: selectedSession.id,
      role: 'user',
      sender: 'Developer',
      content,
      createdAt: new Date().toISOString(),
    };
    setMessages((prev) => [...prev, tempUserMsg]);

    try {
      const agentReply = await HubApi.sendMessage(
        selectedSession.id,
        agentId,
        'Developer',
        content
      );
      // Append agent reply
      setMessages((prev) => {
        if (prev.some((m) => m.id === agentReply.id)) return prev;
        return [...prev, agentReply];
      });
    } catch (err) {
      console.error("Error sending message:", err);
    }
  };

  const handleToggleConnector = async (id: string) => {
    try {
      const updated = await HubApi.toggleConnector(id);
      setConnectors((prev) =>
        prev.map((c) => (c.id === updated.id ? updated : c))
      );
    } catch (err) {
      console.error("Error toggling connector:", err);
    }
  };

  const handleCreateProject = async (name: string, description: string) => {
    try {
      const newProj = await HubApi.createProject(name, description);
      setProjects((prev) => [newProj, ...prev]);
      setSelectedProject(newProj);
    } catch (err) {
      console.error("Error creating project:", err);
    }
  };

  const handleCreateSession = async (title: string, agentIds: string[]) => {
    if (!selectedProject) return;
    try {
      const newSess = await HubApi.createSession(selectedProject.id, title, agentIds);
      setSessions((prev) => [newSess, ...prev]);
      setSelectedSession(newSess);
      setMessages([]);
      setCurrentTab('sessions');
    } catch (err) {
      console.error("Error creating session:", err);
    }
  };

  return (
    <div className="flex h-screen w-screen overflow-hidden bg-[#080c14] text-slate-100 font-sans">
      {/* Barra lateral de navegación principal */}
      <Sidebar
        currentTab={currentTab}
        onSelectTab={setCurrentTab}
        viewMode={viewMode}
        onToggleViewMode={handleToggleViewMode}
        projects={projects}
        selectedProject={selectedProject}
        onSelectProject={(p) => {
          setSelectedProject(p);
        }}
        sessions={sessions}
        selectedSession={selectedSession}
        onSelectSession={(s) => {
          setSelectedSession(s);
          setCurrentTab('sessions');
        }}
        onCreateSession={() => {
          const title = prompt("Título de la nueva sesión:")?.trim();
          if (title) handleCreateSession(title, []);
        }}
        onCreateProject={() => {
          const name = prompt("Nombre del nuevo workspace / proyecto:")?.trim();
          if (name) handleCreateProject(name, "Espacio de trabajo creado desde la interfaz");
        }}
      />

      {/* Viewport Dual-View (GUI o CLI) */}
      <main className="flex-1 flex flex-col h-screen overflow-hidden">
        {viewMode === 'gui' ? (
          <GuiView
            currentTab={currentTab}
            selectedProject={selectedProject}
            selectedSession={selectedSession}
            messages={messages}
            agents={agents}
            connectors={connectors}
            projects={projects}
            onSendMessage={handleSendMessage}
            onToggleConnector={handleToggleConnector}
            onCreateProject={handleCreateProject}
            onCreateSession={handleCreateSession}
            onSelectProject={setSelectedProject}
            onSelectSession={setSelectedSession}
            onSelectTab={setCurrentTab}
          />
        ) : (
          <CliView
            selectedProject={selectedProject}
            selectedSession={selectedSession}
            messages={messages}
            agents={agents}
            connectors={connectors}
            projects={projects}
            sessions={sessions}
            onSendMessage={handleSendMessage}
            onToggleConnector={handleToggleConnector}
            onCreateProject={handleCreateProject}
            onCreateSession={handleCreateSession}
            onSelectProject={setSelectedProject}
            onSelectSession={setSelectedSession}
            onToggleViewMode={handleToggleViewMode}
            onSelectTab={setCurrentTab}
          />
        )}
      </main>
    </div>
  );
};

export default App;
