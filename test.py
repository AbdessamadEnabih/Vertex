import socket

# Use localhost (127.0.0.1) or host.docker.internal for Docker-based setups
HOST = "127.0.0.1"  # or 'host.docker.internal' for Docker on Windows/Mac
PORT = 6380


def send_command(s: socket.socket, command):
    try:
        s.sendall(command.encode() + b"\n")
        response = s.recv(1024).decode().strip()
        print(response.lower())
        print("Response:", response)

        # Check if the command was successful
        if response.lower() == "ok":
            print("Command executed successfully.")
        else:
            print("Command failed:", response)

    except Exception as e:
        print("Error:", e)


# Example interactions
commands = [
    ("SET mykey4 myvalue85665", "Set key-value pair"),
    ("SET mykey4df myv85665", "Set key-value pair"),
    ("SET mykey4qd myvalue865", "Set key-value pair"),
    ("ALL", "Retrieve all key-value pairs"),
]

try:
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        print("Connection established.")
        for cmd, description in commands:
            send_command(s, cmd)
except ConnectionRefusedError:
    print("Connection refused. Make sure the server is running.")
except Exception as e:
    print(f"Error: {e}")
