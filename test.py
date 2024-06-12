import asyncio
import websockets
import sys
import requests

async def send_message(remote_addr):
    # Get list of available users
    try:
        response = requests.get(f'http://localhost:3000/showUsers', headers={'X-Forwarded-For': remote_addr})
        response.raise_for_status()
        print(response.text)
    except Exception as e:
        print(f"Error getting user list: {e}")
        return

    # Connect to WebSocket for sending and receiving messages
    try:
        async with websockets.connect(f'ws://localhost:3000/ws', extra_headers=[('X-Forwarded-For', remote_addr)]) as websocket:
            while True:
                message = input("Enter message to send (type 'quit' to exit): ")
                if message.lower() == 'quit':
                    break
                await websocket.send(message)
                print(f"Sent: {message}")
                response = await websocket.recv()
                print(f"Received: {response}")
    except Exception as e:
        print(f"WebSocket connection error: {e}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python client.py <remote_address>")
        sys.exit(1)

    remote_addr = sys.argv[1]
    asyncio.run(send_message(remote_addr))
