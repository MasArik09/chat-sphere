import { useMessageStore } from '../store/message-store';
import { useConversationStore } from '../store/conversation-store';
import { usePresenceStore } from '../store/presence-store';

class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectTimeout: number | null = null;
  private shouldReconnect = true;
  private token: string | null = null;

  connect(token: string) {
    this.token = token;
    this.shouldReconnect = true;
    
    if (this.socket) {
      this.socket.close();
    }

    const wsUrl = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
    const connectionUrl = `${wsUrl}?token=${token}`;

    try {
      this.socket = new WebSocket(connectionUrl);

      this.socket.onopen = () => {
        console.log('WebSocket connection established.');
        if (this.reconnectTimeout) {
          clearTimeout(this.reconnectTimeout);
          this.reconnectTimeout = null;
        }
      };

      this.socket.onmessage = (event) => {
        try {
          const payload = JSON.parse(event.data);
          this.handleEvent(payload);
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err);
        }
      };

      this.socket.onclose = () => {
        console.log('WebSocket connection closed.');
        if (this.shouldReconnect) {
          this.queueReconnect();
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket error:', error);
      };
    } catch (err) {
      console.error('Failed to create WebSocket:', err);
      this.queueReconnect();
    }
  }

  disconnect() {
    this.shouldReconnect = false;
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout);
      this.reconnectTimeout = null;
    }
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }

  sendEvent(event: string, data: any) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ event, data }));
    }
  }

  private queueReconnect() {
    if (this.reconnectTimeout) return;

    this.reconnectTimeout = window.setTimeout(() => {
      this.reconnectTimeout = null;
      if (this.token && this.shouldReconnect) {
        console.log('Attempting to reconnect WebSocket...');
        this.connect(this.token);
      }
    }, 3000);
  }

  private handleEvent(payload: { event: string; data: any }) {
    const { event, data } = payload;
    switch (event) {
      case 'message.received':
        if (data) {
          useMessageStore.getState().addWebSocketMessage(data);
          useConversationStore.getState().fetchConversations();
        }
        break;
      case 'presence.online':
        if (data && typeof data.user_id === 'number') {
          usePresenceStore.getState().setOnline(data.user_id, true);
        }
        break;
      case 'presence.offline':
        if (data && typeof data.user_id === 'number') {
          usePresenceStore.getState().setOnline(data.user_id, false);
          usePresenceStore.getState().clearUserTyping(data.user_id);
        }
        break;
      case 'typing.start':
        if (data && typeof data.conversation_id === 'number' && typeof data.user_id === 'number') {
          usePresenceStore.getState().setTyping(data.conversation_id, data.user_id, true);
        }
        break;
      case 'typing.stop':
        if (data && typeof data.conversation_id === 'number' && typeof data.user_id === 'number') {
          usePresenceStore.getState().setTyping(data.conversation_id, data.user_id, false);
        }
        break;
      case 'message.read':
        if (data && typeof data.conversation_id === 'number' && typeof data.user_id === 'number' && typeof data.last_read_message_id === 'number') {
          useConversationStore.getState().setConversations(
            useConversationStore.getState().conversations.map((c) => {
              if (c.id === data.conversation_id) {
                return {
                  ...c,
                  participants: c.participants?.map((p) => {
                    if (p.user_id === data.user_id) {
                      return { ...p, last_read_message_id: data.last_read_message_id };
                    }
                    return p;
                  }),
                };
              }
              return c;
            })
          );

          const selected = useConversationStore.getState().selectedConversation;
          if (selected && selected.id === data.conversation_id) {
            useConversationStore.getState().setSelectedConversation({
              ...selected,
              participants: selected.participants?.map((p) => {
                if (p.user_id === data.user_id) {
                  return { ...p, last_read_message_id: data.last_read_message_id };
                }
                return p;
              }),
            });
          }
        }
        break;
      default:
        console.warn('Unhandled WebSocket event:', event);
    }
  }
}

export const webSocketService = new WebSocketService();
