import 'dart:async';
import 'dart:io';
import 'package:flutter/material.dart';
import 'package:yaru/yaru.dart';
import 'services/api_service.dart';

void main() {
  runApp(const MainApp());
}

class MainApp extends StatelessWidget {
  const MainApp({super.key});

  @override
  Widget build(BuildContext context) {
    return YaruTheme(
      builder: (context, yaru, child) {
        return MaterialApp(
          theme: yaru.theme,
          darkTheme: yaru.darkTheme,
          debugShowCheckedModeBanner: false,
          title: 'WSL Plus',
          home: const Scaffold(body: MainScreen()),
        );
      },
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
  final List<Map<String, dynamic>> _distros = [];
  String _defaultDistro = '';
  bool _isInstalling = false;
  String _currentlyInstalling = '';
  Timer? _refreshTimer;

  @override
  void initState() {
    super.initState();
    // Load distros and default on startup
    _loadDistros();

    // Set up periodic refresh every 5 seconds
    _refreshTimer = Timer.periodic(const Duration(seconds: 5), (timer) {
      _loadDistros();
    });
  }

  @override
  void dispose() {
    _refreshTimer?.cancel();
    super.dispose();
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
      final defaultDistro = await ApiService.getDefaultDistro();
      setState(() {
        _distros.clear();
        _distros.addAll(distros);
        _defaultDistro = defaultDistro;
      });
      final runningCount = distros.where((d) => d['running'] as bool).length;
      _addLog('Loaded ${distros.length} distro(s), $runningCount running, default: "$defaultDistro"');
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

  Future<void> _setDefaultDistro(String name) async {
    _addLog('Setting $name as default...');
    try {
      await ApiService.setDefaultDistro(name);
      _addLog('✓ Set $name as default');
      _loadDistros();
    } catch (e) {
      _addLog('✗ Error: $e');
    }
  }

  Future<void> _unregisterDistro(String name) async {
    // Confirm deletion
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Confirm Delete'),
        content: Text('Are you sure you want to unregister $name?\n'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      _deleteDistros([name]);
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

  Future<void> _terminateDistros(List<String> distros) async {
    _addLog('Stopping ${distros.length} distro(s): ${distros.join(", ")}');
    try {
      final results = await ApiService.terminateDistros(distros);

      for (var result in results) {
        if (result['success'] as bool) {
          _addLog('✓ ${result['distro']}: ${result['message']}');
        } else {
          _addLog('✗ ${result['distro']}: ${result['message']}');
        }
      }

      // Refresh distro list to update running states
      _loadDistros();
    } catch (e) {
      _addLog('✗ Error: $e');
    }
  }

  Future<void> _launchDistro(String distro) async {
    _addLog('Launching $distro in new terminal...');
    try {
      await ApiService.launchDistro(distro);
      _addLog('✓ Launched $distro in new terminal');

      // Wait a moment for the distro to start, then refresh
      await Future.delayed(const Duration(milliseconds: 500));
      _loadDistros();
    } catch (e) {
      _addLog('✗ Error: $e');
    }
  }

  Future<void> _showRenameDialog(String oldName) async {
    final controller = TextEditingController();

    final newName = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('Rename $oldName'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: controller,
              decoration: const InputDecoration(
                labelText: 'New name',
                hintText: 'Enter new distribution name',
              ),
              autofocus: true,
            ),
            const SizedBox(height: 16),
            Text(
              'Note: You will need to restart WSL for changes to take effect.',
              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () {
              final name = controller.text.trim();
              if (name.isNotEmpty) {
                Navigator.pop(context, name);
              }
            },
            child: const Text('Rename'),
          ),
        ],
      ),
    );

    if (newName != null && newName.isNotEmpty) {
      _renameDistro(oldName, newName);
    }
  }

  Future<void> _renameDistro(String oldName, String newName) async {
    _addLog('Renaming $oldName to $newName...');
    try {
      final result = await ApiService.renameDistro(oldName, newName);

      if (result['success'] as bool) {
        _addLog('✓ ${result['message']}');

        // Show restart WSL reminder
        if (mounted) {
          showDialog(
            context: context,
            builder: (context) => AlertDialog(
              title: const Text('Restart WSL'),
              content: const Text(
                'The distribution has been renamed successfully.\n\n'
                'To apply the changes, please restart WSL.',
              ),
              actions: [
                TextButton(
                  onPressed: () => Navigator.pop(context),
                  child: const Text('Later'),
                ),
                FilledButton(
                  onPressed: () {
                    Navigator.pop(context);
                    _shutdownWSL();
                  },
                  child: const Text('Shutdown WSL Now'),
                ),
              ],
            ),
          );
        }

        // Refresh list to show new name
        _loadDistros();
      } else {
        _addLog('✗ ${result['message']}');
      }
    } catch (e) {
      _addLog('✗ Error: $e');
    }
  }

  Future<void> _shutdownWSL() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Shutdown WSL'),
        content: const Text(
          'This will terminate all running WSL distributions.\n\n'
          'Are you sure you want to continue?'
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('Shutdown'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      _addLog('Shutting down WSL...');
      try {
        final result = await Process.run('wsl', ['--shutdown']);
        if (result.exitCode == 0) {
          _addLog('✓ WSL shutdown successfully');
          // Refresh list after a short delay to show updated states
          await Future.delayed(const Duration(milliseconds: 500));
          _loadDistros();
        } else {
          _addLog('✗ Failed to shutdown WSL: ${result.stderr}');
        }
      } catch (e) {
        _addLog('✗ Error: $e');
      }
    }
  }

  Future<void> _showWSLInfo() async {
    try {
      final info = await ApiService.getWSLInfo();

      if (mounted) {
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: const Text('WSL System Information'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildInfoRow('Default WSL Version', 'WSL ${info['defaultWslVersion']}'),
                _buildInfoRow('Number of Distros', '${info['numDistros']}'),
                _buildInfoRow('Default Distro', info['defaultDistro'] ?? 'None'),
                _buildInfoRow('Total Disk Usage', info['totalDiskUsage']),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Close'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      _addLog('✗ Error getting WSL info: $e');
    }
  }

  Future<void> _showDistroInfo(String distro) async {
    try {
      final info = await ApiService.getDistroInfo(distro);

      if (mounted) {
        showDialog(
          context: context,
          builder: (context) => AlertDialog(
            title: Text('$distro Information'),
            content: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildInfoRow('Name', info['name']),
                  _buildInfoRow('WSL Version', 'WSL ${info['wslVersion']}'),
                  _buildInfoRow('State', info['state']),
                  _buildInfoRow('Is Default', info['isDefault'] ? 'Yes' : 'No'),
                  _buildInfoRow('Flavor', info['flavor'] ?? 'Unknown'),
                  const Divider(),
                  _buildInfoRow('GUID', info['guid']),
                  _buildInfoRow('Default UID', '${info['defaultUid']}'),
                  _buildInfoRow('Interop Enabled', info['interopEnabled'] ? 'Yes' : 'No'),
                  _buildInfoRow('Drive Mounting', info['driveMounting'] ? 'Yes' : 'No'),
                  _buildInfoRow('Path Appended', info['pathAppended'] ? 'Yes' : 'No'),
                  if (info['isUbuntu'] == true) ...[
                    const Divider(),
                    _buildInfoRow(
                      'Ubuntu Telemetry',
                      info['telemetryEnabled'] == true ? 'Enabled' : 'Disabled',
                    ),
                  ],
                ],
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Close'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      _addLog('✗ Error getting distro info: $e');
    }
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 140,
            child: Text(
              '$label:',
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          Expanded(
            child: Text(value),
          ),
        ],
      ),
    );
  }

  Future<void> _checkAndBackupDistros(List<String> distros) async {
    // Check if any of the distros are running
    final runningDistros = <String>[];
    for (final distroName in distros) {
      final distroData = _distros.firstWhere(
        (d) => d['name'] == distroName,
        orElse: () => <String, dynamic>{},
      );
      if (distroData.isNotEmpty && distroData['running'] == true) {
        runningDistros.add(distroName);
      }
    }

    if (runningDistros.isNotEmpty) {
      // Prompt user to stop running distros first
      final shouldStop = await showDialog<bool>(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('Running Distros Detected'),
          content: Text(
            'The following distro(s) are currently running:\n\n${runningDistros.join(', ')}\n\n'
            'For a consistent backup, distros should be stopped first. Would you like to stop them now?',
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text('Cancel Backup'),
            ),
            FilledButton(
              onPressed: () => Navigator.pop(context, true),
              child: const Text('Stop & Backup'),
            ),
          ],
        ),
      );

      if (shouldStop == true) {
        // Stop the running distros first
        await _terminateDistros(runningDistros);
        // Wait a moment for them to fully stop
        await Future.delayed(const Duration(seconds: 1));
        // Then proceed with backup
        _backupDistros(distros);
      }
    } else {
      // No running distros, proceed with backup
      _backupDistros(distros);
    }
  }

  Future<void> _backupDistros(List<String> distros) async {
    setState(() {
      _isInstalling = true;
      _currentlyInstalling = 'Backing up: ${distros.join(', ')}';
    });

    _addLog('Backing up ${distros.length} distro(s): ${distros.join(", ")}');
    try {
      final results = await ApiService.backupDistros(distros);

      for (var result in results) {
        if (result['success'] as bool) {
          _addLog('✓ ${result['distro']}: ${result['message']}');
          _addLog('  Saved to: ${result['filePath']}');
        } else {
          _addLog('✗ ${result['distro']}: ${result['message']}');
        }
      }
    } catch (e) {
      _addLog('✗ Error: $e');
    } finally {
      setState(() {
        _isInstalling = false;
        _currentlyInstalling = '';
      });
    }
  }

  Future<void> _showBackupDialog() async {
    if (_distros.isEmpty) {
      _addLog('No distros available to backup');
      return;
    }

    final distroNames = _distros.map((d) => d['name'] as String).toList();

    final selected = await showDialog<List<String>>(
      context: context,
      builder: (context) => _BackupDialog(distros: distroNames),
    );

    if (selected != null && selected.isNotEmpty) {
      _checkAndBackupDistros(selected);
    }
  }

  Future<void> _showDeleteDialog() async {
    if (_distros.isEmpty) {
      _addLog('No distros available to delete');
      return;
    }

    final distroNames = _distros.map((d) => d['name'] as String).toList();

    final selected = await showDialog<List<String>>(
      context: context,
      builder: (context) => _DeleteDialog(distros: distroNames),
    );

    if (selected != null && selected.isNotEmpty) {
      _deleteDistros(selected);
    }
  }

  Future<void> _deleteDistros(List<String> distros) async {
    _addLog('Deleting ${distros.length} distro(s)...');
    try {
      final results = await ApiService.unregisterDistros(distros);

      for (final result in results) {
        if (result['success'] as bool) {
          _addLog('✓ ${result['distro']}: ${result['message']}');
        } else {
          _addLog('✗ ${result['distro']}: ${result['message']}');
        }
      }

      _loadDistros();
    } catch (e) {
      _addLog('✗ Error: $e');
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
            border: Border(
              right: BorderSide(
                color: Theme.of(context).colorScheme.outlineVariant,
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
                    _buildButton('Refresh List', _loadDistros),
                    _buildButton('WSL Info', _showWSLInfo),
                    _buildButton('WSL Shutdown', _shutdownWSL),
                    const SizedBox(height: 16),
                    _buildSectionHeader(context, 'Bulk Actions'),
                    _buildButton('Install', _showInstallDialog),
                    _buildButton('Backup', _showBackupDialog),
                    _buildButton('Delete', _showDeleteDialog),
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
                    : Builder(
                        builder: (context) {
                          // Sort so default is first (case-insensitive comparison)
                          final sortedDistros = List<Map<String, dynamic>>.from(_distros);
                          sortedDistros.sort((a, b) {
                            final aName = a['name'] as String;
                            final bName = b['name'] as String;
                            if (aName.toLowerCase() == _defaultDistro.toLowerCase()) return -1;
                            if (bName.toLowerCase() == _defaultDistro.toLowerCase()) return 1;
                            return aName.toLowerCase().compareTo(bName.toLowerCase());
                          });

                          return ListView.builder(
                            padding: const EdgeInsets.all(16.0),
                            itemCount: sortedDistros.length,
                            itemBuilder: (context, index) {
                              final distroData = sortedDistros[index];
                              final distro = distroData['name'] as String;
                              final isRunning = distroData['running'] as bool;
                              final isDefault = distro.toLowerCase() == _defaultDistro.toLowerCase();
                              
                              return Card(
                            margin: const EdgeInsets.only(bottom: 8.0),
                            elevation: isDefault ? 4 : 1,
                            color: isDefault 
                                ? Theme.of(context).colorScheme.primaryContainer
                                : null,
                            child: ListTile(
                              leading: Icon(
                                Icons.computer,
                                color: isDefault 
                                    ? Theme.of(context).colorScheme.primary
                                    : null,
                                size: isDefault ? 28 : 24,
                              ),
                              title: Row(
                                children: [
                                  Text(
                                    distro,
                                    style: isDefault
                                        ? const TextStyle(fontWeight: FontWeight.bold)
                                        : null,
                                  ),
                                  if (isDefault) ...[
                                    const SizedBox(width: 8),
                                    Container(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 8,
                                        vertical: 3,
                                      ),
                                      decoration: BoxDecoration(
                                        color: Theme.of(context).colorScheme.primary,
                                        borderRadius: BorderRadius.circular(6),
                                      ),
                                      child: Text(
                                        'DEFAULT',
                                        style: Theme.of(context).textTheme.labelSmall?.copyWith(
                                              color: Theme.of(context).colorScheme.onPrimary,
                                              fontWeight: FontWeight.bold,
                                              fontSize: 10,
                                            ),
                                      ),
                                    ),
                                  ],
                                  if (isRunning) ...[
                                    const SizedBox(width: 8),
                                    Container(
                                      padding: const EdgeInsets.symmetric(
                                        horizontal: 8,
                                        vertical: 3,
                                      ),
                                      decoration: BoxDecoration(
                                        color: Colors.green,
                                        borderRadius: BorderRadius.circular(6),
                                      ),
                                      child: Row(
                                        mainAxisSize: MainAxisSize.min,
                                        children: [
                                          const Icon(Icons.circle, size: 8, color: Colors.white),
                                          const SizedBox(width: 4),
                                          Text(
                                            'RUNNING',
                                            style: Theme.of(context).textTheme.labelSmall?.copyWith(
                                                  color: Colors.white,
                                                  fontWeight: FontWeight.bold,
                                                  fontSize: 10,
                                                ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                              onTap: () => _launchDistro(distro),
                              trailing: PopupMenuButton<String>(
                                itemBuilder: (context) {
                                  final items = <PopupMenuEntry<String>>[
                                    const PopupMenuItem(
                                      value: 'launch',
                                      child: Text('Launch Terminal'),
                                    ),
                                    const PopupMenuDivider(),
                                    const PopupMenuItem(
                                      value: 'info',
                                      child: Text('View Info'),
                                    ),
                                    if (!isDefault)
                                      const PopupMenuItem(
                                        value: 'set-default',
                                        child: Text('Set as Default'),
                                      ),
                                    if (isRunning)
                                      const PopupMenuItem(
                                        value: 'stop',
                                        child: Text('Stop'),
                                      ),
                                    const PopupMenuItem(
                                      value: 'rename',
                                      child: Text('Rename'),
                                    ),
                                    const PopupMenuItem(
                                      value: 'backup',
                                      child: Text('Backup'),
                                    ),
                                    const PopupMenuDivider(),
                                    const PopupMenuItem(
                                      value: 'delete',
                                      child: Text('Delete'),
                                    ),
                                  ];
                                  return items;
                                },
                                onSelected: (value) {
                                  switch (value) {
                                    case 'launch':
                                      _launchDistro(distro);
                                      break;
                                    case 'info':
                                      _showDistroInfo(distro);
                                      break;
                                    case 'set-default':
                                      _setDefaultDistro(distro);
                                      break;
                                    case 'stop':
                                      _terminateDistros([distro]);
                                      break;
                                    case 'rename':
                                      _showRenameDialog(distro);
                                      break;
                                    case 'backup':
                                      _checkAndBackupDistros([distro]);
                                      break;
                                    case 'delete':
                                      _unregisterDistro(distro);
                                      break;
                                  }
                                },
                              ),
                            ),
                          );
                        },
                      );
                    },
                  ),
              ),
              // Log area at bottom
              Container(
                height: 150,
                decoration: BoxDecoration(
                  border: Border(
                    top: BorderSide(
                      color: Theme.of(context).colorScheme.outlineVariant,
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

  Widget _buildButton(String label, VoidCallback? onPressed) {
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

class _BackupDialog extends StatefulWidget {
  final List<String> distros;

  const _BackupDialog({required this.distros});

  @override
  State<_BackupDialog> createState() => _BackupDialogState();
}

class _BackupDialogState extends State<_BackupDialog> {
  final Set<String> _selected = {};
  String _searchQuery = '';

  @override
  Widget build(BuildContext context) {
    final filtered = widget.distros.where((d) {
      if (_searchQuery.isEmpty) return true;
      return d.toLowerCase().contains(_searchQuery.toLowerCase());
    }).toList();

    return AlertDialog(
      title: const Text('Backup WSL Distros'),
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

                  return CheckboxListTile(
                    title: Text(distro),
                    value: _selected.contains(distro),
                    onChanged: (checked) {
                      setState(() {
                        if (checked == true) {
                          _selected.add(distro);
                        } else {
                          _selected.remove(distro);
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
          child: const Text('Backup'),
        ),
      ],
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

class _DeleteDialog extends StatefulWidget {
  final List<String> distros;

  const _DeleteDialog({required this.distros});

  @override
  State<_DeleteDialog> createState() => _DeleteDialogState();
}

class _DeleteDialogState extends State<_DeleteDialog> {
  final Set<String> _selected = {};
  String _searchQuery = '';

  @override
  Widget build(BuildContext context) {
    final filtered = widget.distros.where((d) {
      if (_searchQuery.isEmpty) return true;
      return d.toLowerCase().contains(_searchQuery.toLowerCase());
    }).toList();

    return AlertDialog(
      title: const Text('Delete WSL Distros'),
      content: SizedBox(
        width: 500,
        height: 400,
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.errorContainer,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(
                children: [
                  Icon(
                    Icons.warning,
                    color: Theme.of(context).colorScheme.onErrorContainer,
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      'WARNING: This will permanently delete the selected distributions and all their data.',
                      style: TextStyle(
                        color: Theme.of(context).colorScheme.onErrorContainer,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 16),
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

                  return CheckboxListTile(
                    title: Text(distro),
                    value: _selected.contains(distro),
                    onChanged: (checked) {
                      setState(() {
                        if (checked == true) {
                          _selected.add(distro);
                        } else {
                          _selected.remove(distro);
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
          style: FilledButton.styleFrom(
            backgroundColor: Theme.of(context).colorScheme.error,
            foregroundColor: Theme.of(context).colorScheme.onError,
          ),
          child: const Text('Delete'),
        ),
      ],
    );
  }
}
