import asyncio
import websockets

async def send_message():
    async with websockets.connect('ws://localhost:3000/ws') as websocket:
        while True:
            message = input("Enter message to send: ")
            await websocket.send(message)
            print(f"Sent: {message}")
            response = await websocket.recv()
            print(f"Received: {response}")

# Run the event loop
asyncio.get_event_loop().run_until_complete(send_message())
