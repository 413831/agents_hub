import React, { useState, useRef, useEffect } from 'react';
import { Project, Session, Message, Connector, Agent, NavTab } from '../types';

interface CliViewProps {
  selectedProject: Project | null;
  selectedSession: Session | null;
  messages: Message[];
  agents: Agent[];
  connectors: Connector[];
  projects: Project[];
  sessions: Session[];
  onSendMessage: (content: string, agentId: string) => Promise<void>;
  onToggleConnector: (id: string) => Promise<void>;
  onCreateProject: (name: string, description: string) => Promise<void>;
  onCreateSession: (title: string, agentIds: string[]) => Promise<void>;
  onSelectProject: (proj: Project) => void;
  onSelectSession: (sess: Session) => void;
  onToggleViewMode: () => void;
  onSelectTab: (tab: NavTab) => void;
}

interface CliLine {
  id: string;
  type: 'system' | 'input' | 'output' | 'error' | 'telemetry';
  text: string;
}

export const CliView: React.FC<CliViewProps> = ({
  selectedProject,
  selectedSession,
  messages,
  agents,
  connectors,
  projects,
  sessions,
  onSendMessage,
  onToggleConnector,
  onCreateProject,
  onCreateSession,
  onSelectProject,
  onSelectSession,
  onToggleViewMode,
  onSelectTab,
}) => {
  const [cliInput, setCliInput] = useState('');
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState<number>(-1);
  const [lines, setLines] = useState<CliLine[]>([]);
  const [selectedAgentId, setSelectedAgentId] = useState<string>('');
  const terminalEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Default agent
  useEffect(() => {
    if (agents.length > 0 && !selectedAgentId) {
      setSelectedAgentId(agents[0].id);
    }
  }, [agents, selectedAgentId]);

  // Initial banner
  useEffect(() => {
    const bannerLines: CliLine[] = [
      { id: 'b1', type: 'system', text: '========================================================================' },
      { id: 'b2', type: 'system', text: ' AGENTS HUB CLI INTERFACE :: Clean Architecture & DDD v1.0' },
      { id: 'b3', type: 'system', text: ` Active Project: ${selectedProject?.name || 'None'} [${selectedProject?.id || '-'}]` },
      { id: 'b4', type: 'system', text: ` Active Session: ${selectedSession?.title || 'None'} [${selectedSession?.id || '-'}]` },
      { id: 'b5', type: 'system', text: ` Active Agent:   ${agents.find(a => a.id === selectedAgentId)?.name || 'Default'}` },
      { id: 'b6', type: 'system', text: ' Type "help" for a list of available commands, or type directly to chat.' },
      { id: 'b7', type: 'system', text: '========================================================================\n' },
    ];
    setLines(bannerLines);
  }, [selectedProject?.id, selectedSession?.id]);

  // Sync existing messages to CLI view
  useEffect(() => {
    if (messages.length === 0) return;
    const msgLines: CliLine[] = [];
    messages.forEach((m) => {
      const isUser = m.role === 'user';
      if (isUser) {
        msgLines.push({
          id: `cli-m-${m.id}-in`,
          type: 'input',
          text: `user@hub:~$ ${m.content}`,
        });
      } else {
        msgLines.push({
          id: `cli-m-${m.id}-out`,
          type: 'output',
          text: `[${m.sender}]:\n${m.content}`,
        });
        if (m.telemetry) {
          msgLines.push({
            id: `cli-m-${m.id}-tel`,
            type: 'telemetry',
            text: `[Model: ${m.telemetry.model} | Latency: ${m.telemetry.totalLatencyMs}ms | TTFT: ${m.telemetry.ttftMs}ms | Tokens: ${m.telemetry.totalTokens} | Finish: ${m.telemetry.finishReason}]`,
          });
        }
      }
    });
    setLines((prev) => {
      const banner = prev.slice(0, 7);
      return [...banner, ...msgLines];
    });
  }, [messages]);

  useEffect(() => {
    terminalEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [lines]);

  const handleCommand = async (cmdText: string) => {
    const trimmed = cmdText.trim();
    if (!trimmed) return;

    // Add to history
    setHistory((prev) => [...prev, trimmed]);
    setHistoryIndex(-1);

    // Append user input line
    const userLine: CliLine = {
      id: `usr-${Date.now()}`,
      type: 'input',
      text: `user@hub:~$ ${trimmed}`,
    };
    setLines((prev) => [...prev, userLine]);

    const parts = trimmed.split(' ');
    const cmd = parts[0].toLowerCase();
    const args = parts.slice(1).join(' ');

    switch (cmd) {
      case 'help':
        setLines((prev) => [
          ...prev,
          {
            id: `help-${Date.now()}`,
            type: 'system',
            text: `Available Commands:
  help                     Show this help guide
  projects                 List all workspace projects
  sessions                 List sessions in current workspace
  connectors               List MCP connectors and exposed tools
  agents                   List all registered agents
  select-agent <agent_id>  Set active agent (e.g. select-agent agent-arch-1)
  create-project <name>    Create a new workspace project
  create-session <title>   Create a new chat session
  toggle-conn <conn_id>    Toggle MCP connector status
  mode gui                 Switch view to Modern GUI Dashboard
  clear                    Clear the terminal output screen
  <prompt>                 Any text without a recognized command is sent to the active agent.`,
          },
        ]);
        break;

      case 'clear':
        setLines([]);
        break;

      case 'mode':
        if (args.toLowerCase() === 'gui') {
          onToggleViewMode();
        } else {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: 'Usage: mode gui' },
          ]);
        }
        break;

      case 'projects':
        const projTable = projects
          .map((p) => `  * [${p.id}] ${p.name.padEnd(28)} Model: ${p.defaultModel || 'N/A'}`)
          .join('\n');
        setLines((prev) => [
          ...prev,
          {
            id: `p-${Date.now()}`,
            type: 'output',
            text: `WORKSPACE PROJECTS (${projects.length}):\n${projTable}`,
          },
        ]);
        break;

      case 'sessions':
        const sessTable = sessions
          .map((s) => `  * [${s.id}] ${s.title.padEnd(35)} Status: ${s.status}`)
          .join('\n');
        setLines((prev) => [
          ...prev,
          {
            id: `s-${Date.now()}`,
            type: 'output',
            text: `SESSIONS IN PROJECT [${selectedProject?.id || 'all'}] (${sessions.length}):\n${sessTable}`,
          },
        ]);
        break;

      case 'connectors':
        const connTable = connectors
          .map(
            (c) =>
              `  * [${c.id}] ${c.name.padEnd(30)} Type: ${c.type.padEnd(8)} Status: ${c.status.padEnd(12)} Tools: [${c.tools.join(', ')}]`
          )
          .join('\n');
        setLines((prev) => [
          ...prev,
          {
            id: `c-${Date.now()}`,
            type: 'output',
            text: `MCP CONNECTORS & SKILLS (${connectors.length}):\n${connTable}`,
          },
        ]);
        break;

      case 'agents':
        const agentTable = agents
          .map((a) => `  * [${a.id}] ${a.name.padEnd(20)} Model: ${a.provider}/${a.model}`)
          .join('\n');
        setLines((prev) => [
          ...prev,
          {
            id: `a-${Date.now()}`,
            type: 'output',
            text: `REGISTERED AGENTS (${agents.length}):\n${agentTable}`,
          },
        ]);
        break;

      case 'select-agent':
        if (!args) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: 'Usage: select-agent <agent_id>' },
          ]);
          return;
        }
        const targetAgent = agents.find((a) => a.id === args);
        if (targetAgent) {
          setSelectedAgentId(targetAgent.id);
          setLines((prev) => [
            ...prev,
            {
              id: `ag-${Date.now()}`,
              type: 'system',
              text: `Active agent updated to: ${targetAgent.name} (${targetAgent.model})`,
            },
          ]);
        } else {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: `Agent '${args}' not found.` },
          ]);
        }
        break;

      case 'create-project':
        if (!args) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: 'Usage: create-project <name>' },
          ]);
          return;
        }
        try {
          await onCreateProject(args, 'Created via CLI');
          setLines((prev) => [
            ...prev,
            { id: `ok-${Date.now()}`, type: 'system', text: `Project '${args}' created successfully.` },
          ]);
        } catch (err: any) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: `Failed to create project: ${err?.message || err}` },
          ]);
        }
        break;

      case 'create-session':
        if (!args) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: 'Usage: create-session <title>' },
          ]);
          return;
        }
        try {
          await onCreateSession(args, selectedAgentId ? [selectedAgentId] : []);
          setLines((prev) => [
            ...prev,
            { id: `ok-${Date.now()}`, type: 'system', text: `Session '${args}' started successfully.` },
          ]);
        } catch (err: any) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: `Failed to create session: ${err?.message || err}` },
          ]);
        }
        break;

      case 'toggle-conn':
        if (!args) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: 'Usage: toggle-conn <connector_id>' },
          ]);
          return;
        }
        try {
          await onToggleConnector(args);
          setLines((prev) => [
            ...prev,
            { id: `ok-${Date.now()}`, type: 'system', text: `Connector '${args}' status toggled.` },
          ]);
        } catch (err: any) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: `Failed to toggle connector: ${err?.message || err}` },
          ]);
        }
        break;

      default:
        // Direct chat prompt!
        try {
          setLines((prev) => [
            ...prev,
            { id: `wait-${Date.now()}`, type: 'system', text: '[Orchestrating request through Clean Architecture LLM port...]' },
          ]);
          await onSendMessage(trimmed, selectedAgentId);
        } catch (err: any) {
          setLines((prev) => [
            ...prev,
            { id: `err-${Date.now()}`, type: 'error', text: `Inference failed: ${err?.message || err}` },
          ]);
        }
        break;
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      const val = cliInput;
      setCliInput('');
      handleCommand(val);
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (history.length > 0) {
        const nextIdx = historyIndex + 1;
        if (nextIdx < history.length) {
          setHistoryIndex(nextIdx);
          setCliInput(history[history.length - 1 - nextIdx]);
        }
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (historyIndex > 0) {
        const nextIdx = historyIndex - 1;
        setHistoryIndex(nextIdx);
        setCliInput(history[history.length - 1 - nextIdx]);
      } else if (historyIndex === 0) {
        setHistoryIndex(-1);
        setCliInput('');
      }
    }
  };

  return (
    <div 
      className="flex-1 bg-[#05070a] text-slate-300 font-mono flex flex-col h-screen overflow-hidden cli-scanlines select-text"
      onClick={() => inputRef.current?.focus()}
    >
      {/* CLI Header Bar */}
      <header className="h-10 border-b border-slate-900 bg-black/60 px-4 flex items-center justify-between text-xs text-slate-400">
        <div className="flex items-center gap-2">
          <span className="w-2.5 h-2.5 rounded-full bg-emerald-500 inline-block animate-pulse"></span>
          <span className="text-slate-300 font-semibold">user@hub:~ (CLI Terminal Mode)</span>
        </div>
        <div className="flex items-center gap-4 text-[11px]">
          <span>Session: {selectedSession?.id || 'none'}</span>
          <span>Agent: {selectedAgentId || 'default'}</span>
          <button
            onClick={onToggleViewMode}
            className="text-emerald-400 hover:text-emerald-300 underline font-mono text-[11px]"
          >
            Switch to GUI
          </button>
        </div>
      </header>

      {/* Terminal Output Area */}
      <div className="flex-1 overflow-y-auto p-4 text-xs space-y-1.5 leading-relaxed">
        {lines.map((l) => {
          if (l.type === 'input') {
            return (
              <div key={l.id} className="text-emerald-400 font-semibold">
                {l.text}
              </div>
            );
          }
          if (l.type === 'telemetry') {
            return (
              <div key={l.id} className="text-purple-400/90 pl-3 border-l-2 border-purple-500/40 text-[11px]">
                {l.text}
              </div>
            );
          }
          if (l.type === 'system') {
            return (
              <div key={l.id} className="text-slate-500 whitespace-pre-wrap">
                {l.text}
              </div>
            );
          }
          if (l.type === 'error') {
            return (
              <div key={l.id} className="text-red-400 whitespace-pre-wrap">
                [ERROR] {l.text}
              </div>
            );
          }
          // output
          return (
            <div key={l.id} className="text-slate-200 whitespace-pre-wrap pl-2 border-l border-slate-800">
              {l.text}
            </div>
          );
        })}
        <div ref={terminalEndRef} />
      </div>

      {/* Monochromatic Command Prompt */}
      <div className="p-3 border-t border-slate-900 bg-black/80 flex items-center gap-2">
        <span className="text-emerald-400 font-bold text-xs select-none">user@hub:~$</span>
        <input
          ref={inputRef}
          type="text"
          value={cliInput}
          onChange={(e) => setCliInput(e.target.value)}
          onKeyDown={handleKeyDown}
          autoFocus
          spellCheck={false}
          placeholder="Escribe un comando o prompt..."
          className="flex-1 bg-transparent text-slate-100 text-xs font-mono focus:outline-none placeholder-slate-700"
        />
      </div>
    </div>
  );
};
