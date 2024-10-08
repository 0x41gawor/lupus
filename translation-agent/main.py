import logging
from flask import Flask, request, jsonify
from kubernetes import client, config
from datetime import datetime

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
        logging.info(f"Got data: {data}")
        update_monitor_status(data)  # Call the function to update the Kubernetes resource
        return jsonify({'message': 'Monitor status updated successfully'}), 200
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=4141)
