import logging
from flask import Flask, request, jsonify
from kubernetes import client, config
from datetime import datetime
import requests

# Load the kubeconfig (works with Minikube or other clusters)
config.load_kube_config()

# Create an API client for the custom resource
k8s_client = client.CustomObjectsApi()

app = Flask(__name__)

# Disable Flask request logging
log = logging.getLogger('werkzeug')
log.setLevel(logging.ERROR)

# Configure the logger
logging.basicConfig(
    format='%(asctime)s %(message)s',  # Specify the log format
    datefmt='%Y/%m/%d %H:%M:%S',       # Specify the date format
    level=logging.INFO                  # Set the log level
)

def update_monitor_status(data):
    # Fetch the custom resource
    monitor = k8s_client.get_namespaced_custom_object(
        group="lupus.gawor.io",
        version="v1",
        namespace='default',
        plural="monitors",
        name='adam'
    )

    # Update the status fields with received data
    monitor['status'] = {
        'gdansk': data.get('Gdansk', monitor['status'].get('gdansk', 0)),
        'krakow': data.get('Krakow', monitor['status'].get('krakow', 0)),
        'poznan': data.get('Poznan', monitor['status'].get('poznan', 0)),
        'warsaw': data.get('Warsaw', monitor['status'].get('warsaw', 0)),
        'lastUpdated': datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')
    }

    # Update the custom resource status
    k8s_client.patch_namespaced_custom_object_status(
        group="lupus.gawor.io",
        version="v1", 
        namespace='default',
        plural="monitors",
        name='adam',
        body=monitor
    )


@app.route('/api/monitor', methods=['POST'])
def monitor_update():
    try:
        data = request.get_json()  # Get the JSON data from the request
        formatted_data = "{" + ", ".join([f"{city}: {value: >2}" for city, value in data.items()]) + "}"
        logging.info(f"Got data: {formatted_data}")
        update_monitor_status(data)  # Call the function to update the Kubernetes resource
        return jsonify({'message': 'Monitor status updated successfully'}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500

def send_move_command(data):
    try:
        # Forward the move command to localhost:4040/api/move
        response = requests.post('http://localhost:4040/api/move', json=data)
        response.raise_for_status()  # Raise an exception for HTTP errors

        # Check if the response is in JSON format
        try:
            return response.json()  # Try to return the response as JSON
        except ValueError:
            # If the response is not JSON, return the text content instead
            return {'message': response.text}  # Return the text content in a dictionary
    except requests.exceptions.RequestException as e:
        logging.error(f"Error forwarding move command: {e}")
        return {'error': str(e)}


@app.route('/api/move', methods=['POST'])
def move_command():
    try:
        data = request.get_json()
        logging.info(f"Got move command: {data}")
        result = send_move_command(data)  # Forward the move command
        return jsonify(result), 200 if 'error' not in result else 500
    except Exception as e:
        logging.error(f"Error in /api/move: {e}")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=4141)
