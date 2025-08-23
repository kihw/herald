// Herald.lol WebSocket Client - Real-time Gaming Updates

import { EventEmitter } from 'events';

export interface WebSocketMessage {
  type: string;
  user_id?: string;
  match_id?: string;
  room_id?: string;
  data: any;
  timestamp: string;
  id: string;
}

export interface ClientMessage {
  action: string;
  data?: any;
}

// Message types
export enum MessageType {
  MATCH_UPDATE = 'match_update',
  PERFORMANCE_UPDATE = 'performance_update',
  RANK_UPDATE = 'rank_update',
  FRIEND_ACTIVITY = 'friend_activity',
  LIVE_MATCH = 'live_match',
  COACHING_SUGGESTION = 'coaching_suggestion',
  CHAMPION_MASTERY = 'champion_mastery',
  SYSTEM_NOTIFICATION = 'system_notification',
  ERROR = 'error',
  PING = 'ping',
  PONG = 'pong'
}

// Client actions
export enum ClientAction {
  SUBSCRIBE = 'subscribe',
  UNSUBSCRIBE = 'unsubscribe',
  JOIN_ROOM = 'join_room',
  LEAVE_ROOM = 'leave_room',
  WATCH_MATCH = 'watch_match',
  UNWATCH_MATCH = 'unwatch_match',
  UPDATE_PREFERENCES = 'update_preferences',
  PONG = 'pong',
  GET_STATS = 'get_stats'
}

// Gaming data interfaces
export interface MatchUpdateData {
  game_id: string;
  status: string;
  game_time: number;
  participants: ParticipantData[];
  team_stats: TeamStatsData;
  event_data?: any;
}

export interface ParticipantData {
  summoner_name: string;
  champion_name: string;
  level: number;
  kills: number;
  deaths: number;
  assists: number;
  cs: number;
  gold: number;
  items: number[];
  kda: number;
}

export interface TeamStatsData {
  blue_team: TeamData;
  red_team: TeamData;
}

export interface TeamData {
  kills: number;
  deaths: number;
  assists: number;
  gold: number;
  dragons: number;
  barons: number;
  towers: number;
  inhibitors: number;
}

export interface PerformanceUpdateData {
  user_id: string;
  current_kda: number;
  average_kda: number;
  cs_per_minute: number;
  vision_score: number;
  damage_share: number;
  gold_efficiency: number;
  improvement_suggestion: string;
}

export interface ClientPreferences {
  match_updates: boolean;
  rank_updates: boolean;
  friend_activity: boolean;
  coaching_suggestions: boolean;
  system_notifications: boolean;
}

export interface WebSocketStats {
  total_connections: number;
  active_connections: number;
  messages_per_second: number;
  last_message_time: string;
  average_response_time_ms: number;
}

export class HeraldWebSocket extends EventEmitter {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string;
  private userId: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private heartbeatInterval: NodeJS.Timeout | null = null;
  private isConnecting = false;
  private isIntentionallyClosed = false;
  
  // Connection state
  private subscriptions = new Set<string>();
  private watchedMatches = new Set<string>();
  private joinedRooms = new Set<string>();
  private preferences: ClientPreferences = {
    match_updates: true,
    rank_updates: true,
    friend_activity: true,
    coaching_suggestions: true,
    system_notifications: true
  };
  
  // Gaming analytics tracking
  private messageCount = 0;
  private lastMessageTime: Date | null = null;
  private responseTimeSum = 0;
  private responseTimeCount = 0;
  private performanceMetrics = {
    messagesReceived: 0,
    averageLatency: 0,
    connectionUptime: 0,
    reconnectCount: 0
  };
  
  constructor(userId: string, token: string, baseUrl?: string) {
    super();
    this.userId = userId;
    this.token = token;
    this.url = this.buildWebSocketUrl(baseUrl || process.env.VITE_WS_URL || 'ws://localhost:8080');
  }
  
  private buildWebSocketUrl(baseUrl: string): string {
    const wsUrl = new URL('/ws', baseUrl.replace('http', 'ws'));
    wsUrl.searchParams.append('user_id', this.userId);
    wsUrl.searchParams.append('token', this.token);
    return wsUrl.toString();
  }
  
  // Connection management
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
        resolve();
        return;
      }
      
      this.isConnecting = true;
      this.isIntentionallyClosed = false;
      
      try {
        this.ws = new WebSocket(this.url);
        this.setupEventHandlers(resolve, reject);
      } catch (error) {
        this.isConnecting = false;
        reject(error);
      }
    });
  }
  
  disconnect(): void {
    this.isIntentionallyClosed = true;
    this.stopHeartbeat();
    
    if (this.ws) {
      this.ws.close(1000, 'Client disconnect');
      this.ws = null;
    }
    
    this.emit('disconnected');
  }
  
  private setupEventHandlers(resolve: () => void, reject: (error: any) => void): void {
    if (!this.ws) return;
    
    this.ws.onopen = () => {
      console.log('🎮 Herald.lol WebSocket connected');
      this.isConnecting = false;
      this.reconnectAttempts = 0;
      this.startHeartbeat();
      this.resubscribeAll();
      this.emit('connected');
      resolve();
    };
    
    this.ws.onmessage = (event) => {
      this.handleMessage(event);
    };
    
    this.ws.onclose = (event) => {
      console.log(`WebSocket closed: ${event.code} - ${event.reason}`);
      this.isConnecting = false;
      this.stopHeartbeat();
      
      if (!this.isIntentionallyClosed && this.shouldReconnect(event.code)) {
        this.scheduleReconnect();
      }
      
      this.emit('disconnected', { code: event.code, reason: event.reason });
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      this.isConnecting = false;
      this.emit('error', error);
      
      if (this.reconnectAttempts === 0) {
        reject(error);
      }
    };
  }
  
  private shouldReconnect(code: number): boolean {
    // Don't reconnect on authentication failures or intentional closes
    return code !== 1000 && code !== 1001 && code !== 1002 && code !== 4001;
  }
  
  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      this.emit('reconnectFailed');
      return;
    }
    
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts);
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts + 1}/${this.maxReconnectAttempts})`);
    
    setTimeout(() => {
      this.reconnectAttempts++;
      this.performanceMetrics.reconnectCount++;
      this.connect().catch((error) => {
        console.error('Reconnection failed:', error);
      });
    }, delay);
  }
  
  private handleMessage(event: MessageEvent): void {
    try {
      const message: WebSocketMessage = JSON.parse(event.data);
      this.messageCount++;
      this.lastMessageTime = new Date();
      this.performanceMetrics.messagesReceived++;
      
      // Calculate response time if this is a pong
      if (message.type === MessageType.PONG) {
        const now = Date.now();
        const sentTime = parseInt(message.data?.timestamp || '0');
        if (sentTime > 0) {
          const responseTime = now - sentTime;
          this.responseTimeSum += responseTime;
          this.responseTimeCount++;
          this.performanceMetrics.averageLatency = this.responseTimeSum / this.responseTimeCount;
        }
      }
      
      // Handle ping-pong
      if (message.type === MessageType.PING) {
        this.sendPong();
        return;
      }
      
      // Emit specific event for message type
      this.emit(message.type, message);
      this.emit('message', message);
      
      // Gaming-specific event handling
      this.handleGamingMessage(message);
      
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
      this.emit('parseError', error);
    }
  }
  
  private handleGamingMessage(message: WebSocketMessage): void {
    switch (message.type) {
      case MessageType.MATCH_UPDATE:
        this.emit('matchUpdate', message.data as MatchUpdateData);
        break;
      
      case MessageType.PERFORMANCE_UPDATE:
        this.emit('performanceUpdate', message.data as PerformanceUpdateData);
        break;
      
      case MessageType.RANK_UPDATE:
        this.emit('rankUpdate', message.data);
        // Show notification for rank changes
        if (message.data.new_rank !== message.data.old_rank) {
          this.showRankNotification(message.data);
        }
        break;
      
      case MessageType.COACHING_SUGGESTION:
        this.emit('coachingSuggestion', message.data);
        break;
      
      case MessageType.SYSTEM_NOTIFICATION:
        this.emit('systemNotification', message.data);
        break;
      
      case MessageType.ERROR:
        console.error('WebSocket error message:', message.data);
        this.emit('serverError', message.data);
        break;
    }
  }
  
  // Messaging methods
  private sendMessage(message: ClientMessage): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('WebSocket not connected, cannot send message');
      return;
    }
    
    try {
      this.ws.send(JSON.stringify(message));
    } catch (error) {
      console.error('Error sending WebSocket message:', error);
      this.emit('sendError', error);
    }
  }
  
  private sendPong(): void {
    this.sendMessage({
      action: ClientAction.PONG,
      data: { timestamp: Date.now() }
    });
  }
  
  // Subscription management
  subscribe(): void {
    this.sendMessage({ action: ClientAction.SUBSCRIBE });
    this.subscriptions.add(this.userId);
  }
  
  unsubscribe(): void {
    this.sendMessage({ action: ClientAction.UNSUBSCRIBE });
    this.subscriptions.delete(this.userId);
  }
  
  // Room management
  joinRoom(roomId: string): void {
    this.sendMessage({
      action: ClientAction.JOIN_ROOM,
      data: { room_id: roomId }
    });
    this.joinedRooms.add(roomId);
  }
  
  leaveRoom(roomId: string): void {
    this.sendMessage({
      action: ClientAction.LEAVE_ROOM,
      data: { room_id: roomId }
    });
    this.joinedRooms.delete(roomId);
  }
  
  // Match tracking
  watchMatch(matchId: string): void {
    this.sendMessage({
      action: ClientAction.WATCH_MATCH,
      data: { match_id: matchId }
    });
    this.watchedMatches.add(matchId);
  }
  
  unwatchMatch(matchId: string): void {
    this.sendMessage({
      action: ClientAction.UNWATCH_MATCH,
      data: { match_id: matchId }
    });
    this.watchedMatches.delete(matchId);
  }
  
  // Preferences
  updatePreferences(preferences: Partial<ClientPreferences>): void {
    this.preferences = { ...this.preferences, ...preferences };
    this.sendMessage({
      action: ClientAction.UPDATE_PREFERENCES,
      data: this.preferences
    });
  }
  
  // Stats
  getStats(): void {
    this.sendMessage({ action: ClientAction.GET_STATS });
  }
  
  // Heartbeat
  private startHeartbeat(): void {
    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.sendMessage({
          action: ClientAction.PONG,
          data: { timestamp: Date.now() }
        });
      }
    }, 30000); // 30 second heartbeat
  }
  
  private stopHeartbeat(): void {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }
  
  // Resubscribe to all subscriptions after reconnect
  private resubscribeAll(): void {
    if (this.subscriptions.size > 0) {
      this.subscribe();
    }
    
    for (const roomId of this.joinedRooms) {
      this.joinRoom(roomId);
    }
    
    for (const matchId of this.watchedMatches) {
      this.watchMatch(matchId);
    }
    
    if (Object.values(this.preferences).some(Boolean)) {
      this.updatePreferences(this.preferences);
    }
  }
  
  // Notifications
  private showRankNotification(data: any): void {
    if ('Notification' in window && Notification.permission === 'granted') {
      const isPromotion = this.isRankPromotion(data.old_rank, data.new_rank);
      const title = isPromotion ? 'Rank Promotion! 🎉' : 'Rank Change';
      const body = `${data.old_rank} → ${data.new_rank} (${data.lp} LP)`;
      
      new Notification(title, {
        body,
        icon: '/herald-logo.png',
        badge: '/herald-badge.png',
        tag: 'rank-update',
        renotify: true
      });
    }
  }
  
  private isRankPromotion(oldRank: string, newRank: string): boolean {
    const ranks = ['IRON', 'BRONZE', 'SILVER', 'GOLD', 'PLATINUM', 'EMERALD', 'DIAMOND', 'MASTER', 'GRANDMASTER', 'CHALLENGER'];
    const divisions = ['IV', 'III', 'II', 'I'];
    
    const parseRank = (rank: string) => {
      const parts = rank.split(' ');
      const tier = ranks.indexOf(parts[0]);
      const division = divisions.indexOf(parts[1] || 'I');
      return tier * 4 + (3 - division);
    };
    
    return parseRank(newRank) > parseRank(oldRank);
  }
  
  // Getters
  get isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
  
  get connectionState(): number {
    return this.ws?.readyState || WebSocket.CLOSED;
  }
  
  get metrics() {
    return {
      ...this.performanceMetrics,
      connectionUptime: this.lastMessageTime ? Date.now() - this.lastMessageTime.getTime() : 0,
      messageCount: this.messageCount,
      averageResponseTime: this.performanceMetrics.averageLatency
    };
  }
}

// Export singleton instance for global use
let heraldWebSocket: HeraldWebSocket | null = null;

export function createWebSocketInstance(userId: string, token: string, baseUrl?: string): HeraldWebSocket {
  if (heraldWebSocket) {
    heraldWebSocket.disconnect();
  }
  
  heraldWebSocket = new HeraldWebSocket(userId, token, baseUrl);
  return heraldWebSocket;
}

export function getWebSocketInstance(): HeraldWebSocket | null {
  return heraldWebSocket;
}

// React hook for WebSocket connection
export function useWebSocket(userId?: string, token?: string) {
  const [isConnected, setIsConnected] = React.useState(false);
  const [connectionError, setConnectionError] = React.useState<Error | null>(null);
  const [reconnectCount, setReconnectCount] = React.useState(0);
  
  React.useEffect(() => {
    if (!userId || !token) return;
    
    const ws = createWebSocketInstance(userId, token);
    
    const handleConnect = () => {
      setIsConnected(true);
      setConnectionError(null);
    };
    
    const handleDisconnect = () => {
      setIsConnected(false);
    };
    
    const handleError = (error: Error) => {
      setConnectionError(error);
    };
    
    const handleReconnect = () => {
      setReconnectCount(prev => prev + 1);
    };
    
    ws.on('connected', handleConnect);
    ws.on('disconnected', handleDisconnect);
    ws.on('error', handleError);
    ws.on('reconnectFailed', handleReconnect);
    
    // Connect
    ws.connect().catch(setConnectionError);
    
    return () => {
      ws.off('connected', handleConnect);
      ws.off('disconnected', handleDisconnect);
      ws.off('error', handleError);
      ws.off('reconnectFailed', handleReconnect);
      ws.disconnect();
    };
  }, [userId, token]);
  
  return {
    isConnected,
    connectionError,
    reconnectCount,
    webSocket: heraldWebSocket
  };
}

declare global {
  var React: typeof import('react');
}