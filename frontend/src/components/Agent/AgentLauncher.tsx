import { useEffect, useRef, useState } from "react";
import { Badge, Flex, IconButton, Text, Tooltip } from "@radix-ui/themes";
import {
  Bot,
  MessageSquarePlus,
  MessagesSquare,
  PanelLeftClose,
  PanelLeftOpen,
  SendHorizontal,
  Sparkles,
  Trash2,
} from "lucide-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

import CButton from "../common/CButton";
import CModel from "../common/CModel";
import AgentService from "../../services/agent/agent.service";
import type {
  AgentSessionSummary,
  AgentStoredMessage,
  AgentStreamEvent,
} from "../../types/agent";
import { useErrorHandler } from "../../hooks/useErrorHandler";
import "./style.css";

type TimelineItem = {
  id: string;
  kind: "user" | "assistant" | "activity";
  content: string;
  meta?: string;
  isError?: boolean;
};

const SESSION_STORAGE_KEY = "grubzo-agent-session-id";

function getErrorMessage(error: unknown) {
  if (typeof error === "string") return error;
  if (error instanceof Error) return error.message;
  return "Something went wrong";
}

function formatDateLabel(value: string) {
  try {
    return new Date(value).toLocaleString([], {
      month: "short",
      day: "numeric",
      hour: "numeric",
      minute: "2-digit",
    });
  } catch {
    return value;
  }
}

function messagePreview(session: AgentSessionSummary) {
  if (session.last_message?.trim()) {
    return session.last_message;
  }
  return "No messages yet";
}

function sessionTitle(session: AgentSessionSummary) {
  if (session.title?.trim()) {
    return session.title;
  }
  return `Chat ${formatDateLabel(session.updated_at)}`;
}

function MarkdownContent({
  content,
  preview = false,
}: {
  content: string;
  preview?: boolean;
}) {
  return (
    <div className={`agent-markdown ${preview ? "preview" : ""}`}>
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          p: ({ children }) => <Text size="2">{children}</Text>,
          h1: ({ children }) => (
            <Text size={preview ? "2" : "4"} weight="bold">
              {children}
            </Text>
          ),
          h2: ({ children }) => (
            <Text size={preview ? "2" : "3"} weight="bold">
              {children}
            </Text>
          ),
          h3: ({ children }) => (
            <Text size="2" weight="bold">
              {children}
            </Text>
          ),
          code(props) {
            const { inline, children, className, ...rest } = props as {
              inline?: boolean;
              children?: React.ReactNode;
              className?: string;
            };
            if (inline) {
              return (
                <code className="agent-inline-code" {...rest}>
                  {children}
                </code>
              );
            }
            return (
              <pre className="agent-code-block">
                <code className={className} {...rest}>
                  {children}
                </code>
              </pre>
            );
          },
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}

function toTimelineMessage(message: AgentStoredMessage): TimelineItem | null {
  if (message.Kind === "text") {
    return {
      id: `${message.Kind}-${message.ID}`,
      kind: message.Role === "assistant" ? "assistant" : "user",
      content: message.Content,
    };
  }
  if (message.Kind === "tool_call") {
    return {
      id: `${message.Kind}-${message.ID}`,
      kind: "activity",
      content: `Using ${message.ToolName}...`,
    };
  }
  if (message.Kind === "tool_result") {
    const [headline] = message.Content.split("\n");
    return {
      id: `${message.Kind}-${message.ID}`,
      kind: "activity",
      content: headline || `${message.ToolName} finished`,
      meta: message.ToolName,
      isError: message.IsError,
    };
  }
  return null;
}

function eventToTimelineItem(event: AgentStreamEvent): TimelineItem | null {
  switch (event.type) {
    case "tool_call":
      return {
        id: `tool-call-${Date.now()}-${event.tool_name}`,
        kind: "activity",
        content: `Using ${event.tool_name}...`,
      };
    case "tool_result": {
      const summary =
        event.payload?.content?.[0]?.text ||
        event.payload?.structuredContent?.note ||
        `${event.tool_name} finished`;
      return {
        id: `tool-result-${Date.now()}-${event.tool_name}`,
        kind: "activity",
        content: summary,
        meta: event.tool_name,
        isError: Boolean(event.payload?.isError),
      };
    }
    case "assistant":
      return {
        id: `assistant-${Date.now()}`,
        kind: "assistant",
        content: event.message || "",
      };
    case "error":
      return {
        id: `error-${Date.now()}`,
        kind: "activity",
        content: event.message || "Something went wrong",
        isError: true,
      };
    default:
      return null;
  }
}

const AgentLauncher = () => {
  const { showError } = useErrorHandler();
  const [open, setOpen] = useState(false);
  const [sessions, setSessions] = useState<AgentSessionSummary[]>([]);
  const [selectedSessionId, setSelectedSessionId] = useState<string | null>(
    null,
  );
  const [timeline, setTimeline] = useState<TimelineItem[]>([]);
  const [draft, setDraft] = useState("");
  const [loadingSessions, setLoadingSessions] = useState(false);
  const [loadingMessages, setLoadingMessages] = useState(false);
  const [sessionsCollapsed, setSessionsCollapsed] = useState(false);
  const [sending, setSending] = useState(false);
  const threadRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!threadRef.current) return;
    threadRef.current.scrollTop = threadRef.current.scrollHeight;
  }, [timeline, open]);

  useEffect(() => {
    if (!open) return;
    const previousOverflow = document.body.style.overflow;
    const previousRootOverflow = document.documentElement.style.overflow;
    document.body.style.overflow = "hidden";
    document.documentElement.style.overflow = "hidden";
    return () => {
      document.body.style.overflow = previousOverflow;
      document.documentElement.style.overflow = previousRootOverflow;
    };
  }, [open]);

  useEffect(() => {
    if (!open) return;
    void loadSessions();
  }, [open]);

  useEffect(() => {
    if (!open || !selectedSessionId || sending) return;
    void loadMessages(selectedSessionId);
  }, [open, selectedSessionId, sending]);

  const loadSessions = async (preferredSessionId?: string | null) => {
    setLoadingSessions(true);
    try {
      const nextSessions = await AgentService.listSessions();
      setSessions(nextSessions);

      const storedSessionId = localStorage.getItem(SESSION_STORAGE_KEY);
      const targetSessionId = preferredSessionId !== undefined
        ? preferredSessionId
        : selectedSessionId || storedSessionId || nextSessions[0]?.id || null;

      if (
        targetSessionId &&
        nextSessions.some((s) => s.id === targetSessionId)
      ) {
        setSelectedSessionId(targetSessionId);
        localStorage.setItem(SESSION_STORAGE_KEY, targetSessionId);
      } else if (!sending) {
        setSelectedSessionId(null);
        setTimeline([]);
        localStorage.removeItem(SESSION_STORAGE_KEY);
      }
    } catch (error) {
      setTimeline((prev) => [
        ...prev,
        {
          id: `session-load-error-${Date.now()}`,
          kind: "activity",
          content: getErrorMessage(error),
          isError: true,
        },
      ]);
      showError(error);
    } finally {
      setLoadingSessions(false);
    }
  };

  const loadMessages = async (sessionId: string) => {
    setLoadingMessages(true);
    try {
      const messages = await AgentService.listMessages(sessionId);
      setTimeline(
        messages.map(toTimelineMessage).filter(Boolean) as TimelineItem[],
      );
    } catch (error) {
      setTimeline((prev) => [
        ...prev,
        {
          id: `message-load-error-${Date.now()}`,
          kind: "activity",
          content: getErrorMessage(error),
          isError: true,
        },
      ]);
      showError(error);
    } finally {
      setLoadingMessages(false);
    }
  };

  const handleNewSession = () => {
    setSelectedSessionId(null);
    setTimeline([]);
    localStorage.removeItem(SESSION_STORAGE_KEY);
  };

  const handleSelectSession = (sessionId: string) => {
    setSelectedSessionId(sessionId);
    localStorage.setItem(SESSION_STORAGE_KEY, sessionId);
  };

  const handleDeleteSession = async (sessionId: string) => {
    try {
      await AgentService.deleteSession(sessionId);

      if (selectedSessionId === sessionId) {
        setSelectedSessionId(null);
        setTimeline([]);
        localStorage.removeItem(SESSION_STORAGE_KEY);
      }

      await loadSessions(null);
    } catch (error) {
      showError(error);
    }
  };

  const handleSend = async () => {
    const message = draft.trim();
    if (!message || sending) return;

    setDraft("");
    setSending(true);
    setTimeline((prev) => [
      ...prev,
      {
        id: `user-${Date.now()}`,
        kind: "user",
        content: message,
      },
    ]);

    let nextSessionId = selectedSessionId;

    try {
      await AgentService.streamChat(
        {
          session_id: selectedSessionId ?? undefined,
          message,
        },
        (event) => {
          if (event.type === "session" && event.session_id) {
            nextSessionId = event.session_id;
            setSelectedSessionId(event.session_id);
            localStorage.setItem(SESSION_STORAGE_KEY, event.session_id);
          }

          const item = eventToTimelineItem(event);
          if (item) {
            setTimeline((prev) => [...prev, item]);
          }
        },
      );

      if (nextSessionId) {
        await loadSessions(nextSessionId);
        await loadMessages(nextSessionId);
      } else {
        await loadSessions();
      }
    } catch (error) {
      setTimeline((prev) => [
        ...prev,
        {
          id: `send-error-${Date.now()}`,
          kind: "activity",
          content: getErrorMessage(error),
          meta: "Agent",
          isError: true,
        },
      ]);
      showError(error);
    } finally {
      setSending(false);
    }
  };

  return (
    <>
      <Tooltip content="Open Grubzo assistant">
        <IconButton
          radius="small"
          variant="solid"
          onClick={() => setOpen(true)}
          className="agent-launcher-btn"
        >
          <Bot size={18} />
        </IconButton>
      </Tooltip>

      <CModel
        open={open}
        onClose={() => setOpen(false)}
        title="Grubzo Assistant"
        size="full"
        anchor="center"
        closeOnBackdrop
      >
        <div
          className={`agent-shell ${
            sessionsCollapsed ? "sessions-collapsed" : ""
          }`}
        >
          <div className="agent-sidebar">
            <div className="agent-toolbar">
              <Flex direction="column" gap="3">
                <Flex
                  justify="between"
                  align="center"
                  className="agent-session-count"
                >
                  {!sessionsCollapsed ? (
                    <Text size="2" weight="bold">
                      Sessions
                    </Text>
                  ) : null}
                  <Flex align="center" gap="2">
                    {!sessionsCollapsed ? (
                      <Badge variant="soft" color="violet" size="3">
                        {sessions.length}
                      </Badge>
                    ) : null}
                    <IconButton
                      radius="small"
                      variant="soft"
                      className="agent-toolbar-icon"
                      onClick={() =>
                        setSessionsCollapsed((current) => !current)
                      }
                    >
                      {sessionsCollapsed ? (
                        <PanelLeftOpen size={15} />
                      ) : (
                        <PanelLeftClose size={15} />
                      )}
                    </IconButton>
                  </Flex>
                </Flex>
                {!sessionsCollapsed ? (
                  <CButton
                    label="New Session"
                    startIcon={<MessageSquarePlus size={16} />}
                    onClick={handleNewSession}
                    fullWidth
                    styles={{ justifyContent: "center" }}
                  />
                ) : (
                  <IconButton
                    radius="small"
                    variant="soft"
                    className="agent-toolbar-icon"
                    onClick={handleNewSession}
                  >
                    <MessageSquarePlus size={16} />
                  </IconButton>
                )}
              </Flex>
            </div>

            <Flex className="agent-session-list" direction="column" gap="2">
              {loadingSessions ? (
                <Text size="2" color="gray">
                  Loading sessions...
                </Text>
              ) : sessions.length === 0 ? (
                <Text size="2" color="gray">
                  No previous sessions yet.
                </Text>
              ) : (
                sessions.map((session) => (
                  <button
                    key={session.id}
                    type="button"
                    className={`agent-session-btn ${
                      session.id === selectedSessionId ? "active" : ""
                    }`}
                    onClick={() => handleSelectSession(session.id)}
                  >
                    <Flex direction="column" gap="2">
                      <Flex justify="between" align="center" gap="2">
                        <div className="agent-session-meta">
                          <Text size="2" weight="bold" className="agent-session-title">
                            {sessionTitle(session)}
                          </Text>
                        </div>
                        <Flex align="center" gap="2" className="agent-session-actions">
                          <div className="agent-session-meta">
                            <Text size="1" color="gray">
                              {formatDateLabel(session.updated_at)}
                            </Text>
                          </div>
                          <IconButton
                            type="button"
                            radius="small"
                            variant="ghost"
                            color="gray"
                            className="agent-session-delete"
                            onClick={(event) => {
                              event.preventDefault();
                              event.stopPropagation();
                              void handleDeleteSession(session.id);
                            }}
                          >
                            <Trash2 size={14} />
                          </IconButton>
                        </Flex>
                      </Flex>
                      <div className="agent-session-meta agent-session-preview">
                        <MarkdownContent
                          content={messagePreview(session)}
                          preview
                        />
                      </div>
                    </Flex>
                  </button>
                ))
              )}
            </Flex>
          </div>

          <div className="agent-main">
            <div className="agent-thread" ref={threadRef}>
              {loadingMessages ? (
                <Text size="2" color="gray">
                  Loading messages...
                </Text>
              ) : timeline.length === 0 ? (
                <div className="agent-empty">
                  <Flex direction="column" gap="3" align="center">
                    <Sparkles size={28} />
                    <Text size="5" weight="bold">
                      Start a fresh Grubzo session
                    </Text>
                    <Text size="2" color="gray">
                      Ask about food, track an order, or place a new one. Your
                      previous sessions are listed on the left.
                    </Text>
                  </Flex>
                </div>
              ) : (
                timeline.map((item) => (
                  <div
                    key={item.id}
                    className={`agent-bubble ${item.kind} ${
                      item.isError ? "error" : ""
                    }`}
                  >
                    <Flex direction="column" gap="2">
                      {item.meta && (
                        <Flex align="center" gap="2">
                          <MessagesSquare size={14} />
                          <Text size="1" weight="bold">
                            {item.meta}
                          </Text>
                        </Flex>
                      )}
                      {item.kind === "assistant" ? (
                        <MarkdownContent content={item.content} />
                      ) : (
                        <Text
                          size="2"
                          color={
                            item.kind === "activity" && item.isError
                              ? "red"
                              : undefined
                          }
                        >
                          {item.content}
                        </Text>
                      )}
                    </Flex>
                  </div>
                ))
              )}
            </div>

            <div>
              <Flex className="agent-composer" direction="column" gap="3">
                <Flex className="agent-input-wrap">
                  <textarea
                    className="agent-input"
                    value={draft}
                    onChange={(event) => setDraft(event.target.value)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter" && !event.shiftKey) {
                        event.preventDefault();
                        void handleSend();
                      }
                    }}
                    disabled={sending}
                  />
                  <IconButton
                    radius="small"
                    variant="solid"
                    className="agent-send-btn"
                    onClick={() => void handleSend()}
                    disabled={!draft.trim() || sending}
                  >
                    <SendHorizontal size={16} />
                  </IconButton>
                </Flex>
                <Flex
                  justify="center"
                  align="center"
                  className="agent-composer-hint"
                >
                  <Text size="1" color="gray">
                    {selectedSessionId
                      ? "Resume the selected session or start a new one from the left."
                      : "A new session will start when you send your first message."}
                  </Text>
                </Flex>
              </Flex>
            </div>
          </div>
        </div>
      </CModel>
    </>
  );
};

export default AgentLauncher;
