from kubernetes import client, config
import argparse
import time
import json
from datetime import datetime
import subprocess  # Ensure this line is included
from kubernetes.client import CustomObjectsApi  # Import CustomObjectsApi
from flask import Flask, request, jsonify
from threading import Thread

# Load Kubernetes configuration
config.load_kube_config()

# Kubernetes API clients
core_v1_api = client.CoreV1Api()
apps_v1_api = client.AppsV1Api()

global interval
round = 0

def get_pods_by_deployment_prefix(prefix):
    try:
        # List all deployments in the current namespace
        deployments = apps_v1_api.list_deployment_for_all_namespaces(watch=False)
        matching_deployments = [d for d in deployments.items if d.metadata.name.startswith(prefix)]
        
        # Find all pods associated with matching deployments
        pods = core_v1_api.list_pod_for_all_namespaces(watch=False)
        deployment_pod_map = {}
        for deployment in matching_deployments:
            filtered_pods = []
            for pod in pods.items:
                if pod.metadata.owner_references:
                    for owner in pod.metadata.owner_references:
                        if owner.kind == "ReplicaSet" and owner.name.startswith(deployment.metadata.name):
                            filtered_pods.append(pod)
            deployment_pod_map[deployment.metadata.name] = filtered_pods
        
        return deployment_pod_map
    except client.rest.ApiException as e:
        print(f"Exception when calling Kubernetes API: {e}")
        return {}

def get_pod_resources(pod):
    containers = pod.spec.containers
    resource_info = []
    for container in containers:
        resources = container.resources
        resource_info.append({
            "container_name": container.name,
            "requests": resources.requests or {},
            "limits": resources.limits or {},
        })
    return resource_info

def get_pod_actual_usage(pod_name, namespace):
    try:
        # Use kubectl top to fetch actual usage (requires metrics-server installed)
        result = subprocess.run(
            ["kubectl", "top", "pod", pod_name, "-n", namespace, "--no-headers"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        if result.returncode == 0:
            usage_data = result.stdout.strip().split()
            if len(usage_data) >= 3:
                return {
                    "cpu": usage_data[1],
                    "memory": usage_data[2],
                }
        else:
            print(f"Error fetching actual usage for pod {pod_name}: {result.stderr}")
            return {}
    except Exception as e:
        print(f"Error running kubectl top for pod {pod_name}: {e}")
        return {}

def get_json_data():
    deployment_prefix = "open5gs-upf"
    deployment_pod_map = get_pods_by_deployment_prefix(deployment_prefix)
    metrics = {}

    for deployment_name, pods in deployment_pod_map.items():
        # Initialize deployment entry in metrics
        metrics[deployment_name] = {
            "requests": {},
            "limits": {},
            "actual": {}
        }

        for pod in pods:
            pod_name = pod.metadata.name
            namespace = pod.metadata.namespace

            # Resource requests and limits
            resource_info = get_pod_resources(pod)
            for container_info in resource_info:
                metrics[deployment_name]["requests"] = container_info["requests"]
                metrics[deployment_name]["limits"] = container_info["limits"]

            # Actual resource usage
            actual_usage = get_pod_actual_usage(pod_name, namespace)
            metrics[deployment_name]["actual"] = actual_usage
    
    return json.dumps(metrics)

def send_to_kube(state):
    custom_objects_api = CustomObjectsApi()  # Use CustomObjectsApi to interact with CRDs
    try:
        # Retrieve the custom resource
        observe = custom_objects_api.get_namespaced_custom_object(
            group="lupus.gawor.io",
            version="v1",
            namespace='default',
            plural="elements",
            name='lola-lola'
        )
        
        # Update the `status.input` field with the state
        observe_status = observe.get('status', {})
        observe_status['input'] = json.loads(state)  # Convert JSON string to an object
        observe_status['lastUpdated'] = datetime.utcnow().isoformat() + "Z"  # Proper ISO 8601 format

        # Update the custom resource's status
        custom_objects_api.patch_namespaced_custom_object_status(
            group="lupus.gawor.io",
            version="v1",
            namespace='default',
            plural="elements",
            name='lola-lola',
            body={"status": observe_status}  # Send only the `status` field
        )
        print("Updated Kubernetes custom resource status successfully.")
    except Exception as e:
        print(f"Error updating custom resource: {e}")

def periodic_task():
    json_data = get_json_data()
    timestamp = datetime.utcnow().strftime('%Y/%m/%d %H:%M:%S')
    global round 
    round = round + 1
    print(timestamp + " Round: " + str(round) + "\n" + json_data)
    send_to_kube(json_data)

app = Flask(__name__)

@app.route('/api/interval', methods=['POST'])
def update_interval():
    global interval
    try:
        data = request.get_json()
        if "value" in data and isinstance(data["value"], int) and data["value"] > 0:
            interval = data["value"]
            return jsonify({"message": "Interval updated successfully", "new_interval": interval}), 200
        else:
            return jsonify({"error": "Invalid input. 'value' must be a positive integer."}), 400
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Periodic K8s metrics fetcher')
    parser.add_argument('--interval', type=int, default=60, help='Interval in seconds for periodic task')
    args = parser.parse_args()
    
    # Declare that you are modifying the global variable
    interval = args.interval

    # Start Flask server in a separate thread
    flask_thread = Thread(target=lambda: app.run(host='0.0.0.0', port=9000, debug=False, use_reloader=False))
    flask_thread.daemon = True
    flask_thread.start()
    
    while True:
        periodic_task()
        time.sleep(interval)