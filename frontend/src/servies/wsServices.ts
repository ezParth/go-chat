let socket: WebSocket | null = null

export const initWebSocket = (url: string): WebSocket => {
  if (!socket || socket.readyState === WebSocket.CLOSED) {
    socket = new WebSocket(url)
  }
  return socket
}

export const getWebSocket = (): WebSocket | null => {
  return socket
}

export const closeWebSocket = () => {
  if (socket) {
    socket.close()
    socket = null
  }
}
