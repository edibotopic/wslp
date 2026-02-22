import 'dart:convert';
import 'package:http/http.dart' as http;

class ApiService {
  static const String baseUrl = 'http://localhost:8080';

  static Future<List<String>> getDistros() async {
    final response = await http.get(Uri.parse('$baseUrl/api/distros'));
    
    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      final distros = data['distros'] as List;
      return distros.map((d) => d['name'] as String).toList();
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

  static Future<void> unregisterDistro(String name) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/unregister'),
      headers: {'Content-Type': 'application/json'},
      body: json.encode({'name': name}),
    );
    
    if (response.statusCode != 200) {
      throw Exception('Failed to unregister distro');
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
}
