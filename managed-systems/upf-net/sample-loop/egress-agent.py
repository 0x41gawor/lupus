# command_forwarder.py
from flask import Flask, request, jsonify
import requests

# Initialize Flask app
app = Flask(__name__)

# Function to handle /api/commands endpoint
@app.route('/api/commands', methods=['POST'])
def handle_commands():
    try:
        data = request.json
        commands = data.get('commands', [])
        for command in commands:
            # Forward each command to the external API
            response = requests.post('http://192.168.56.111:5000/api/move', json=command)
            if response.status_code != 200:
                print(f"Failed to send command: {command}")
        return jsonify({"status": "success"}), 200
    except Exception as e:
        print(f"Error handling commands: {e}")
        return jsonify({"status": "error", "message": str(e)}), 500

if __name__ == "__main__":
    # Run the Flask server
    app.run(host='0.0.0.0', port=5001)