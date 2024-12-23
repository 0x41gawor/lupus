from flask import Flask, request, jsonify

app = Flask(__name__)

# Hardcoded default values
DEFAULT_VALUES = {
    "requests": {
        "memory": "128Mi",
        "cpu": "100m"
    },
    "limits": {
        "memory": "256Mi",
        "cpu": "250m"
    }
}

# ------------------ Parsing Helpers ------------------ #
def parse_cpu(cpu_str: str) -> int:
    """
    Convert a CPU string to an integer representing millicores.
    
    Examples:
      "300m"  -> 300
      "2"     -> 2000  (interpreted as 2 CPU cores -> 2000 millicores)
      "1.5"   -> 1500
    """
    cpu_str = cpu_str.strip().lower()
    if cpu_str.endswith("m"):
        # e.g. "300m" -> 300
        numeric_part = cpu_str[:-1]  # remove "m"
        return int(float(numeric_part))
    else:
        # e.g. "2" -> 2000, "1.5" -> 1500
        return int(float(cpu_str) * 1000)


def parse_memory(mem_str: str) -> int:
    """
    Convert a memory string to an integer representing megabytes (MB).
    
    Examples:
      "300Mi" -> 300
      "1Gi"   -> 1024
      "1024Ki"-> 1
      "512"   -> 512   (no unit, assume MB)
    """
    mem_str = mem_str.strip()
    lower_str = mem_str.lower()

    if lower_str.endswith("mi"):
        # e.g. "300Mi"
        numeric_part = lower_str.replace("mi", "")
        return int(float(numeric_part))
    elif lower_str.endswith("gi"):
        # e.g. "2Gi" -> 2 * 1024 = 2048
        numeric_part = lower_str.replace("gi", "")
        return int(float(numeric_part) * 1024)
    elif lower_str.endswith("ki"):
        # e.g. "1024Ki" -> 1024 / 1024 = 1
        numeric_part = lower_str.replace("ki", "")
        return int(float(numeric_part) / 1024)
    else:
        # e.g. "512" -> 512
        return int(float(lower_str))

# ------------------ Comparison Helpers ------------------ #

def is_higher_cpu(actual_cpu_str, default_cpu_str) -> bool:
    """Return True if actual_cpu_str is higher than default_cpu_str (in millicores)."""
    return parse_cpu(actual_cpu_str) > parse_cpu(default_cpu_str)

def is_higher_memory(actual_mem_str, default_mem_str) -> bool:
    """Return True if actual_mem_str is higher than default_mem_str (in MB)."""
    return parse_memory(actual_mem_str) > parse_memory(default_mem_str)

def is_default_cpu(actual_cpu_str, default_cpu_str) -> bool:
    """Return True if actual == default (in millicores)."""
    return parse_cpu(actual_cpu_str) == parse_cpu(default_cpu_str)

def is_default_memory(actual_mem_str, default_mem_str) -> bool:
    """Return True if actual == default (in MB)."""
    return parse_memory(actual_mem_str) == parse_memory(default_mem_str)

# ------------------ Core Logic ------------------ #
def determine_point(data):
    """
    Determine the operational point type:
    
    1. NORMAL
       - actual < default_values['requests'] for both cpu, memory

    2. NORMAL_TO_CRITICAL
       - requests == default_values['requests']
       - limits == default_values['limits']
       - actual > default_values['requests'] for cpu or memory

    3. CRITICAL
       - actual, requests, limits are ALL higher than the defaults
         (at least one field of each is higher than the default)

    4. CRITICAL_TO_NORMAL
       - requests, limits are above their defaults
       - actual is below the default (requests) for both cpu, memory
    """
    # Extract input values
    req_cpu_str    = data["requests"]["cpu"]
    req_mem_str    = data["requests"]["memory"]
    lim_cpu_str    = data["limits"]["cpu"]
    lim_mem_str    = data["limits"]["memory"]
    act_cpu_str    = data["actual"]["cpu"]
    act_mem_str    = data["actual"]["memory"]

    # Extract defaults
    def_req_cpu = DEFAULT_VALUES["requests"]["cpu"]
    def_req_mem = DEFAULT_VALUES["requests"]["memory"]
    def_lim_cpu = DEFAULT_VALUES["limits"]["cpu"]
    def_lim_mem = DEFAULT_VALUES["limits"]["memory"]

    # 1. NORMAL: actual < default_values.requests (cpu & mem)
    condition_normal = (not is_higher_cpu(act_cpu_str, def_req_cpu) and
                        not is_higher_memory(act_mem_str, def_req_mem))

    # 2. NORMAL_TO_CRITICAL:
    #    - requests == default (cpu & mem)
    #    - limits == default (cpu & mem)
    #    - actual > default.requests for cpu or memory
    condition_normal_to_critical = (
        is_default_cpu(req_cpu_str, def_req_cpu) and
        is_default_memory(req_mem_str, def_req_mem) and
        is_default_cpu(lim_cpu_str, def_lim_cpu) and
        is_default_memory(lim_mem_str, def_lim_mem) and
        (
            is_higher_cpu(act_cpu_str, def_req_cpu) or
            is_higher_memory(act_mem_str, def_req_mem)
        )
    )

    # 3. CRITICAL:
    #    - actual is higher than default requests (cpu or mem)
    #    - requests is higher than default requests (cpu or mem)
    #    - limits is higher than default limits (cpu or mem)
    condition_critical = (
        (is_higher_cpu(act_cpu_str, def_req_cpu) or is_higher_memory(act_mem_str, def_req_mem)) and
        (is_higher_cpu(req_cpu_str, def_req_cpu) or is_higher_memory(req_mem_str, def_req_mem)) and
        (is_higher_cpu(lim_cpu_str, def_lim_cpu) or is_higher_memory(lim_mem_str, def_lim_mem))
    )

    # 4. CRITICAL_TO_NORMAL:
    #    - requests, limits are above defaults
    #    - actual is below default.requests for cpu & mem
    #    => requests > default, limits > default, actual < default
    condition_critical_to_normal = (
        (is_higher_cpu(req_cpu_str, def_req_cpu) or is_higher_memory(req_mem_str, def_req_mem)) and
        (is_higher_cpu(lim_cpu_str, def_lim_cpu) or is_higher_memory(lim_mem_str, def_lim_mem)) and
        (not is_higher_cpu(act_cpu_str, def_req_cpu)) and
        (not is_higher_memory(act_mem_str, def_req_mem))
    )

    # Decide which point applies in priority order
    if condition_critical:
        return "CRITICAL"
    elif condition_normal_to_critical:
        return "NORMAL_TO_CRITICAL"
    elif condition_critical_to_normal:
        return "CRITICAL_TO_NORMAL"
    elif condition_normal:
        return "NORMAL"
    else:
        # Fallback if no exact condition matched
        return "NORMAL"

def _int_mul(value: int, factor: float) -> int:
    """
    Multiplies an integer by a float factor and returns an int.
    E.g. _int_mul(100, 1.2) -> 120
    """
    return int(value * factor)

def _cpu_to_str(millicores: int) -> str:
    """
    Format an integer millicore value back into a string with the 'm' suffix.
    E.g. 120 -> '120m'
    """
    return f"{millicores}m"

def _mem_to_str(megabytes: int) -> str:
    """
    Format an integer MB value back into a string with the 'Mi' suffix.
    E.g. 256 -> '256Mi'
    """
    return f"{megabytes}Mi"

def generate_spec(actual: dict) -> dict:
    """
    Given something like:
      actual = {"cpu": "110m", "memory": "18Mi"}
    Return a dict of the form:
      {
        "spec": {
          "requests": {"cpu": "...", "memory": "..."},
          "limits":   {"cpu": "...", "memory": "..."}
        }
      }
    Where each CPU/memory is either scaled from actual (if above defaults)
    or uses the default value (if at or below defaults).
    """
    # Parse default values
    def_req_cpu = parse_cpu(DEFAULT_VALUES["requests"]["cpu"])
    def_req_mem = parse_memory(DEFAULT_VALUES["requests"]["memory"])
    def_lim_cpu = parse_cpu(DEFAULT_VALUES["limits"]["cpu"])
    def_lim_mem = parse_memory(DEFAULT_VALUES["limits"]["memory"])

    # Parse actual values
    actual_cpu = parse_cpu(actual["cpu"])
    actual_mem = parse_memory(actual["memory"])

    # ----- CPU logic -----
    if actual_cpu > def_req_cpu:
        requests_cpu = _int_mul(actual_cpu, 1.2)
        limits_cpu   = _int_mul(actual_cpu, 2.4)
    else:
        requests_cpu = def_req_cpu
        limits_cpu   = def_lim_cpu

    # ----- Memory logic -----
    if actual_mem > def_req_mem:
        requests_mem = _int_mul(actual_mem, 1.2)
        limits_mem   = _int_mul(actual_mem, 2.4)
    else:
        requests_mem = def_req_mem
        limits_mem   = def_lim_mem

    return {
        "result": {
            "requests": {
                "cpu": _cpu_to_str(requests_cpu),
                "memory": _mem_to_str(requests_mem)
            },
            "limits": {
                "cpu": _cpu_to_str(limits_cpu),
                "memory": _mem_to_str(limits_mem)
            }
        }
    }

# ------------------ Flask Endpoints ------------------ #
@app.route("/v1/data/policy/point", methods=["POST"])
def logic_endpoint():
    data = request.get_json(force=True)
    point = determine_point(data['input'])
    return jsonify({"result": point})

@app.route("/v1/data/policy/spec", methods=["POST"])
def spec_endpoint():
    data = request.get_json(force=True)
    actual = data['input']["actual"]
    spec_obj = generate_spec(actual)
    return jsonify(spec_obj)

@app.route("/v1/data/policy/interval", methods=["POST"])
def interval_endpoint():
    data = request.get_json(force=True)
    point_value = data['input']["point"]

    if point_value in ("NORMAL_TO_CRITICAL", "CRITICAL"):
        interval = "HIGH"
    else:
        interval = "LOW"

    return jsonify({"result": interval})


if __name__ == "__main__":
    # Run the Flask app. You can also set a different port, debug mode, etc.
    app.run(host="0.0.0.0", port=9500, debug=True)
