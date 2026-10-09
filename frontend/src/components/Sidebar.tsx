import React from 'react';
import { 
  FolderKanban, 
  MessageSquare, 
  Cpu, 
  Terminal, 
  LayoutDashboard, 
  Plus, 
  Layers, 
  Radio
} from 'lucide-react';
import { Project, Session, ViewMode, NavTab } from '../types';

interface SidebarProps {
  currentTab: NavTab;
  onSelectTab: (tab: NavTab) => void;
  viewMode: ViewMode;
  onToggleViewMode: () => void;
  projects: Project[];
  selectedProject: Project | null;
  onSelectProject: (proj: Project) => void;
  sessions: Session[];
  selectedSession: Session | null;
  onSelectSession: (sess: Session) => void;
  onCreateSession: () => void;
  onCreateProject: () => void;
}

export const Sidebar: React.FC<SidebarProps> = ({
  currentTab,
  onSelectTab,
  viewMode,
  onToggleViewMode,
  projects,
  selectedProject,
  onSelectProject,
  sessions,
  selectedSession,
  onSelectSession,
  onCreateSession,
  onCreateProject,
}) => {
  return (
    <aside className="w-72 bg-[#0c101c] border-r border-slate-800/80 flex flex-col h-screen select-none">
      {/* Brand & Dual-View Toggle Header */}
      <div className="p-4 border-b border-slate-800/80 flex flex-col gap-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-600 to-indigo-600 flex items-center justify-center text-white font-bold text-sm shadow-lg shadow-blue-500/20">
              <Layers className="w-4 h-4" />
            </div>
            <div>
              <div className="text-sm font-semibold tracking-tight text-white flex items-center gap-1.5">
                Agents Hub
                <span className="text-[10px] px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-400 font-mono font-medium border border-blue-500/20">
                  v1.0
                </span>
              </div>
              <p className="text-[11px] text-slate-400">Clean Arch & DDD</p>
            </div>
          </div>
        </div>

        {/* Dual-View Mode Switcher */}
        <div className="flex items-center p-1 bg-slate-900/90 rounded-lg border border-slate-800 text-xs">
          <button
            id="toggle-gui-mode"
            onClick={() => { if (viewMode !== 'gui') onToggleViewMode(); }}
            className={`flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-md transition-all font-medium ${
              viewMode === 'gui'
                ? 'bg-blue-600 text-white shadow-sm'
                : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50'
            }`}
          >
            <LayoutDashboard className="w-3.5 h-3.5" />
            Modo GUI
          </button>
          <button
            id="toggle-cli-mode"
            onClick={() => { if (viewMode !== 'cli') onToggleViewMode(); }}
            className={`flex-1 flex items-center justify-center gap-1.5 py-1.5 rounded-md transition-all font-medium ${
              viewMode === 'cli'
                ? 'bg-emerald-600 text-white shadow-sm font-mono'
                : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/50 font-mono'
            }`}
          >
            <Terminal className="w-3.5 h-3.5" />
            Modo CLI
          </button>
        </div>
      </div>

      {/* Main Navigation Tabs */}
      <div className="px-3 pt-3 flex flex-col gap-1">
        <button
          id="nav-projects-tab"
          onClick={() => onSelectTab('projects')}
          className={`flex items-center justify-between px-3 py-2 rounded-lg text-xs font-medium transition-all ${
            currentTab === 'projects'
              ? 'bg-blue-600/15 text-blue-400 border border-blue-500/30'
              : 'text-slate-300 hover:bg-slate-800/50 hover:text-white'
          }`}
        >
          <div className="flex items-center gap-2.5">
            <FolderKanban className="w-4 h-4" />
            <span>Proyectos</span>
          </div>
          <span className="text-[11px] bg-slate-800/80 text-slate-400 px-1.5 py-0.5 rounded-md font-mono">
            {projects.length}
          </span>
        </button>

        <button
          id="nav-sessions-tab"
          onClick={() => onSelectTab('sessions')}
          className={`flex items-center justify-between px-3 py-2 rounded-lg text-xs font-medium transition-all ${
            currentTab === 'sessions'
              ? 'bg-blue-600/15 text-blue-400 border border-blue-500/30'
              : 'text-slate-300 hover:bg-slate-800/50 hover:text-white'
          }`}
        >
          <div className="flex items-center gap-2.5">
            <MessageSquare className="w-4 h-4" />
            <span>Sesiones Activas</span>
          </div>
          <span className="text-[11px] bg-slate-800/80 text-slate-400 px-1.5 py-0.5 rounded-md font-mono">
            {sessions.length}
          </span>
        </button>

        <button
          id="nav-connectors-tab"
          onClick={() => onSelectTab('connectors')}
          className={`flex items-center justify-between px-3 py-2 rounded-lg text-xs font-medium transition-all ${
            currentTab === 'connectors'
              ? 'bg-blue-600/15 text-blue-400 border border-blue-500/30'
              : 'text-slate-300 hover:bg-slate-800/50 hover:text-white'
          }`}
        >
          <div className="flex items-center gap-2.5">
            <Cpu className="w-4 h-4" />
            <span>Conectores (MCP)</span>
          </div>
          <span className="text-[10px] text-emerald-400 bg-emerald-500/10 px-1.5 py-0.5 rounded border border-emerald-500/20 font-mono">
            LIVE
          </span>
        </button>
      </div>

      {/* Dynamic Sub-List depending on context */}
      <div className="flex-1 overflow-y-auto px-3 py-3 mt-1 flex flex-col gap-3">
        {/* Workspace selector dropdown or list */}
        <div className="flex flex-col gap-1.5">
          <div className="flex items-center justify-between px-1">
            <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
              Workspace Activo
            </span>
            <button 
              onClick={onCreateProject}
              title="Nuevo Proyecto"
              className="text-slate-400 hover:text-white transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
            </button>
          </div>
          <div className="flex flex-col gap-1">
            {projects.map((p) => (
              <button
                key={p.id}
                onClick={() => onSelectProject(p)}
                className={`text-left px-2.5 py-1.5 rounded-md text-xs transition-colors flex items-center justify-between ${
                  selectedProject?.id === p.id
                    ? 'bg-slate-800 text-white font-medium border-l-2 border-blue-500'
                    : 'text-slate-400 hover:bg-slate-800/40 hover:text-slate-200'
                }`}
              >
                <span className="truncate pr-2">{p.name}</span>
                <span className="text-[10px] text-slate-400 font-mono flex-shrink-0">
                  {p.defaultModel?.split('-')[0] || 'ia'}
                </span>
              </button>
            ))}
          </div>
        </div>

        {/* Sessions list */}
        <div className="flex flex-col gap-1.5 pt-2 border-t border-slate-800/60">
          <div className="flex items-center justify-between px-1">
            <span className="text-[11px] font-semibold text-slate-400 uppercase tracking-wider">
              Historial de Chats
            </span>
            <button 
              onClick={onCreateSession}
              title="Nueva Sesión"
              className="text-slate-400 hover:text-white transition-colors"
            >
              <Plus className="w-3.5 h-3.5" />
            </button>
          </div>
          <div className="flex flex-col gap-1">
            {sessions.map((s) => (
              <button
                key={s.id}
                onClick={() => {
                  onSelectSession(s);
                  if (currentTab !== 'sessions') onSelectTab('sessions');
                }}
                className={`text-left px-2.5 py-2 rounded-md text-xs transition-colors flex flex-col gap-0.5 ${
                  selectedSession?.id === s.id
                    ? 'bg-blue-600/20 text-blue-300 font-medium border border-blue-500/30'
                    : 'text-slate-400 hover:bg-slate-800/40 hover:text-slate-200'
                }`}
              >
                <div className="truncate text-slate-200">{s.title}</div>
                <div className="text-[10px] text-slate-400 flex items-center gap-1.5 font-mono">
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 inline-block"></span>
                  {s.agentIds?.length ? `${s.agentIds.length} agente(s)` : 'Orchestrator'}
                </div>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Footer System Status */}
      <div className="p-3 border-t border-slate-800/80 bg-slate-950/40 flex items-center justify-between text-xs text-slate-400">
        <div className="flex items-center gap-2">
          <Radio className="w-3.5 h-3.5 text-emerald-400 animate-pulse" />
          <span className="text-[11px]">Wails v2 Core</span>
        </div>
        <span className="text-[10px] font-mono bg-slate-800/60 text-slate-400 px-2 py-0.5 rounded">
          Go 1.25
        </span>
      </div>
    </aside>
  );
};
