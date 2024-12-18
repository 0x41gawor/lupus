from flask import Flask, request, jsonify
from kubernetes import client, config
from kubernetes.client.rest import ApiException

app = Flask(__name__)

# Load Kubernetes config (use kubeconfig if running locally, in-cluster config if running in a cluster)
try:
    config.load_incluster_config()
except:
    config.load_kube_config()

# Function to patch deployment resources
def patch_deployment_resources(namespace, deployment_name, resource_type, cpu, memory):
    try:
        api_instance = client.AppsV1Api()

        # Define the patch
        patch = {
            "spec": {
                "template": {
                    "spec": {
                        "containers": [
                            {
                                "name": "upf",
                                "resources": {
                                    resource_type.split(".")[-1]: {  # Extract "requests" or "limits"
                                        "cpu": cpu,
                                        "memory": memory
                                    }
                                }
                            }
                        ]
                    }
                }
            }
        }

        # Patch the deployment
        api_response = api_instance.patch_namespaced_deployment(name=deployment_name, namespace=namespace, body=patch)
        return api_response

    except ApiException as e:
        return str(e)

@app.route('/api/set-requests', methods=['POST'])
def set_requests():
    data = request.get_json()
    deployment_name = data.get('name')  # Deployment name is passed in 'name'
    namespace = data.get('namespace', 'open5gs')  # Default namespace if not provided
    cpu = data.get('cpu')
    memory = data.get('memory')

    if not all([deployment_name, cpu, memory]):
        return jsonify({"error": "Missing required fields: name, cpu, memory"}), 400

    response = patch_deployment_resources(namespace, deployment_name, "resources.requests", cpu, memory)
    return jsonify({"response": str(response)})

@app.route('/api/set-limits', methods=['POST'])
def set_limits():
    data = request.get_json()
    deployment_name = data.get('name')  # Deployment name is passed in 'name'
    namespace = data.get('namespace', 'open5gs')  # Default namespace if not provided
    cpu = data.get('cpu')
    memory = data.get('memory')

    if not all([deployment_name, cpu, memory]):
        return jsonify({"error": "Missing required fields: name, cpu, memory"}), 400

    response = patch_deployment_resources(namespace, deployment_name, "resources.limits", cpu, memory)
    return jsonify({"response": str(response)})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=9001)