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
  bool _isInstalling = false;
  String _currentlyInstalling = '';

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

  Future<void> _showDefaultDistro() async {
    _addLog('Getting default distro...');
    try {
      final defaultDistro = await ApiService.getDefaultDistro();
      _addLog('Default distro: $defaultDistro');
      if (mounted) {
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: const Text('Default WSL Distro'),
            content: Text(defaultDistro),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('OK'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      _addLog('Error: $e');
    }
  }

  Future<void> _showInstallDialog() async {
    _addLog('Loading available distros...');
    try {
      final available = await ApiService.getAvailableDistros();
      _addLog('Loaded ${available.length} available distros');
      
      if (!mounted) return;

      final selected = await showDialog<List<String>>(
        context: context,
        builder: (context) => _InstallDialog(available: available),
      );

      if (selected != null && selected.isNotEmpty) {
        _installDistros(selected);
      }
    } catch (e) {
      _addLog('Error: $e');
    }
  }

  Future<void> _installDistros(List<String> distros) async {
    setState(() {
      _isInstalling = true;
      _currentlyInstalling = distros.join(', ');
    });
    
    _addLog('Installing ${distros.length} distro(s): ${distros.join(", ")}');
    try {
      final results = await ApiService.installDistros(distros);
      
      for (var result in results) {
        if (result['success'] as bool) {
          _addLog('✓ ${result['distro']}: ${result['message']}');
        } else {
          _addLog('✗ ${result['distro']}: ${result['message']}');
        }
      }
      
      // Refresh the distro list
      _loadDistros();
    } catch (e) {
      _addLog('✗ Error: $e');
    } finally {
      setState(() {
        _isInstalling = false;
        _currentlyInstalling = '';
      });
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
                    _buildButton('WSL Default', _showDefaultDistro),
                    const SizedBox(height: 16),
                    _buildSectionHeader(context, 'Distro Actions'),
                    _buildButton('Install', _showInstallDialog),
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
                              itemCount: _logs.length + (_isInstalling ? 1 : 0),
                              reverse: true,
                              itemBuilder: (context, index) {
                                if (index == 0 && _isInstalling) {
                                  return Padding(
                                    padding: const EdgeInsets.symmetric(vertical: 2.0),
                                    child: Row(
                                      children: [
                                        SizedBox(
                                          width: 12,
                                          height: 12,
                                          child: CircularProgressIndicator(
                                            strokeWidth: 2,
                                            valueColor: AlwaysStoppedAnimation<Color>(
                                              Theme.of(context).colorScheme.primary,
                                            ),
                                          ),
                                        ),
                                        const SizedBox(width: 8),
                                        Expanded(
                                          child: Text(
                                            'Installing $_currentlyInstalling...',
                                            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                                  fontFamily: 'monospace',
                                                  color: Theme.of(context).colorScheme.primary,
                                                ),
                                          ),
                                        ),
                                      ],
                                    ),
                                  );
                                }
                                
                                final logIndex = _isInstalling ? index - 1 : index;
                                final log = _logs[_logs.length - 1 - logIndex];
                                Color? textColor = Theme.of(context).colorScheme.onSurfaceVariant;
                                
                                if (log.startsWith('✓')) {
                                  textColor = Colors.green;
                                } else if (log.startsWith('✗')) {
                                  textColor = Colors.red;
                                }
                                
                                return Padding(
                                  padding: const EdgeInsets.symmetric(vertical: 2.0),
                                  child: Text(
                                    log,
                                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                          fontFamily: 'monospace',
                                          color: textColor,
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

class _InstallDialog extends StatefulWidget {
  final List<Map<String, String>> available;

  const _InstallDialog({required this.available});

  @override
  State<_InstallDialog> createState() => _InstallDialogState();
}

class _InstallDialogState extends State<_InstallDialog> {
  final Set<String> _selected = {};
  String _searchQuery = '';

  @override
  Widget build(BuildContext context) {
    final filtered = widget.available.where((d) {
      if (_searchQuery.isEmpty) return true;
      final name = d['name']!.toLowerCase();
      final friendly = d['friendlyName']!.toLowerCase();
      final query = _searchQuery.toLowerCase();
      return name.contains(query) || friendly.contains(query);
    }).toList();

    return AlertDialog(
      title: const Text('Install WSL Distros'),
      content: SizedBox(
        width: 500,
        height: 400,
        child: Column(
          children: [
            TextField(
              decoration: const InputDecoration(
                labelText: 'Search',
                prefixIcon: Icon(Icons.search),
              ),
              onChanged: (value) {
                setState(() {
                  _searchQuery = value;
                });
              },
            ),
            const SizedBox(height: 16),
            Expanded(
              child: ListView.builder(
                itemCount: filtered.length,
                itemBuilder: (context, index) {
                  final distro = filtered[index];
                  final name = distro['name']!;
                  final friendlyName = distro['friendlyName']!;
                  
                  return CheckboxListTile(
                    title: Text(friendlyName),
                    subtitle: Text(name),
                    value: _selected.contains(name),
                    onChanged: (checked) {
                      setState(() {
                        if (checked == true) {
                          _selected.add(name);
                        } else {
                          _selected.remove(name);
                        }
                      });
                    },
                  );
                },
              ),
            ),
            if (_selected.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(top: 8.0),
                child: Text(
                  '${_selected.length} distro(s) selected',
                  style: Theme.of(context).textTheme.bodySmall,
                ),
              ),
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('Cancel'),
        ),
        FilledButton(
          onPressed: _selected.isEmpty
              ? null
              : () => Navigator.pop(context, _selected.toList()),
          child: const Text('Install'),
        ),
      ],
    );
  }
}
