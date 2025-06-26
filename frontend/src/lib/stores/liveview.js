import { writable, derived, get } from 'svelte/store';
import { browser } from '$app/environment';
import { page } from '$app/stores';

/**
 * LiveView-like real-time store for SvelteKit + WebSockets
 * Provides Phoenix LiveView style real-time updates with automatic reconnection
 */

// Connection states
const CONNECTION_STATES = {
  DISCONNECTED: 'disconnected',
  CONNECTING: 'connecting',
  CONNECTED: 'connected',
  RECONNECTING: 'reconnecting',
  ERROR: 'error'
};

// Default configuration
const DEFAULT_CONFIG = {
  url: 'ws://localhost:8081/ws',
  reconnectInterval: 1000,
  maxReconnectAttempts: 10,
  heartbeatInterval: 30000,
  debug: false
};

// Internal stores
const connectionState = writable(CONNECTION_STATES.DISCONNECTED);
const lastError = writable(null);
const reconnectAttempts = writable(0);
const liveData = writable({});
const subscribers = writable(new Map());

// WebSocket connection
let ws = null;
let heartbeatTimer = null;
let reconnectTimer = null;
let config = { ...DEFAULT_CONFIG };

/**
 * Enhanced LiveView store with Phoenix-like capabilities
 */
function createLiveView() {
  // Public reactive stores
  const stores = {
    // Connection status
    connected: derived(connectionState, $state => $state === CONNECTION_STATES.CONNECTED),
    connecting: derived(connectionState, $state => $state === CONNECTION_STATES.CONNECTING),
    connectionState,

    // Error handling
    error: lastError,

    // Live data - reactive to all real-time updates
    data: liveData,

    // Individual topic subscriptions
    topics: writable(new Map())
  };

  /**
   * Initialize WebSocket connection
   */
  function connect(userConfig = {}) {
    if (!browser) return;

    config = { ...DEFAULT_CONFIG, ...userConfig };

    if (ws?.readyState === WebSocket.OPEN) {
      if (config.debug) console.log('LiveView: Already connected');
      return;
    }

    connectionState.set(CONNECTION_STATES.CONNECTING);
    lastError.set(null);

    try {
      ws = new WebSocket(config.url);
      setupWebSocketHandlers();

      if (config.debug) console.log('LiveView: Connecting to', config.url);
    } catch (error) {
      handleConnectionError(error);
    }
  }

  /**
   * Setup WebSocket event handlers
   */
  function setupWebSocketHandlers() {
    ws.onopen = () => {
      connectionState.set(CONNECTION_STATES.CONNECTED);
      reconnectAttempts.set(0);
      lastError.set(null);

      startHeartbeat();
      resubscribeToTopics();

      if (config.debug) console.log('LiveView: Connected');
    };

    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        handleMessage(message);
      } catch (error) {
        console.error('LiveView: Failed to parse message', error);
      }
    };

    ws.onclose = (event) => {
      connectionState.set(CONNECTION_STATES.DISCONNECTED);
      stopHeartbeat();

      if (config.debug) console.log('LiveView: Disconnected', event.code, event.reason);

      // Auto-reconnect unless it was a manual close
      if (event.code !== 1000) {
        scheduleReconnect();
      }
    };

    ws.onerror = (error) => {
      handleConnectionError(error);
    };
  }

  /**
   * Handle incoming WebSocket messages
   */
  function handleMessage(message) {
    const { type, topic, event, payload, timestamp } = message;

    if (config.debug) console.log('LiveView: Message received', message);

    switch (type) {
      case 'heartbeat':
        send({ type: 'heartbeat_ack' });
        break;

      case 'broadcast':
        handleBroadcast(topic, event, payload, timestamp);
        break;

      case 'update':
        handleUpdate(topic, payload, timestamp);
        break;

      case 'error':
        lastError.set(payload);
        break;

      default:
        // Custom event handlers
        notifySubscribers(type, message);
    }
  }

  /**
   * Handle broadcast messages (LiveView-style)
   */
  function handleBroadcast(topic, event, payload, timestamp) {
    // Update live data store
    liveData.update(data => ({
      ...data,
      [topic]: {
        ...data[topic],
        lastEvent: event,
        lastPayload: payload,
        lastUpdate: timestamp || Date.now()
      }
    }));

    // Notify topic subscribers
    const topicSubs = get(subscribers).get(topic) || new Set();
    topicSubs.forEach(callback => {
      try {
        callback({ event, payload, timestamp, topic });
      } catch (error) {
        console.error('LiveView: Subscriber error', error);
      }
    });
  }

  /**
   * Handle data updates
   */
  function handleUpdate(topic, payload, timestamp) {
    liveData.update(data => ({
      ...data,
      [topic]: {
        ...data[topic],
        ...payload,
        lastUpdate: timestamp || Date.now()
      }
    }));
  }

  /**
   * Notify event subscribers
   */
  function notifySubscribers(event, message) {
    const eventSubs = get(subscribers).get(event) || new Set();
    eventSubs.forEach(callback => {
      try {
        callback(message);
      } catch (error) {
        console.error('LiveView: Event subscriber error', error);
      }
    });
  }

  /**
   * Send message to server
   */
  function send(message) {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
      if (config.debug) console.log('LiveView: Message sent', message);
      return true;
    } else {
      if (config.debug) console.warn('LiveView: Cannot send, not connected', message);
      return false;
    }
  }

  /**
   * Subscribe to a topic (LiveView-style)
   */
  function subscribe(topic, callback) {
    subscribers.update(subs => {
      if (!subs.has(topic)) {
        subs.set(topic, new Set());
      }
      subs.get(topic).add(callback);
      return subs;
    });

    // Send subscription message to server
    send({
      type: 'subscribe',
      topic,
      timestamp: Date.now()
    });

    // Return unsubscribe function
    return () => unsubscribe(topic, callback);
  }

  /**
   * Unsubscribe from topic
   */
  function unsubscribe(topic, callback) {
    subscribers.update(subs => {
      const topicSubs = subs.get(topic);
      if (topicSubs) {
        topicSubs.delete(callback);
        if (topicSubs.size === 0) {
          subs.delete(topic);
          // Notify server
          send({
            type: 'unsubscribe',
            topic,
            timestamp: Date.now()
          });
        }
      }
      return subs;
    });
  }

  /**
   * Push event to server (LiveView-style)
   */
  function push(topic, event, payload = {}) {
    return send({
      type: 'push',
      topic,
      event,
      payload,
      timestamp: Date.now()
    });
  }

  /**
   * Join a topic (LiveView-style presence)
   */
  function join(topic, payload = {}) {
    return send({
      type: 'join',
      topic,
      payload,
      timestamp: Date.now()
    });
  }

  /**
   * Leave a topic
   */
  function leave(topic) {
    send({
      type: 'leave',
      topic,
      timestamp: Date.now()
    });

    // Clean up local subscriptions
    subscribers.update(subs => {
      subs.delete(topic);
      return subs;
    });
  }

  /**
   * Handle connection errors
   */
  function handleConnectionError(error) {
    connectionState.set(CONNECTION_STATES.ERROR);
    lastError.set(error.message || 'Connection failed');

    if (config.debug) console.error('LiveView: Connection error', error);

    scheduleReconnect();
  }

  /**
   * Schedule reconnection attempt
   */
  function scheduleReconnect() {
    const attempts = get(reconnectAttempts);

    if (attempts >= config.maxReconnectAttempts) {
      if (config.debug) console.log('LiveView: Max reconnect attempts reached');
      return;
    }

    const delay = Math.min(config.reconnectInterval * Math.pow(2, attempts), 30000);

    connectionState.set(CONNECTION_STATES.RECONNECTING);
    reconnectAttempts.update(n => n + 1);

    if (config.debug) console.log(`LiveView: Reconnecting in ${delay}ms (attempt ${attempts + 1})`);

    reconnectTimer = setTimeout(() => {
      connect();
    }, delay);
  }

  /**
   * Resubscribe to all topics after reconnection
   */
  function resubscribeToTopics() {
    const subs = get(subscribers);
    subs.forEach((callbacks, topic) => {
      if (callbacks.size > 0) {
        send({
          type: 'subscribe',
          topic,
          timestamp: Date.now()
        });
      }
    });
  }

  /**
   * Start heartbeat to keep connection alive
   */
  function startHeartbeat() {
    stopHeartbeat();
    heartbeatTimer = setInterval(() => {
      send({ type: 'heartbeat', timestamp: Date.now() });
    }, config.heartbeatInterval);
  }

  /**
   * Stop heartbeat
   */
  function stopHeartbeat() {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer);
      heartbeatTimer = null;
    }
  }

  /**
   * Disconnect WebSocket
   */
  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }

    stopHeartbeat();

    if (ws) {
      ws.close(1000, 'Manual disconnect');
      ws = null;
    }

    connectionState.set(CONNECTION_STATES.DISCONNECTED);
  }

  /**
   * Get current data for a topic
   */
  function getTopicData(topic) {
    return get(liveData)[topic] || {};
  }

  /**
   * Create a derived store for a specific topic
   */
  function createTopicStore(topic, initialData = {}) {
    return derived(liveData, $data => $data[topic] || initialData);
  }

  // Public API
  return {
    // Stores
    ...stores,

    // Connection management
    connect,
    disconnect,

    // Messaging
    send,
    push,

    // Topic management
    subscribe,
    unsubscribe,
    join,
    leave,

    // Data access
    getTopicData,
    createTopicStore,

    // Utilities
    isConnected: () => get(connectionState) === CONNECTION_STATES.CONNECTED,

    // Configuration
    configure: (newConfig) => {
      config = { ...config, ...newConfig };
    }
  };
}

// Create singleton instance
export const liveview = createLiveView();

// Auto-connect in browser
if (browser) {
  // Auto-connect when the store is imported
  liveview.connect();

  // Disconnect on page unload
  window.addEventListener('beforeunload', () => {
    liveview.disconnect();
  });
}

// Helper functions for common patterns
export function useLiveTopic(topic, initialData = {}) {
  return liveview.createTopicStore(topic, initialData);
}

export function useLiveSubscription(topic, callback) {
  let unsubscribe;

  return {
    subscribe: () => {
      unsubscribe = liveview.subscribe(topic, callback);
    },
    unsubscribe: () => {
      if (unsubscribe) unsubscribe();
    }
  };
}

// Export for direct usage
export default liveview;
