import axios from "axios";
import type {
  AgentChatRequest,
  AgentSessionSummary,
  AgentStoredMessage,
  AgentStreamEvent,
} from "../../types/agent";

const LIST_SESSIONS_URL = "/api/chat/sessions";
const LIST_MESSAGES_URL = (sessionId: string) => `/api/chat/sessions/${sessionId}/messages`;
const DELETE_SESSION_URL = (sessionId: string) => `/api/chat/sessions/${sessionId}`;
const CHAT_URL = "/api/chat";

function parseSSEEvent(block: string): AgentStreamEvent | null {
  const lines = block.split("\n").map((line) => line.trim());
  let eventType = "message";
  const dataLines: string[] = [];

  for (const line of lines) {
    if (line.startsWith("event:")) {
      eventType = line.slice(6).trim();
    }
    if (line.startsWith("data:")) {
      dataLines.push(line.slice(5).trim());
    }
  }

  if (dataLines.length === 0) {
    return null;
  }

  const payload = dataLines.join("\n");

  try {
    return {
      type: eventType as AgentStreamEvent["type"],
      ...JSON.parse(payload),
    };
  } catch {
    return {
      type: eventType as AgentStreamEvent["type"],
      message: payload,
    };
  }
}

const AgentService = {
  async listSessions(): Promise<AgentSessionSummary[]> {
    const response = await axios.get<{ sessions: AgentSessionSummary[] }>(
      LIST_SESSIONS_URL
    );
    return response.data.sessions ?? [];
  },

  async listMessages(sessionId: string): Promise<AgentStoredMessage[]> {
    const response = await axios.get<{ messages: AgentStoredMessage[] }>(
      LIST_MESSAGES_URL(sessionId)
    );
    return response.data.messages ?? [];
  },

  async deleteSession(sessionId: string): Promise<void> {
    await axios.delete(DELETE_SESSION_URL(sessionId));
  },

  async streamChat(
    body: AgentChatRequest,
    onEvent: (event: AgentStreamEvent) => void
  ): Promise<void> {
    const response = await fetch(CHAT_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "text/event-stream",
      },
      credentials: "include",
      body: JSON.stringify(body),
    });

    if (!response.ok || !response.body) {
      const message = await response.text();
      throw new Error(message || "Unable to start agent session");
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const blocks = buffer.split("\n\n");
      buffer = blocks.pop() ?? "";

      for (const block of blocks) {
        const event = parseSSEEvent(block);
        if (event) {
          onEvent(event);
        }
      }
    }

    if (buffer.trim()) {
      const event = parseSSEEvent(buffer);
      if (event) {
        onEvent(event);
      }
    }
  },
};

export default AgentService;
