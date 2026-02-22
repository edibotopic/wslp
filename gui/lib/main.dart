import 'package:flutter/material.dart';
import 'services/api_service.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
        colorSchemeSeed: const Color.fromRGBO(200, 200, 200, 1.0),
      ),
      darkTheme: ThemeData(brightness: Brightness.dark),
      title: 'WSL Plus',
      home: const Scaffold(body: MainScreen()),
    );
  }
}

class MainScreen extends StatefulWidget {
  const MainScreen({super.key});

  @override
  State<MainScreen> createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  final List<String> _logs = [];
  final List<String> _distros = [];

  @override
  void initState() {
    super.initState();
    // Load distros on startup
    _loadDistros();
  }

  void _addLog(String message) {
    setState(() {
      _logs.add(message);
    });
  }

  Future<void> _loadDistros() async {
    _addLog('Loading distros...');
    try {
      final distros = await ApiService.getDistros();
      setState(() {
        _distros.clear();
        _distros.addAll(distros);
      });
      _addLog('Loaded ${distros.length} distro(s)');
    } catch (e) {
      _addLog('Error: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        // Sidebar with buttons
        Container(
          width: 200,
          decoration: BoxDecoration(
            color: Theme.of(context).colorScheme.surfaceContainerHighest,
            border: Border(
              right: BorderSide(
                color: Theme.of(context).colorScheme.outline,
                width: 1,
              ),
            ),
          ),
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Text(
                  'WSL Plus',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
              ),
              const Divider(height: 1),
              Expanded(
                child: ListView(
                  padding: const EdgeInsets.all(8.0),
                  children: [
                    _buildSectionHeader(context, 'WSL Actions'),
                    _buildButton('Refresh List', _loadDistros),
                    _buildButton('WSL Info', () => _addLog('WSL Info clicked')),
                    _buildButton('WSL Default', () => _addLog('WSL Default clicked')),
                    const SizedBox(height: 16),
                    _buildSectionHeader(context, 'Distro Actions'),
                    _buildButton('Install', () => _addLog('Install clicked')),
                    _buildButton('Rename', () => _addLog('Rename clicked')),
                    _buildButton('Backup', () => _addLog('Backup clicked')),
                  ],
                ),
              ),
            ],
          ),
        ),
        // Main content area
        Expanded(
          child: Column(
            children: [
              // Distro list area
              Expanded(
                child: _distros.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.computer,
                              size: 64,
                              color: Theme.of(context).colorScheme.secondary,
                            ),
                            const SizedBox(height: 16),
                            Text(
                              'No distros loaded',
                              style: Theme.of(context).textTheme.titleMedium,
                            ),
                            const SizedBox(height: 8),
                            Text(
                              'Click "Refresh List" to reload',
                              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                                  ),
                            ),
                          ],
                        ),
                      )
                    : ListView.builder(
                        padding: const EdgeInsets.all(16.0),
                        itemCount: _distros.length,
                        itemBuilder: (context, index) {
                          return Card(
                            margin: const EdgeInsets.only(bottom: 8.0),
                            child: ListTile(
                              leading: Icon(
                                Icons.computer,
                                color: Theme.of(context).colorScheme.primary,
                              ),
                              title: Text(_distros[index]),
                              trailing: PopupMenuButton(
                                itemBuilder: (context) => [
                                  const PopupMenuItem(
                                    value: 'info',
                                    child: Text('View Info'),
                                  ),
                                  const PopupMenuItem(
                                    value: 'rename',
                                    child: Text('Rename'),
                                  ),
                                  const PopupMenuItem(
                                    value: 'backup',
                                    child: Text('Backup'),
                                  ),
                                ],
                                onSelected: (value) {
                                  _addLog('$value: ${_distros[index]}');
                                },
                              ),
                            ),
                          );
                        },
                      ),
              ),
              // Log area at bottom
              Container(
                height: 150,
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.surfaceContainerHighest,
                  border: Border(
                    top: BorderSide(
                      color: Theme.of(context).colorScheme.outline,
                      width: 1,
                    ),
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Row(
                        children: [
                          Icon(
                            Icons.terminal,
                            size: 16,
                            color: Theme.of(context).colorScheme.secondary,
                          ),
                          const SizedBox(width: 8),
                          Text(
                            'Activity Log',
                            style: Theme.of(context).textTheme.labelLarge,
                          ),
                          const Spacer(),
                          if (_logs.isNotEmpty)
                            IconButton(
                              icon: const Icon(Icons.clear_all, size: 20),
                              onPressed: () {
                                setState(() {
                                  _logs.clear();
                                });
                              },
                              tooltip: 'Clear log',
                            ),
                        ],
                      ),
                    ),
                    const Divider(height: 1),
                    Expanded(
                      child: _logs.isEmpty
                          ? Center(
                              child: Text(
                                'No activity yet',
                                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                                    ),
                              ),
                            )
                          : ListView.builder(
                              padding: const EdgeInsets.all(8.0),
                              itemCount: _logs.length,
                              reverse: true,
                              itemBuilder: (context, index) {
                                return Padding(
                                  padding: const EdgeInsets.symmetric(vertical: 2.0),
                                  child: Text(
                                    _logs[_logs.length - 1 - index],
                                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                          fontFamily: 'monospace',
                                          color: Theme.of(context).colorScheme.onSurfaceVariant,
                                        ),
                                  ),
                                );
                              },
                            ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildSectionHeader(BuildContext context, String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 4.0),
      child: Text(
        title,
        style: Theme.of(context).textTheme.labelSmall?.copyWith(
              fontWeight: FontWeight.bold,
              color: Theme.of(context).colorScheme.primary,
            ),
      ),
    );
  }

  Widget _buildButton(String label, VoidCallback onPressed) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2.0),
      child: SizedBox(
        width: double.infinity,
        child: FilledButton.tonal(
          onPressed: onPressed,
          child: Text(label),
        ),
      ),
    );
  }
}
