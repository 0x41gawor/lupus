# periodic_k8s_updater.py
import requests
import time
from kubernetes import client, config
import argparse
from datetime import datetime

# Global state shared across threads
state = None

# Function for periodic HTTP requests and K8s custom resource update
def periodic_task(interval, k8s_client):
    global state
    while True:
        try:
            # Step 1: Fetch data from the external API
            response = requests.get('http://127.0.0.1:7000/api/data')
            if response.status_code == 200:
                state = response.json()
                print(f"Received state: {state}")

                # Step 2: Update Kubernetes custom resource
                observe = k8s_client.get_namespaced_custom_object(
                    group="lupus.gawor.io",
                    version="v1",
                    namespace='default',
                    plural="elements",
                    name='lola-observe1'
                )
                
                # Update the `status.input` field with the state
                observe_status = observe.get('status', {})
                observe_status['input'] = state
                observe_status['lastUpdated'] = datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')

                observe['status'] = observe_status
                

                # Patch the custom resource status
                k8s_client.patch_namespaced_custom_object_status(
                    group="lupus.gawor.io",
                    version="v1",
                    namespace='default',
                    plural="elements",
                    name='lola-observe1',
                    body=observe
                )
                print("Updated Kubernetes custom resource status successfully.")
            else:
                print(f"Failed to fetch data: {response.status_code}")
        except Exception as e:
            print(f"Error during periodic task: {e}")

        # Wait for the specified interval before repeating
        time.sleep(interval)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Periodic HTTP and K8s updater')
    parser.add_argument('-interval', type=int, default=60, help='Interval in seconds for periodic task')
    args = parser.parse_args()

    # Load Kubernetes config and create client
    config.load_kube_config()
    k8s_client = client.CustomObjectsApi()

    # Run the periodic task with the specified interval
    periodic_task(args.interval, k8s_client)