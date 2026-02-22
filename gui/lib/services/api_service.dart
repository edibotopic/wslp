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
}
