import React, { useState, useRef, useEffect } from 'react';
import { 
  Send, 
  Bot, 
  User, 
  Clock, 
  Zap, 
  CheckCircle2, 
  Activity, 
  FolderKanban, 
  Plus, 
  Layers, 
  ShieldCheck, 
  Power,
  SlidersHorizontal,
  Sparkles
} from 'lucide-react';
import { Project, Session, Message, Connector, Agent, NavTab } from '../types';

interface GuiViewProps {
  currentTab: NavTab;
  selectedProject: Project | null;
  selectedSession: Session | null;
  messages: Message[];
  agents: Agent[];
  connectors: Connector[];
  projects: Project[];
  onSendMessage: (content: string, agentId: string) => Promise<void>;
  onToggleConnector: (id: string) => Promise<void>;
  onCreateProject: (name: string, description: string) => Promise<void>;
  onCreateSession: (title: string, agentIds: string[]) => Promise<void>;
  onSelectProject: (proj: Project) => void;
  onSelectSession: (sess: Session) => void;
  onSelectTab: (tab: NavTab) => void;
}

export const GuiView: React.FC<GuiViewProps> = ({
  currentTab,
  selectedProject,
  selectedSession,
  messages,
  agents,
  connectors,
  projects,
  onSendMessage,
  onToggleConnector,
  onCreateProject,
  onCreateSession,
  onSelectProject,
  onSelectSession,
  onSelectTab,
}) => {
  const [inputText, setInputText] = useState('');
  const [selectedAgentId, setSelectedAgentId] = useState<string>('');
  const [isSending, setIsSending] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Modals / forms state
  const [showNewProjModal, setShowNewProjModal] = useState(false);
  const [newProjName, setNewProjName] = useState('');
  const [newProjDesc, setNewProjDesc] = useState('');

  const [showNewSessModal, setShowNewSessModal] = useState(false);
  const [newSessTitle, setNewSessTitle] = useState('');

  useEffect(() => {
    if (agents.length > 0 && !selectedAgentId) {
      setSelectedAgentId(agents[0].id);
    }
  }, [agents, selectedAgentId]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSubmitMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim() || isSending) return;
    const text = inputText;
    setInputText('');
    setIsSending(true);
    try {
      await onSendMessage(text, selectedAgentId);
    } finally {
      setIsSending(false);
    }
  };

  const handleCreateProjectSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newProjName.trim()) return;
    await onCreateProject(newProjName, newProjDesc);
    setNewProjName('');
    setNewProjDesc('');
    setShowNewProjModal(false);
  };

  const handleCreateSessionSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newSessTitle.trim()) return;
    await onCreateSession(newSessTitle, selectedAgentId ? [selectedAgentId] : []);
    setNewSessTitle('');
    setShowNewSessModal(false);
  };

  // ----------------------------------------------------
  // SUB-VIEW 1: PROYECTOS (Espacios de Trabajo)
  // ----------------------------------------------------
  if (currentTab === 'projects') {
    return (
      <div className="flex-1 bg-[#090d16] flex flex-col h-screen overflow-y-auto p-8">
        <div className="max-w-6xl w-full mx-auto flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-white tracking-tight flex items-center gap-2">
                <FolderKanban className="w-6 h-6 text-blue-500" />
                Espacios de Trabajo y Proyectos
              </h1>
              <p className="text-sm text-slate-400 mt-1">
                Aislamiento de contexto DDD para modelos, agentes asignados y logs de ejecución.
              </p>
            </div>
            <button
              onClick={() => setShowNewProjModal(true)}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-sm font-medium flex items-center gap-2 shadow-lg shadow-blue-500/20 transition-all"
            >
              <Plus className="w-4 h-4" />
              Nuevo Proyecto
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            {projects.map((proj) => {
              const isSelected = selectedProject?.id === proj.id;
              return (
                <div
                  key={proj.id}
                  className={`p-5 rounded-xl border transition-all flex flex-col justify-between ${
                    isSelected
                      ? 'bg-slate-900/90 border-blue-500/60 shadow-lg shadow-blue-500/10'
                      : 'bg-slate-900/40 border-slate-800 hover:border-slate-700 hover:bg-slate-900/70'
                  }`}
                >
                  <div className="flex flex-col gap-2">
                    <div className="flex items-center justify-between">
                      <span className="text-xs font-mono text-blue-400 bg-blue-500/10 px-2 py-0.5 rounded border border-blue-500/20">
                        {proj.id}
                      </span>
                      {isSelected && (
                        <span className="text-xs bg-emerald-500/15 text-emerald-400 px-2 py-0.5 rounded-full font-medium flex items-center gap-1 border border-emerald-500/30">
                          <CheckCircle2 className="w-3 h-3" /> Activo
                        </span>
                      )}
                    </div>
                    <h3 className="text-base font-semibold text-white mt-1">{proj.name}</h3>
                    <p className="text-xs text-slate-400 leading-relaxed line-clamp-3">
                      {proj.description || 'Sin descripción.'}
                    </p>
                  </div>

                  <div className="mt-6 pt-4 border-t border-slate-800/80 flex items-center justify-between text-xs">
                    <span className="text-slate-400 font-mono">
                      Modelo: <strong className="text-slate-200">{proj.defaultModel || 'Multi-LLM'}</strong>
                    </span>
                    <button
                      onClick={() => {
                        onSelectProject(proj);
                        onSelectTab('sessions');
                      }}
                      className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-blue-400 rounded-md font-medium transition-colors"
                    >
                      Abrir Sesiones &rarr;
                    </button>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Modal Crear Proyecto */}
        {showNewProjModal && (
          <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-[#101524] border border-slate-800 rounded-xl p-6 max-w-md w-full shadow-2xl">
              <h2 className="text-lg font-bold text-white mb-2">Crear Nuevo Espacio de Trabajo</h2>
              <form onSubmit={handleCreateProjectSubmit} className="flex flex-col gap-4">
                <div>
                  <label className="text-xs text-slate-400 block mb-1">Nombre del Proyecto</label>
                  <input
                    type="text"
                    required
                    value={newProjName}
                    onChange={(e) => setNewProjName(e.target.value)}
                    placeholder="ej: Microservicios Orchestrator"
                    className="w-full bg-slate-900 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
                <div>
                  <label className="text-xs text-slate-400 block mb-1">Descripción</label>
                  <textarea
                    rows={3}
                    value={newProjDesc}
                    onChange={(e) => setNewProjDesc(e.target.value)}
                    placeholder="Propósito, agentes involucrados y contexto..."
                    className="w-full bg-slate-900 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                  />
                </div>
                <div className="flex items-center justify-end gap-3 mt-2">
                  <button
                    type="button"
                    onClick={() => setShowNewProjModal(false)}
                    className="px-4 py-2 text-xs text-slate-400 hover:text-white"
                  >
                    Cancelar
                  </button>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-xs font-medium"
                  >
                    Crear Proyecto
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}
      </div>
    );
  }

  // ----------------------------------------------------
  // SUB-VIEW 2: CONECTORES (MCP, Skills, Plugins)
  // ----------------------------------------------------
  if (currentTab === 'connectors') {
    return (
      <div className="flex-1 bg-[#090d16] flex flex-col h-screen overflow-y-auto p-8">
        <div className="max-w-6xl w-full mx-auto flex flex-col gap-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-white tracking-tight flex items-center gap-2">
                <SlidersHorizontal className="w-6 h-6 text-indigo-500" />
                Conectores de Ecosistema & MCP
              </h1>
              <p className="text-sm text-slate-400 mt-1">
                Servidores Model Context Protocol (MCP), Skills nativas y Plugins externos para tool-calling.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {connectors.map((c) => {
              const isConnected = c.status === 'connected';
              return (
                <div
                  key={c.id}
                  className="bg-slate-900/50 border border-slate-800/80 hover:border-slate-700/80 rounded-xl p-5 flex flex-col justify-between transition-all"
                >
                  <div className="flex flex-col gap-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-xs uppercase font-mono font-bold tracking-wider px-2 py-0.5 rounded bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                          {c.type}
                        </span>
                        <h3 className="text-base font-semibold text-white">{c.name}</h3>
                      </div>
                      <button
                        onClick={() => onToggleConnector(c.id)}
                        className={`px-3 py-1 rounded-full text-xs font-medium flex items-center gap-1.5 transition-all ${
                          isConnected
                            ? 'bg-emerald-500/15 text-emerald-400 border border-emerald-500/30 hover:bg-red-500/10 hover:text-red-400 hover:border-red-500/30'
                            : 'bg-slate-800 text-slate-400 border border-slate-700 hover:bg-emerald-500/10 hover:text-emerald-400 hover:border-emerald-500/30'
                        }`}
                      >
                        <Power className="w-3 h-3" />
                        {isConnected ? 'Activo' : 'Desactivado'}
                      </button>
                    </div>

                    <p className="text-xs text-slate-400 leading-relaxed">
                      {c.description}
                    </p>

                    <div className="flex flex-col gap-1.5 mt-2">
                      <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
                        Tools / Herramientas Expuestas:
                      </span>
                      <div className="flex flex-wrap gap-1.5">
                        {c.tools.map((t) => (
                          <span
                            key={t}
                            className="text-[11px] font-mono bg-slate-950/80 text-blue-300 px-2 py-0.5 rounded border border-slate-800"
                          >
                            {t}
                          </span>
                        ))}
                      </div>
                    </div>
                  </div>

                  <div className="mt-4 pt-3 border-t border-slate-800/70 flex items-center justify-between text-xs text-slate-400 font-mono">
                    <span className="truncate max-w-[280px]">Endpoint: {c.endpoint}</span>
                    <span className="text-[10px] text-slate-400">ID: {c.id}</span>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    );
  }

  // ----------------------------------------------------
  // SUB-VIEW 3: SESIONES / CHAT PRINCIPAL (Default)
  // ----------------------------------------------------
  return (
    <div className="flex-1 bg-[#090d16] flex flex-col h-screen overflow-hidden">
      {/* Session Top Bar */}
      <header className="h-16 border-b border-slate-800/80 px-6 flex items-center justify-between bg-[#0b0f1a]/80 backdrop-blur-md">
        <div className="flex items-center gap-3">
          <div className="flex flex-col">
            <div className="flex items-center gap-2">
              <h2 className="text-sm font-semibold text-white">
                {selectedSession?.title || 'Selecciona o crea una sesión'}
              </h2>
              <span className="text-[11px] font-mono px-2 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/20">
                {selectedProject?.name || 'Workspace'}
              </span>
            </div>
            <span className="text-[11px] text-slate-400">
              {messages.length} mensajes registrados &bull; Telemetría en tiempo real
            </span>
          </div>
        </div>

        {/* Top Controls: Agent Selector */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 bg-slate-900 border border-slate-800 px-3 py-1.5 rounded-lg text-xs">
            <Bot className="w-3.5 h-3.5 text-blue-400" />
            <span className="text-slate-400">Agente:</span>
            <select
              value={selectedAgentId}
              onChange={(e) => setSelectedAgentId(e.target.value)}
              className="bg-transparent text-white font-medium focus:outline-none cursor-pointer"
            >
              {agents.map((ag) => (
                <option key={ag.id} value={ag.id} className="bg-slate-900 text-white">
                  {ag.name} ({ag.model})
                </option>
              ))}
            </select>
          </div>

          <button
            onClick={() => setShowNewSessModal(true)}
            className="px-3 py-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-medium rounded-lg flex items-center gap-1.5 transition-all"
          >
            <Plus className="w-3.5 h-3.5" />
            Nueva Sesión
          </button>
        </div>
      </header>

      {/* Messages Stream */}
      <div className="flex-1 overflow-y-auto p-6 flex flex-col gap-5">
        {messages.length === 0 ? (
          <div className="flex-1 flex flex-col items-center justify-center text-center p-8 text-slate-400 max-w-md mx-auto my-auto">
            <div className="w-12 h-12 rounded-xl bg-blue-600/10 border border-blue-500/20 flex items-center justify-center text-blue-400 mb-4">
              <Sparkles className="w-6 h-6" />
            </div>
            <h3 className="text-base font-semibold text-slate-200">Comienza la interacción</h3>
            <p className="text-xs text-slate-400 mt-1.5 leading-relaxed">
              Envía un mensaje o prompt. El Hub enrutará la solicitud al agente seleccionado y registrará la telemetría correspondiente.
            </p>
          </div>
        ) : (
          messages.map((m) => {
            const isUser = m.role === 'user';
            return (
              <div
                key={m.id}
                className={`flex gap-3 max-w-4xl ${
                  isUser ? 'ml-auto flex-row-reverse' : 'mr-auto'
                }`}
              >
                {/* Avatar Icon */}
                <div
                  className={`w-8 h-8 rounded-lg flex-shrink-0 flex items-center justify-center text-xs font-bold ${
                    isUser
                      ? 'bg-blue-600 text-white shadow-md shadow-blue-500/20'
                      : 'bg-indigo-600 text-white shadow-md shadow-indigo-500/20'
                  }`}
                >
                  {isUser ? <User className="w-4 h-4" /> : <Bot className="w-4 h-4" />}
                </div>

                {/* Message Body */}
                <div className="flex flex-col gap-1.5 max-w-2xl">
                  <div
                    className={`flex items-center gap-2 text-[11px] ${
                      isUser ? 'justify-end text-slate-400' : 'text-slate-400'
                    }`}
                  >
                    <span className="font-semibold text-slate-200">{m.sender}</span>
                    <span className="font-mono text-slate-400">
                      {new Date(m.createdAt || Date.now()).toLocaleTimeString()}
                    </span>
                  </div>

                  <div
                    className={`p-4 rounded-xl text-sm leading-relaxed whitespace-pre-wrap ${
                      isUser
                        ? 'bg-blue-600/90 text-white rounded-tr-none shadow-md'
                        : 'bg-[#121829] border border-slate-800 text-slate-100 rounded-tl-none shadow-md'
                    }`}
                  >
                    {m.content}
                  </div>

                  {/* Telemetry Footer Badge (Fase 1 / Fase 3 preview) */}
                  {m.telemetry && (
                    <div className="mt-1 flex flex-wrap items-center gap-2 text-[11px] text-slate-400 font-mono bg-slate-900/60 border border-slate-800/80 px-2.5 py-1 rounded-md">
                      <span className="text-blue-400 flex items-center gap-1 font-medium">
                        <Activity className="w-3 h-3" />
                        {m.telemetry.model}
                      </span>
                      <span className="text-slate-400">&bull;</span>
                      <span className="text-emerald-400 flex items-center gap-1">
                        <Clock className="w-3 h-3" />
                        {m.telemetry.totalLatencyMs}ms
                      </span>
                      <span className="text-slate-400">&bull;</span>
                      <span className="text-purple-400 flex items-center gap-1">
                        <Zap className="w-3 h-3" />
                        TTFT: {m.telemetry.ttftMs}ms
                      </span>
                      <span className="text-slate-400">&bull;</span>
                      <span className="text-amber-400">
                        Tokens: {m.telemetry.totalTokens}
                      </span>
                    </div>
                  )}
                </div>
              </div>
            );
          })
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* Input Bar */}
      <footer className="p-4 border-t border-slate-800/80 bg-[#0b0f1a]/90 backdrop-blur-md">
        <form
          onSubmit={handleSubmitMessage}
          className="max-w-4xl mx-auto flex items-center gap-3 bg-slate-900/90 border border-slate-800 rounded-xl p-2 px-3 focus-within:border-blue-500/80 transition-all shadow-inner"
        >
          <input
            type="text"
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            disabled={isSending}
            placeholder={
              isSending
                ? 'El agente está procesando...'
                : 'Escribe una instrucción para el Hub de Agentes (Enter para enviar)...'
            }
            className="flex-1 bg-transparent text-sm text-white placeholder-slate-400 focus:outline-none"
          />
          <button
            type="submit"
            disabled={!inputText.trim() || isSending}
            className="p-2 bg-blue-600 hover:bg-blue-500 disabled:bg-slate-800 disabled:text-slate-600 text-white rounded-lg transition-colors shadow-sm"
          >
            <Send className="w-4 h-4" />
          </button>
        </form>
      </footer>

      {/* Modal Nueva Sesión */}
      {showNewSessModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-[#101524] border border-slate-800 rounded-xl p-6 max-w-md w-full shadow-2xl">
            <h2 className="text-lg font-bold text-white mb-2">Crear Nueva Sesión de Chat</h2>
            <form onSubmit={handleCreateSessionSubmit} className="flex flex-col gap-4">
              <div>
                <label className="text-xs text-slate-400 block mb-1">Título de la Sesión</label>
                <input
                  type="text"
                  required
                  value={newSessTitle}
                  onChange={(e) => setNewSessTitle(e.target.value)}
                  placeholder="ej: Refactorización Concurrente de Workers"
                  className="w-full bg-slate-900 border border-slate-800 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-blue-500"
                />
              </div>
              <div className="flex items-center justify-end gap-3 mt-2">
                <button
                  type="button"
                  onClick={() => setShowNewSessModal(false)}
                  className="px-4 py-2 text-xs text-slate-400 hover:text-white"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg text-xs font-medium"
                >
                  Iniciar Sesión
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};
