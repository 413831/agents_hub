// frontend/src/main.js

document.addEventListener('DOMContentLoaded', () => {
  // 1. Tab Navigation
  const navButtons = document.querySelectorAll('.nav-btn');
  const viewContainers = document.querySelectorAll('.view-container');

  navButtons.forEach(btn => {
    btn.addEventListener('click', () => {
      const targetView = btn.getAttribute('data-view');
      
      navButtons.forEach(b => b.classList.remove('active'));
      viewContainers.forEach(v => v.classList.remove('active'));

      btn.classList.add('active');
      const targetEl = document.getElementById(targetView);
      if (targetEl) {
        targetEl.classList.add('active');
      }
    });
  });

  // 2. Chat View - Form handling & Wails Bridge
  const sendBtn = document.getElementById('chat-send-btn');
  const chatInput = document.getElementById('chat-input');
  const chatBox = document.getElementById('chat-messages');
  const providerSelect = document.getElementById('provider-select');
  const modelSelect = document.getElementById('model-select');

  function appendMessage(role, text) {
    const bubble = document.createElement('div');
    bubble.className = `msg-bubble ${role}`;
    bubble.innerHTML = `<strong>${role === 'user' ? 'Tú' : 'AI Assistant'}:</strong> <div>${text}</div>`;
    chatBox.appendChild(bubble);
    chatBox.scrollTop = chatBox.scrollHeight;
  }

  async function handleSend() {
    const text = chatInput.value.trim();
    if (!text) return;

    appendMessage('user', text);
    chatInput.value = '';

    const req = {
      provider: providerSelect.value,
      model: modelSelect.value,
      messages: [{ role: 'user', content: text }],
      temperature: 0.7,
      max_tokens: 1024
    };

    // Invocar Wails si está disponible en window.go, o emular respuesta limpia
    if (window.go && window.go.wails && window.go.wails.LLMHandler) {
      try {
        const resp = await window.go.wails.LLMHandler.GenerateCompletion(req);
        appendMessage('assistant', resp.content);
      } catch (err) {
        appendMessage('assistant', `[Error]: ${err}`);
      }
    } else {
      // Mock para vista desacoplada inicial
      setTimeout(() => {
        appendMessage('assistant', `[${req.provider.toUpperCase()} Prototipo]: Recibido prompt "${text}". Conectado a Clean Architecture.`);
      }, 500);
    }
  }

  if (sendBtn && chatInput) {
    sendBtn.addEventListener('click', handleSend);
    chatInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') handleSend();
    });
  }

  // 3. Task Queue Trigger
  const enqueueBtn = document.getElementById('enqueue-job-btn');
  if (enqueueBtn) {
    enqueueBtn.addEventListener('click', async () => {
      if (window.go && window.go.wails && window.go.wails.TaskHandler) {
        try {
          await window.go.wails.TaskHandler.EnqueueLLMTask('task:llm:inference', { prompt: 'Tarea pesada background' });
        } catch (e) {
          console.error(e);
        }
      } else {
        const tbody = document.querySelector('#tasks-table tbody');
        if (tbody) {
          const tr = document.createElement('tr');
          const id = 'job_' + Math.floor(Math.random() * 10000);
          tr.innerHTML = `
            <td><code>${id}</code></td>
            <td>task:llm:inference</td>
            <td><span class="badge running">running</span></td>
            <td>50%</td>
            <td>Procesando en Message Broker...</td>
          `;
          tbody.prepend(tr);
        }
      }
    });
  }
});
