export interface AgentSessionSummary {
  id: string;
  title?: string;
  provider: string;
  model: string;
  last_message: string;
  created_at: string;
  updated_at: string;
}

export interface AgentStoredMessage {
  ID: number;
  SessionID: string;
  Role: string;
  Kind: string;
  Content: string;
  ToolName: string;
  ToolCallID: string;
  MetaJSON: string;
  IsError: boolean;
  CreatedAt: string;
}

export interface AgentStreamEvent {
  type:
    | "session"
    | "tool_call"
    | "tool_result"
    | "assistant"
    | "done"
    | "error";
  session_id?: string;
  message?: string;
  tool_name?: string;
  payload?: any;
}

export interface AgentChatRequest {
  session_id?: string;
  message: string;
  model?: string;
}
