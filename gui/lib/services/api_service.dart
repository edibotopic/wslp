import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  static const String baseUrl = 'http://localhost:8080';

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
}
