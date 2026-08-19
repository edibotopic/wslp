import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  // Defaults to 8080 (the wslp `serve` command's default port), but can be
  // pointed at a different port via [configurePort] — see main() for how
  // this is wired up to a --port=<n> command-line argument / WSLP_PORT
  // environment variable, so the GUI can talk to a server started with a
  // custom port.
  static String _port = '8080';

  static String get baseUrl => 'http://localhost:$_port';

  static void configurePort(String port) {
    _port = port;
  }

  static Future<List<Map<String, dynamic>>> getDistros() async {
    final response = await http.get(Uri.parse('$baseUrl/api/distros'));

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final distros = data['distros'] as List;
      return distros.map((d) => {
        'name': d['name'] as String,
        'state': d['state'] as String,
        'running': d['running'] as bool,
      }).toList();
    } else {
      throw Exception('Failed to load distros');
    }
  }

  static Future<String> getDefaultDistro() async {
    final response = await http.get(Uri.parse('$baseUrl/api/default'));
    
    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return data['default'] as String;
    } else {
      throw Exception('Failed to get default distro');
    }
  }

  static Future<List<Map<String, String>>> getAvailableDistros() async {
    final response = await http.get(Uri.parse('$baseUrl/api/available'));
    
    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final available = data['available'] as List;
      return available.map((d) => {
        'name': d['name'] as String,
        'friendlyName': d['friendlyName'] as String,
      }).toList();
    } else {
      throw Exception('Failed to load available distros');
    }
  }

  static Future<List<Map<String, dynamic>>> installDistros(List<String> distros) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/install'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'distros': distros}),
    );
    
    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final results = data['results'] as List;
      return results.map((r) => {
        'distro': r['distro'] as String,
        'success': r['success'] as bool,
        'message': r['message'] as String,
        'registered': r['registered'] as bool? ?? false,
      }).toList();
    } else {
      throw Exception('Failed to install distros');
    }
  }

  static Future<List<Map<String, dynamic>>> unregisterDistros(List<String> distros) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/unregister'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'distros': distros}),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final results = data['results'] as List;
      return results.map((r) => {
        'distro': r['distro'] as String,
        'success': r['success'] as bool,
        'message': r['message'] as String,
      }).toList();
    } else {
      throw Exception('Failed to unregister distros');
    }
  }

  static Future<void> setDefaultDistro(String name) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/set-default'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'name': name}),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to set default distro');
    }
  }

  static Future<List<Map<String, dynamic>>> backupDistros(
    List<String> distros, {
    String? customName,
    String? backupDir,
  }) async {
    final body = {
      'distros': distros,
      if (customName != null && customName.isNotEmpty) 'customName': customName,
      if (backupDir != null && backupDir.isNotEmpty) 'backupDir': backupDir,
    };

    final response = await http.post(
      Uri.parse('$baseUrl/api/backup'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode(body),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final results = data['results'] as List;
      return results.map((r) => {
        'distro': r['distro'] as String,
        'success': r['success'] as bool,
        'message': r['message'] as String,
        'filePath': r['filePath'] as String? ?? '',
      }).toList();
    } else {
      throw Exception('Failed to backup distros');
    }
  }

  static Future<List<Map<String, dynamic>>> terminateDistros(List<String> distros) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/terminate'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'distros': distros}),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final results = data['results'] as List;
      return results.map((r) => {
        'distro': r['distro'] as String,
        'success': r['success'] as bool,
        'message': r['message'] as String,
      }).toList();
    } else {
      throw Exception('Failed to terminate distros');
    }
  }

  static Future<void> launchDistro(String name) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/launch'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'name': name}),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to launch distro: ${response.body}');
    }
  }

  static Future<Map<String, dynamic>> renameDistro(String oldName, String newName) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/rename'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'oldName': oldName, 'newName': newName}),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return {
        'oldName': data['oldName'] as String,
        'newName': data['newName'] as String,
        'success': data['success'] as bool,
        'message': data['message'] as String,
      };
    } else {
      throw Exception('Failed to rename distro: ${response.body}');
    }
  }

  static Future<Map<String, dynamic>> getWSLInfo() async {
    final response = await http.get(Uri.parse('$baseUrl/api/wsl-info'));

    if (response.statusCode == 200) {
      return json.decode(response.body) as Map<String, dynamic>;
    } else {
      throw Exception('Failed to get WSL info');
    }
  }

  static Future<Map<String, dynamic>> getDistroInfo(String name) async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/distro-info?name=${Uri.encodeComponent(name)}'),
    );

    if (response.statusCode == 200) {
      return json.decode(response.body) as Map<String, dynamic>;
    } else {
      throw Exception('Failed to get distro info');
    }
  }

  /// Lists Canonical Workshop (canonical/workshop) environments running
  /// inside [name]. Returns an empty list if Workshop isn't installed in
  /// that distro — this is expected, not an error condition.
  static Future<List<Map<String, dynamic>>> getWorkshops(String name) async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/workshops?distro=${Uri.encodeComponent(name)}'),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final workshops = data['workshops'] as List;
      return workshops.map((w) => {
        'project': w['project'] as String,
        'name': w['name'] as String,
        'status': w['status'] as String,
      }).toList();
    } else {
      throw Exception('Failed to load workshops');
    }
  }

  /// Starts or stops the named workshop (project/name identify it, as
  /// returned by [getWorkshops]). Throws with the server's error message
  /// (e.g. "workshop is already started") on failure.
  static Future<void> workshopAction({
    required String distro,
    required String project,
    required String name,
    required String action,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/workshop-action'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'distro': distro,
        'project': project,
        'name': name,
        'action': action,
      }),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to $action workshop: ${response.body}');
    }

    final data = json.decode(response.body);
    if (data['success'] != true) {
      throw Exception(data['message'] as String? ?? 'Failed to $action workshop');
    }
  }

  /// Opens an interactive `workshop shell` session for the named workshop
  /// in a new terminal window (non-blocking).
  static Future<void> shellWorkshop({
    required String distro,
    required String project,
    required String name,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/workshop-shell'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({
        'distro': distro,
        'project': project,
        'name': name,
      }),
    );

    if (response.statusCode != 200) {
      throw Exception('Failed to open workshop shell: ${response.body}');
    }
  }

  /// Gracefully stops the wslp server. The server acknowledges the request
  /// and then shuts itself down asynchronously, so a successful response
  /// here doesn't guarantee it's fully stopped yet, only that it accepted
  /// the request.
  static Future<void> shutdownServer() async {
    final response = await http.post(Uri.parse('$baseUrl/api/shutdown'));

    if (response.statusCode != 200) {
      throw Exception('Failed to stop server: ${response.body}');
    }
  }

  static Future<bool> getUbuntuTelemetry() async {
    final response = await http.get(Uri.parse('$baseUrl/api/ubuntu-telemetry'));

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return data['enabled'] as bool;
    } else {
      throw Exception('Failed to get Ubuntu telemetry status');
    }
  }

  static Future<bool> setUbuntuTelemetry(bool enabled) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/ubuntu-telemetry'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'enabled': enabled}),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      return data['enabled'] as bool;
    } else {
      throw Exception('Failed to set Ubuntu telemetry status');
    }
  }
}
